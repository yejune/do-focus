---
name: do-foundation-claude
aliases: [do-foundation-claude]
category: foundation
description: Claude Code 공식 Skill 작성 키트 - Skills, sub-agent, slash command, hook, memory, settings, IAM 권한 관리
version: 2.0.0
modularized: false
allowed-tools: Read, Write, Edit, Grep, Glob
tags:
  - foundation
  - claude-code
  - skills
  - sub-agents
  - slash-commands
  - orchestration
  - hooks
  - memory
  - settings
  - iam
  - best-practices
user-invocable: true
---

# Claude Code 작성 키트

Claude Code Skills, sub-agent, 커스텀 slash command, hook, memory, settings, IAM 권한에 대한 종합 참고서

공식 문서 참조:
- [Skills 가이드](reference/claude-code-skills-official.md) - Agent Skills 생성 및 관리
- [Sub-agent 가이드](reference/claude-code-sub-agents-official.md) - Sub-agent 개발 및 위임
- [Slash Command](reference/claude-code-custom-slash-commands-official.md) - 커스텀 명령어 생성
- [Hook 시스템](reference/claude-code-hooks-official.md) - 이벤트 기반 자동화
- [Memory 관리](reference/claude-code-memory-official.md) - 컨텍스트 및 지식 유지
- [Settings 설정](reference/claude-code-settings-official.md) - 설정 계층 및 관리
- [IAM 권한](reference/claude-code-iam-official.md) - 접근 제어 및 보안
- [전체 설정 가이드](reference/complete-configuration-guide.md) - 종합 설정

## 빠른 참고 (30초)

Skills: ~/.claude/skills/ (개인) 또는 .claude/skills/ (프로젝트), 500줄 이하, progressive disclosure
Sub-agent: Task(subagent_type="...") 위임 전용, nested spawning 불가
Command: $ARGUMENTS/$1/$2 매개변수, @file 참조, .claude/commands/ 저장
Hook: settings.json의 이벤트 기반 자동화, PreToolUse/PostToolUse 이벤트
Memory: Enterprise, Project, User 계층, @import.md 문법
Settings: 4단계 계층 (Enterprise, User, Project, Local)
IAM: 계층화된 권한 (Read, Bash, Write, Admin), 도구별 규칙

## 구현 가이드 (5분)

### 기능

- 종합적인 Claude Code 작성 참고서
- Skills, sub-agent, command, hook, settings 지원
- IAM 권한 및 보안 모범 사례
- Progressive disclosure 문서 구조
- MCP 통합 패턴 및 예제

### 사용 시점

- 공식 표준을 따르는 새로운 Claude Code skill 생성
- 적절한 위임 패턴으로 커스텀 sub-agent 개발
- 매개변수 처리가 있는 커스텀 slash command 구현
- 이벤트 기반 자동화를 위한 hook 설정
- IAM 권한 및 접근 제어 구성

### 핵심 패턴

패턴 1: Skill 생성 (500줄 이하)

```yaml
---
name: skill-name
description: 간단한 설명 (최대 200자)
version: 1.0.0
updated: 2025-11-26
status: active
---

## 빠른 참고 (30초)
## 구현 가이드 (5분)
## 고급 구현 (10+ 분)
```

패턴 2: Sub-agent 위임

순차 위임:
```
result1 = Task(subagent_type: "workflow-spec", prompt: "분석")
result2 = Task(subagent_type: "code-backend", prompt: "구현", context: result1)
```

병렬 위임:
```
results = await Promise.all([
  Task(subagent_type: "code-backend", prompt: "Backend"),
  Task(subagent_type: "code-frontend", prompt: "Frontend")
])
```

패턴 3: Hook 통합 Command

```json
{
  "hooks": {
    "PreToolUse": [{
      "matcher": "Write",
      "hooks": [{"type": "command", "command": "validate-write"}]
    }]
  }
}
```

## Do Skills 및 Sub-agent 디렉토리

### 빠른 접근 패턴

