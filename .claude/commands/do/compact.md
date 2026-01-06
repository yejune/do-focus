---
description: 빠른 컨텍스트 정리 - 요약을 파일에 저장 후 /clear
allowed-tools: Write, Read
argument-hint: [요약할 내용]
---

# /do:compact - 빠른 컨텍스트 정리

## 실행 단계

1. **현재 세션 요약 생성** (아래 형식):

```markdown
# 세션 컨텍스트

## 완료된 작업
-

## 현재 상태
- 브랜치:
- 진행 상황:

## 다음 할 일
-

## 중요 결정사항
-
```

2. **`.do/context.md` 파일에 저장**

3. **사용자에게 안내**:
```
컨텍스트가 .do/context.md에 저장됨.

다음 단계:
1. /clear 실행
2. "컨텍스트 복원해줘" 또는 ".do/context.md 읽어줘" 입력
```

## 추가 컨텍스트

$ARGUMENTS
