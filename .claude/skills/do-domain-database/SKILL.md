---
name: do-domain-database
description: PostgreSQL, MongoDB, Redis 및 현대적 애플리케이션을 위한 고급 데이터 패턴을 다루는 Database 전문가
version: 1.0.0
category: domain
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
tags:
  - database
  - postgresql
  - mongodb
  - redis
  - data-patterns
  - performance
updated: 2025-01-06
status: active
author: Do Team
---

# Database Domain Specialist

## 빠른 참조 (30초)

PostgreSQL, MongoDB, Redis 및 확장 가능한 현대적 애플리케이션을 위한 포괄적인 데이터베이스 패턴과 구현 가이드

핵심 역량:
- PostgreSQL: 고급 관계형 패턴, 최적화 및 확장
- MongoDB: Document 모델링, aggregation 및 NoSQL 성능 튜닝
- Redis: 인메모리 캐싱, 실시간 분석 및 분산 시스템
- Multi-Database: 하이브리드 아키텍처 및 데이터 통합 패턴
- Performance: 쿼리 최적화, 인덱싱 전략 및 확장
- Operations: 연결 관리, 마이그레이션 및 모니터링

사용 시점:
- 데이터베이스 스키마 및 데이터 모델 설계
- 캐싱 전략 및 성능 최적화 구현
- 확장 가능한 데이터 아키텍처 구축
- 멀티 데이터베이스 시스템 작업
- 데이터베이스 쿼리 및 성능 최적화

---

## 구현 가이드 (5분)

### 빠른 시작 워크플로우

Database Stack 초기화:

```
// PHP/JS 스타일 슈도 코드

// DatabaseManager 인스턴스 생성
dbManager = new DatabaseManager()

// PostgreSQL 설정 (관계형 데이터용)
postgresql = dbManager.setupPostgresql({
  connectionString: "postgresql://...",
  connectionPoolSize: 20,
  enableQueryLogging: true
})

// MongoDB 설정 (문서 저장용)
mongodb = dbManager.setupMongodb({
  connectionString: "mongodb://...",
  databaseName: "app_data",
  enableSharding: true
})

// Redis 설정 (캐싱 및 실시간 기능용)
redis = dbManager.setupRedis({
  connectionString: "redis://...",
  maxConnections: 50,
  enableClustering: true
})

// 통합 데이터베이스 인터페이스 사용
userData = dbManager.getUserWithProfile(userId)
analytics = dbManager.getUserAnalytics(userId, "30d")
```

### 핵심 구성 요소

1. PostgreSQL (`modules/postgresql.md`)
   - 고급 스키마 설계 및 제약 조건
   - 복잡한 쿼리 최적화 및 인덱싱
   - Window 함수 및 CTE
   - 파티셔닝 및 Materialized View
   - Connection pooling 및 성능 튜닝

2. MongoDB (`modules/mongodb.md`)
   - Document 모델링 및 스키마 설계
   - 분석용 Aggregation pipeline
   - 인덱싱 전략 및 성능
   - 샤딩 및 확장 패턴
   - 데이터 일관성 및 검증

3. Redis (`modules/redis.md`)
   - 다중 계층 캐싱 전략
   - 실시간 분석 및 카운팅
   - 분산 락 및 조정
   - Pub/sub 메시징 및 스트림
   - 고급 데이터 구조 (HyperLogLog, Geo)

---

## 고급 패턴 (10분 이상)

### Multi-Database 아키텍처

Polyglot Persistence 패턴:

```
// DataRouter 클래스
class DataRouter {
  constructor() {
    this.postgresql = new PostgreSQLConnection()
    this.mongodb = new MongoDBConnection()
    this.redis = new RedisConnection()
  }

  getUserProfile(userId) {
    // PostgreSQL에서 구조화된 사용자 데이터 조회
    user = this.postgresql.getUser(userId)

    // MongoDB에서 유연한 프로필 데이터 조회
    profile = this.mongodb.getUserProfile(userId)

    // Redis에서 실시간 상태 조회
    status = this.redis.getUserStatus(userId)

    return this.mergeUserData(user, profile, status)
  }

  updateUserData(userId, data) {
    // 각 데이터 타입을 적절한 데이터베이스로 라우팅

    // 구조화된 데이터가 있는 경우
    if (data.structuredData) {
      this.postgresql.updateUser(userId, data.structuredData)
    }

    // 프로필 데이터가 있는 경우
    if (data.profileData) {
      this.mongodb.updateUserProfile(userId, data.profileData)
    }

    // 실시간 데이터가 있는 경우
    if (data.realTimeData) {
      this.redis.setUserStatus(userId, data.realTimeData)
    }

    // 데이터베이스 간 캐시 무효화
    this.invalidateUserCache(userId)
  }
}
```

데이터 동기화:

