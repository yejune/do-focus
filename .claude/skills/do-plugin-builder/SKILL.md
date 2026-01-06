---
name: do-plugin-builder
description: Claude Code plugin 개발 패턴, 템플릿 및 모범 사례. Plugin 생성, 구성 요소 정의, 문제 해결 시 사용.
version: 1.1.0
category: foundation
allowed-tools: Read, Write, Edit, Grep, Glob, Bash, TodoWrite
tags:
  - plugin
  - claude-code
  - development
  - templates
  - hooks
  - commands
  - agents
  - skills
  - mcp
  - lsp
updated: 2026-01-06
status: active
author: Do Team
---

# Claude Code Plugin Builder

## 빠른 참조

Plugin 개발 필수 사항: 올바른 구조와 구성 요소로 Claude Code plugin을 빌드.

디렉토리 구조:
- `.claude-plugin/plugin.json` - Plugin 매니페스트 (필수)
- `commands/` - Slash commands (plugin 루트)
- `agents/` - Custom agents (plugin 루트)
- `skills/` - Agent skills (plugin 루트)
- `hooks/` - Event handlers (plugin 루트)
- `.mcp.json` - MCP 서버 설정
- `.lsp.json` - LSP 서버 설정

핵심 제약: 구성 요소 디렉토리는 반드시 plugin 루트 수준에 배치. `.claude-plugin/` 내부가 아님.

사용 시점:
- 새 Claude Code plugin 생성
- Plugin 구성 요소 정의 (commands, agents, skills, hooks)
- MCP 또는 LSP 서버 설정
- Plugin 로딩 문제 해결
- 독립 설정을 plugin 형식으로 마이그레이션

---

## 구현 가이드

### Plugin 디렉토리 구조

올바른 Plugin 레이아웃:
```
my-plugin/
  .claude-plugin/
    plugin.json          # 필수: Plugin 메타데이터만
  commands/              # 루트 수준
    my-command.md
  agents/                # 루트 수준
    my-agent.md
  skills/                # 루트 수준
    my-skill/
      SKILL.md
  hooks/                 # 루트 수준
    hooks.json
  .mcp.json              # 루트 수준
  .lsp.json              # 루트 수준
  LICENSE
  CHANGELOG.md
  README.md
```

피해야 할 실수: 구성 요소 디렉토리를 .claude-plugin 폴더 내에 배치하지 말 것.

### plugin.json Schema

최소 구성:
```json
{
  "name": "my-plugin"
}
```

완전한 구성:
```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "description": "Plugin 목적 설명",
  "author": {
    "name": "Author Name",
    "email": "author@example.com",
    "url": "https://example.com"
  },
  "homepage": "https://github.com/user/my-plugin",
  "repository": "https://github.com/user/my-plugin",
  "license": "MIT",
  "keywords": ["claude-code", "automation"],
  "commands": ["./commands"],
  "agents": ["./agents"],
  "skills": ["./skills"],
  "hooks": ["./hooks/hooks.json"],
  "mcpServers": ["./.mcp.json"],
  "lspServers": ["./.lsp.json"],
  "outputStyles": ["./output-styles"]
}
```

필수 필드:
- name: Kebab-case 고유 식별자 (문자, 숫자, 하이픈만 허용)

선택 필드:
- version: 시맨틱 버전 (MAJOR.MINOR.PATCH)
- description: Plugin 목적 설명
- author: name, email, url을 포함하는 객체
- homepage, repository, license, keywords
- 구성 요소 경로 참조 (commands, agents, skills, hooks)
- 서버 설정 (mcpServers, lspServers)
- 출력 스타일 참조 (outputStyles)

### 경로 설정 규칙

경로 규칙: plugin 루트 기준 상대 경로, `./`로 시작, 배열 지원

환경 변수: `${CLAUDE_PLUGIN_ROOT}` (plugin 디렉토리), `${CLAUDE_PROJECT_DIR}` (프로젝트 디렉토리)

### Slash Commands

Command 파일 구조 (commands/my-command.md):
```markdown
---
description: 발견용 Command 설명
---

Command 지시사항과 프롬프트 내용.

Arguments: $ARGUMENTS (전체), $1, $2 (위치별)
파일 참조: @path/to/file.md
```

Frontmatter 필드:
- description (필수): 도움말 표시용 Command 목적

Argument 처리:
- `$ARGUMENTS` - 전체 argument를 단일 문자열로
- `$1`, `$2`, `$3` - 개별 위치 argument
- `@file.md` - 파일 내용 주입

Command 네임스페이싱: `/plugin-name:command-name`으로 접근

### Custom Agents

Agent 파일 구조 (agents/my-agent.md):
```markdown
---
name: my-agent
description: Agent 목적과 기능
tools: Read, Write, Edit, Grep, Glob, Bash
model: sonnet
permissionMode: default
skills:
  - skill-name-one
  - skill-name-two
---

Agent system prompt와 지시사항.
```

Frontmatter: name (필수), description, tools, model (sonnet/opus/haiku/inherit), permissionMode (default/bypassPermissions/plan/passthrough), skills

사용 가능한 Tools: Read, Write, Edit, Grep, Glob, Bash, WebFetch, WebSearch, Task, TodoWrite

### Agent Skills

