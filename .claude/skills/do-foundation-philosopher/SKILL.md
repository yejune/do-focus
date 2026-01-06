---
name: do-foundation-philosopher
description: First Principles Analysis, Stanford Design Thinking, MIT Systems Engineering을 통합한 전략적 사고 프레임워크로 더 깊이 있는 문제 해결과 의사결정 지원
version: 1.0.0
modularized: true
updated: 2025-01-06
status: active
tags:
  - foundation
  - strategic-thinking
  - first-principles
  - trade-off-analysis
  - cognitive-bias
  - decision-making
allowed-tools: Read, Grep, Glob
---

# Do Foundation Philosopher

빠른 계산보다 깊은 분석을 촉진하는 전략적 사고 프레임워크. 세 가지 검증된 방법론을 통합해 체계적인 문제 해결 지원.

핵심 철학: 행동 전에 깊이 생각한다. 가정에 의문을 제기한다. 대안을 고려한다. Trade-off를 명시화한다. 인지 편향을 점검한다.

## Quick Reference (30초)

Philosopher Framework란?

복잡한 의사결정을 위한 구조화된 접근법:
- First Principles Analysis: 문제를 근본적 진실로 분해
- Stanford Design Thinking: 발산-수렴적 솔루션 생성
- MIT Systems Engineering: 체계적 위험 평가 및 검증

5단계 사고 프로세스:
1. Assumption Audit: 당연하게 여기는 것에 의문 제기
2. First Principles Decomposition: 근본 원인으로 분해
3. Alternative Generation: 여러 솔루션 옵션 생성
4. Trade-off Analysis: 옵션을 체계적으로 비교
5. Cognitive Bias Check: 사고 품질 검증

활성화 시점:
- 5개 이상 파일에 영향을 주는 아키텍처 의사결정
- 기술 선택 (라이브러리, 프레임워크, 데이터베이스)
- 성능 vs 유지보수성 Trade-off
- 리팩토링 범위 의사결정
- Breaking changes 고려
- 장기적 영향이 큰 모든 의사결정

Quick Access:
- 가정 질문 기법: [Assumption Matrix Module](modules/assumption-matrix.md)
- 근본 원인 분석: [First Principles Module](modules/first-principles.md)
- 옵션 비교: [Trade-off Analysis Module](modules/trade-off-analysis.md)
- 편향 방지: [Cognitive Bias Module](modules/cognitive-bias.md)

---

## Implementation Guide (5분)

### Phase 1: Assumption Audit

목적: 숨겨진 가정이 blind spot이 되기 전에 표면화.

5가지 핵심 질문:
- 증거 없이 참이라고 가정하는 것은 무엇인가?
- 이 가정이 틀리면 어떻게 되는가?
- 이것은 hard constraint인가, 단순한 선호인가?
- 이 가정을 뒷받침하는 증거는?
- 누가 이 가정을 검증해야 하는가?

Assumption 범주:
- Technical Assumptions: 기술 역량, 성능 특성, 호환성
- Business Assumptions: 사용자 행동, 시장 상황, 예산 가용성
- Team Assumptions: 스킬 수준, 가용성, 도메인 지식
- Timeline Assumptions: 배포 기대치, 의존성 일정

Assumption 문서화 형식:
- Assumption statement: 가정에 대한 명확한 설명
- Confidence level: 증거 기반으로 High, Medium, Low
- Evidence basis: 가정을 뒷받침하는 근거
- Risk if wrong: 가정이 틀릴 경우의 결과
- Validation method: 커밋 전 검증 방법

WHY: 검토되지 않은 가정은 프로젝트 실패와 재작업의 주요 원인.

### Phase 2: First Principles Decomposition

목적: 복잡성을 뚫고 근본 원인과 기본 요구사항 발견.

Five Whys 기법:
- Surface Problem: 사용자나 시스템이 관찰하는 것
- First Why: 즉각적인 원인 분석
- Second Why: 기저 원인 조사
- Third Why: 시스템적 요인 식별
- Fourth Why: 조직 또는 프로세스 요인
- Fifth Why (Root Cause): 해결해야 할 근본 이슈

