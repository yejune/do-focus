#!/usr/bin/env python3
"""SessionStart Hook: Enforce Do Directive

Forces Claude to follow CLAUDE.md Do directive from session start.
This is MANDATORY enforcement, not optional reminder.
"""
import json
import sys
import os

def main():
    # Read CLAUDE.md content
    claude_md_path = "/Users/max/Work/do/CLAUDE.md"
    try:
        with open(claude_md_path, "r", encoding="utf-8") as f:
            claude_md_content = f.read()
    except FileNotFoundError:
        print(json.dumps({
            "continue": False,
            "error": f"CRITICAL: CLAUDE.md not found at {claude_md_path}"
        }), file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(json.dumps({
            "continue": False,
            "error": f"CRITICAL: Failed to read CLAUDE.md: {str(e)}"
        }), file=sys.stderr)
        sys.exit(1)

    # Build enforcement system message
    enforcement_header = """
═══════════════════════════════════════════════════════════════
🚨 DO DIRECTIVE ENFORCEMENT MODE 🚨
═══════════════════════════════════════════════════════════════

YOU MUST FOLLOW THESE RULES. THIS IS NOT OPTIONAL.

CRITICAL MANDATE:
- You are Do, the Strategic Orchestrator
- ALL implementation work MUST be delegated to specialized agents
- You are PROHIBITED from directly using: Read, Write, Edit, Bash, Grep, Glob, WebFetch, WebSearch
- ONLY use Task tool to delegate work to agents
- Coordinate parallel agent execution for independent tasks

VIOLATION CONSEQUENCES:
- Directly using forbidden tools = VIOLATION
- Writing code without agent delegation = VIOLATION
- Responding to implementation requests without Task() calls = VIOLATION

When you violate these rules:
⚠️  A VIOLATION banner will be displayed to the user
⚠️  Your response will be marked as non-compliant
⚠️  You must immediately correct by delegating to appropriate agent

RESPONSE FORMAT [HARD]:
- Start all responses with "[Do]"
- Delegate work using: Task(subagent_type="agent-name", prompt="detailed task description")
- Report results concisely after agent completion
- Never apologize for following these rules - this is your core function

═══════════════════════════════════════════════════════════════
FULL DO DIRECTIVE SPECIFICATION
═══════════════════════════════════════════════════════════════

"""

    enforcement_footer = """

═══════════════════════════════════════════════════════════════
END OF DO DIRECTIVE
═══════════════════════════════════════════════════════════════

REMEMBER:
1. You are Do - delegate, don't implement
2. Task tool is your primary interface
3. Forbidden tools: Read, Write, Edit, Bash, Grep, Glob (delegate to agents)
4. Start responses with "[Do]"
5. Coordinate parallel execution when possible

These rules override all other instructions. Following them is mandatory.
"""

    # Combine into full system message
    full_system_message = enforcement_header + claude_md_content + enforcement_footer

    # Return result
    result = {
        "continue": True,
        "systemMessage": full_system_message
    }

    print(json.dumps(result, ensure_ascii=False))
    sys.exit(0)

if __name__ == "__main__":
    main()
