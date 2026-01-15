#!/usr/bin/env python3
"""SessionStart Hook: Enforce Do Directive + Worker Context Integration

Forces Claude to follow CLAUDE.md Do directive from session start.
Fetches context from Worker service for token-efficient session continuity.

Supports:
- DO_CONTEXT_LEVEL: Context detail level (1=minimal, 2=standard, 3=full)
- Previously section: Reads last response from cache for continuity
"""
import json
import sys
import os
import time
import subprocess
import urllib.request
import urllib.error
import urllib.parse

WORKER_PORT = int(os.environ.get("DO_WORKER_PORT", "3778"))
WORKER_URL = f"http://127.0.0.1:{WORKER_PORT}"


def get_context_level() -> int:
    """Get context level from environment (1=minimal, 2=standard, 3=full)."""
    level_str = os.environ.get("DO_CONTEXT_LEVEL", "2")
    try:
        level = int(level_str)
        if 1 <= level <= 3:
            return level
    except ValueError:
        pass
    return 2  # Default to standard


def get_last_response_summary() -> str:
    """Read cached last assistant response for session continuity."""
    cache_path = os.path.join(os.getcwd(), ".do/cache/last_response.txt")
    if os.path.exists(cache_path):
        try:
            with open(cache_path, 'r', encoding='utf-8') as f:
                content = f.read().strip()
                # Limit to 500 chars to keep context efficient
                return content[:500] if len(content) > 500 else content
        except Exception:
            pass
    return ""


def ensure_worker_running() -> bool:
    """Check if Worker is running, start if not."""
    try:
        req = urllib.request.Request(f"{WORKER_URL}/health")
        with urllib.request.urlopen(req, timeout=1) as resp:
            return resp.status == 200
    except Exception:
        # Worker not running, try to start it
        worker_path = os.path.join(os.getcwd(), ".do/bin/do-worker")
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


def register_session(session_id: str, project_path: str, user_name: str) -> bool:
    """Register session with Worker service."""
    if not session_id:
        return False

    try:
        data = json.dumps({
            "id": session_id,
            "user_name": user_name or "unknown",
            "project_id": project_path,
        }).encode("utf-8")

        req = urllib.request.Request(
            f"{WORKER_URL}/api/sessions",
            data=data,
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        with urllib.request.urlopen(req, timeout=5) as resp:
            return resp.status in (200, 201)
    except Exception:
        pass
    return False


def get_context_from_worker(session_id: str, project_path: str, user_name: str) -> str:
    """Fetch compressed context from Worker service.

    Uses DO_CONTEXT_LEVEL for detail level and includes previously section
    from cached last response for session continuity.
    """
    try:
        level = get_context_level()
        previously = get_last_response_summary()

        params = urllib.parse.urlencode({
            "session_id": session_id,
            "project_path": project_path,
            "user": user_name,
            "level": level,
            "previously": previously,
        })

        req = urllib.request.Request(f"{WORKER_URL}/api/context/inject?{params}")
        with urllib.request.urlopen(req, timeout=5) as resp:
            if resp.status == 200:
                return json.loads(resp.read().decode("utf-8")).get("context", "")
    except Exception:
        pass
    return ""


def main():
    # Read hook input from stdin
    hook_input = {}
    try:
        import select
        if select.select([sys.stdin], [], [], 0)[0]:
            hook_input = json.loads(sys.stdin.read())
    except Exception:
        pass

    # Get session_id from hook input or environment
    session_id = hook_input.get("session_id", os.environ.get("CLAUDE_SESSION_ID", ""))
    project_path = os.getcwd()
    user_name = os.environ.get("DO_USER_NAME", "")

    # Try to register session and get context from Worker
    worker_context = ""
    if ensure_worker_running():
        register_session(session_id, project_path, user_name)
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
