# Philosopher Framework Examples

실제 개발 시나리오에서 Philosopher Framework의 실용적 적용 사례.

## Example 1: 기술 선택

### 시나리오
팀이 React 애플리케이션의 상태 관리 솔루션을 선택해야 함.

### Phase 1: Assumption Audit

식별된 가정:
- Assumption 1: 애플리케이션의 복잡성이 크게 증가할 것
  - Confidence: Medium
  - Risk if wrong: 무거운 솔루션으로 over-engineering

- Assumption 2: 팀이 Redux 패턴에 익숙함
  - Confidence: Low (확인 필요)
  - Risk if wrong: 가파른 학습 곡선, 느린 배포

- Assumption 3: 이 애플리케이션에서 성능이 중요함
  - Confidence: High (요구사항 기반)
  - Risk if wrong: 조기 최적화

AskUserQuestion 적용:
- 팀에 Redux 경험자 2명, 미경험자 3명 확인
- 공유 상태가 필요한 컴포넌트가 50개 이상 될 것으로 확인
- 성능 요구사항이 초기 로드에 대한 것이지 상태 업데이트가 아님을 명확히 함

### Phase 2: First Principles Decomposition

Surface Problem: React 앱에 상태 관리 필요

First Why: 여러 컴포넌트가 동일한 데이터를 공유하고 업데이트해야 함
Second Why: 규모가 커지면 컴포넌트 prop drilling이 다루기 힘들어짐
Third Why: 애플리케이션이 깊이 중첩된 컴포넌트 계층 구조를 가짐
Root Cause: 중앙화되고 예측 가능한 상태 접근 패턴 필요

Constraint 분석:
- Hard Constraint: React 18과 호환되어야 함
- Soft Constraint: 팀의 익숙한 패턴 선호
- Degree of Freedom: 어떤 상태 관리 접근법이든 선택 가능

### Phase 3: Alternative Generation

Option A - Redux Toolkit:
- Pros: 업계 표준, 광범위한 에코시스템, DevTools
- Cons: 보일러플레이트, 학습 곡선, 잠재적 over-engineering

Option B - Zustand:
- Pros: 최소 보일러플레이트, 배우기 쉬움, 경량
- Cons: 작은 에코시스템, 대규모 앱에서 덜 구조화됨

Option C - React Context + useReducer:
- Pros: 내장, 의존성 없음, 익숙한 React 패턴
- Cons: 규모에서 성능 우려, DevTools 없음

Option D - Jotai:
- Pros: 원자적 접근, 최소 리렌더링, 간단한 API
- Cons: 다른 멘탈 모델, 작은 커뮤니티

### Phase 4: Trade-off Analysis

기준 및 가중치 (AskUserQuestion으로 확인):
- Learning Curve: 25% (팀 경험 혼합)
- Scalability: 25% (앱 성장 예상)
- Performance: 20% (중요하지만 핵심은 아님)
- Ecosystem: 15% (도구가 중요)
- Bundle Size: 15% (초기 로드가 우선)

점수:
- Redux Toolkit: Learning 5, Scale 9, Perf 7, Eco 9, Size 5 = 7.0
- Zustand: Learning 9, Scale 7, Perf 8, Eco 6, Size 9 = 7.8
- Context: Learning 8, Scale 5, Perf 5, Eco 5, Size 10 = 6.5
- Jotai: Learning 7, Scale 8, Perf 9, Eco 5, Size 8 = 7.4

### Phase 5: Cognitive Bias Check

Anchoring: 초기 본능은 친숙함으로 인해 Redux - 대안을 진정으로 고려했는지 확인
Confirmation: Zustand 한계와 Redux 장점을 적극적으로 검색
Sunk Cost: 의사결정에 영향을 주는 이전 투자 없음
Overconfidence: 확장성 예측의 불확실성 인정

### 권고

선택: Zustand
근거: 팀 구성을 고려할 때 학습 곡선과 확장성의 최적 균형
수용한 Trade-off: 더 빠른 팀 생산성을 위해 작은 에코시스템
완화: 복잡성이 Zustand 역량을 초과할 경우 Redux로 마이그레이션 경로 계획
Review Trigger: 상태 로직이 20개 store를 초과하거나 복잡한 미들웨어가 필요할 경우

---

## Example 2: 성능 최적화 의사결정

### 시나리오
API 엔드포인트가 느림 (P95: 2초, 목표: 200ms).

### Phase 1: Assumption Audit

식별된 가정:
- Assumption: 데이터베이스 쿼리가 병목
  - Confidence: Medium (직감 기반)
  - Risk if wrong: 잘못된 컴포넌트 최적화

