# Do

Claude Code용 프로젝트별 AI 개발 환경

---

## 빠른 시작

1. Claude Code 설치: `brew install anthropics/tap/claude`
2. godo 설치: `brew install yejune/tap/godo`
3. 프로젝트에서 초기화: `godo sync` (설치/업데이트 자동 감지)
4. 초기 설정: `/do:setup`

---

## 설정

`/do:setup` 실행하여 개인 설정:

```json
// .claude/settings.local.json (gitignore)
{
  "env": {
    "DO_USER_NAME": "이름",
    "DO_LANGUAGE": "ko",
    "DO_COMMIT_LANGUAGE": "en",
    "DO_AI_FOOTER": "false",
    "DO_CONFIRM_CHANGES": "true"
  }
}
```

| 환경변수 | 설명 | 옵션 |
|---------|------|------|
| DO_USER_NAME | 사용자 이름 | 자유 입력 |
| DO_LANGUAGE | 대화 언어 | ko, en, ja, zh |
| DO_COMMIT_LANGUAGE | 커밋 언어 | ko, en |
| DO_AI_FOOTER | AI 푸터 | true, false |
| DO_CONFIRM_CHANGES | 수정 확인 | true, false |

---

## 요구사항

- Claude Code CLI
- Python 3.10+ (훅 실행용)
- uv (Python 패키지 관리자)

---

## godo 명령어

```bash
godo sync         # 설치 또는 업데이트
godo selfupdate   # godo CLI 자체 업데이트
godo version      # 버전 확인
```

---

## 주요 기능

### 27개 전문 에이전트

- **expert-***: backend, frontend, database, security, testing, debug, performance, devops, uiux
- **manager-***: tdd, spec, docs, quality, project, git, strategy, claude-code
- **builder-***: agent, skill, command, plugin
- **mcp-***: context7, playwright, figma, notion, sequential-thinking
- **ai-***: nano-banana (이미지 생성)

### 47개 스킬

- **lang-***: python, typescript, javascript, go, rust, java, kotlin, swift, ruby, php, elixir, scala, r, cpp, csharp, flutter
- **platform-***: vercel, supabase, firebase-auth, firestore, neon, railway, convex, clerk, auth0
- **workflow-***: project, spec, testing, templates, worktree, jit-docs
- **domain-***: backend, frontend, database, uiux
- **foundation-***: core, claude, quality, context, philosopher
- **library-***: shadcn, mermaid, nextra

### 자동 백그라운드 작업 추적

`SubagentStop` 훅으로 병렬 실행 작업 자동 추적

### 한국어 지원

대화, 문서, 주석 모두 한국어 설정 가능

---

## 디렉토리 구조

```
.claude/
  agents/do/       # 에이전트 정의 (27개)
  skills/do-*/     # 스킬 정의 (47개)
  hooks/do/        # 훅 스크립트
  commands/do/     # 슬래시 커맨드

.do/
  cache/           # 캐시 및 상태 파일
  memory/          # 세션 상태 저장
```

---

## 라이선스

MIT
