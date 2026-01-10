---
name: do-workflow-spec
description: EARS 형식, 요구사항 명확화, Plan-Run-Sync 통합을 통한 SPEC workflow 오케스트레이션
version: 2.1.0
category: workflow
tags:
  - workflow
  - spec
  - ears
  - requirements
  - planning
updated: 2026-01-06
status: active
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
user-invocable: true
---

# SPEC Workflow 관리

## 빠른 참조

SPEC-First 개발 방법론: EARS 형식으로 체계적 요구사항 정의, Plan-Run-Sync workflow 통합.

핵심 기능:
- EARS Format Specifications: 5가지 요구사항 패턴으로 명확성 확보
- Requirement Clarification: 4단계 체계적 프로세스
- SPEC Document Templates: 일관된 구조의 표준화된 문서
- Plan-Run-Sync Integration: 원활한 workflow 연결
- Parallel Development: Git Worktree 기반 SPEC 격리
- Quality Gates: TRUST 5 framework 검증

사용 시점:
- 기능 계획 및 요구사항 정의
- SPEC 문서 생성 및 유지보수
- 병렬 기능 개발 조율
- 품질 보증 및 검증 계획

빠른 명령어:
- SPEC 생성: /do:1-plan "사용자 인증 시스템"
- Worktree 병렬 SPEC: /do:1-plan "로그인" "가입" --worktree
- 새 branch와 SPEC: /do:1-plan "결제 처리" --branch
- 기존 SPEC 업데이트: /do:1-plan SPEC-001 "OAuth 지원 추가"

---

## EARS 5 패턴

### Ubiquitous - 항상 활성

용도: 시스템 전체 품질 속성
예시: 로깅, 입력 검증, 에러 처리
테스트 전략: 모든 기능 테스트에서 공통 검증으로 포함

```pseudo
// PHP/JS 스타일 슈도 코드
RULE: "시스템은 항상 [동작]을 수행해야 한다"

EXAMPLE_REQUIREMENTS = [
  "시스템은 항상 입력을 검증해야 한다",
  "시스템은 항상 에러를 로깅해야 한다",
  "시스템은 항상 요청을 감사해야 한다"
]
```

### Event-Driven - 트리거-응답

용도: 사용자 상호작용, 시스템 간 통신
예시: 버튼 클릭, 파일 업로드, 결제 완료
테스트 전략: 이벤트 시뮬레이션으로 예상 응답 검증

```pseudo
RULE: "WHEN [이벤트] THEN [동작]"

EXAMPLE_REQUIREMENTS = [
  "WHEN 사용자가 로그인하면 THEN JWT 토큰을 발급한다",
  "WHEN 파일 업로드 완료 THEN 썸네일을 생성한다",
  "WHEN 결제 성공 THEN 확인 이메일을 발송한다"
]
```

### State-Driven - 조건부 동작

용도: 접근 제어, 상태 기계, 조건부 비즈니스 로직
예시: 계정 상태 확인, 재고 검증, 권한 확인
테스트 전략: 상태 설정 후 조건부 동작 검증

```pseudo
RULE: "IF [조건] THEN [동작]"

EXAMPLE_REQUIREMENTS = [
  "IF 계정이 active이면 THEN 로그인을 허용한다",
  "IF 재고가 충분하면 THEN 주문을 진행한다",
  "IF 권한이 있으면 THEN 수정을 허용한다"
]
```

### Unwanted - 금지 동작

용도: 보안 취약점 방지, 데이터 무결성 보호
예시: 평문 비밀번호 금지, 무단 접근 차단, PII 로그 금지
테스트 전략: 금지 동작 시도의 차단 검증

```pseudo
RULE: "시스템은 [금지 동작]하지 않아야 한다"

EXAMPLE_REQUIREMENTS = [
  "시스템은 평문 비밀번호를 저장하지 않아야 한다",
  "시스템은 SQL injection을 허용하지 않아야 한다",
  "시스템은 민감 정보를 로그에 기록하지 않아야 한다"
]
```

### Optional - 선택 기능

용도: MVP scope 정의, 기능 우선순위화
예시: OAuth 로그인, 다크 모드, 오프라인 모드
테스트 전략: 구현 상태에 따라 조건부 테스트 실행

```pseudo
RULE: "가능하면 [기능]을 제공한다"

EXAMPLE_REQUIREMENTS = [
  "가능하면 소셜 로그인을 제공한다",
  "가능하면 다크 모드를 지원한다",
  "가능하면 오프라인 캐싱을 제공한다"
]
```

---

## 요구사항 명확화 프로세스

### Step 0 - 가정 분석 (Philosopher Framework)

Scope 정의 전 AskUserQuestion으로 기저 가정 표면화 및 검증.

