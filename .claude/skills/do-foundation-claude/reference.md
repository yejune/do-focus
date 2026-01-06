# do-foundation-claude 참조 문서

## API 참조

### Skill 정의 API

Frontmatter 필드:
- `name` (필수): kebab-case 형식의 Skill 식별자, 최대 64자
- `description` (필수): 한 줄 설명, 최대 1024자
- `version`: 시맨틱 버전 (예: "2.0.0")
- `tools`: 허용된 도구의 쉼표 구분 목록
- `modularized`: 모듈화된 파일 구조 여부 (boolean)
- `category`: Skill 카테고리 (foundation, domain, workflow, library, integration)
- `tags`: 검색 가능한 키워드 배열
- `aliases`: Skill 호출을 위한 대체 이름

### Sub-agent 위임 API

Task 호출:
- `Task(subagent_type, prompt)`: 전문 sub-agent 호출
- `Task(subagent_type, prompt, context)`: 이전 작업의 context와 함께 호출
- 체이닝을 위한 구조화된 결과 객체 반환

사용 가능한 Sub-agent 유형:
- `spec-builder`: EARS 형식 specification 생성
- `tdd-implementer`: RED-GREEN-REFACTOR TDD 실행
- `backend-expert`: Backend 아키텍처 및 API 개발
- `frontend-expert`: Frontend UI 구현
- `security-expert`: 보안 분석 및 검증
- `docs-manager`: 기술 문서 생성
- `quality-gate`: TRUST 5 검증
- `agent-factory`: 새 sub-agent 생성
- `skill-factory`: 규정 준수 skill 생성

### Command 매개변수 API

매개변수 유형:
- `$1`, `$2`, `$3`: 위치 인자
- `$ARGUMENTS`: 모든 인자를 단일 문자열로
- `@filename`: 파일 내용 주입

Command 위치:
- 개인: `~/.claude/commands/`
- 프로젝트: `.claude/commands/`

---

## 설정 옵션

### Settings 계층

우선순위 (높은 순서):
1. Enterprise 설정 (`/etc/claude/settings.json`)
2. User 설정 (`~/.claude/settings.json`)
3. Project 설정 (`.claude/settings.json`)
4. Local 설정 (`.claude/settings.local.json`)

### 도구 권한

권한 수준:
- `Read, Grep, Glob`: 분석용 읽기 전용 접근
- `Read, Write, Edit, Grep, Glob`: 전체 파일 조작
- `Bash`: 시스템 명령 실행 (명시적 허가 필요)
- `WebFetch, WebSearch`: 외부 웹 접근

### Memory 설정

Memory 파일 위치:
- Enterprise: `/etc/claude/CLAUDE.md`
- User: `~/.claude/CLAUDE.md`
- Project: `./CLAUDE.md` 또는 `.claude/CLAUDE.md`

Memory Import 문법:
```markdown
@import path/to/file.md
```

---

## 통합 패턴

### Command-Agent-Skill 오케스트레이션

순차 패턴:
1. Command가 `$ARGUMENTS`로 사용자 입력 수신
2. Command가 `Skill("skill-name")`으로 관련 Skill 로드
3. Command가 `Task(subagent_type, prompt)`로 sub-agent에 위임
4. Sub-agent가 로드된 skill context로 실행
5. 결과가 command로 반환되어 표시

병렬 패턴:
- 여러 독립 `Task()` 호출이 동시 실행
- 모든 완료 후 결과 집계
- 작업 간 의존성이 없을 때 사용

### Hook 통합

PreToolUse Hook:
- 모든 도구 호출 전에 실행
- 도구 실행을 차단하거나 수정 가능
- 검증, 로깅, 보안 검사에 사용

PostToolUse Hook:
- 도구 완료 후 실행
- 결과를 처리하거나 수정 가능
- 백업, 감사, 알림에 사용

Hook 설정 (settings.json):
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "ToolName",
        "hooks": [{"type": "command", "command": "hook-script"}]
      }
    ]
  }
}
```

### MCP 서버 통합

Context7 통합:
- 실시간 문서 조회에 사용
- 2단계 패턴: library ID 해석 후 문서 가져오기
- token 제한 응답 지원

MCP 도구 호출:
- MCP 제공 기능에는 `mcp__` 접두사
- settings.json에서 서버 설정

---

## 문제 해결

### Skill 로드 안됨

증상: Skill 인식 안됨, context 누락

해결책:
1. 파일 위치 확인 (`~/.claude/skills/` 또는 `.claude/skills/`)
2. SKILL.md frontmatter 문법 확인 (유효한 YAML)
3. 이름이 kebab-case, 최대 64자인지 확인
4. 파일 크기가 500줄 이하인지 확인

### Sub-agent 위임 실패

증상: Task() 오류 반환, 불완전한 결과

해결책:
1. subagent_type이 유효한지 확인
2. prompt의 명확성과 구체성 확인
3. 필요한 context가 제공되었는지 확인
4. token 예산 검토 (각 Task()는 200K 할당)

### Hook 실행 안됨

증상: PreToolUse/PostToolUse 트리거 안됨

해결책:
1. matcher 패턴이 도구 이름과 정확히 일치하는지 확인
2. hook 스크립트가 존재하고 실행 가능한지 확인
3. settings.json 문법 검토
4. command 권한 확인

### Memory 파일 문제

증상: CLAUDE.md 내용 적용 안됨

해결책:
1. 파일이 올바른 계층 위치에 있는지 확인
2. 파일 인코딩 확인 (UTF-8 필수)
3. @import 경로 검토 (파일 기준 상대 경로)
4. 파일 읽기 권한 확인

---

## 외부 리소스

### 공식 문서

- [Claude Code Skills 가이드](https://docs.anthropic.com/claude-code/skills)
- [Sub-agent 문서](https://docs.anthropic.com/claude-code/agents)
- [Custom Commands 참조](https://docs.anthropic.com/claude-code/commands)
- [Hook 시스템 가이드](https://docs.anthropic.com/claude-code/hooks)
- [Memory 관리](https://docs.anthropic.com/claude-code/memory)
- [Settings 설정](https://docs.anthropic.com/claude-code/settings)
- [IAM 및 권한](https://docs.anthropic.com/claude-code/iam)

### 모범 사례

- SKILL.md를 500줄 이하로 유지
- Progressive disclosure 사용 (빠른 참고, 구현, 고급)
- 최소 권한 도구 권한 적용
- description에 트리거 시나리오 문서화
- 각 패턴에 대한 작동 예제 포함

### 관련 Skills

- `do-foundation-core`: 핵심 실행 패턴 및 SPEC 워크플로우
- `do-foundation-context`: Token 예산 및 세션 관리
- `do-workflow-project`: 프로젝트 초기화 및 설정
- `do-docs-generation`: 문서 자동화

---

Version: 2.0.0
Last Updated: 2025-12-06
