---
name: do-foundation-core
description: Do의 기초 원칙 - TRUST 5, SPEC-First TDD, 위임 패턴, 토큰 최적화, 점진적 공개, 모듈 아키텍처를 포함한 AI 기반 개발 워크플로우의 핵심 체계
version: 2.3.0
modularized: true
updated: 2025-12-06
status: active
tags:
  - foundation
  - core
  - orchestration
  - agents
  - commands
  - trust-5
  - spec-first-tdd
allowed-tools: Read, Grep, Glob
---

# Do Foundation Core

AI 기반 개발 워크플로우의 품질, 효율성, 확장성을 보장하는 6가지 핵심 원칙

핵심 철학: 품질 우선, 테스트 주도, 모듈화, 검증된 패턴과 자동화를 통한 효율적 개발

## 빠른 참고 (30초)

Do Foundation Core란?
품질, 효율, 확장성을 위한 6가지 필수 원칙:

1. TRUST 5 Framework - 품질 게이트 (Test-first, Readable, Unified, Secured, Trackable)
2. SPEC-First TDD - 명세 기반 테스트 주도 개발 워크플로우
3. Delegation Patterns - 전문 에이전트 위임 (직접 실행 금지)
4. Token Optimization - 200K 예산 관리 및 컨텍스트 효율화
5. Progressive Disclosure - 3단계 지식 전달 (30초, 5분, 10분+)
6. Modular System - 확장 가능한 파일 분할 아키텍처

빠른 접근:
- 품질 표준: [TRUST 5 Module](modules/trust-5-framework.md)
- 개발 흐름: [SPEC-First TDD Module](modules/spec-first-tdd.md)
- 에이전트 조정: [Delegation Patterns Module](modules/delegation-patterns.md)
- 예산 관리: [Token Optimization Module](modules/token-optimization.md)
- 콘텐츠 구조: [Progressive Disclosure Module](modules/progressive-disclosure.md)
- 파일 조직: [Modular System Module](modules/modular-system.md)
- 에이전트 목록: [Agents Reference Module](modules/agents-reference.md)
- 명령어 참조: [Commands Reference Module](modules/commands-reference.md)
- 보안 규칙: [Execution Rules Module](modules/execution-rules.md)

활용 사례:
- 새 에이전트 생성 시 품질 표준 적용
- 새 스킬 개발 시 구조 가이드라인 참조
- 복잡한 워크플로우 조율
- 토큰 예산 계획 및 최적화
- 문서 아키텍처 설계
- 품질 게이트 구성

---

## 구현 가이드 (5분)

### 1. TRUST 5 Framework - 품질 보증 체계

목적: 코드 품질, 보안, 유지보수성을 보장하는 자동화된 품질 게이트

5가지 기둥:

Test-first 기둥:
- 요구사항: 테스트 커버리지 85% 이상 유지
- 검증: 커버리지 리포트 자동 실행
- 실패 시: 머지 차단, 누락 테스트 생성
- WHY: 높은 커버리지로 코드 신뢰성 확보, 프로덕션 결함 감소
- IMPACT: 버그 조기 발견, 디버깅 시간 60-70% 단축

Readable 기둥:
- 요구사항: 명확하고 설명적인 네이밍 규칙
- 검증: 린터 체크 실행
- 실패 시: 경고 및 리팩토링 제안
- WHY: 명확한 네이밍으로 코드 이해도와 팀 협업 개선
- IMPACT: 온보딩 시간 40% 단축, 유지보수 속도 향상

Unified 기둥:
- 요구사항: 일관된 포맷팅과 임포트 패턴
- 검증: 포맷터 및 정렬 도구 체크
- 실패 시: 자동 포맷 또는 경고
- WHY: 일관성으로 스타일 논쟁과 머지 충돌 제거
- IMPACT: 코드 리뷰 시간 30% 단축, 가독성 향상