Constraint 분석:
- Hard Constraints: 협상 불가 (보안, 규정, 물리적 한계, 예산)
- Soft Constraints: 협상 가능한 선호 (일정, 기능 범위, 도구)
- Self-Imposed Constraints: 요구사항으로 위장한 가정
- Degrees of Freedom: 창의적 솔루션이 가능한 영역

분해 질문:
- 이 요청 뒤에 있는 실제 목표는?
- 정말로 해결하려는 문제는?
- 제약이 없다면 솔루션은 어떤 모습일까?
- 최소 실행 가능 솔루션은?
- 목표 달성하면서 무엇을 제거할 수 있는가?

WHY: 대부분의 문제는 잘못된 추상화 수준에서 해결됨.

### Phase 3: Alternative Generation

목적: 차선의 솔루션으로 조기 수렴 방지.

생성 규칙:
- 최소 3개의 서로 다른 대안 필수
- 최소 1개의 비전통적 옵션 포함
- 항상 "아무것도 안 함"을 baseline으로 포함
- 단기 vs 장기 영향 고려
- 점진적 접근과 변혁적 접근 모두 탐색

Alternative 범주:
- Conservative: 낮은 위험, 점진적 개선, 익숙한 기술
- Balanced: 중간 위험, 상당한 개선, 약간의 혁신
- Aggressive: 높은 위험, 변혁적 변화, 최신 접근법
- Radical: 근본적 가정에 도전, 완전히 다른 접근

창의성 기법:
- Inversion: 문제를 악화시키는 것은? 그 반대를 수행
- Analogy: 다른 도메인은 유사 문제를 어떻게 해결하는가?
- Constraint Removal: 예산, 시간, 기술에 제한이 없다면?
- Simplification: 가장 단순한 가능 솔루션은?

WHY: 첫 번째 솔루션이 최선의 솔루션인 경우는 드묾.

### Phase 4: Trade-off Analysis

목적: 암묵적 Trade-off를 명시적이고 비교 가능하게 전환.

표준 평가 기준:
- Performance: 속도, 처리량, 지연시간, 리소스 사용
- Maintainability: 코드 명확성, 문서화, 팀 친숙도
- Implementation Cost: 개발 시간, 복잡성, 학습 곡선
- Risk Level: 기술 위험, 실패 확률, 롤백 난이도
- Scalability: 성장 용량, 유연성, 미래 대비
- Security: 취약점 표면, 규정 준수, 데이터 보호

가중 점수 방법:
- 프로젝트 우선순위 기반으로 기준에 가중치 할당 (총 100%)
- 각 옵션에 기준별 1-10점 부여
- 가중 복합 점수 계산
- 각 점수에 대한 근거 문서화
- 가중치 변경에 대한 점수 민감도 파악

Trade-off 문서화:
- What we gain: 선택한 접근법의 주요 이점
- What we sacrifice: 수용한 명시적 비용과 제한
- Why acceptable: 이러한 Trade-off 수용 근거
- Mitigation plan: 단점 해결 방법

WHY: 암묵적 Trade-off는 후회와 재검토로 이어짐.

### Phase 5: Cognitive Bias Check

목적: 일반적인 사고 오류를 점검해 권고 품질 보장.

주요 모니터링 편향:
- Anchoring Bias: 처음 접한 정보에 과도한 의존
- Confirmation Bias: 기존 믿음을 지지하는 증거만 탐색
- Sunk Cost Fallacy: 과거 투자 때문에 계속 진행
- Availability Heuristic: 최근 또는 기억하기 쉬운 사건에 과도한 비중
- Overconfidence Bias: 자신의 판단에 과도한 확신

편향 탐지 질문:
- 이 솔루션에 집착하는 이유가 먼저 생각했기 때문인가?
- 내 선호에 반하는 증거를 적극적으로 찾았는가?
- 이전 투자가 없는 상태에서 시작한다면 이것을 권고할까?
- 적용되지 않을 수 있는 최근 경험에 영향 받고 있지 않은가?
- 이 권고에 대한 내 생각을 바꿀 요인은?

완화 전략:
- Pre-mortem: 의사결정이 실패했다고 상상. 무엇이 잘못되었나?
- Devil's advocate: 자신의 권고에 반대 논증
- Outside view: Base rate가 성공에 대해 말하는 바는?
- Disagreement seeking: 도전할 사람과 상담
- Reversal test: 반대가 제안되면 무슨 말을 할 것인가?

