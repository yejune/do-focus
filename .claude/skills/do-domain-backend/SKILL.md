---
name: do-domain-backend
description: API 설계, 데이터베이스 통합, 마이크로서비스 아키텍처, 최신 백엔드 패턴을 다루는 백엔드 개발 전문가
version: 1.0.0
category: domain
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
tags:
  - backend
  - api
  - database
  - microservices
  - architecture
updated: 2026-01-06
status: active
author: Do Team
---

# 백엔드 개발 전문가

## 요약

API 설계, 데이터베이스 통합, 마이크로서비스 아키텍처를 다루는 백엔드 개발 전문가 스킬

핵심 기능:
- API 설계: REST, GraphQL, gRPC with OpenAPI 3.1
- 데이터베이스: PostgreSQL, MongoDB, Redis, 캐싱 전략
- 마이크로서비스: 분산 패턴, 이벤트 기반 아키텍처
- 보안: 인증, 권한부여, OWASP 준수
- 성능: 캐싱, 최적화, 모니터링, 확장

사용 시기:
- 백엔드 API 개발 및 아키텍처 설계
- 데이터베이스 설계 및 최적화
- 마이크로서비스 구현
- 성능 최적화 및 확장성 개선
- 백엔드 시스템 보안 통합

---

## API 설계 패턴

### RESTful API 구조

기본 앱 설정:
- FastAPI 인스턴스 생성 (title, version, lifespan 설정)
- CORS 미들웨어 추가
- 라이프사이클에서 DB/캐시 연결 초기화

Pydantic 모델 정의:
- 요청 모델: 입력 검증 (UserCreate, OrderRequest)
- 응답 모델: 출력 형식 지정 (UserResponse, TokenResponse)
- EmailStr, validator 등 내장 검증 활용

라우터 설계:
- GET /users - 목록 조회 (페이지네이션, 필터링)
- GET /users/{id} - 단건 조회
- POST /users - 생성 (요청 본문 검증)
- PUT /users/{id} - 전체 수정
- PATCH /users/{id} - 부분 수정
- DELETE /users/{id} - 삭제

### GraphQL 구조

Strawberry 타입 정의:
- @strawberry.type 데코레이터로 타입 클래스 정의
- Query 타입: 조회 작업 (@strawberry.field)
- Mutation 타입: 변경 작업
- Subscription 타입: 실시간 업데이트

리졸버 패턴:
- async 함수로 비동기 데이터 로딩
- DataLoader로 N+1 문제 해결
- 인증 정보는 context로 전달

---

## 데이터베이스 통합

### PostgreSQL (SQLAlchemy)

비동기 엔진 설정:
- create_async_engine 사용
- pool_size: 20 (기본 연결 수)
- max_overflow: 30 (추가 허용 연결)
- pool_pre_ping: True (연결 상태 확인)
- pool_recycle: 3600 (1시간마다 재연결)

세션 관리:
- AsyncSession으로 비동기 트랜잭션 처리
- 컨텍스트 매니저로 자동 커밋/롤백
- Depends()로 요청 단위 세션 주입

### MongoDB (Motor)

클라이언트 설정:
- AsyncIOMotorClient로 비동기 연결
- maxPoolSize: 50 (최대 연결)
- minPoolSize: 10 (최소 유지 연결)
- waitQueueTimeoutMS: 5000 (대기 타임아웃)

인덱스 최적화:
- 자주 조회하는 필드에 단일 인덱스
- 복합 조회에 복합 인덱스
- unique=True로 중복 방지

### Redis

연결 풀 설정:
- ConnectionPool.from_url 사용
- max_connections: 50
- decode_responses: True (자동 문자열 변환)

캐싱 패턴:
- setex로 TTL 설정 (기본 3600초)
- get/set으로 단순 캐시
- 캐시 키 네이밍: "entity:id" 형식

---

## 마이크로서비스 아키텍처

### 이벤트 기반 아키텍처

EventBus 구현:
- RabbitMQ 연결 (aio-pika)
- ExchangeType.TOPIC으로 라우팅
- 메시지 영속성: delivery_mode=2

이벤트 발행:
- 이벤트 타입 문자열 (order.created, user.registered)
- JSON 페이로드 (데이터 + 타임스탬프)
- routing_key로 구독자 필터링

이벤트 구독:
- 큐 선언 후 exchange에 바인딩
- 핸들러 함수로 메시지 처리
- message.process()로 확인(ack) 관리

### 서비스 디스커버리

Consul 등록:
- 서비스명, ID, 포트 등록
- 헬스체크 엔드포인트 설정 (/health)
- 주기적 상태 확인 (interval: 10s)

서비스 탐색:
- 서비스명으로 건강한 인스턴스 조회
- 로드 밸런싱: 랜덤 또는 라운드로빈 선택
- 장애 인스턴스 자동 제외

---

## 인증 및 보안

### JWT 인증

토큰 생성:
- jwt.encode()로 토큰 발급
- 페이로드: sub (사용자 ID), exp (만료시간)
- 알고리즘: HS256 또는 RS256
- 만료 시간: 액세스 30분, 리프레시 7일