Secured 기둥:
- 요구사항: OWASP 보안 표준 준수
- 검증: 보안 전문 에이전트 분석
- 실패 시: 머지 차단, 보안 리뷰 필수
- WHY: 보안 취약점은 치명적 비즈니스 및 법적 위험 초래
- IMPACT: 일반 보안 취약점 95%+ 방지

Trackable 기둥:
- 요구사항: 명확하고 구조화된 커밋 메시지
- 검증: Git 커밋 메시지 패턴 매칭
- 실패 시: 적절한 커밋 메시지 형식 제안
- WHY: 명확한 이력으로 디버깅, 감사, 협업 지원
- IMPACT: 이슈 조사 시간 50% 단축

통합 지점:
- Pre-commit hooks: 자동 검증
- CI/CD 파이프라인: 품질 게이트 강제
- 에이전트 워크플로우: core-quality 검증
- 문서화: 품질 메트릭

상세 참조: [TRUST 5 Framework Module](modules/trust-5-framework.md)

---

### 2. SPEC-First TDD - 개발 워크플로우

목적: 구현 전 명확한 요구사항 정의를 보장하는 명세 기반 개발

3단계 워크플로우:

```
Phase 1: SPEC (/do:1-plan)
  workflow-spec → EARS format
  출력: .do/specs/SPEC-XXX/spec.md
  /clear 실행 (45-50K 토큰 절약)

Phase 2: TDD (/do:2-run)
  RED: 실패 테스트
  GREEN: 통과 코드
  REFACTOR: 최적화
  검증: 85%+ 커버리지

Phase 3: Docs (/do:3-sync)
  API 문서
  아키텍처 다이어그램
  프로젝트 리포트
```

EARS 형식:
- Ubiquitous: 시스템 전체 (항상 활성)
- Event-driven: 트리거 기반 (X 발생 시 Y 실행)
- State-driven: 조건부 (X 상태일 때 Y 실행)
- Unwanted: 금지 사항 (X 금지)
- Optional: 선택 사항 (가능하면 X 실행)

토큰 예산: SPEC 30K, TDD 180K, Docs 40K, 총합 250K

핵심 관행: Phase 1 완료 후 /clear 실행하여 컨텍스트 초기화

상세 참조: [SPEC-First TDD Module](modules/spec-first-tdd.md)

---

### 3. Delegation Patterns - 에이전트 위임

목적: 직접 실행 없이 전문 에이전트에 작업 위임

핵심 원칙 [HARD]: Alfred는 모든 작업을 Task()를 통해 전문 에이전트에 위임해야 함

WHY: 직접 실행은 전문화, 품질 게이트, 토큰 최적화를 우회
IMPACT: 적절한 위임으로 작업 성공률 40% 향상, 병렬 실행 가능

위임 구문 (슈도 코드):
```
result = Task(
  agent_type: "specialized_agent",
  prompt: "명확하고 구체적인 작업",
  context: { relevant: "data" }
)
```

3가지 패턴:

순차 패턴 (의존성 있을 때):
```
design = Task(agent: "api-designer", prompt: "API 설계")
code = Task(agent: "backend-expert", prompt: "구현", context: { design: design })
```

병렬 패턴 (독립적 작업):
```
results = Promise.all([
  Task(agent: "backend-expert", prompt: "Backend"),
  Task(agent: "frontend-expert", prompt: "Frontend")
])
```

조건부 패턴 (분석 기반):
```
analysis = Task(agent: "debug-helper", prompt: "분석")
if (analysis.type == "security")
  Task(agent: "security-expert", prompt: "수정")
```

에이전트 선택:
- 간단 (1 파일): 1-2 에이전트 순차
- 중간 (3-5 파일): 2-3 에이전트 순차
- 복잡 (10+ 파일): 5+ 에이전트 혼합

상세 참조: [Delegation Patterns Module](modules/delegation-patterns.md)

---

### 4. Token Optimization - 예산 관리

