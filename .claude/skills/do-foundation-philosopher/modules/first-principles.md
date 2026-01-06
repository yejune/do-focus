# First Principles Module

근본 원인 분석 및 기본적 분해 기법에 대한 심층 가이드.

## Five Whys 기법

### 적용 프로세스

관찰된 문제로 시작하여 반복적으로 "왜"를 물음:

Level 1 - Surface Problem:
- 사용자나 시스템이 경험하는 것은?
- 보이는 증상은?
- 언제, 어디서 발생하는가?

Level 2 - First Why:
- 즉각적인 원인은?
- 증상을 촉발하는 것은?
- 직접 관련된 컴포넌트는?

Level 3 - Second Why:
- 그 즉각적인 원인이 존재하는 이유는?
- 가능하게 하는 조건은?
- 기여하는 상류 요인은?

Level 4 - Third Why:
- 그 조건이 존재하는 이유는?
- 작용하는 시스템적 요인은?
- 여기로 이끈 프로세스나 설계 의사결정은?

Level 5 - Fourth/Fifth Why (Root Cause):
- 그 의사결정이 내려진 이유는?
- 이끈 근본적 제약이나 가정은?
- 재발을 방지하려면 무엇이 변해야 하는가?

### 일반적인 함정

너무 일찍 멈춤:
- 증상을 원인으로 받아들임
- 시스템적 요인을 파고들지 않음
- 시스템 대신 개인을 탓함

너무 멀리 감:
- 철학적이거나 변경 불가능한 원인에 도달
- 실행 가능한 구체성을 잃음
- 프로젝트 통제를 벗어난 범위로 확장

분기 혼란:
- 각 레벨에서 여러 유효한 답변
- 가장 영향력 높은 분기를 먼저 탐색
- 나중을 위해 대안 분기 문서화

## Constraint 분석

### Hard Constraints

정의: 변경할 수 없는 협상 불가 요구사항.

예시:
- 보안 규정 요구사항 (SOC2, GDPR, HIPAA)
- 물리적 제한 (네트워크 지연, 스토리지 용량)
- 법적 요구사항 (데이터 보존, 접근성)
- 예산 상한 (승인된 자금 한도)
- 호환성 요구사항 (기존 시스템 통합)

처리: 이러한 제약 내에서 작동하는 솔루션 설계.

### Soft Constraints

정의: Trade-off가 허용 가능하면 조정할 수 있는 선호.

예시:
- 일정 선호 (원하지만 협상 가능한 날짜)
- 기능 범위 (있으면 좋은 것 vs 필수)
- 기술 선호 (익숙한 것 vs 최적인 것)
- 품질 수준 (충분히 좋은 것 vs 완벽)
- 팀 선호 (방법 vs 결과)

처리: AskUserQuestion으로 정말 협상 가능한 것을 명확히 함.

### Self-Imposed Constraints

정의: 요구사항으로 위장한 가정.

일반적인 self-imposed constraints:
- "Technology X를 사용해야 함" (대안이 존재할 때)
- "이 방식으로 해야 함" (다른 접근법이 작동할 때)
- "...할 시간이 없음" (우선순위 조정이 도움이 될 때)
- "항상 그렇게 해왔음" (레거시 프로세스 관성)

처리: 실제 제약인지 습관인지 의문 제기.

## 분해 패턴

### Goal Decomposition

고수준 목표를 실행 가능한 하위 목표로 분해:

고수준 목표: 애플리케이션 성능 개선

Sub-Goal 1: 페이지 로드 시간 단축
- Metric: Time to first contentful paint
- Target: 1.5초 미만

Sub-Goal 2: API 응답 시간 개선
- Metric: P95 응답 시간
- Target: 200ms 미만

Sub-Goal 3: 리소스 소비 감소
- Metric: 메모리 및 CPU 사용
- Target: 20% 감소

### Solution Space Mapping

모든 가능한 솔루션 방향 식별:

문제: 데이터베이스 쿼리가 느림

Solution Space:
- 쿼리 최적화 (인덱스, 쿼리 재작성)
- 캐싱 레이어 (Redis, 인메모리)
- 데이터베이스 스케일링 (읽기 복제본, 샤딩)
- 아키텍처 변경 (CQRS, 이벤트 소싱)
- 데이터 모델 재설계 (비정규화, 집계)
- 하드웨어 업그레이드 (더 빠른 디스크, 더 많은 메모리)

AskUserQuestion으로 제약을 고려해 어떤 방향이 실행 가능한지 탐색.

## AskUserQuestion 통합

문제를 분해할 때:
- AskUserQuestion으로 문제 이해 확인
- AskUserQuestion으로 특정 제약이 존재하는 이유 탐색
- AskUserQuestion으로 근본 원인 식별 확인
- AskUserQuestion으로 분해 완성도 검증

---

Version: 1.0.0
Parent Skill: do-foundation-philosopher