가정 분류:
- 기술 가정: 기술 역량, API 가용성, 성능 특성
- 비즈니스 가정: 사용자 행동, 시장 요구사항, 일정 실현 가능성
- 팀 가정: 기술 가용성, 자원 배분, 지식 격차
- 통합 가정: 서드파티 서비스 신뢰성, 호환성 기대치

가정 문서화:
- 가정 내용: 가정하는 바의 명확한 기술
- 신뢰 수준: 증거 기반 High, Medium, Low
- 증거 근거: 가정을 뒷받침하는 정보
- 오류 시 위험: 가정이 틀릴 경우의 결과
- 검증 방법: 중요 노력 투입 전 확인 방법

### Step 0.5 - 근본 원인 분석

기능 요청 또는 문제 기반 SPEC에 Five Whys 적용:
- 표면 문제: 사용자가 관찰하거나 요청하는 것
- 첫 번째 Why: 이 요청을 유발하는 즉각적 필요
- 두 번째 Why: 그 필요를 만드는 근본적 문제
- 세 번째 Why: 기여하는 시스템적 요인
- 근본 원인: 솔루션이 해결해야 할 근본적 이슈

### Step 1 - Scope 정의

- 지원 인증 방법 식별
- 검증 규칙 및 제약 정의
- 실패 처리 전략 결정
- 세션 관리 접근 방식 수립

### Step 2 - 제약 추출

- 성능 요구사항: 응답 시간 목표
- 보안 요구사항: OWASP 준수, 암호화 표준
- 호환성 요구사항: 지원 브라우저 및 디바이스
- 확장성 요구사항: 동시 사용자 목표

### Step 3 - 성공 기준 정의

- 테스트 커버리지: 최소 퍼센트 목표
- 응답 시간: 백분위 목표 (P50, P95, P99)
- 기능 완료: 모든 정상 시나리오 검증 통과
- 품질 게이트: 린터 경고 0, 보안 취약점 0

### Step 4 - 테스트 시나리오 생성

- 정상 케이스: 유효 입력으로 예상 출력
- 에러 케이스: 유효하지 않은 입력의 에러 처리
- 경계 케이스: 경계 조건 및 코너 케이스
- 보안 케이스: 인젝션 공격, 권한 상승 시도

---

## Plan-Run-Sync Workflow 통합

### PLAN Phase (/do:1-plan)

- manager-spec agent가 사용자 입력 분석
- EARS 형식 요구사항 생성
- 사용자 상호작용으로 요구사항 명확화
- .do/specs/ 디렉토리에 SPEC 문서 생성
- Git branch 생성 (선택적 --branch 플래그)
- Git Worktree 설정 (선택적 --worktree 플래그)

### RUN Phase (/do:2-run)

- manager-tdd agent가 SPEC 문서 로드
- RED-GREEN-REFACTOR TDD cycle 실행
- do-workflow-testing skill 참조로 테스트 패턴 활용
- Domain Expert agent 위임 (expert-backend, expert-frontend 등)
- manager-quality agent로 품질 검증

### SYNC Phase (/do:3-sync)

- manager-docs agent가 문서 동기화
- SPEC에서 API 문서 생성
- README 및 아키텍처 문서 업데이트
- CHANGELOG entry 생성
- SPEC 참조로 버전 관리 커밋

---

## Git Worktree 병렬 개발

### Worktree 개념

- 여러 branch를 위한 독립적 작업 디렉토리
- 각 SPEC에 격리된 개발 환경 제공
- 병렬 작업에 branch 전환 불필요
- 기능 격리로 merge conflict 감소

### Worktree 생성

명령어 `/do:1-plan "로그인 기능" "가입 기능" --worktree`로 여러 SPEC 생성:
- project-worktrees 디렉토리에 SPEC별 하위 디렉토리 생성

### Worktree 이점

- 병렬 개발: 여러 기능 동시 개발
- 팀 협업: SPEC별 명확한 소유권 경계
- 의존성 격리: 기능별 다른 라이브러리 버전
- 위험 감소: 불안정 코드가 다른 기능에 영향 없음

---

## SPEC 문서 구조

### SPEC Metadata Schema

필수 필드:
- SPEC ID: 순번 (SPEC-001, SPEC-002 등)
- Title: 영문 기능명
- Created: ISO 8601 timestamp
- Status: Planned, In Progress, Completed, Blocked
- Priority: High, Medium, Low
- Assigned: 구현 담당 Agent

선택 필드:
- Related SPECs: 의존성 및 관련 기능
- Epic: 상위 기능 그룹
- Estimated Effort: 시간 또는 스토리 포인트
- Labels: 분류용 태그

### Lifecycle Level

- spec-first: 구현 전 SPEC 작성, 완료 후 폐기 (일회성, 프로토타입)
- spec-anchored: 구현과 함께 SPEC 유지 (핵심 기능, API 계약)
- spec-as-source: SPEC이 single source of truth (규제 환경, 코드 생성)