목적: 전략적 컨텍스트 관리로 200K 토큰 예산 효율화

예산 배분:

SPEC Phase:
- 토큰 예산: 30K
- 전략: 요구사항만 로드, 완료 후 /clear 실행
- WHY: 명세 단계는 요구사항 분석에 최소 컨텍스트만 필요
- IMPACT: 구현 단계를 위해 45-50K 토큰 절약

TDD Phase:
- 토큰 예산: 180K
- 전략: 선택적 파일 로딩, 구현 관련 파일만 로드
- WHY: 구현은 깊은 컨텍스트 필요하지만 전체 코드베이스는 불필요
- IMPACT: 예산 내 70% 더 큰 구현 가능

Docs Phase:
- 토큰 예산: 40K
- 전략: 결과 캐싱 및 템플릿 재사용
- WHY: 문서화는 완료된 작업 산출물 기반
- IMPACT: 중복 파일 읽기 60% 감소

총 예산:
- 결합 예산: 모든 단계에서 250K 토큰
- 전략: 단계 간 컨텍스트 리셋으로 Phase 분리
- WHY: 깨끗한 컨텍스트 경계로 토큰 비대화 방지
- IMPACT: 동일 예산 내 2-3배 더 큰 프로젝트 가능

토큰 절약 전략:

1. Phase 분리: 단계 간 /clear
   - /do:1-plan 후 (45-50K 절약)
   - 컨텍스트 150K 초과 시
   - 50+ 메시지 후

2. 선택적 로딩: 필요한 파일만 로드

3. 컨텍스트 최적화: 20-30K 토큰 목표

4. 모델 선택: Sonnet (품질), Haiku (속도/비용)

모니터링: /context 명령, 예산 추적, /clear 제안

비용 절감: Haiku 70% 저렴 → 총 60-70% 절감

상세 참조: [Token Optimization Module](modules/token-optimization.md)

---

### 5. Progressive Disclosure - 콘텐츠 아키텍처

목적: 가치와 깊이의 균형을 맞춘 3단계 지식 전달

3가지 레벨:

빠른 참고 레벨:
- 시간 투자: 30초
- 내용: 핵심 원칙과 필수 개념
- 토큰 사용: 약 1,000 토큰
- WHY: 시간 제약 사용자에게 빠른 가치 전달
- IMPACT: 5% 시간에 80% 이해

구현 레벨:
- 시간 투자: 5분
- 내용: 워크플로우, 실용 예시, 통합 패턴
- 토큰 사용: 약 3,000 토큰
- WHY: 개념에서 실행으로 연결하는 실행 가이드
- IMPACT: 깊은 전문 지식 없이 즉시 생산적 작업 가능

고급 레벨:
- 시간 투자: 10분+
- 내용: 심층 기술, 엣지 케이스, 최적화 기법
- 토큰 사용: 약 5,000 토큰
- WHY: 복잡한 시나리오를 위한 마스터리 레벨 지식 제공
- IMPACT: 포괄적 커버리지로 에스컬레이션 70% 감소

SKILL.md 구조 (300-400 줄):
```
## 빠른 참고 (30초)
## 구현 가이드 (5분)
## 고급 패턴 (10분+)
## 연관 항목
```

모듈 아키텍처:
- SKILL.md: 진입점, 상호 참조
- modules/: 심층 다이브, 무제한
- examples.md: 동작 샘플
- reference.md: 외부 링크

파일 분할 (400줄 초과 시):
```
SKILL.md (300-400 줄)
  빠른 참고 (60-80)
  구현 가이드 (140-200)
  고급 패턴 (60-100)
  참조 (10-20)

오버플로우 → modules/[topic].md
```

상세 참조: [Progressive Disclosure Module](modules/progressive-disclosure.md)

---

### 6. Modular System - 파일 조직

목적: 무제한 콘텐츠를 위한 확장 가능한 파일 구조

