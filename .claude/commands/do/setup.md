# /do:setup

Do 환경 설정을 진행합니다.

## allowed-tools
Read, Write, Edit, AskUserQuestion

## Step 1: 현재 설정 확인

`.claude/settings.local.json` 파일을 읽어서 현재 설정 확인.
없으면 기본값 사용.

## Step 2-1: 기본 설정 (첫 번째 AskUserQuestion)

**먼저 시스템 사용자명을 가져옵니다:**
- `whoami` 명령 실행하여 시스템 사용자명 확인

AskUserQuestion으로 4개 질문:

1. **이름** (DO_USER_NAME)
   - 옵션 1: "{시스템사용자명} 사용 (추천)" - 설명: "현재 시스템 로그인 이름을 사용합니다"
   - 옵션 2: "설정 안 함" - 설명: "이름 없이 진행합니다"
   - 옵션 3 (Other): "직접 입력 (아래 입력칸에 원하는 이름 작성)"
   - 질문 텍스트: "사용자 이름을 선택하거나, 직접 입력하려면 맨 아래 입력칸을 사용하세요"

2. **대화 언어** (DO_LANGUAGE)
   - 한국어 (ko), English (en), 日本語 (ja), 中文 (zh)

3. **커밋 언어** (DO_COMMIT_LANGUAGE)
   - 한국어 (ko), English (en)

4. **응답 스타일** (style - env 아님)
   - sprint, pair, direct

**처리 로직:**
- 옵션 1 선택: DO_USER_NAME = 시스템사용자명
- 옵션 2 선택: DO_USER_NAME = "" (빈 문자열)
- 옵션 3 (직접 입력): DO_USER_NAME = 입력된 값

## Step 2-2: 추가 설정 (두 번째 AskUserQuestion)

AskUserQuestion으로 2개 질문:

1. **AI 푸터** (DO_AI_FOOTER)
   - 예 (true), 아니오 (false)

2. **에이전트 수정 확인** (DO_CONFIRM_CHANGES)
   - 예 (true), 아니오 (false)

## Step 3: 설정 저장

`.claude/settings.local.json` 파일 업데이트:

```json
{
  "style": "{선택한 스타일}",
  "env": {
    "DO_USER_NAME": "{이름}",
    "DO_LANGUAGE": "{언어코드}",
    "DO_COMMIT_LANGUAGE": "{커밋언어}",
    "DO_AI_FOOTER": "{true/false}",
    "DO_CONFIRM_CHANGES": "{true/false}"
  }
}
```

기존 settings.local.json 내용 유지하면서 env 필드만 업데이트.

## Step 4: 완료 메시지

설정 완료!
- 이름: {이름}
- 대화 언어: {언어}
- 커밋 언어: {커밋언어}
- AI 푸터: {예/아니오}
- 에이전트 확인: {예/아니오}
- 스타일: {스타일}
