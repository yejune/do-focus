# /do:setup

Do 환경 설정을 진행합니다.

## allowed-tools
Read, Write, Edit, AskUserQuestion

## Step 1: 현재 설정 확인

`.claude/settings.local.json` 파일을 읽어서 현재 설정 확인.
없으면 기본값 사용.

## Step 2: 사용자에게 질문

AskUserQuestion으로 6개 질문:

1. **이름** (DO_USER_NAME)
   - 현재 값 또는 빈 값

2. **대화 언어** (DO_LANGUAGE)
   - 한국어 (ko), English (en), 日本語 (ja), 中文 (zh)

3. **커밋 언어** (DO_COMMIT_LANGUAGE)
   - 한국어 (ko), English (en)

4. **AI 푸터** (DO_AI_FOOTER)
   - 예 (true), 아니오 (false)

5. **에이전트 수정 확인** (DO_CONFIRM_CHANGES)
   - 예 (true), 아니오 (false)

6. **응답 스타일** (style - env 아님)
   - sprint, pair, direct

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
