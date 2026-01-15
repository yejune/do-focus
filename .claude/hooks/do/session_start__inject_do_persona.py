#!/usr/bin/env python3
"""SessionStart Hook: Enforce Do Directive + Worker Context Integration

Forces Claude to follow CLAUDE.md Do directive from session start.
Fetches context from Worker service for token-efficient session continuity.
"""
import json
import sys
import os
import time
import subprocess

# Optional: requests for HTTP calls
try:
    import requests
    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False

WORKER_PORT = int(os.environ.get("DO_WORKER_PORT", "3778"))
WORKER_URL = f"http://127.0.0.1:{WORKER_PORT}"


def ensure_worker_running() -> bool:
    """Check if Worker is running, start if not."""
    if not HAS_REQUESTS:
        return False

    try:
        resp = requests.get(f"{WORKER_URL}/health", timeout=1)
        return resp.status_code == 200
    except Exception:
        # Worker not running, try to start it
        worker_path = os.path.join(os.getcwd(), ".do/worker/bin/do-worker")
        if os.path.exists(worker_path):
            try:
                subprocess.Popen(
                    [worker_path],
                    stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL,
                    start_new_session=True
                )
                time.sleep(0.5)  # Wait for startup
                return True
            except Exception:
                pass
        return False


def get_context_from_worker(session_id: str, project_path: str, user_name: str) -> str:
    """Fetch compressed context from Worker service."""
    if not HAS_REQUESTS:
        return ""

    try:
        resp = requests.get(
            f"{WORKER_URL}/api/context/inject",
            params={
                "session_id": session_id,
                "project_path": project_path,
                "user_name": user_name,
                "level": 1  # Progressive disclosure level
            },
            timeout=5
        )
        if resp.status_code == 200:
            return resp.json().get("context", "")
    except Exception:
        pass
    return ""


def main():
    # Get environment variables
    session_id = os.environ.get("CLAUDE_SESSION_ID", "")
    project_path = os.getcwd()
    user_name = os.environ.get("DO_USER_NAME", "")

    # Try to get context from Worker
    worker_context = ""
    if ensure_worker_running():
        worker_context = get_context_from_worker(session_id, project_path, user_name)

    # Base enforcement message
    base_message = """
DO DIRECTIVE ENFORCEMENT

YOU ARE DO - THE STRATEGIC ORCHESTRATOR

MANDATORY RULES:
- ALL implementation work -> delegate via Task tool
- FORBIDDEN direct use: Read, Write, Edit, Bash, Grep, Glob, WebFetch, WebSearch
- Start responses with "[Do]"
- Execute independent tasks in parallel

VIOLATIONS:
- Direct tool use without delegation = VIOLATION
- Code writing without Task() = VIOLATION

Follow CLAUDE.md specification in project root.
"""

    # Combine base message with worker context
    if worker_context:
        system_message = f"{base_message.strip()}\n\n---\n\n{worker_context}"
    else:
        system_message = base_message.strip()

    # Return result
    result = {
        "continue": True,
        "systemMessage": system_message
    }

    print(json.dumps(result, ensure_ascii=False))
    sys.exit(0)


if __name__ == "__main__":
    main()
