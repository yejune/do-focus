#!/usr/bin/env python3
"""
Background Task Tracker Utility

Provides functions to start/stop tracking and register tasks.
Used by main agent before launching background tasks.
"""

import json
import os
import uuid
from datetime import datetime
from pathlib import Path
from typing import Optional


def get_project_root() -> Path:
    """Get project root from environment or current directory."""
    if "CLAUDE_PROJECT_DIR" in os.environ:
        return Path(os.environ["CLAUDE_PROJECT_DIR"])
    return Path.cwd()


PROJECT_ROOT = get_project_root()
CACHE_DIR = PROJECT_ROOT / ".do" / "cache"
STATE_FILE = CACHE_DIR / "background-tasks.json"


def start_tracking(total_tasks: int, session_id: Optional[str] = None) -> str:
    """
    Start a new background task tracking session.

    Args:
        total_tasks: Total number of tasks to track
        session_id: Optional session ID (auto-generated if not provided)

    Returns:
        Session ID for this tracking session
    """
    CACHE_DIR.mkdir(parents=True, exist_ok=True)

    session_id = session_id or str(uuid.uuid4())[:8]

    state = {
        "active": True,
        "session_id": session_id,
        "total": total_tasks,
        "completed": 0,
        "started_at": datetime.now().isoformat(),
        "tasks": []
    }

    with open(STATE_FILE, "w", encoding="utf-8") as f:
        json.dump(state, f, indent=2, ensure_ascii=False)

    return session_id


def register_task(agent_id: str, description: str) -> None:
    """
    Register a task that was just launched.

    Args:
        agent_id: The agent/task ID returned by Task tool
        description: Short description of the task
    """
    if not STATE_FILE.exists():
        return

    with open(STATE_FILE, "r", encoding="utf-8") as f:
        state = json.load(f)

    state["tasks"].append({
        "agent_id": agent_id,
        "description": description,
        "status": "running",
        "started_at": datetime.now().isoformat()
    })

    with open(STATE_FILE, "w", encoding="utf-8") as f:
        json.dump(state, f, indent=2, ensure_ascii=False)


def stop_tracking() -> dict:
    """
    Stop tracking and return final state.

    Returns:
        Final state with all task results
    """
    if not STATE_FILE.exists():
        return {}

    with open(STATE_FILE, "r", encoding="utf-8") as f:
        state = json.load(f)

    state["active"] = False
    state["stopped_at"] = datetime.now().isoformat()

    with open(STATE_FILE, "w", encoding="utf-8") as f:
        json.dump(state, f, indent=2, ensure_ascii=False)

    return state


def get_status() -> dict:
    """
    Get current tracking status.

    Returns:
        Current state dict
    """
    if not STATE_FILE.exists():
        return {"active": False}

    with open(STATE_FILE, "r", encoding="utf-8") as f:
        return json.load(f)


def is_tracking_active() -> bool:
    """Check if tracking is currently active."""
    status = get_status()
    return status.get("active", False)


# CLI interface for testing
if __name__ == "__main__":
    import sys

    if len(sys.argv) < 2:
        print("Usage: background_task_tracker.py <command> [args]")
        print("Commands:")
        print("  start <total_tasks>  - Start tracking session")
        print("  register <id> <desc> - Register a task")
        print("  status               - Get current status")
        print("  stop                 - Stop tracking")
        sys.exit(1)

    command = sys.argv[1]

    if command == "start" and len(sys.argv) >= 3:
        total = int(sys.argv[2])
        session_id = start_tracking(total)
        print(f"Tracking started: session={session_id}, total={total}")

    elif command == "register" and len(sys.argv) >= 4:
        register_task(sys.argv[2], sys.argv[3])
        print(f"Task registered: {sys.argv[2]}")

    elif command == "status":
        status = get_status()
        print(json.dumps(status, indent=2, ensure_ascii=False))

    elif command == "stop":
        state = stop_tracking()
        print("Tracking stopped")
        print(json.dumps(state, indent=2, ensure_ascii=False))

    else:
        print(f"Unknown command: {command}")
        sys.exit(1)
