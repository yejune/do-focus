#!/usr/bin/env python3
"""
SubagentStop Hook: Background Task Tracker

Automatically tracks background task completion and outputs summary when all done.

Trigger: Fires when any subagent completes
State: Stored in .do/cache/background-tasks.json
"""

import json
import os
import sys
from datetime import datetime
from pathlib import Path

# Project root detection
def get_project_root() -> Path:
    """Get project root from environment or current directory."""
    if "CLAUDE_PROJECT_DIR" in os.environ:
        return Path(os.environ["CLAUDE_PROJECT_DIR"])
    return Path.cwd()

PROJECT_ROOT = get_project_root()
CACHE_DIR = PROJECT_ROOT / ".do" / "cache"
STATE_FILE = CACHE_DIR / "background-tasks.json"


def load_state() -> dict:
    """Load current task tracking state."""
    if not STATE_FILE.exists():
        return {"active": False, "session_id": None, "total": 0, "completed": 0, "tasks": []}

    try:
        with open(STATE_FILE, "r", encoding="utf-8") as f:
            return json.load(f)
    except (json.JSONDecodeError, IOError):
        return {"active": False, "session_id": None, "total": 0, "completed": 0, "tasks": []}


def save_state(state: dict) -> None:
    """Save task tracking state."""
    CACHE_DIR.mkdir(parents=True, exist_ok=True)
    with open(STATE_FILE, "w", encoding="utf-8") as f:
        json.dump(state, f, indent=2, ensure_ascii=False)


def parse_hook_input() -> dict:
    """Parse input from Claude Code hook system."""
    try:
        input_data = sys.stdin.read()
        if input_data:
            return json.loads(input_data)
    except (json.JSONDecodeError, IOError):
        pass
    return {}


def format_duration(seconds: float) -> str:
    """Format duration in human readable form."""
    if seconds < 60:
        return f"{seconds:.1f}s"
    elif seconds < 3600:
        minutes = int(seconds // 60)
        secs = int(seconds % 60)
        return f"{minutes}m {secs}s"
    else:
        hours = int(seconds // 3600)
        minutes = int((seconds % 3600) // 60)
        return f"{hours}h {minutes}m"


def output_progress(state: dict) -> None:
    """Output progress update."""
    completed = state["completed"]
    total = state["total"]
    percentage = (completed / total * 100) if total > 0 else 0

    # Progress bar
    bar_width = 20
    filled = int(bar_width * completed / total) if total > 0 else 0
    bar = "█" * filled + "░" * (bar_width - filled)

    print(f"\n[{bar}] {completed}/{total} ({percentage:.0f}%)", file=sys.stderr)


def output_summary(state: dict) -> None:
    """Output final summary when all tasks complete."""
    total = state["total"]
    tasks = state["tasks"]

    # Calculate stats
    successful = len([t for t in tasks if t.get("status") == "completed"])
    failed = len([t for t in tasks if t.get("status") == "failed"])

    total_time = sum(t.get("duration", 0) for t in tasks)

    print("\n" + "=" * 50, file=sys.stderr)
    print("Background Tasks Complete", file=sys.stderr)
    print("=" * 50, file=sys.stderr)
    print(f"Total: {total}", file=sys.stderr)
    print(f"Success: {successful}", file=sys.stderr)
    if failed > 0:
        print(f"Failed: {failed}", file=sys.stderr)
    print(f"Duration: {format_duration(total_time)}", file=sys.stderr)
    print("=" * 50, file=sys.stderr)

    # List failed tasks if any
    if failed > 0:
        print("\nFailed Tasks:", file=sys.stderr)
        for task in tasks:
            if task.get("status") == "failed":
                print(f"  - {task.get('description', 'Unknown')}", file=sys.stderr)


def main():
    """Main hook entry point."""
    # Parse hook input
    hook_input = parse_hook_input()

    # Load current state
    state = load_state()

    # Check if tracking is active
    if not state.get("active"):
        return  # No active tracking session

    # Extract task info from hook input
    agent_id = hook_input.get("agent_id", "unknown")
    status = hook_input.get("status", "completed")
    description = hook_input.get("description", "")

    # Find and update task
    task_found = False
    for task in state["tasks"]:
        if task.get("agent_id") == agent_id:
            task["status"] = status
            task["completed_at"] = datetime.now().isoformat()
            if task.get("started_at"):
                start = datetime.fromisoformat(task["started_at"])
                task["duration"] = (datetime.now() - start).total_seconds()
            task_found = True
            break

    # If task not found, add it (for tasks started before tracking)
    if not task_found:
        state["tasks"].append({
            "agent_id": agent_id,
            "description": description,
            "status": status,
            "completed_at": datetime.now().isoformat(),
            "duration": 0
        })

    # Update completed count
    state["completed"] = len([t for t in state["tasks"] if t.get("status") in ["completed", "failed"]])

    # Save state
    save_state(state)

    # Output progress
    output_progress(state)

    # Check if all tasks complete
    if state["completed"] >= state["total"]:
        output_summary(state)
        # Reset state
        state["active"] = False
        save_state(state)


if __name__ == "__main__":
    main()
