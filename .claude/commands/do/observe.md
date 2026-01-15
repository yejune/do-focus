# /do:observe

수동 관찰 기록 커맨드. 의사결정, 버그 수정, 학습 내용 등을 명시적으로 기록.

## 사용법

```
/do:observe decision "JWT 대신 세션 기반 인증 선택 - 보안 요구사항 충족"
/do:observe bugfix "race condition 수정 - 동시 요청 시 데이터 충돌"
/do:observe learning "Redis 캐싱 패턴 학습"
```

## allowed-tools
Bash

## Step 1: 인자 파싱

사용자 입력에서 타입과 내용 추출:
- 첫 번째 단어: observation 타입 (decision, bugfix, feature, refactor, learning, docs)
- 나머지: 내용

## Step 2: Worker에 전송

```bash
curl -X POST http://127.0.0.1:3778/api/observations \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "$CLAUDE_SESSION_ID",
    "type": "{타입}",
    "content": "{내용}",
    "user_name": "$DO_USER_NAME"
  }'
```

## Step 3: 확인 메시지

```
✓ 관찰 기록됨: [{타입}] {내용 요약}
```

## 예시

입력: `/do:observe decision API 버전 관리를 URL path 방식으로 결정`
출력: `✓ 관찰 기록됨: [decision] API 버전 관리를 URL path 방식으로 결정`
