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
LOG_FILE="$HOME/.do/claude-session-${CURRENT_SESSION}.log"

echo "🔗 Current Session"
echo ""
echo "Session ID: ${CURRENT_SESSION}"
echo "Log File:   ${LOG_FILE/#$HOME/~}"

# Check if log file exists
if [ -f "$LOG_FILE" ]; then
    LINE_COUNT=$(wc -l < "$LOG_FILE" | tr -d ' ')
    FILE_SIZE=$(du -h "$LOG_FILE" | cut -f1)
    echo "Status:     Active (${LINE_COUNT} lines, ${FILE_SIZE})"
else
    echo "Status:     No log file found"
fi

echo ""
echo "Recent Sessions:"

# List recent session logs
if [ -d "$HOME/.do" ]; then
    ls -lht "$HOME/.do/claude-session-"*.log 2>/dev/null | head -5 | while read -r line; do
        filename=$(echo "$line" | awk '{print $NF}')
        basename=$(basename "$filename")
        size=$(echo "$line" | awk '{print $5}')

        # Check if this is the current session
        if [ "$filename" = "$LOG_FILE" ]; then
            lines=$(wc -l < "$filename" | tr -d ' ')
            echo "  📄 ${basename} (current, ${lines} lines)"
        else
            # Convert size to human readable
            if [ -f "$filename" ]; then
                file_size=$(du -h "$filename" | cut -f1)
                echo "  📄 ${basename} (${file_size})"
            fi
        fi
    done

    # Check if no files were found
    if ! ls "$HOME/.do/claude-session-"*.log >/dev/null 2>&1; then
        echo "  (No session logs found)"
    fi
else
    echo "  (No .do directory found)"
fi

echo ""
echo "💡 Tip: Use 'tail -f ${LOG_FILE/#$HOME/~}' to follow current session in real-time"
```
