---
name: do-lang-go
description: Go 1.23+ 개발 전문가로 Fiber, Gin, GORM 및 동시성 프로그래밍 패턴을 다룬다. 고성능 마이크로서비스, CLI 도구 또는 클라우드 네이티브 애플리케이션 구축 시 사용.
version: 1.0.0
category: language
tags: [go, golang, fiber, gin, concurrency, microservices]
context7-libraries: [/gofiber/fiber, /gin-gonic/gin, /go-gorm/gorm]
related-skills: [do-lang-rust, do-domain-backend]
updated: 2026-01-06
status: active
---

## 빠른 참조

Go 1.23+ 개발 전문가로 고성능 백엔드 시스템과 CLI 애플리케이션 구축.

자동 트리거: `.go`, `go.mod`, `go.sum`, goroutine, channel, Fiber, Gin, GORM, Echo, Chi

핵심 사용 사례:
- 고성능 REST API 및 마이크로서비스
- 동시성 및 병렬 처리 시스템
- CLI 도구 및 시스템 유틸리티
- 클라우드 네이티브 컨테이너 서비스

빠른 패턴:

```go
// Fiber API
app := fiber.New()
app.Get("/api/users/:id", func(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"id": c.Params("id")})
})
app.Listen(":3000")

// Gin API
r := gin.Default()
r.GET("/api/users/:id", func(c *gin.Context) {
    c.JSON(200, gin.H{"id": c.Param("id")})
})
r.Run(":3000")

// Errgroup
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return processUsers(ctx) })
g.Go(func() error { return processOrders(ctx) })
if err := g.Wait(); err != nil { log.Fatal(err) }
```

---

## Go 1.23 언어 기능

새로운 기능:
- 정수 범위 반복: `for i := range 10 { fmt.Println(i) }`
- Profile-Guided Optimization (PGO) 2.0
- 개선된 제네릭과 타입 추론

```go
// 제네릭
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice { result[i] = fn(v) }
    return result
}

// 오류 처리
var ErrNotFound = errors.New("not found")
if err != nil { return fmt.Errorf("fetch user %d: %w", id, err) }
if errors.Is(err, ErrNotFound) { /* 처리 */ }
```

---

## 웹 프레임워크: Fiber v3

```go
app := fiber.New(fiber.Config{ErrorHandler: customErrorHandler, Prefork: true})
app.Use(recover.New(), logger.New(), cors.New())
app.Use(limiter.New(limiter.Config{Max: 100, Expiration: time.Minute}))

api := app.Group("/api/v1")
api.Get("/users", listUsers)
api.Get("/users/:id", getUser)
api.Post("/users", createUser)
app.Listen(":3000")

func getUser(c fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil { return fiber.NewError(fiber.StatusBadRequest, "Invalid ID") }
    return c.JSON(fiber.Map{"id": id})
}
```

---

## 웹 프레임워크: Gin

```go
r := gin.Default()
r.Use(cors.Default())

api := r.Group("/api/v1")
api.GET("/users", listUsers)
api.POST("/users", createUser)
r.Run(":3000")

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2"`
    Email string `json:"email" binding:"required,email"`
}

func createUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, gin.H{"id": 1, "name": req.Name})
}
```

---

## ORM: GORM 1.25

```go
type User struct {
    gorm.Model
    Name  string `gorm:"uniqueIndex;not null"`
    Email string `gorm:"uniqueIndex;not null"`
    Posts []Post `gorm:"foreignKey:AuthorID"`
}

db.First(&user, 1)
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(10)
}).First(&user, 1)

db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil { return err }
    return tx.Create(&profile).Error
})
```

---

## PostgreSQL: pgx

```go
config, _ := pgxpool.ParseConfig(connString)
config.MaxConns = 25
config.MinConns = 5
pool, _ := pgxpool.NewWithConfig(ctx, config)

err := pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).
    Scan(&user.ID, &user.Name)

rows, _ := pool.Query(ctx, "SELECT id, name FROM users LIMIT $1", 10)
defer rows.Close()
for rows.Next() { rows.Scan(&u.ID, &u.Name) }

