#!/usr/bin/env python3
"""
PostToolUse Hook: Check tobrew after git commit & Remove AI footer

Detects git commit commands and:
1. Checks for AI footer in commit message and removes it (3rd defense layer)
2. Prompts user to release if tobrew files exist

Event: PostToolUse
Trigger: Successful git commit command
Action: Clean AI footer + Check for tobrew.* files and suggest release

Usage:
    Automatically triggered after git commit via PostToolUse hook
"""

import os
import sys
import json
import subprocess
from pathlib import Path


def check_and_remove_ai_footer(project_dir: Path) -> bool:
    """Check latest commit for AI footer and remove it via amend.

    Args:
        project_dir: Project root directory

    Returns:
        True if AI footer was found and removed, False otherwise
    """
    try:
        # Get latest commit message
        result = subprocess.run(
            ["git", "log", "-1", "--pretty=%B"],
            cwd=project_dir,
            capture_output=True,
            text=True,
            check=True
        )
        commit_msg = result.stdout.strip()

        # Check for AI footer patterns
        ai_footer_patterns = [
            "Co-Authored-By: Claude Code",
            "Co-authored-by: Claude Code",
        ]

        has_ai_footer = any(pattern in commit_msg for pattern in ai_footer_patterns)

        if not has_ai_footer:
            return False

        # Remove AI footer lines
        lines = commit_msg.split('\n')
        cleaned_lines = [
            line for line in lines
            if not any(pattern in line for pattern in ai_footer_patterns)
            and line.strip() != '---'
        ]

        # Remove trailing empty lines
        while cleaned_lines and not cleaned_lines[-1].strip():
            cleaned_lines.pop()

        cleaned_msg = '\n'.join(cleaned_lines)

        # Amend commit with cleaned message
        subprocess.run(
            ["git", "commit", "--amend", "-m", cleaned_msg],
            cwd=project_dir,
            capture_output=True,
            text=True,
            check=True
        )

        return True

    except subprocess.CalledProcessError:
        # Git command failed, ignore
        return False
    except Exception:
        # Any other error, ignore
        return False


def find_tobrew_files(project_dir: Path) -> list[str]:
    """Find tobrew configuration files in project directory.

    Args:
        project_dir: Project root directory

    Returns:
        List of found tobrew file names
    """
    tobrew_patterns = ["tobrew.lock", "tobrew.yaml", "tobrew.yml", ".tobrew"]
    found_files = []

    for pattern in tobrew_patterns:
        file_path = project_dir / pattern
        if file_path.exists():
            found_files.append(pattern)

    return found_files


def create_release_prompt(tobrew_files: list[str]) -> dict:
    """Create system message with release prompt.

    Args:
        tobrew_files: List of found tobrew files

    Returns:
        System message dict for Claude Code
    """
    files_list = ", ".join(tobrew_files)

    message = f"""✅ 커밋 완료
📦 tobrew 파일 발견: {files_list}

릴리즈가 필요할 수 있습니다.

다음 중 선택하세요:
1. "예, 릴리즈" - git push && tobrew release --patch 실행
2. "나중에" - 커밋만 유지 (나중에 수동으로 릴리즈)

릴리즈 하시겠습니까?"""

    return {
        "type": "systemMessage",
        "content": message
    }


def main():
    """Main hook execution function."""
    try:
        # Get environment variables from hook context
        tool_name = os.environ.get("TOOL_NAME", "")
        command = os.environ.get("COMMAND", "")
        project_dir = os.environ.get("CLAUDE_PROJECT_DIR", os.getcwd())

        # Check if this was a git commit command
        if tool_name != "Bash" or "git commit" not in command:
            # Not a git commit, exit silently
            sys.exit(0)

        project_path = Path(project_dir)

        # 3rd defense layer: Check and remove AI footer
        ai_footer_removed = check_and_remove_ai_footer(project_path)
        if ai_footer_removed:
            print(json.dumps({
                "type": "systemMessage",
                "content": "🛡️ 3차 방어선 작동: AI 푸터 자동 제거됨"
            }, ensure_ascii=False))

        # Check for tobrew files
        tobrew_files = find_tobrew_files(project_path)

        if not tobrew_files:
            # No tobrew files found, exit silently
            sys.exit(0)

        # tobrew files found, create release prompt
        prompt_message = create_release_prompt(tobrew_files)

        # Output as JSON for Claude Code to process
        print(json.dumps(prompt_message, ensure_ascii=False))
        sys.exit(0)

    except Exception as e:
        # Log error but don't block workflow
        print(f"Error in tobrew check hook: {e}", file=sys.stderr)
        sys.exit(0)


if __name__ == "__main__":
    main()
