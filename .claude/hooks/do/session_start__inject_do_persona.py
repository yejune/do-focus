#!/usr/bin/env python3
"""SessionStart Hook: Inject Do Persona

Forces Claude to follow CLAUDE.md Do directive from session start.
"""
import json
import sys

def main():
    # Read CLAUDE.md content
    try:
        with open("/Users/max/Work/do/CLAUDE.md", "r", encoding="utf-8") as f:
            claude_md = f.read()
    except:
        claude_md = ""
    
    # Force Do persona
    do_reminder = """
🤖 [Do Mode Active]

You MUST follow these rules from CLAUDE.md:
1. Start all responses with "[Do]"
2. Delegate ALL work to agents via Task tool
3. NEVER use: Read, Write, Edit, Bash, Grep, Glob directly
4. Coordinate parallel agent execution
5. Report results concisely

VIOLATION = Not being Do
"""
    
    result = {
        "continue": True,
        "systemMessage": do_reminder
    }
    
    print(json.dumps(result, ensure_ascii=False))
    sys.exit(0)

if __name__ == "__main__":
    main()
