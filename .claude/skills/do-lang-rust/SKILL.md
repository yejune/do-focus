---
name: do-lang-rust
description: Rust 1.92+ 개발 전문가. Axum, Tokio, SQLx 및 메모리 안전 시스템 프로그래밍 지원. 고성능, 메모리 안전 애플리케이션이나 WebAssembly 구축에 사용.
version: 1.1.0
updated: 2026-01-06
status: active
user-invocable: false
---

## 빠른 참조

Rust 1.92+ 고성능, 메모리 안전 애플리케이션 개발 전문가.

자동 트리거: `.rs`, `Cargo.toml`, async/await, Tokio, Axum, SQLx, serde, lifetimes, traits

핵심 사용 사례:
- 고성능 REST API 및 마이크로서비스
- 메모리 안전 동시성 시스템
- CLI 도구 및 시스템 유틸리티
- WebAssembly 애플리케이션
- 저지연 네트워킹 서비스

빠른 패턴:

Axum REST API:
```rust
let app = Router::new()
    .route("/api/users/:id", get(get_user))
    .with_state(app_state);
let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
axum::serve(listener, app).await?;
```

SQLx 비동기 핸들러:
```rust
async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> Result<Json<User>, AppError> {
    let user = sqlx::query_as!(User, "SELECT * FROM users WHERE id = $1", id)
        .fetch_optional(&state.db).await?
        .ok_or(AppError::NotFound)?;
    Ok(Json(user))
}
```

---

## Rust 1.92 기능

최신 Rust 기능:
- Rust 2024 에디션 사용 가능 (Rust 1.85에서 릴리스)
- 비동기 traits 안정화 (async-trait 크레이트 불필요)
- const generics로 컴파일 타임 배열 크기 지정
- let-else 패턴으로 조기 반환
- polonius로 개선된 borrow checker

비동기 Traits (안정화):
```rust
trait AsyncRepository {
    async fn get(&self, id: i64) -> Result<User, Error>;
    async fn create(&self, user: CreateUser) -> Result<User, Error>;
}

impl AsyncRepository for PostgresRepository {
    async fn get(&self, id: i64) -> Result<User, Error> {
        sqlx::query_as!(User, "SELECT * FROM users WHERE id = $1", id)
            .fetch_one(&self.pool).await
    }
}
```

Let-Else 패턴:
```rust
fn get_user(id: Option<i64>) -> Result<User, Error> {
    let Some(id) = id else { return Err(Error::MissingId); };
    let Ok(user) = repository.find(id) else { return Err(Error::NotFound); };
    Ok(user)
}
```

---

## 웹 프레임워크: Axum 0.8

설치:
```toml
[dependencies]
axum = "0.8"
tokio = { version = "1.48", features = ["full"] }
tower-http = { version = "0.6", features = ["cors", "trace"] }
```

완전한 API 설정:
```rust
use axum::{extract::{Path, State, Query}, routing::{get, post}, Router, Json};
use tower_http::cors::CorsLayer;

#[derive(Clone)]
struct AppState { db: PgPool }

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let pool = PgPoolOptions::new()
        .max_connections(25)
        .connect(&std::env::var("DATABASE_URL")?).await?;

    let app = Router::new()
        .route("/api/v1/users", get(list_users).post(create_user))
        .route("/api/v1/users/:id", get(get_user))
        .layer(CorsLayer::permissive())
        .with_state(AppState { db: pool });

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
    axum::serve(listener, app).await?;
    Ok(())
}
```

핸들러 패턴:
```rust
async fn list_users(
    State(state): State<AppState>,
    Query(params): Query<ListParams>,
) -> Result<Json<Vec<User>>, AppError> {
    let users = sqlx::query_as!(User,
        "SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
        params.limit.unwrap_or(10), params.offset.unwrap_or(0))
        .fetch_all(&state.db).await?;
    Ok(Json(users))
}
```

미들웨어:
```rust
use axum::middleware::{self, Next};
use axum::http::Request;

async fn auth_middleware<B>(
    State(state): State<AppState>,
    request: Request<B>,
    next: Next<B>,
) -> Result<Response, AppError> {
    let token = request.headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "))
        .ok_or(AppError::Unauthorized)?;
    let claims = verify_token(token, &state.jwt_secret)?;
    Ok(next.run(request).await)
}

let protected = Router::new()
    .route("/users/me", get(get_current_user))
    .layer(middleware::from_fn_with_state(state.clone(), auth_middleware));
```

---

## 비동기 런타임: Tokio 1.48

태스크 관리:
```rust
use tokio::task::{JoinHandle, JoinSet};

let handle: JoinHandle<i32> = tokio::spawn(async { 42 });

let mut set = JoinSet::new();
for i in 0..10 {
    set.spawn(async move { process(i).await });
}
while let Some(result) = set.join_next().await {
    println!("완료: {:?}", result);
}
```

채널:
```rust
use tokio::sync::{mpsc, oneshot, broadcast};

// MPSC 채널
let (tx, mut rx) = mpsc::channel::<Message>(100);
tokio::spawn(async move {
    while let Some(msg) = rx.recv().await { process(msg).await; }
});
tx.send(Message::new()).await?;

// 단일 응답용 oneshot
let (tx, rx) = oneshot::channel::<Response>();
tokio::spawn(async move { let _ = tx.send(compute().await); });
let response = rx.await?;
```

