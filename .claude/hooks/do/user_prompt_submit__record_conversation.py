#!/usr/bin/env python3
"""
User Prompt Submit Hook: 사용자 대화 기록
사용자 입력을 Worker DB에 저장
"""
import json
import os
import sys
import urllib.request
import urllib.error

WORKER_PORT = int(os.environ.get("DO_WORKER_PORT", "3778"))
WORKER_URL = f"http://127.0.0.1:{WORKER_PORT}"


def send_to_worker(session_id: str, content: str, prompt_number: int = 1):
    """Worker에 대화 전송 (observations + user_prompts 둘 다)"""
    # 1. observations 저장 (기존)
    try:
        data = {
            "session_id": session_id,
            "type": "conversation",
            "content": content[:2000],
            "importance": 3,
            "tags": ["user", "prompt"]
        }
        req = urllib.request.Request(
            f"{WORKER_URL}/api/observations",
            data=json.dumps(data).encode('utf-8'),
            headers={"Content-Type": "application/json"},
            method="POST"
        )
        urllib.request.urlopen(req, timeout=2)
    except Exception:
        pass

    # 2. user_prompts 저장 (신규)
    try:
        import time
        data = {
            "session_id": session_id,
            "prompt_number": prompt_number,
            "prompt_text": content[:5000],  # 5000자 제한
            "created_at_epoch": int(time.time())
        }
        req = urllib.request.Request(
            f"{WORKER_URL}/api/prompts",
            data=json.dumps(data).encode('utf-8'),
            headers={"Content-Type": "application/json"},
            method="POST"
        )
        urllib.request.urlopen(req, timeout=2)
    except Exception:
        pass


def main():
    # Read hook input from stdin
    hook_input = {}
    try:
        raw = sys.stdin.read()
        if raw.strip():
            hook_input = json.loads(raw)
    except Exception:
        pass

    session_id = hook_input.get("session_id", os.environ.get("CLAUDE_SESSION_ID", ""))
    prompt = hook_input.get("prompt", "")

    if session_id and prompt:
        send_to_worker(session_id, prompt)

    # Always continue (don't block)
    print(json.dumps({"continue": True}))


if __name__ == "__main__":
    main()