WHY: 전문가도 시간 압박 하에서 인지 편향에 빠짐.

---

## Advanced Implementation

### Do 워크플로우 통합

SPEC Phase 통합:
- /do:1-plan 중 Assumption Audit 적용
- spec.md Problem Analysis 섹션에 가정 문서화
- plan.md에 고려한 대안 포함
- acceptance.md에 검증 기준 정의

TDD Phase 통합:
- First Principles로 핵심 테스트 시나리오 식별
- Edge case에 대한 테스트 대안 생성
- 테스트 커버리지 의사결정에 Trade-off Analysis 적용

Quality Phase 통합:
- 코드 리뷰 과정에 Cognitive Bias Check 포함
- 구현 후 가정이 여전히 유효한지 검증
- 최종 문서에 수용한 Trade-off 기록

### 시간 배분 가이드

복잡한 의사결정에 권장 노력 분배:
- Assumption Audit: 분석 시간의 15%
- First Principles Decomposition: 분석 시간의 25%
- Alternative Generation: 분석 시간의 20%
- Trade-off Analysis: 분석 시간의 25%
- Cognitive Bias Check: 분석 시간의 15%

총 분석 vs 구현 비율:
- 단순 의사결정 (1-2 파일): 분석 10%, 구현 90%
- 중간 의사결정 (3-10 파일): 분석 25%, 구현 75%
- 복잡 의사결정 (10+ 파일): 분석 40%, 구현 60%
- 아키텍처 의사결정: 분석 50%, 구현 50%

### 의사결정 문서화 템플릿

Strategic Decision Record:

Decision Title: 의사결정 내용에 대한 명확한 진술

Context: 이 의사결정이 필요한 이유

Assumptions Examined:
- Assumption 1: 신뢰도 및 검증 상태
- Assumption 2: 신뢰도 및 검증 상태

Root Cause Analysis:
- 식별된 Surface problem
- Five Whys를 통해 결정된 Root cause

Alternatives Considered:
- Option A: 장점, 단점, 점수
- Option B: 장점, 단점, 점수
- Option C: 장점, 단점, 점수

Trade-offs Accepted:
- 선택한 접근법으로 얻는 것
- 희생하는 것과 그 이유

Bias Check Completed:
- 편향 완화 단계 수행 확인

Final Decision: 주요 근거와 함께 선택된 옵션

Success Criteria: 의사결정 정확성 측정 방법

Review Trigger: 재검토를 유발하는 조건

---

## Works Well With

Agents:
- manager-strategy: SPEC 분석 및 계획의 주요 소비자
- expert-backend: 기술 선택 의사결정
- expert-frontend: 아키텍처 및 프레임워크 선택
- expert-database: 스키마 설계 Trade-off
- manager-quality: 코드 리뷰 편향 점검

Skills:
- do-foundation-core: TRUST 5 및 SPEC 워크플로우 통합
- do-workflow-spec: SPEC 형식의 Assumption 문서화

Commands:
- /do:1-plan: 명세 중 Philosopher Framework 적용
- /do:2-run: 구현 중 문서화된 Trade-off 참조

---

## Quick Decision Matrix

어떤 Phase를 사용할지:

Simple Bug Fix: Philosopher 건너뜀 (직접 구현)
Feature Addition: Phase 1, 3, 4 (가정, 대안, Trade-off)
Refactoring: Phase 1, 2, 4 (가정, 근본 원인, Trade-off)
Technology Selection: 전체 5 Phase (전체 분석 필수)
Architecture Change: 전체 5 Phase + 확장 문서화

---

Module Deep Dives:
- [Assumption Matrix](modules/assumption-matrix.md)
- [First Principles](modules/first-principles.md)
- [Trade-off Analysis](modules/trade-off-analysis.md)
- [Cognitive Bias](modules/cognitive-bias.md)

Examples: [examples.md](examples.md)
External Resources: [reference.md](reference.md)

---

Version: 1.0.0
Last Updated: 2025-01-06
Status: Active
Origin: Claude Code Philosopher Ignition 프레임워크에서 영감