Skill 구조 (skills/my-skill/SKILL.md):
```markdown
---
name: my-skill
description: Skill 목적과 사용 시점
allowed-tools: Read, Grep, Glob
---

# Skill 이름

## 빠른 참조

간단한 개요와 핵심 개념.

## 구현 가이드

상세한 구현 패턴.

## 고급 패턴

전문가 수준 지식.
```

Skill 발견: 모델이 컨텍스트 관련성에 따라 호출. Task 컨텍스트가 skill 설명과 일치하면 자동 로드.

### Hooks 설정

hooks.json 구조:
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "bash ./hooks/validate-write.sh"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "prompt",
            "prompt": "작업이 성공적으로 완료되었는지 확인"
          }
        ]
      }
    ]
  }
}
```

Hook Events: PreToolUse, PostToolUse, PostToolUseFailure, PermissionRequest, UserPromptSubmit, Notification, Stop, SubagentStart, SubagentStop, SessionStart, SessionEnd, PreCompact

Hook 유형:
- command: Bash 명령 실행
- prompt: LLM 프롬프트 전송
- agent: 처리를 위한 에이전트 호출

Matcher 패턴:
- 정확한 도구 이름: "Write", "Bash"
- 와일드카드: "*"는 모든 도구에 일치
- 도구 이름 기반 도구별 필터링

### MCP 서버 설정

.mcp.json 구조:
```json
{
  "mcpServers": {
    "my-server": {
      "command": "npx",
      "args": ["-y", "@my-org/mcp-server"],
      "env": {
        "API_KEY": "${API_KEY}"
      }
    }
  }
}
```

Transport 유형: stdio (기본), http, sse

필드: command (필수), args, env, type, url (http/sse용)

### LSP 서버 설정

.lsp.json 구조:
```json
{
  "lspServers": {
    "python": {
      "command": "pylsp",
      "args": [],
      "extensionToLanguage": {
        ".py": "python",
        ".pyi": "python"
      },
      "env": {
        "PYTHONPATH": "${CLAUDE_PROJECT_DIR}"
      }
    }
  }
}
```

필수 필드:
- command: LSP 서버 실행 파일
- extensionToLanguage: 파일 확장자와 언어 매핑

선택 필드: args, env, transport, initializationOptions, settings, workspaceFolder, startupTimeout, shutdownTimeout, restartOnCrash, maxRestarts, loggingConfig

---

## 고급 패턴

### 개발 워크플로우

로컬 개발:
```bash
# 단일 plugin 테스트
claude --plugin-dir ./my-plugin

# 다중 plugin 테스트
claude --plugin-dir ./plugin-one --plugin-dir ./plugin-two
```

구성 요소 테스트:
- Commands: `/plugin-name:command-name` 호출
- Agents: `/agents`로 목록 확인 후 이름으로 호출
- Skills: skill 도메인 관련 질문
- Hooks: 이벤트 트리거 및 디버그 로그 확인

디버깅:
```bash
# 디버그 모드 활성화
claude --debug

# plugin 구조 유효성 검사
claude plugin validate

# plugin 오류 확인
/plugin errors
```

### 보안 모범 사례

경로 보안:
- plugin 상대 경로에 항상 `${CLAUDE_PLUGIN_ROOT}` 사용
- 절대 경로 하드코딩 금지
- hook 스크립트에서 모든 입력 검증
- 경로 탐색 공격 방지

권한 가이드라인:
- 도구 접근에 최소 권한 적용
- 에이전트 권한을 필요한 작업으로 제한
- hook 명령 입력 검증
- 환경 변수 정제

### 배포와 설치

Plugin 설치 범위:
- user: `~/.claude/settings.json` (개인, 기본)
- project: `.claude/settings.json` (팀, 버전 관리)
- local: `.claude/settings.local.json` (개발자, gitignore)
- managed: `managed-settings.json` (기업, 읽기 전용)

CLI 명령: `claude plugin install/uninstall/list/enable/disable/update <plugin-name>`

---

## 문제 해결

일반적인 문제:

Plugin 로드 실패:
- `.claude-plugin/plugin.json` 존재 확인
- plugin.json 구문 유효성 검증
- name 필드가 kebab-case인지 확인
- 구성 요소 디렉토리가 루트 수준인지 확인

Commands 미발견:
- command 파일에 .md 확장자 확인
- description이 있는 YAML frontmatter 검증
- plugin.json의 commands 경로 확인
- `/plugin-name:command-name`으로 테스트

Hooks 미트리거:
- hooks.json 구문 검증
- matcher 패턴이 도구 이름과 일치하는지 확인
- hook 명령이 실행 가능한지 확인
- 상세 로그를 위해 디버그 모드 활성화

MCP 서버 실패:
- command가 PATH에 존재하는지 확인
- 환경 변수가 올바르게 설정되었는지 확인
- transport 유형이 서버와 일치하는지 확인
- 먼저 서버를 독립적으로 테스트

---

## 관련 Skill

- do-foundation-claude - Claude Code 설정과 패턴
- do-foundation-core - 핵심 개발 워크플로우
- do-workflow-project - 프로젝트 초기화
- do-domain-backend - 백엔드 plugin 개발
- do-domain-frontend - 프론트엔드 plugin 개발

---

상태: Production Ready
최종 업데이트: 2026-01-06
관리: Do Team
버전 변경: v1.1.0 - PostToolUseFailure, SubagentStart hook 이벤트 추가; agent hook 유형 추가; LSP 고급 옵션 추가; managed 설치 범위 추가