Select와 타임아웃:
```rust
async fn timeout_operation() -> Result<Data, Error> {
    tokio::select! {
        result = fetch_data() => result,
        _ = tokio::time::sleep(Duration::from_secs(5)) => Err(Error::Timeout),
    }
}
```

---

## 데이터베이스: SQLx 0.8

타입 안전 쿼리:
```rust
#[derive(Debug, sqlx::FromRow)]
struct User { id: i64, name: String, email: String }

async fn user_operations(pool: &PgPool) -> Result<(), sqlx::Error> {
    let user = sqlx::query_as!(User,
        "SELECT id, name, email FROM users WHERE id = $1", 1i64)
        .fetch_one(pool).await?;

    let mut tx = pool.begin().await?;
    sqlx::query!("INSERT INTO users (name, email) VALUES ($1, $2)", "John", "john@example.com")
        .execute(&mut *tx).await?;
    tx.commit().await?;
    Ok(())
}
```

연결 풀:
```rust
let pool = PgPoolOptions::new()
    .max_connections(25)
    .min_connections(5)
    .acquire_timeout(Duration::from_secs(5))
    .connect(&std::env::var("DATABASE_URL")?).await?;
```

---

## 직렬화: Serde 1.0

```rust
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct User {
    id: i64,
    #[serde(rename = "userName")]
    name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    profile_url: Option<String>,
    #[serde(default)]
    is_active: bool,
}
```

---

## 에러 처리

thiserror:
```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("database error: {0}")]
    Database(#[from] sqlx::Error),
    #[error("not found: {0}")]
    NotFound(String),
    #[error("unauthorized")]
    Unauthorized,
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match &self {
            AppError::NotFound(_) => (StatusCode::NOT_FOUND, self.to_string()),
            AppError::Unauthorized => (StatusCode::UNAUTHORIZED, self.to_string()),
            AppError::Database(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Internal error".into()),
        };
        (status, Json(json!({"error": message}))).into_response()
    }
}
```

---

## CLI 개발: clap

```rust
use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "myapp", version, about)]
struct Cli {
    #[arg(short, long, global = true)]
    config: Option<PathBuf>,
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Serve { #[arg(short, long, default_value = "3000")] port: u16 },
    Migrate,
}

fn main() {
    let cli = Cli::parse();
    match cli.command {
        Commands::Serve { port } => serve(port),
        Commands::Migrate => migrate(),
    }
}
```

---

## 테스트 패턴

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_create_user() {
        let pool = setup_test_db().await;
        let result = create_user(&pool, "John", "john@example.com").await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().name, "John");
    }
}
```

---

## 고급 패턴

### 동시성 제어

Rate-Limited 작업:
```rust
use tokio::sync::Semaphore;

async fn rate_limited(items: Vec<String>, max: usize) -> Vec<String> {
    let sem = std::sync::Arc::new(Semaphore::new(max));
    let handles: Vec<_> = items.into_iter().map(|item| {
        let sem = sem.clone();
        tokio::spawn(async move {
            let _permit = sem.acquire().await.unwrap();
            process_item(item).await
        })
    }).collect();
    futures::future::join_all(handles).await.into_iter().filter_map(|r| r.ok()).collect()
}
```

### 성능 최적화

Release 빌드:
```toml
[profile.release]
lto = true
codegen-units = 1
panic = "abort"
strip = true
```

### 배포

최소 컨테이너 (5-15MB):
```dockerfile
FROM rust:1.92-alpine AS builder
WORKDIR /app
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main(){}" > src/main.rs && cargo build --release
COPY src ./src
RUN touch src/main.rs && cargo build --release

FROM alpine:latest
COPY --from=builder /app/target/release/app /app
EXPOSE 3000
CMD ["/app"]
```

---

## Context7 통합

라이브러리 문서 접근:
- `/rust-lang/rust` - Rust 언어 및 stdlib
- `/tokio-rs/tokio` - Tokio 비동기 런타임
- `/tokio-rs/axum` - Axum 웹 프레임워크
- `/launchbadge/sqlx` - SQLx 비동기 SQL
- `/serde-rs/serde` - 직렬화 프레임워크
- `/dtolnay/thiserror` - 에러 derive
- `/clap-rs/clap` - CLI 파서

---

## 함께 사용할 스킬

- `do-lang-go` - Go 시스템 프로그래밍 패턴
- `do-domain-backend` - REST API 아키텍처 및 마이크로서비스 패턴
- `do-foundation-quality` - Rust 애플리케이션 보안 강화
- `do-workflow-tdd` - 테스트 주도 개발 워크플로우

---

## 문제 해결

일반 이슈:
- Cargo 에러: `cargo clean && cargo build`
- 버전 확인: `rustc --version && cargo --version`
- 의존성 문제: `cargo update && cargo tree`
- 컴파일 타임 SQL 확인: `cargo sqlx prepare`

성능 특성:
- 시작 시간: 50-100ms
- 메모리 사용량: 5-20MB 기본
- 처리량: 100k-200k req/s
- 지연: p99 5ms 미만
- 컨테이너 크기: 5-15MB (alpine)

---

## 추가 리소스

완전한 언어 참조 및 Context7 라이브러리 매핑은 [reference.md](reference.md) 참조.

프로덕션 준비 코드 예제는 [examples.md](examples.md) 참조.

---

Last Updated: 2026-01-06
Version: 1.1.0
