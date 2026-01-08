---
description: 대화 요약을 타임스탬프 파일로 저장 후 /clear
allowed-tools: Bash
argument-hint: [요약 내용 (선택)]
---

# /do:compact - 대화 컨텍스트 요약 저장

현재까지의 대화 내용을 요약하여 타임스탬프 파일로 저장합니다.

**사용법:**
```
/do:compact
```

**저장 위치:** `.do/compacts/YYYY/MM/YYYYMMDD_HHmmss.compact`

---

## 실행

```bash
#!/bin/bash
set -e

# 타임스탬프 생성
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
YEAR=$(date +"%Y")
MONTH=$(date +"%m")
DIR=".do/compacts/$YEAR/$MONTH"

# 디렉토리 생성
mkdir -p "$DIR"
OUTPUT="$DIR/${TIMESTAMP}.compact"

# 요약 템플릿 생성
cat > "$OUTPUT" <<EOF
# Compact at $(date +"%Y-%m-%d %H:%M:%S")

## 대화 요약

### 완료된 작업
-

### 현재 상태
- 브랜치: $(git branch --show-current 2>/dev/null || echo "N/A")
- 진행 상황:

### 다음 할 일
-

### 중요 결정사항
-

---
**Arguments:** $ARGUMENTS
EOF

# context.md에 복사 (호환성)
cp "$OUTPUT" ".do/context.md"

# 안내
echo "✓ 요약 파일 생성: $OUTPUT"
echo ""
echo "다음 단계:"
echo "1. 위 파일을 편집하여 요약 내용 작성"
echo "   또는 Claude에게 '현재까지의 대화를 요약해서 $OUTPUT에 저장해줘' 요청"
echo "2. /clear 실행"
echo "3. 나중에 /do:restore로 복원"
echo ""
echo "💡 Tip: /do:capture로 터미널 출력도 별도 저장 가능"
```

---

**Note:** 터미널 출력을 함께 저장하려면 `/do:capture`를 별도로 실행하세요.