표준 구조:
```
.claude/skills/skill-name/
  SKILL.md           # 핵심 (300-400 줄)
  modules/           # 확장 (무제한)
    patterns.md
  examples.md        # 동작 샘플
  reference.md       # 외부 링크
  scripts/           # 유틸리티 (선택)
  templates/         # 템플릿 (선택)
```

파일 원칙:

1. SKILL.md: 300-400줄, 점진적 공개, 상호 참조
2. modules/: 주제별, 무제한, 독립적
3. examples.md: 복사/붙여넣기 가능, 주석 포함
4. reference.md: API 문서, 리소스

상호 참조 구문:
```
상세: [Module](modules/patterns.md)
예시: [Examples](examples.md#auth)
외부: [Reference](reference.md#api)
```

발견 흐름: SKILL.md → 주제 → modules/[topic].md → 심층 다이브

상세 참조: [Modular System Module](modules/modular-system.md)

---

## 고급 구현 (10분+)

고급 패턴(교차 모듈 통합, 품질 검증, 오류 처리 등)은 상세 모듈 참조에서 확인

주요 고급 주제:
- 교차 모듈 통합: TRUST 5 + SPEC-First TDD 결합
- 토큰 최적화 위임: 컨텍스트 리셋과 병렬 실행
- 점진적 에이전트 워크플로우: 에스컬레이션 패턴
- 품질 검증: 실행 전/후 검증
- 오류 처리: 위임 실패 복구

상세 참조: [examples.md](examples.md) - 동작 코드 샘플

---

## 연관 항목

에이전트:
- agent-factory - 기초 원칙으로 에이전트 생성
- skill-factory - 모듈 아키텍처로 스킬 생성
- core-quality - 자동화된 TRUST 5 검증
- workflow-spec - EARS 형식 명세
- workflow-tdd - RED-GREEN-REFACTOR 실행
- workflow-docs - 점진적 공개 문서화

스킬:
- do-foundation-claude - CLAUDE.md와 기초 패턴
- do-foundation-context - 토큰 최적화
- do-workflow-project - 프로젝트 관리

도구:
- AskUserQuestion - 모든 사용자 상호작용 및 확인에 직접 사용

명령어:
- /do:1-plan - SPEC-First Phase 1
- /do:2-run - TDD Phase 2
- /do:3-sync - Documentation Phase 3
- /do:9-feedback - 지속적 개선
- /clear - 토큰 관리

기초 모듈 (확장 문서):
- [Agents Reference](modules/agents-reference.md) - 에이전트 카탈로그 및 7계층 구조
- [Commands Reference](modules/commands-reference.md) - 6개 핵심 명령어 워크플로우
- [Execution Rules](modules/execution-rules.md) - 보안, Git 전략, 준수

---

## 빠른 결정 매트릭스

시나리오별 주요 원칙:

새 에이전트: TRUST 5, Delegation, Token Opt, Modular
새 스킬: Progressive, Modular, TRUST 5, Token Opt
워크플로우: Delegation Patterns, SPEC-First, Token Opt
품질: TRUST 5 Framework, SPEC-First TDD
예산: Token Optimization, Progressive, Modular
문서화: Progressive, Modular, Token Optimization

모듈 심층 다이브:
- [TRUST 5 Framework](modules/trust-5-framework.md)
- [SPEC-First TDD](modules/spec-first-tdd.md)
- [Delegation Patterns](modules/delegation-patterns.md)
- [Token Optimization](modules/token-optimization.md)
- [Progressive Disclosure](modules/progressive-disclosure.md)
- [Modular System](modules/modular-system.md)
- [Agents Reference](modules/agents-reference.md)
- [Commands Reference](modules/commands-reference.md)
- [Execution Rules](modules/execution-rules.md)

전체 예시: [examples.md](examples.md)
외부 리소스: [reference.md](reference.md)

---

Version: 2.3.0
Last Updated: 2025-12-06
Status: Active (320-420줄 최적화)
