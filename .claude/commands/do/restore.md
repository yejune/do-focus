---
description: 저장된 컨텍스트 복원 - 최신 compact 파일 또는 context.md 읽기
allowed-tools: Bash
---

# /do:restore - 컨텍스트 복원

최신 `.do/compacts/*.compact` 파일을 찾아 복원하거나, 없으면 `.do/context.md`를 fallback으로 사용합니다.

## 실행

1. `.do/compacts/` 디렉토리에서 최신 compact 파일 검색 (타임스탬프 기반)
2. compact 파일이 있으면 해당 파일 표시
3. 없으면 `.do/context.md` 표시
4. 둘 다 없으면 안내 메시지 출력

## 구현

```bash
# 최신 compact 파일 찾기 (타임스탬프 역순 정렬)
LATEST=$(find .do/compacts -name "*.compact" 2>/dev/null | sort -r | head -1)

if [ -n "$LATEST" ]; then
    echo "✓ 컨텍스트 복원: $LATEST"
    echo ""
    cat "$LATEST"
elif [ -f ".do/context.md" ]; then
    echo "✓ 컨텍스트 복원: .do/context.md"
    echo ""
    cat ".do/context.md"
else
    echo "✗ 저장된 컨텍스트가 없습니다."
    echo ""
    echo "사용법:"
    echo "  /do:compact - 현재 컨텍스트 저장"
    exit 1
fi
```

## 출력 예시

### Compact 파일 복원:
```
✓ 컨텍스트 복원: .do/compacts/2026/01/20260108_163045.compact

# Compact at 2026-01-08 16:30:45

## 대화 요약
- restore.md 슬래시 커맨드 수정 요청
- compact 파일 우선, context.md fallback

## 터미널 출력
...
```

### Context.md Fallback:
```
✓ 컨텍스트 복원: .do/context.md

# 프로젝트 컨텍스트
...
```

### 파일 없음:
```
✗ 저장된 컨텍스트가 없습니다.

사용법:
  /do:compact - 현재 컨텍스트 저장
```