### 디렉토리 구조

- .do/specs/: SPEC 문서 파일 (SPEC-001-feature-name.md)
- .do/memory/: 세션 상태 파일 (last-session-state.json)
- .do/docs/: 생성된 문서 (api-documentation.md)

---

## Constitution 정렬

SPEC 생성 전 `.do/project/tech.md`의 프로젝트 constitution 확인.

Constitution 구성요소:
- Technology Stack: 필수 버전 및 framework
- Naming Conventions: 변수, 함수, 파일 명명 표준
- Forbidden Libraries: 금지 라이브러리 및 대안
- Architectural Patterns: 계층화 규칙 및 의존성 방향
- Security Standards: 인증 패턴 및 암호화 요구사항
- Logging Standards: 로그 형식 및 구조화 로깅 요구사항

Constitution 검증:
- 모든 SPEC 기술 선택이 Constitution stack 버전과 정렬
- SPEC이 금지 라이브러리나 패턴 도입 금지
- SPEC이 Constitution 명명 규칙 준수
- SPEC이 아키텍처 경계 및 계층화 존중

---

## 품질 지표

### SPEC 품질 지표

- 요구사항 명확성: 모든 EARS 패턴 적절히 사용
- 테스트 커버리지: 모든 요구사항에 대응 테스트 시나리오
- 제약 완전성: 기술 및 비즈니스 제약 정의
- 성공 기준 측정 가능성: 정량화된 완료 지표

### 검증 체크리스트

- 모든 EARS 요구사항 테스트 가능
- 모호한 언어 없음 (should, might, usually)
- 모든 에러 케이스 문서화
- 성능 목표 정량화
- 보안 요구사항 OWASP 준수

---

## 통합 예제

### 순차 Workflow

```pseudo
// 단일 기능 순차 처리
USER_REQUEST
  -> /do:1-plan "기능 설명"
  -> manager-spec이 SPEC-001 생성
  -> /clear (토큰 최적화)
  -> /do:2-run SPEC-001
  -> manager-tdd가 RED-GREEN-REFACTOR 실행
  -> /do:3-sync SPEC-001
  -> manager-docs가 문서 업데이트
  -> FEATURE_COMPLETE
```

### 병렬 Workflow

```pseudo
// 여러 기능 병렬 처리
USER_REQUEST
  -> /do:1-plan "feature1" "feature2" "feature3" --worktree
  -> manager-spec이 SPEC-001, SPEC-002, SPEC-003 생성
  -> Git Worktree 병렬 개발 설정
  -> /clear (토큰 최적화)
  -> Session1: /do:2-run SPEC-001
  -> Session2: /do:2-run SPEC-002
  -> Session3: /do:2-run SPEC-003
  -> Worktree merge to main
  -> /do:3-sync SPEC-001 SPEC-002 SPEC-003
  -> ALL_FEATURES_COMPLETE
```

### 의존성 체인

```pseudo
// 순차 의존성 개발
/do:1-plan "database schema" --branch
  -> SPEC-001 생성 (foundation)
  -> /do:2-run SPEC-001
  -> DB schema 구현
/do:1-plan "backend API" --branch
  -> SPEC-002 생성 (SPEC-001 의존)
  -> /do:2-run SPEC-002
  -> Backend API 구현
/do:1-plan "frontend UI" --branch
  -> SPEC-003 생성 (SPEC-002 의존)
  -> /do:2-run SPEC-003
  -> Frontend UI 구현
/do:3-sync SPEC-001 SPEC-002 SPEC-003
  -> FULL_STACK_FEATURE_COMPLETE
```

---

## 토큰 관리

### 세션 전략

- PLAN phase: 세션 토큰의 약 30% 사용
- RUN phase: 세션 토큰의 약 60% 사용
- SYNC phase: 세션 토큰의 약 10% 사용

### 컨텍스트 최적화

- SPEC 문서는 .do/specs/ 디렉토리에 영속
- 세션 메모리는 .do/memory/에 크로스 세션 컨텍스트
- SPEC ID 참조로 최소 컨텍스트 전달
- Agent 위임으로 토큰 오버헤드 감소

---

## 관련 Skills

- do-foundation-core: SPEC-First TDD 방법론, TRUST 5 framework
- do-workflow-testing: TDD 구현 및 테스트 자동화
- do-workflow-project: 프로젝트 초기화 및 설정
- do-workflow-worktree: Git Worktree 병렬 개발 관리
- manager-spec: SPEC 생성 및 요구사항 분석 Agent
- manager-tdd: SPEC 요구사항 기반 TDD 구현
- manager-quality: TRUST 5 품질 검증 및 게이트 강제

---

Version: 2.1.0 (Full Recovery)
Last Updated: 2026-01-06
Integration Status: Complete - Full Plan-Run-Sync workflow
