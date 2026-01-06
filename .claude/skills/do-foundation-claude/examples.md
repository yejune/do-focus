# Claude Code 작성 키트 예제

Skills, sub-agent, 커스텀 slash command, hook, memory, settings에 대한 실용적인 Claude Code 패턴 예제

---

## 예제 1: 커스텀 Skill 생성

**시나리오**: Python 테스팅 모범 사례를 위한 skill 생성

**입력**:
```yaml
# 필요: 다음 조건을 충족하는 Python 테스팅 skill
# - Claude Code 공식 표준 준수
# - Progressive disclosure 사용
# - 500줄 이하 유지
# - 실용적인 예제 포함
```

**출력**:
```yaml
# 파일: .claude/skills/python-testing/SKILL.md
---
name: python-testing
description: pytest, coverage, TDD 패턴을 활용한 Python 테스팅 모범 사례
version: 1.0.0
updated: 2025-11-26
status: active
tools: Read, Write, Edit, Bash, Grep, Glob
---

# Python 테스팅 전문가

## 빠른 참고 (30초)

Python Testing Mastery - pytest, coverage 분석, TDD 워크플로우를 포함한 종합 테스팅 패턴

핵심 기능:
- 단위 테스트: pytest fixture, 파라미터화, mocking
- 통합 테스트: 데이터베이스, API, 서비스 테스팅
- Coverage: pytest-cov 설정 및 리포팅
- TDD 워크플로우: RED-GREEN-REFACTOR 사이클 구현

사용 시점:
- Python 모듈용 단위 테스트 작성
- pytest 설정 구성
- 새 기능에 TDD 적용
- coverage 목표 달성

---

## 구현 가이드 (5분)

### pytest 설정

```toml
# pyproject.toml
[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py", "*_test.py"]
python_functions = ["test_*"]
addopts = "-v --tb=short --strict-markers"
markers = [
    "slow: marks tests as slow",
    "integration: marks integration tests"
]

[tool.coverage.run]
source = ["src"]
branch = true
omit = ["*/tests/*", "*/__init__.py"]

[tool.coverage.report]
fail_under = 80
show_missing = true
```

---

## 함께 잘 작동하는 것들

- do-lang-python - Python 3.13+ 패턴
- do-domain-backend - Backend 테스팅 전략
- do-workflow-tdd - TDD 워크플로우 통합
```

**설명**: 이 skill은 progressive disclosure (빠른 참고/구현/고급)를 사용한 Claude Code 공식 표준을 따르며, 500줄 이하를 유지하고 실용적인 pytest 예제를 포함

---

## 예제 2: Sub-agent 위임 패턴

**시나리오**: SPEC 워크플로우를 위한 효율적인 sub-agent 위임 구현

**입력**:
```
// 필요: SPEC 구현을 위한 여러 agent 조정
// 요구사항:
// - 의존 작업에는 순차 실행
// - 독립 작업에는 병렬 실행
// - 적절한 에러 처리
// - Agent 간 context 전달
```

**출력**:

### 순차 위임 패턴

```
// Phase 1: spec-builder로 분석
analysis = Task(
    subagent_type: "spec-builder",
    prompt: "다음 요구사항을 분석하고 SPEC 생성:
    요구사항: JWT 토큰 기반 사용자 인증 시스템
    출력 형식: SPEC ID, 요구사항 목록, 아키텍처 개요, 테스트 기준",
    context: {
        project_type: "web_api",
        language: "python",
        framework: "fastapi"
    }
)

// Phase 2: tdd-implementer로 구현 (analysis에 의존)
implementation = Task(
    subagent_type: "tdd-implementer",
    prompt: "TDD 접근법으로 SPEC 구현:
    SPEC ID: " + analysis.spec_id + "
    RED-GREEN-REFACTOR 사이클 따르기",
    context: {
        spec_id: analysis.spec_id,
        architecture: analysis.architecture
    }
)

// Phase 3: quality-gate로 검증 (implementation에 의존)
validation = Task(
    subagent_type: "quality-gate",
    prompt: "구현 검증:
    SPEC ID: " + implementation.spec_id + "
    확인: 모든 테스트 통과, Coverage 80% 이상, 보안 이슈 없음",
    context: {
        implementation: implementation,
        original_spec: analysis
    }
)
```

### 병렬 위임 패턴

```
// 독립 작업은 동시 실행으로 3배 빠른 처리

// 세 작업 모두 병렬 실행 가능
results = await Promise.all([
    // Backend 구현
    Task(
        subagent_type: "backend-expert",
        prompt: "SPEC-001용 API 엔드포인트 구현",
        context: {spec_id: "SPEC-001", focus: "api"}
    ),
    // Frontend 구현
    Task(
        subagent_type: "frontend-expert",
        prompt: "SPEC-001용 UI 컴포넌트 구현",
        context: {spec_id: "SPEC-001", focus: "ui"}
    ),
    // 문서 생성
    Task(
        subagent_type: "docs-manager",
        prompt: "SPEC-001용 API 문서 생성",
        context: {spec_id: "SPEC-001", focus: "docs"}
    )
])

// 결과 집계
backend_result = results[0]
frontend_result = results[1]
docs_result = results[2]

// 통합 검증 (순차 - 모든 병렬 작업에 의존)
integration = Task(
    subagent_type: "quality-gate",
    prompt: "모든 컴포넌트 통합 검증",
    context: {
        backend: backend_result,
        frontend: frontend_result,
        docs: docs_result
    }
)
```

**설명**: 이 패턴은 의존 작업에 순차 위임, 독립 작업에 병렬 실행, 탄력적 워크플로우를 위한 적절한 에러 처리를 보여줌

---

## 예제 3: Hook이 있는 커스텀 Slash Command

**시나리오**: 검증용 pre/post hook이 있는 배포 command 생성

**입력**:
```
// 필요: 다음 기능의 배포 command
// - 배포 전 검증
// - 배포 전 테스트 실행
// - 완료 시 알림
// - 롤백 지원
```

**출력**:
```markdown
# 파일: .claude/commands/deploy.md
---
description: 검증, 테스팅, 알림이 포함된 애플리케이션 배포
allowed-tools: Bash, Read, Write, Glob
---

# 배포 Command

지정된 환경에 애플리케이션 배포

## 인자

- $1: 환경 (staging/production)
- $2: 버전 태그 (선택사항, 기본값 latest)

## 실행 흐름

### Step 1: 배포 전 검증

배포 요구사항 확인

### Step 2: 테스트 실행

배포 전 전체 테스트 스위트 실행

### Step 3: 빌드 및 배포

환경에 빌드 및 푸시

### Step 4: 헬스 체크

배포 성공 확인

### Step 5: 알림

팀에 배포 알림
```

```json
// 파일: .claude/settings.json (hooks 섹션)
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "validate-bash-command",
            "description": "실행 전 bash 명령 검증"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "git add $FILE",
            "description": "작성된 파일 자동 스테이징"
          }
        ]
      }
    ]
  }
}
```

**설명**: 이 패턴은 검증, 테스팅, 알림을 위한 hook과 커스텀 slash command를 결합하여 완전한 배포 워크플로우 생성

---

## 공통 패턴

### 패턴 1: Memory 파일 구성

효율적인 context 로딩을 위한 memory 구성:

```markdown
# 파일: .claude/CLAUDE.md (프로젝트 레벨 memory)

## 프로젝트 개요
- 이름: MyApp
- 유형: Web API
- 스택: Python 3.13, FastAPI, PostgreSQL

## 개발 가이드라인
- 모든 새 기능에 TDD 적용
- 최소 80% 테스트 coverage
- 모든 곳에 타입 힌트 사용

## 활성 SPEC
- SPEC-001: 사용자 인증 (진행 중)
- SPEC-002: API 속도 제한 (계획됨)

@import architecture.md
@import coding-standards.md
```

### 패턴 2: Settings 계층

적절한 레벨에서 설정 구성:

```json
// ~/.claude/settings.json (사용자 레벨)
{
  "preferences": {
    "outputStyle": "concise",
    "codeStyle": "modern"
  },
  "permissions": {
    "allowedTools": ["Read", "Write", "Edit", "Bash"]
  }
}
```

```json
// .claude/settings.json (프로젝트 레벨)
{
  "model": "claude-sonnet-4-5-20250929",
  "permissions": {
    "allow": ["Read", "Write", "Edit"],
    "deny": ["Bash dangerous commands"]
  },
  "hooks": {
    "PreToolUse": [...]
  }
}
```

### 패턴 3: IAM 권한 계층

Agent 역할에 따른 권한 정의:

## 권한 계층

### Tier 1: 읽기 전용 Agent
- 도구: Read, Grep, Glob
- 용도: 코드 분석, 문서 검토
- 예시 agent: code-analyzer, doc-reviewer

### Tier 2: 쓰기 제한 Agent
- 도구: Read, Write, Edit, Grep, Glob
- 제한: 프로덕션 파일 수정 불가
- 용도: 코드 생성, 리팩토링
- 예시 agent: code-generator, refactorer

### Tier 3: 전체 접근 Agent
- 도구: Bash 포함 전체
- 제한: 위험한 명령은 승인 필요
- 용도: 배포, 시스템 관리
- 예시 agent: deployer, admin

### Tier 4: 관리자 Agent
- 도구: 상승된 권한으로 전체
- 용도: 시스템 설정, 보안
- 예시 agent: security-auditor, config-manager

---

## 안티패턴 (피해야 할 패턴)

### 안티패턴 1: 모놀리식 Skills

**문제**: 500줄을 초과하는 skill은 유지보수와 로딩이 어려움

**해결책**: 교차 참조가 있는 집중된 skill로 분할

### 안티패턴 2: 중첩 Sub-agent 스폰

**문제**: Sub-agent가 다른 sub-agent를 스폰하면 context 문제 발생

**해결책**: 모든 sub-agent 위임은 메인 스레드에서만

### 안티패턴 3: Skill에 하드코딩된 경로

**문제**: 하드코딩된 경로는 이식성을 해침

**해결책**: 상대 경로와 프로젝트 참조 사용

---

## 통합 예제

### 완전한 SPEC 워크플로우

```
// 완전한 SPEC-First TDD 워크플로우

// Step 1: Plan - SPEC 생성
plan_result = Task(
    subagent_type: "spec-builder",
    prompt: "SPEC 생성: 아바타 업로드가 있는 사용자 프로필 관리",
    context: {project: "@CLAUDE.md"}
)

// Step 2: context 정리 (plan 후)
// /clear

// Step 3: Run - TDD로 구현
run_result = Task(
    subagent_type: "tdd-implementer",
    prompt: "SPEC 구현: " + plan_result.spec_id,
    context: {spec: plan_result}
)

// Step 4: Sync - 문서 생성
sync_result = Task(
    subagent_type: "docs-manager",
    prompt: "문서 생성: " + run_result.spec_id,
    context: {implementation: run_result}
)
```

### Hook 기반 품질 보증

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write",
        "hooks": [{"type": "command", "command": "lint-check $FILE"}]
      },
      {
        "matcher": "Bash",
        "hooks": [{"type": "command", "command": "validate-command $COMMAND"}]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Write",
        "hooks": [{"type": "command", "command": "run-tests --affected $FILE"}]
      }
    ]
  }
}
```

---

*전체 참조 문서는 reference/ 디렉토리 참조*
