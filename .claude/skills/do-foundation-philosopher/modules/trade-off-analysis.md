# Trade-off Analysis Module

체계적 옵션 비교 및 의사결정 프레임워크에 대한 심층 가이드.

## 가중 점수 방법

### Step 1: 평가 기준 정의

기술 의사결정을 위한 표준 기준:

Performance 기준:
- 응답 시간 및 지연
- 처리량 및 용량
- 리소스 효율성 (CPU, 메모리, 네트워크)
- 부하 시 확장성

Quality 기준:
- 코드 유지보수성 및 가독성
- 테스트 커버리지 및 테스트 가능성
- 문서 완성도
- 에러 처리 견고성

Cost 기준:
- 구현 노력 (인일)
- 학습 곡선 및 교육 필요성
- 운영 비용 (인프라, 라이선싱)
- 도입되는 기술 부채

Risk 기준:
- 구현 복잡성
- 외부 요인 의존성
- 실패 모드 및 복구
- 롤백 난이도

Strategic 기준:
- 아키텍처 비전 정렬
- 미래 유연성 및 확장성
- 팀 스킬 개발
- 산업 표준 준수

### Step 2: 가중치 할당

가중치 할당 프로세스:

1. 의사결정과 관련된 모든 기준 나열
2. AskUserQuestion으로 사용자 우선순위 이해
3. 우선순위에 따라 기준에 100% 분배
4. 가중치 할당 근거 문서화

가중치 분배 예시:

Performance-Critical 프로젝트:
- Performance: 35%
- Quality: 20%
- Cost: 15%
- Risk: 20%
- Strategic: 10%

Maintainability-Focused 프로젝트:
- Performance: 15%
- Quality: 35%
- Cost: 20%
- Risk: 15%
- Strategic: 15%

Rapid Delivery 프로젝트:
- Performance: 15%
- Quality: 15%
- Cost: 40%
- Risk: 20%
- Strategic: 10%

### Step 3: 옵션 점수 부여

점수 가이드 (1-10 척도):

Score 9-10: 우수, 명확히 우월
Score 7-8: 좋음, 평균 이상
Score 5-6: 적절, 요구사항 충족
Score 3-4: 평균 이하, 우려 있음
Score 1-2: 부족, 상당한 문제

점수 부여 요구사항:
- 각 점수에 구체적 근거 제공
- 증거나 경험 참조
- 점수의 불확실성 고려
- 옵션 간 일관성 유지

### Step 4: 복합 점수 계산

계산 방법:
- 각 점수에 기준 가중치를 곱함
- 가중 점수 합산으로 총점 산출
- 옵션 간 총점 비교
- 가중치 변경에 대한 점수 민감도 분석

## Trade-off 문서화

### Trade-off Record 형식

각 중요한 Trade-off에 대해:

Trade-off ID: T-001
Decision Context: 이 Trade-off가 필요한 의사결정
What We Gain: 선택한 접근법의 이점
What We Sacrifice: 수용한 비용 또는 제한
Why Acceptable: 이 Trade-off 수용 근거
Mitigation Plan: 단점 영향 감소 조치
Review Trigger: 재검토를 유발하는 조건

### 일반적인 Trade-off 패턴

Speed vs Quality:
- 더 빠른 배포 vs 더 철저한 테스트
- 빠른 수정 vs 적절한 솔루션
- MVP vs 전체 기능 세트

Performance vs Maintainability:
- 최적화된 코드 vs 가독성 높은 코드
- 커스텀 솔루션 vs 표준 라이브러리
- 인라인 로직 vs 추상화 레이어

Flexibility vs Simplicity:
- 설정 가능 vs 하드코딩
- 제네릭 vs 구체적
- 플러그인 아키텍처 vs 모놀리식

Cost vs Capability:
- 빌드 vs 구매
- 오픈 소스 vs 상용
- 클라우드 vs 온프레미스

## AskUserQuestion 통합

Trade-off 분석 시:
- AskUserQuestion으로 기준 가중치가 우선순위와 일치하는지 확인
- AskUserQuestion으로 점수와 함께 옵션 제시
- AskUserQuestion으로 Trade-off 수용 가능성 검증
- AskUserQuestion으로 다른 가중치에 대한 민감도 탐색

AskUserQuestion Trade-off 확인 예시:

Question: Trade-off 분석 결과 Option B가 최고 점수. 약간의 성능 (6점)을 희생하고 더 나은 유지보수성 (9점)을 얻음. 이 Trade-off가 수용 가능한가?

Options:
- 인정된 Trade-off와 함께 Option B 수용
- 유지보수성보다 성능 우선
- 성능 영향에 대한 더 깊은 분석 요청
- 요소를 결합한 하이브리드 접근 탐색

---

Version: 1.0.0
Parent Skill: do-foundation-philosopher
