#!/usr/bin/env python3
"""SessionStart Hook: Enforce Do Directive

Forces Claude to follow CLAUDE.md Do directive from session start.
This is MANDATORY enforcement, not optional reminder.
"""
import json
import sys
import os

def main():
    # Concise enforcement message
    # CLAUDE.md is already loaded as project context, no need to duplicate
    system_message = """
🚨 DO DIRECTIVE ENFORCEMENT

YOU ARE DO - THE STRATEGIC ORCHESTRATOR

MANDATORY RULES:
- ALL implementation work → delegate via Task tool
- FORBIDDEN direct use: Read, Write, Edit, Bash, Grep, Glob, WebFetch, WebSearch
- Start responses with "[Do]"
- Execute independent tasks in parallel

VIOLATIONS:
- Direct tool use without delegation = VIOLATION
- Code writing without Task() = VIOLATION

Follow CLAUDE.md specification in project root.
"""

    # Return result
    result = {
        "continue": True,
        "systemMessage": system_message.strip()
    }

    print(json.dumps(result, ensure_ascii=False))
    sys.exit(0)

if __name__ == "__main__":
    main()
