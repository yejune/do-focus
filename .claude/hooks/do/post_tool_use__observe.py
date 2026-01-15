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
    """도구 사용을 관찰 타입으로 분류. (type, content, agent_name, importance, tags, result_preview) 반환"""

    # tool_output의 첫 500자를 result_preview로 추출
    result_preview = str(tool_output)[:500] if tool_output else None

    # Task 도구 → delegation
    if tool_name == "Task":
        agent_type = tool_input.get("subagent_type", "unknown")
        description = tool_input.get("description", "")
        return "delegation", f"[{agent_type}] {description}", agent_type, 4, ["agent", agent_type], result_preview

    # Write/Edit → feature or plan
    if tool_name in ("Write", "Edit", "MultiEdit"):
        file_path = tool_input.get("file_path", "")
        if ".do/plans/" in file_path or "/.claude/plans/" in file_path:
            return "plan", f"플랜 생성: {file_path}", None, 5, ["plan"], result_preview
        if "test" in file_path.lower():
            return "test", f"테스트: {file_path}", None, 3, ["test"], result_preview
        return "feature", f"수정: {file_path}", None, 3, ["code"], result_preview

    # Bash
    if tool_name == "Bash":
        command = tool_input.get("command", "")[:100]  # 100자 제한
        if "git commit" in command:
            return "commit", f"커밋", None, 4, ["git"], result_preview
        if "git push" in command:
            return "push", f"푸시", None, 4, ["git"], result_preview
        return "bash", f"실행: {command}", None, 2, ["bash"], result_preview

    # Read
    if tool_name == "Read":
        file_path = tool_input.get("file_path", "")
        return "read", f"읽기: {file_path}", None, 1, ["read"], result_preview

    # Grep/Glob
    if tool_name in ("Grep", "Glob"):
        pattern = tool_input.get("pattern", "")
        return "search", f"검색: {pattern}", None, 1, ["search"], result_preview

    # WebFetch/WebSearch
    if tool_name == "WebFetch":
        url = tool_input.get("url", "")
        return "web", f"웹: {url}", None, 2, ["web"], result_preview
    if tool_name == "WebSearch":
        query = tool_input.get("query", "")
        return "web", f"검색: {query}", None, 2, ["web"], result_preview

    # 기타 모든 도구
    return "tool", f"{tool_name}", None, 1, [tool_name.lower()], result_preview

def send_to_worker(session_id: str, obs_type: str, content: str,
                   agent_name: str = None, importance: int = 3, tags: list = None,
                   result_preview: str = None):
    """Worker에 관찰 전송 (urllib 사용)"""
    try:
        data = {
            "session_id": session_id,
            "type": obs_type,
            "content": content,
            "importance": importance,
            "tags": tags or [],
            "result_preview": result_preview[:500] if result_preview else None
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

    obs_type, content, agent_name, importance, tags, result_preview = classify_observation(tool_name, tool_input, tool_output)

    if obs_type and content:
        send_to_worker(session_id, obs_type, content, agent_name, importance, tags, result_preview)

    print(json.dumps({"continue": True}))

if __name__ == "__main__":
    main()