Core Skills: Skill("do-foundation-claude") - 이 종합 작성 키트
Agent 생성: Task(subagent_type: "agent-factory") - 표준화된 sub-agent 생성
Skill 생성: Task(subagent_type: "skill-factory") - 규정 준수 skill 생성
품질 검증: Task(subagent_type: "quality-gate") - TRUST 5 검증
문서화: Task(subagent_type: "docs-manager") - 기술 문서 생성

핵심 전문 Skills:
- do-lang-python - Python 3.13+ (FastAPI, async 패턴)
- do-lang-typescript - TypeScript 5.9+ (React 19, Next.js 16)
- do-domain-backend - 현대적 backend 아키텍처 패턴
- do-domain-frontend - React 19, Next.js 16, Vue 3.5 프레임워크
- do-quality-security - OWASP Top 10, 위협 모델링
- do-essentials-debug - AI 기반 디버깅 (Context7 활용)

필수 Sub-agent:
- spec-builder - EARS 형식 specification 생성
- tdd-implementer - RED-GREEN-REFACTOR TDD 실행
- security-expert - 보안 분석 및 검증
- backend-expert - Backend 아키텍처 및 API 개발
- frontend-expert - Frontend UI 구현
- performance-engineer - 성능 최적화 및 분석

## 필수 구현 패턴

### Command, Agent, Skill 오케스트레이션

순차 워크플로우:
```
// Phase 1: 분석 -> spec-builder -> analysis
analysis = Task(subagent_type: "spec-builder", prompt: "Analyze: $ARGUMENTS")
// Phase 2: 구현 -> tdd-implementer -> code + tests
implementation = Task(subagent_type: "tdd-implementer", prompt: "Implement: " + analysis.spec_id)
// Phase 3: 검증 -> quality-gate -> approval
validation = Task(subagent_type: "quality-gate", prompt: "Validate: " + implementation)
```

병렬 개발:
```
// 독립 작업 동시 실행
results = await Promise.all([
  Task(subagent_type: "backend-expert", prompt: "Backend: $1"),
  Task(subagent_type: "frontend-expert", prompt: "Frontend: $1"),
  Task(subagent_type: "docs-manager", prompt: "Docs: $1")
])
// 모든 결과 통합
Task(subagent_type: "quality-gate", prompt: "Integrate", context: {results: results})
```

### Token 세션 관리

독립적인 200K 세션: 각 Task()는 새로운 200K token context 생성

```
Task(subagent_type: "backend-expert", prompt: "복잡한 작업")  // 200K 세션
Task(subagent_type: "frontend-expert", prompt: "UI 작업")     // 새로운 200K 세션
```

### 파일 참조 패턴

매개변수 처리:
- 위치: $1, $2, $3
- 전체 인자: $ARGUMENTS
- 파일 참조: @config.yaml, @path/to/file.md

### Hook 통합

Pre/Post 도구 실행 (Hook 가이드 참조):
```json
{
  "hooks": {
    "PreToolUse": [{"matcher": "Bash", "hooks": [{"type": "command", "command": "validate-command"}]}],
    "PostToolUse": [{"matcher": "Write", "hooks": [{"type": "command", "command": "backup-file"}]}]
  }
}
```

## 핵심 참고 가이드

종합 문서:
- [Commands 가이드](reference/claude-code-custom-slash-commands-official.md) - 완전한 command 생성
- [Hook 시스템](reference/claude-code-hooks-official.md) - 이벤트 기반 자동화
- [Memory 관리](reference/claude-code-memory-official.md) - 컨텍스트 유지
- [Settings 설정](reference/claude-code-settings-official.md) - 설정 계층
- [IAM 권한](reference/claude-code-iam-official.md) - 접근 제어
- [전체 설정](reference/complete-configuration-guide.md) - 종합 설정

## 함께 잘 작동하는 것들

- do-core-agent-factory - 적절한 표준으로 새 sub-agent 생성
- do-cc-commands - 커스텀 slash command 생성
- do-cc-hooks - Claude Code hook 구현
- do-cc-configuration - Claude Code settings 관리
- do-quality-gate - 코드 품질 표준 검증
- do-essentials-debug - skill 로딩 문제 디버깅