tx, _ := pool.Begin(ctx)
defer tx.Rollback(ctx)
tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", name)
tx.Commit(ctx)
```

---

## 동시성 패턴

```go
// Errgroup
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { users, err = fetchUsers(ctx); return err })
g.Go(func() error { orders, err = fetchOrders(ctx); return err })
if err := g.Wait(); err != nil { return nil, err }

// Semaphore 속도 제한
var sem = semaphore.NewWeighted(10)
func processWithLimit(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)
    for _, item := range items {
        item := item
        g.Go(func() error {
            if err := sem.Acquire(ctx, 1); err != nil { return err }
            defer sem.Release(1)
            return processItem(ctx, item)
        })
    }
    return g.Wait()
}

// Worker Pool
func workerPool(ctx context.Context, jobs <-chan Job, n int) <-chan Result {
    results := make(chan Result, 100)
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case <-ctx.Done(): return
                default: results <- processJob(job)
                }
            }
        }()
    }
    go func() { wg.Wait(); close(results) }()
    return results
}

// Context 타임아웃
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
```

---

## 테스트 패턴

```go
// 테이블 기반 테스트
tests := []struct {
    name    string
    input   CreateUserInput
    wantErr bool
}{
    {"valid", CreateUserInput{Name: "John"}, false},
    {"empty", CreateUserInput{Name: ""}, true},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        _, err := svc.Create(tt.input)
        if tt.wantErr { require.Error(t, err) }
    })
}

// HTTP 테스트
app := fiber.New()
app.Get("/users/:id", getUser)
req := httptest.NewRequest("GET", "/users/1", nil)
resp, _ := app.Test(req)
assert.Equal(t, 200, resp.StatusCode)
```

---

## CLI: Cobra + Viper

```go
var rootCmd = &cobra.Command{Use: "myapp", Short: "Description"}

func init() {
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
    viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
    viper.SetEnvPrefix("MYAPP")
    viper.AutomaticEnv()
}
```

---

## 성능 및 배포

PGO 빌드:
```bash
GODEBUG=pgo=1 ./myapp -cpuprofile=default.pgo
go build -pgo=default.pgo -o myapp
```

객체 풀링:
```go
var bufferPool = sync.Pool{New: func() interface{} { return make([]byte, 4096) }}
buf := bufferPool.Get().([]byte)
defer bufferPool.Put(buf)
```

컨테이너 배포 (10-20MB):
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]
```

Graceful Shutdown:
```go
go func() { app.Listen(":3000") }()
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
app.Shutdown()
```

---

## Context7 라이브러리

```
/golang/go        - Go 언어 및 표준 라이브러리
/gofiber/fiber    - Fiber 웹 프레임워크
/gin-gonic/gin    - Gin 웹 프레임워크
/labstack/echo    - Echo 웹 프레임워크
/go-chi/chi       - Chi 라우터
/go-gorm/gorm     - GORM ORM
/sqlc-dev/sqlc    - 타입 안전 SQL
/jackc/pgx        - PostgreSQL 드라이버
/spf13/cobra      - CLI 프레임워크
/spf13/viper      - 설정 관리
/stretchr/testify - 테스트 툴킷
```

---

## 관련 스킬

- `do-domain-backend` - REST API 아키텍처와 마이크로서비스
- `do-lang-rust` - 시스템 프로그래밍 동반자
- `do-workflow-testing` - 테스트 주도 개발

---

## 문제 해결

일반 이슈:
- 모듈 오류: `go mod tidy && go mod verify`
- 버전 확인: `go version && go env GOVERSION`
- 빌드 이슈: `go clean -cache && go build -v`

성능 진단:
- CPU 프로파일링: `go test -cpuprofile=cpu.prof -bench=.`
- 메모리 프로파일링: `go test -memprofile=mem.prof -bench=.`
- 경쟁 조건 탐지: `go test -race ./...`

---

Last Updated: 2026-01-06
Version: 1.0.0
