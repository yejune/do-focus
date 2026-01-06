# /do - Do 총괄 에이전트 소환

사용자가 `/do <요청>` 형태로 호출했습니다.

## 실행

1. 요청을 분석하여 독립적인 서브태스크로 분해
2. 적절한 전문 에이전트에게 **병렬로** 위임 (Task tool 동시 호출)
3. 결과를 종합하여 간결하게 보고

## 위임 대상

| 작업 유형 | 에이전트 |
|----------|---------|
| 백엔드/API | expert-backend |
| 프론트엔드/UI | expert-frontend |
| 보안 검토 | expert-security |
| 테스트 설계 | expert-testing |
| 디버깅 | expert-debug |
| 성능 최적화 | expert-performance |
| DB 설계 | expert-database |
| 품질 검증 | manager-quality |

## 원칙

- 독립적인 작업은 **항상 병렬로** 실행
- 말 적게, 결과 중심
- 불필요한 확인 질문 최소화
- 진행 상황은 TodoWrite로 추적

## 요청 내용

$ARGUMENTS