AskUserQuestion 적용:
- 진행 전 프로파일링 데이터 요청
- 발견: DB 60%, 직렬화 30%, 네트워크 10%

### Phase 2: First Principles Decomposition

Surface Problem: API가 느림 (2초)

First Why: 응답 생성에 시간이 너무 오래 걸림
Second Why: 데이터베이스 쿼리가 너무 많은 데이터 반환
Third Why: 5개만 필요한데 쿼리가 모든 컬럼을 가져옴
Fourth Why: 최적화 없이 ORM 기본 select 사용
Root Cause: 이 사용 사례에 최적화되지 않은 ORM 사용 패턴

### Phase 3: Alternative Generation

Option A - 기존 쿼리 최적화:
- 컬럼 선택, 인덱스 추가
- Effort: Low, Risk: Low

Option B - 캐싱 레이어 추가:
- 자주 접근하는 데이터에 Redis 캐시
- Effort: Medium, Risk: Low

Option C - 데이터 비정규화:
- 읽기 최적화된 뷰 생성
- Effort: High, Risk: Medium

Option D - 아키텍처 변경:
- CQRS 패턴 구현
- Effort: Very High, Risk: High

### Phase 4: Trade-off Analysis

근본 원인이 쿼리 최적화이므로, Option A가 근본 원인을 직접 해결.

권고: Option A로 시작, 불충분하면 Option B 추가

### Phase 5: Cognitive Bias Check

Availability bias 회피: 최근 프로젝트 때문에 캐싱으로 바로 점프하지 않음
Overconfidence 회피: 원인을 가정하기 전에 프로파일링 데이터 요청

---

## Example 3: 리팩토링 범위 의사결정

### 시나리오
레거시 모듈 업데이트 필요. 팀이 전체 재작성 vs 점진적 개선을 논쟁.

### Phase 1: Assumption Audit

식별된 가정:
- Assumption: 전체 재작성에 3개월 소요
  - Confidence: Low (추정치가 종종 틀림)
  - Risk if wrong: 프로젝트 초과

- Assumption: 레거시 코드가 유지보수 불가
  - Confidence: Medium (메트릭 필요)
  - Risk if wrong: 불필요한 재작성

AskUserQuestion 적용:
- 코드 메트릭 요청 (복잡도, 테스트 커버리지, 변경 빈도)
- 결과: 40% 테스트 커버리지, 순환 복잡도 25 (높음), 월 3건 버그

### Phase 2: First Principles Decomposition

Surface Problem: 레거시 모듈이 문제

First Why: 버그가 자주 발생하고 수정이 어려움
Second Why: 코드가 복잡하고 테스트가 부족
Third Why: 원래 설계가 현재 요구사항을 예상하지 못함
Root Cause: 시간이 지나며 축적된 설계-요구사항 불일치

핵심 통찰: 모든 부분이 똑같이 문제인 것은 아님. 핵심 알고리즘은 견고, 인터페이스 레이어가 지저분.

### Phase 3: Alternative Generation

Option A - 전체 재작성:
- 처음부터 시작, 현대적 패턴
- Risk: Second-system effect, 기능 동등성 문제

Option B - Strangler 패턴:
- 점진적으로 부분 교체
- Risk: 장기화된 하이브리드 상태, 복잡성

Option C - 타겟 리팩토링:
- 가장 영향력 높은 영역만 수정
- Risk: 손대지 않은 영역에 기술 부채 잔존

Option D - 인터페이스 래퍼:
- 레거시 내부 위에 깨끗한 인터페이스
- Risk: 문제를 숨기지 해결하지 않음

### Phase 4: Trade-off Analysis

AskUserQuestion으로 결정: 팀에 3개월이 아닌 6주 가용.

시간 제약을 고려해 Option C (타겟 리팩토링) 선택.
집중 영역: 인터페이스 레이어 (버그의 60%, 코드의 30%), 먼저 테스트 추가.

### 권고

접근법: 테스트 우선 접근의 타겟 리팩토링
범위: 인터페이스 레이어만 (버그의 60%, 코드의 30%)
Trade-off: 핵심 알고리즘 기술 부채 잔존
완화: 핵심 알고리즘 주변 포괄적 테스트
Review Trigger: 2개월 내 버그율이 50% 감소하지 않으면

---

## 핵심 교훈

1. 중요한 의사결정 전에 항상 데이터로 가정 검증
2. 근본 원인 분석이 종종 더 단순한 솔루션을 드러냄
3. 시간과 리소스 제약이 정당하게 옵션을 좁힘
4. 향후 참조를 위해 Trade-off 문서화
5. 의사결정에 대한 명확한 review trigger 설정

---

Version: 1.0.0
Parent Skill: do-foundation-philosopher
