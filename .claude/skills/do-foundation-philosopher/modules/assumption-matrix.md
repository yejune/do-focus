# Assumption Matrix Module

체계적인 가정 식별 및 검증에 대한 심층 가이드.

## Assumption 범주

### Technical Assumptions

검토할 일반적인 기술 가정:
- Technology X가 우리의 규모 요구사항을 처리할 수 있음
- API가 하위 호환성을 유지할 것임
- 선택한 접근법으로 성능이 허용 가능할 것임
- 기존 시스템과의 통합이 단순할 것임
- 개발팀이 필요한 전문 지식을 보유함

검증 방법:
- 개념 증명 구현
- 부하 테스트 및 벤치마킹
- API 문서 및 변경 이력 검토
- 경험 많은 엔지니어와 아키텍처 리뷰
- 스킬 평가 및 교육 갭 분석

### Business Assumptions

검토할 일반적인 비즈니스 가정:
- 사용자가 새 기능을 채택할 것임
- 예산이 계속 가용할 것임
- 일정이 현실적임
- 요구사항이 안정적이고 잘 이해됨
- 이해관계자 우선순위가 정렬됨

검증 방법:
- 사용자 리서치 및 피드백 세션
- 재무팀과 예산 확인
- 과거 프로젝트 데이터 비교
- 요구사항 리뷰 미팅
- 이해관계자 정렬 워크숍

### Team Assumptions

검토할 일반적인 팀 가정:
- 팀이 이 작업에 여력이 있음
- 스킬이 사내에 가용함
- 커뮤니케이션 채널이 효과적임
- 다른 팀에 대한 의존성이 관리 가능함
- 지식 이전이 필요 없음

검증 방법:
- 용량 계획 리뷰
- 스킬 매트릭스 평가
- 팀 회고 피드백
- 의존성 매핑 연습
- 문서 감사

## Assumption 문서화 템플릿

각 중요한 가정에 대해:

Assumption ID: A-001
Statement: 가정에 대한 명확한 설명
Category: Technical, Business, 또는 Team
Confidence: High, Medium, 또는 Low
Evidence: 이 가정을 뒷받침하는 것
Risk if Wrong: 가정이 실패할 경우 프로젝트에 미치는 영향
Validation Plan: 검증 방법과 시기
Owner: 검증 책임자
Status: Unvalidated, Validated, 또는 Invalidated

## AskUserQuestion 통합

가정을 표면화할 때 AskUserQuestion 사용:
- 검증해야 할 핵심 가정 확인
- 사용자의 도메인 지식에서 증거 수집
- 가정 검증 노력 우선순위 결정
- 검증되지 않은 가정에 대한 허용 가능한 위험 수준 결정

---

Version: 1.0.0
Parent Skill: do-foundation-philosopher