토큰 검증:
- jwt.decode()로 페이로드 추출
- ExpiredSignatureError: 토큰 만료
- InvalidTokenError: 잘못된 토큰

비밀번호 처리:
- CryptContext(schemes=["bcrypt"])
- hash(): 해싱
- verify(): 검증

### OAuth2 미들웨어

의존성 주입:
- OAuth2PasswordBearer(tokenUrl="auth/login")
- Depends()로 보호 라우트에 적용
- HTTPException 401: 인증 실패 응답

---

## 캐싱 전략

### 캐시 어사이드 패턴

조회 흐름:
- 캐시에서 먼저 조회
- 캐시 히트: 즉시 반환
- 캐시 미스: DB 조회 후 캐시 저장

캐시 무효화:
- 데이터 변경 시 해당 키 삭제
- 목록 캐시: 와일드카드 삭제 (users:list:*)
- TTL 기반 자동 만료

### 캐시 키 설계

단일 엔티티: "{entity}:{id}"
- 예시: user:123, order:456

목록 캐시: "{entity}:list:{params}"
- 예시: users:list:page=1&limit=10

의존 캐시: "{parent}:{id}:{child}"
- 예시: user:123:orders

---

## 성능 최적화

### 연결 풀 튜닝

풀 크기 설정:
- pool_size: 동시 요청의 80%
- max_overflow: 피크 시 추가 여유 (50%)
- pool_recycle: 연결 수명 (1시간)

모니터링 포인트:
- 대기 시간 (pool wait time)
- 활성 연결 수
- 오버플로우 빈도

### 쿼리 최적화

N+1 문제 해결:
- selectinload(): 관련 엔티티 일괄 로드
- joinedload(): JOIN으로 한번에 조회
- 복잡한 관계: subqueryload()

페이지네이션:
- OFFSET/LIMIT: 단순하지만 느림
- 커서 기반: 대용량에 적합
- keyset: ID 또는 timestamp 기준

### 느린 쿼리 추적

이벤트 리스너 설정:
- before_cursor_execute: 시작 시간 기록
- after_cursor_execute: 실행 시간 계산
- 임계값 초과 시 로깅 (0.1초 이상)

---

## 서킷 브레이커

### 상태 관리

상태 종류:
- CLOSED: 정상 동작
- OPEN: 호출 차단 (장애 상태)
- HALF_OPEN: 복구 테스트 중

상태 전이 조건:
- CLOSED에서 OPEN: 연속 실패 임계값 도달
- OPEN에서 HALF_OPEN: 복구 대기 시간 경과
- HALF_OPEN에서 CLOSED: 테스트 성공

### 구현 요소

설정값:
- failure_threshold: 5 (연속 실패 횟수)
- recovery_timeout: 30초 (복구 대기)
- success_threshold: 3 (복구 확인 횟수)

재시도 전략:
- tenacity.retry 데코레이터 활용
- stop_after_attempt(3): 최대 3회
- wait_exponential: 지수 백오프 (1, 2, 4초)

---

## 문제 해결

### 연결 풀 고갈

증상: 요청 타임아웃, "too many connections" 오류

해결 방법:
- pool_size와 max_overflow 증가
- 연결 누수 점검 (컨텍스트 매니저 사용 확인)
- pool_pre_ping으로 끊어진 연결 감지

### 느린 쿼리

증상: 응답 지연, DB CPU 급등

해결 방법:
- EXPLAIN ANALYZE로 쿼리 분석
- 적절한 인덱스 추가
- 결과 캐싱 적용
- 읽기 전용 레플리카 활용

### 비동기 컨텍스트 메모리 누수

증상: 점진적 메모리 증가, OOM 발생

해결 방법:
- async 컨텍스트 매니저 올바르게 사용
- 라이프스팬 핸들러에서 정리 작업 수행
- 태스크 취소 및 정리 모니터링

### CORS 오류

증상: 브라우저에서 크로스 오리진 요청 차단

해결 방법:
- allow_origins에 클라이언트 도메인 추가
- 쿠키 인증 시 allow_credentials=True
- OPTIONS preflight 요청 처리 확인

---

## 기술 스택

언어: Python 3.13+, Node.js 22+, Go 1.23

프레임워크: FastAPI, Django, Express.js, Gin

데이터베이스: PostgreSQL 16+, MongoDB 7+, Redis 7+

메시징: RabbitMQ, Apache Kafka, Redis Pub/Sub

배포: Docker, Kubernetes

모니터링: Prometheus, Grafana, OpenTelemetry

통합 패턴:
- RESTful APIs with OpenAPI 3.1
- GraphQL with Apollo Federation
- gRPC for high-performance services
- Event-driven architecture with CQRS
- Circuit breakers and resilience patterns

---

## 연관 스킬

- do-domain-frontend - Full-stack 개발 통합
- do-domain-database - 고급 데이터베이스 패턴
- do-foundation-core - 핵심 아키텍처 원칙
- do-platform-* - 특정 플랫폼 통합

---

Status: Production Ready
Last Updated: 2026-01-06
Maintained by: Do Backend Team
