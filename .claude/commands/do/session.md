---
description: Show current Claude session information
allowed-tools: ["Bash"]
argument-hint: (no arguments)
---

# Session Information

Shows current Claude session ID, log file location, and recent session history.

## Output

```bash
# Get current session ID
CURRENT_SESSION="${CLAUDE_SESSION_ID:-unknown}"

# Find Git root or use current directory
GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo "$PWD")

# Parse session ID to get directory structure
if [ "$CURRENT_SESSION" != "unknown" ]; then
    # YYYYMMDD-HHmmss -> YYYY/MM/DD
    SESSION_YEAR="${CURRENT_SESSION:0:4}"
    SESSION_MONTH="${CURRENT_SESSION:4:2}"
    SESSION_DAY="${CURRENT_SESSION:6:2}"
    SESSION_DIR="${GIT_ROOT}/.do/claude-session/${SESSION_YEAR}/${SESSION_MONTH}/${SESSION_DAY}"
    LOG_FILE="${SESSION_DIR}/${CURRENT_SESSION}.session"
else
    LOG_FILE="unknown"
fi

echo "🔗 Current Session"
echo ""
echo "Session ID: ${CURRENT_SESSION}"
echo "Log File:   ${LOG_FILE/#$HOME/~}"

# Check if log file exists
if [ "$LOG_FILE" != "unknown" ] && [ -f "$LOG_FILE" ]; then
    LINE_COUNT=$(wc -l < "$LOG_FILE" | tr -d ' ')
    FILE_SIZE=$(du -h "$LOG_FILE" | cut -f1)
    echo "Status:     Active (${LINE_COUNT} lines, ${FILE_SIZE})"
else
    echo "Status:     No log file found"
fi

echo ""
echo "Recent Sessions:"

# List recent session logs using find (project directory)
SESSION_BASE="${GIT_ROOT}/.do/claude-session"
if [ -d "$SESSION_BASE" ]; then
    find "$SESSION_BASE" -name "*.session" -type f 2>/dev/null | \
        xargs ls -lt 2>/dev/null | head -5 | while read -r line; do
        filename=$(echo "$line" | awk '{print $NF}')
        # Extract relative path from project .do/claude-session/
        relative_path="${filename#$SESSION_BASE/}"
        size=$(du -h "$filename" | cut -f1)

        # Check if this is the current session
        if [ "$filename" = "$LOG_FILE" ]; then
            lines=$(wc -l < "$filename" | tr -d ' ')
            echo "  📄 ${relative_path} (current, ${lines} lines)"
        else
            echo "  📄 ${relative_path} (${size})"
        fi
    done

    # Check if no files were found
    if ! find "$SESSION_BASE" -name "*.session" -type f 2>/dev/null | grep -q .; then
        echo "  (No session logs found)"
    fi
else
    echo "  (No session directory found)"
fi

echo ""
echo "💡 Tip: Use 'tail -f ${LOG_FILE/#$HOME/~}' to follow current session in real-time"
```
