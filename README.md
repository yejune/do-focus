# Do

Claude Code용 프로젝트별 AI 개발 환경

---

## 설치

### 방법 1: 프로젝트에 직접 복사

```bash
# Do 클론
git clone https://github.com/yejune/do /tmp/do

# 원하는 프로젝트로 복사
cp -r /tmp/do/.claude /your/project/
cp -r /tmp/do/.do /your/project/
cp /tmp/do/CLAUDE.md /your/project/

# 정리
rm -rf /tmp/do
```

### 방법 2: Git subtree (업데이트 용이)

```bash
cd /your/project
git subtree add --prefix=.claude https://github.com/yejune/do.git main --squash
```

### 방법 3: Git submodule

```bash
cd /your/project
git submodule add https://github.com/yejune/do.git .do-kit
ln -s .do-kit/.claude .claude
ln -s .do-kit/.do .do
cp .do-kit/CLAUDE.md .
```

---

## 설정

프로젝트별 설정은 `.do/config/sections/`에서 관리:

```yaml
# .do/config/sections/language.yaml
conversation_language: ko  # 대화 언어 (ko, en, ja, zh)
code_comments: en          # 코드 주석 언어

# .do/config/sections/user.yaml
name: "사용자 이름"
```

---

## 요구사항

- Claude Code CLI
- Python 3.10+ (훅 실행용)
- uv (Python 패키지 관리자)

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
  config/          # 프로젝트 설정
```

---

## 라이선스

MIT
