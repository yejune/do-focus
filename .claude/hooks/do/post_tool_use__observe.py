#!/usr/bin/env python3
"""
Post Tool Use Hook: 작업 관찰 자동 기록
urllib 사용 (표준 라이브러리)
"""
import json
import os
import sys
import urllib.request
import urllib.error

WORKER_PORT = int(os.environ.get("DO_WORKER_PORT", "3778"))
WORKER_URL = f"http://127.0.0.1:{WORKER_PORT}"

def classify_observation(tool_name: str, tool_input: dict, tool_output: str) -> tuple:
    """도구 사용을 관찰 타입으로 분류. (type, content, agent_name, importance, tags) 반환"""

    # Task 도구 → delegation
    if tool_name == "Task":
        agent_type = tool_input.get("subagent_type", "unknown")
        description = tool_input.get("description", "")
        return "delegation", f"[{agent_type}] {description}", agent_type, 4, ["agent", agent_type]

    # Write/Edit → feature
    if tool_name in ("Write", "Edit", "MultiEdit"):
        file_path = tool_input.get("file_path", "")
        if "test" in file_path.lower():
            return "test", f"테스트 파일 수정: {file_path}", None, 3, ["test"]
        return "feature", f"파일 수정: {file_path}", None, 3, ["code"]

    # Bash git commit
    if tool_name == "Bash":
        command = tool_input.get("command", "")
        if "git commit" in command:
            return "commit", f"커밋 실행", None, 4, ["git"]

    return None, None, None, None, None

def send_to_worker(session_id: str, obs_type: str, content: str,
                   agent_name: str = None, importance: int = 3, tags: list = None):
    """Worker에 관찰 전송 (urllib 사용)"""
    try:
        data = {
            "session_id": session_id,
            "type": obs_type,
            "content": content,
            "importance": importance,
            "tags": tags or []
        }
        if agent_name:
            data["agent_name"] = agent_name

        req = urllib.request.Request(
            f"{WORKER_URL}/api/observations",
            data=json.dumps(data).encode('utf-8'),
            headers={"Content-Type": "application/json"},
            method="POST"
        )
        urllib.request.urlopen(req, timeout=1)
    except:
        pass  # 실패해도 무시

def main():
    hook_input = json.loads(sys.stdin.read())

    session_id = hook_input.get("session_id", os.environ.get("CLAUDE_SESSION_ID", ""))
    tool_name = hook_input.get("tool_name", "")
    tool_input = hook_input.get("tool_input", {})
    tool_output = hook_input.get("tool_output", "")

    obs_type, content, agent_name, importance, tags = classify_observation(tool_name, tool_input, tool_output)

    if obs_type and content:
        send_to_worker(session_id, obs_type, content, agent_name, importance, tags)

    print(json.dumps({"continue": True}))

if __name__ == "__main__":
    main()