```
class DataSyncManager {
  syncUserData(userId) {
    // PostgreSQL에서 MongoDB로 검색용 동기화
    pgUser = this.postgresql.getUser(userId)
    searchDocument = this.createSearchDocument(pgUser)
    this.mongodb.upsertUserSearch(userId, searchDocument)

    // Redis 캐시 업데이트
    cacheData = this.createCacheDocument(pgUser)
    this.redis.setUserCache(userId, cacheData, { ttl: 3600 })
  }
}
```

### 성능 최적화

쿼리 성능 분석:

```
// PostgreSQL 쿼리 최적화
function analyzeQueryPerformance(query) {
  explainResult = postgresql.execute("EXPLAIN (ANALYZE, BUFFERS) " + query)
  return QueryAnalyzer(explainResult).getOptimizationSuggestions()
}

// MongoDB aggregation 최적화
function optimizeAggregationPipeline(pipeline) {
  optimizer = new AggregationOptimizer()
  return optimizer.optimizePipeline(pipeline)
}

// Redis 성능 모니터링
function monitorRedisPerformance() {
  metrics = redis.info()
  return PerformanceAnalyzer(metrics).getRecommendations()
}
```

확장 전략:

```
// PostgreSQL Read Replicas
readReplicas = postgresql.setupReadReplicas([
  "postgresql://replica1...",
  "postgresql://replica2..."
])

// MongoDB 샤딩
mongodb.setupSharding({
  shardKey: "user_id",
  numShards: 4
})

// Redis 클러스터링
redis.setupCluster([
  "redis://node1:7000",
  "redis://node2:7000",
  "redis://node3:7000"
])
```

---

## 사용 예제

### 데이터베이스 작업

```
// PostgreSQL 고급 쿼리
users = postgresql.query(
  "SELECT * FROM users WHERE created_at > $1 ORDER BY activity_score DESC LIMIT 100",
  [thirtyDaysAgo]
)

// MongoDB 분석
analytics = mongodb.aggregate("events", [
  { $match: { timestamp: { $gte: startDate } } },
  { $group: { _id: "$type", count: { $sum: 1 } } },
  { $sort: { count: -1 } }
])

// Redis 캐싱 작업
async function getUserData(userId) {
  cacheKey = "user:" + userId
  data = await redis.get(cacheKey)

  if (!data) {
    data = await fetchFromDatabase(userId)
    await redis.setex(cacheKey, 3600, JSON.stringify(data))
  }

  return JSON.parse(data)
}
```

### Multi-Database 트랜잭션

```
async function createUserWithProfile(userData, profileData) {
  try {
    // 데이터베이스 간 트랜잭션 시작
    await transactionManager.begin()

    // PostgreSQL에 사용자 생성
    userId = await postgresql.insertUser(userData)

    // MongoDB에 프로필 생성
    await mongodb.insertUserProfile(userId, profileData)

    // Redis에 초기 캐시 설정
    await redis.setUserCache(userId, {
      id: userId,
      status: "active",
      createdAt: new Date().toISOString()
    })

    await transactionManager.commit()
    return userId

  } catch (error) {
    // 데이터베이스 간 자동 롤백
    await transactionManager.rollback()
    logger.error("User creation failed: " + error)
    throw error
  }
}
```

---

## 보완 Skill

- `do-domain-backend` - API 통합 및 비즈니스 로직
- `do-foundation-core` - 데이터베이스 마이그레이션 및 스키마 관리
- `do-workflow-project` - 데이터베이스 프로젝트 설정 및 구성
- `do-platform-supabase` - Supabase 데이터베이스 통합 패턴
- `do-platform-neon` - Neon 데이터베이스 통합 패턴
- `do-platform-firestore` - Firestore 데이터베이스 통합 패턴

기술 통합:
- ORM 및 ODM (SQLAlchemy, Mongoose, TypeORM)
- Connection pooling (PgBouncer, connection pools)
- 마이그레이션 도구 (Alembic, Flyway)
- 모니터링 (pg_stat_statements, MongoDB Atlas)
- 캐시 무효화 및 동기화

---

## 기술 스택

관계형 데이터베이스:
- PostgreSQL 14+ (주요)
- MySQL 8.0+ (대안)
- Connection pooling (PgBouncer, SQLAlchemy)

NoSQL 데이터베이스:
- MongoDB 6.0+ (주요)
- Document 모델링 및 검증
- Aggregation 프레임워크
- 샤딩 및 복제

인메모리 데이터베이스:
- Redis 7.0+ (주요)
- 고급 기능용 Redis Stack
- 클러스터링 및 고가용성
- 고급 데이터 구조

지원 도구:
- 마이그레이션 도구 (Alembic, Flyway)
- 모니터링 (Prometheus, Grafana)
- ORM/ODM (SQLAlchemy, Mongoose)
- 연결 관리

성능 기능:
- 쿼리 최적화 및 분석
- 인덱스 관리 및 전략
- 캐싱 계층 및 무효화
- 로드 밸런싱 및 페일오버

---

상세 구현 패턴 및 데이터베이스별 최적화는 `modules/` 디렉토리 참조
