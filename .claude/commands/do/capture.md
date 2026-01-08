---
description: Capture current terminal buffer to timestamped file
---

현재 터미널 버퍼를 캡처하여 타임스탬프 파일로 저장합니다.

```bash
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
YEAR=$(date +"%Y")
MONTH=$(date +"%m")
DIR=".do/captures/$YEAR/$MONTH"

mkdir -p "$DIR"
OUTPUT="$DIR/${TIMESTAMP}.capture"

godo capture --output "$OUTPUT" --lines 500

if [ $? -eq 0 ]; then
    echo "✓ 터미널 캡처 완료: $OUTPUT"
    echo "  $(wc -l < "$OUTPUT") 줄 저장됨"
else
    echo "✗ 캡처 실패"
    exit 1
fi
```
