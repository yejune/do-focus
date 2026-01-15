#!/usr/bin/env python3
"""
Stop Hook: Generate session summary from transcript
Triggered when Claude Code session stops.
"""
import json
import os
import sys
import re
import urllib.request
import urllib.error
from pathlib import Path

WORKER_PORT = int(os.environ.get("DO_WORKER_PORT", "3778"))
WORKER_URL = f"http://127.0.0.1:{WORKER_PORT}"


def strip_system_reminders(text: str) -> str:
    """Remove <system-reminder> tags from text"""
    text = re.sub(r'<system-reminder>[\s\S]*?</system-reminder>', '', text)
    text = re.sub(r'\n{3,}', '\n\n', text).strip()
    return text


def parse_transcript_jsonl(transcript_path: str) -> dict:
    """Parse transcript JSONL file and extract assistant messages

    Returns:
        {
            "assistant_messages": List of text responses,
            "tool_uses": List of tool names,
            "total_entries": int
        }
    """
    result = {
        "assistant_messages": [],
        "tool_uses": [],
        "total_entries": 0
    }

    if not transcript_path or not Path(transcript_path).exists():
        return result

    try:
        with open(transcript_path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                try:
                    entry = json.loads(line)
                    result["total_entries"] += 1

                    if entry.get("type") == "assistant":
                        message = entry.get("message", {})
                        content = message.get("content", [])

                        if isinstance(content, list):
                            for item in content:
                                if item.get("type") == "text":
                                    text = strip_system_reminders(item.get("text", ""))
                                    if text.strip():
                                        result["assistant_messages"].append(text)
                                elif item.get("type") == "tool_use":
                                    result["tool_uses"].append(item.get("name", ""))
                        elif isinstance(content, str):
                            text = strip_system_reminders(content)
                            if text.strip():
                                result["assistant_messages"].append(text)
                except json.JSONDecodeError:
                    continue
    except Exception:
        pass

    return result


def get_last_assistant_message(transcript_path: str) -> str:
    """Get the last assistant text message"""
    parsed = parse_transcript_jsonl(transcript_path)
    if parsed["assistant_messages"]:
        return parsed["assistant_messages"][-1]
    return ""


def get_last_response_with_tools(transcript_path: str) -> str:
    """Get the last assistant response including tool_use blocks.

    Returns:
        Full response with tool_use info (max 100KB)
    """
    if not transcript_path or not Path(transcript_path).exists():
        return ""

    last_response = ""
    try:
        with open(transcript_path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                try:
                    entry = json.loads(line)
                    if entry.get("type") == "assistant":
                        message = entry.get("message", {})
                        content = message.get("content", [])
                        parts = []

                        if isinstance(content, list):
                            for item in content:
                                if item.get("type") == "text":
                                    text = strip_system_reminders(item.get("text", ""))
                                    if text.strip():
                                        parts.append(text)
                                elif item.get("type") == "tool_use":
                                    tool_name = item.get("name", "unknown")
                                    tool_input = item.get("input", {})
                                    input_str = json.dumps(tool_input, ensure_ascii=False)
                                    if len(input_str) > 500:
                                        input_str = input_str[:500] + "..."
                                    parts.append(f"\n[Tool: {tool_name}]\n{input_str}")

                        if parts:
                            last_response = "\n".join(parts)
                except json.JSONDecodeError:
                    continue
    except Exception:
        pass

    return last_response[:100000] if last_response else ""


def update_prompt_response(session_id: str, response: str) -> bool:
    """Update the latest user_prompt with the assistant response."""
    if not response or not session_id:
        return False

    try:
        data = json.dumps({
            "session_id": session_id,
            "response": response
        }).encode("utf-8")

        req = urllib.request.Request(
            f"{WORKER_URL}/api/prompts/latest/response",
            data=data,
            headers={"Content-Type": "application/json"},
            method="PUT"
        )

        with urllib.request.urlopen(req, timeout=5) as resp:
            return resp.status == 200
    except Exception:
        return False


def request_summary_generation(session_id: str, transcript_path: str) -> bool:
    """Request Worker to generate session summary"""
    last_message = get_last_assistant_message(transcript_path)
    parsed = parse_transcript_jsonl(transcript_path)

    try:
        data = json.dumps({
            "session_id": session_id,
            "last_assistant_message": last_message[:2000] if last_message else "",
            "transcript_stats": {
                "total_entries": parsed["total_entries"],
                "assistant_count": len(parsed["assistant_messages"]),
                "tool_count": len(parsed["tool_uses"])
            }
        }).encode("utf-8")

        req = urllib.request.Request(
            f"{WORKER_URL}/api/summaries/generate",
            data=data,
            headers={"Content-Type": "application/json"},
            method="POST"
        )

        with urllib.request.urlopen(req, timeout=5) as resp:
            return resp.status in (200, 201, 202)
    except Exception:
        return False


def main():
    try:
        hook_input = json.loads(sys.stdin.read())
    except json.JSONDecodeError:
        hook_input = {}

    session_id = hook_input.get("session_id", os.environ.get("CLAUDE_SESSION_ID", ""))
    transcript_path = hook_input.get("transcript_path", "")

    summary_requested = False
    response_saved = False

    if session_id:
        # 1. Save assistant response to user_prompts
        response = get_last_response_with_tools(transcript_path)
        if response:
            response_saved = update_prompt_response(session_id, response)

        # 2. Request summary generation (existing)
        summary_requested = request_summary_generation(session_id, transcript_path)

    # Stop hook은 차단하지 않음
    print(json.dumps({"continue": True}))


if __name__ == "__main__":
    main()
