---
name: do-platform-railway
description: Railway 컨테이너 배포 전문가 - Docker, 다중 서비스 아키텍처, 지속적 볼륨, 자동 확장을 다룸. 컨테이너화된 풀스택 애플리케이션 배포 시 사용
version: 1.1.0
category: platform
tags: [railway, docker, containers, multi-service, auto-scaling]
context7-libraries: [railway]
related-skills: [do-platform-vercel, do-domain-backend]
updated: 2025-12-30
status: active
allowed-tools: Read, Write, Bash, Grep, Glob
user-invocable: false
---

# do-platform-railway: 컨테이너 배포 전문가

## 빠른 참조

Railway 플랫폼 핵심: Docker 및 Nixpacks 빌드, 다중 서비스 아키텍처, 지속적 볼륨, 프라이빗 네트워킹, 자동 확장 기능을 갖춘 컨테이너 중심 배포 플랫폼

### Railway 최적 사용 사례

컨테이너 워크로드:
- 커스텀 런타임이 필요한 풀스택 컨테이너 애플리케이션
- 서비스 간 통신이 필요한 다중 서비스 아키텍처
- 지속적 연결이 필요한 백엔드 서비스 (WebSocket, gRPC)
- 커스텀 런타임 요구사항 (Python, Go, Rust, Elixir)
- 관리형 PostgreSQL, MySQL, Redis를 사용하는 데이터베이스 기반 애플리케이션

인프라 요구사항:
- 상태 저장 워크로드를 위한 지속적 볼륨 스토리지
- 보안 서비스 메시를 위한 프라이빗 네트워킹
- 글로벌 가용성을 위한 다중 지역 배포
- CPU, 메모리, 요청 메트릭 기반 자동 확장

### 빌드 전략 선택

Docker 빌드: 커스텀 시스템 의존성, 다단계 빌드, 특정 베이스 이미지가 필요한 경우
Nixpacks 빌드: 표준 런타임(Node.js, Python, Go), 제로 설정, 빠른 빌드가 필요한 경우

참고: Nixpacks는 유지보수 모드로 더 이상 권장되지 않음. 신규 서비스는 Railpack이 기본값. 기존 서비스는 Nixpacks로 계속 동작. 마이그레이션하려면 railway.toml에서 builder = "RAILPACK" 설정

### 주요 CLI 명령어

```bash
railway login && railway init && railway link
railway up                    # 현재 디렉토리 배포
railway up --detach          # 로그 없이 배포
railway up --service api     # 특정 서비스 배포
railway variables --set KEY=value
railway logs --service api
```

---

## 구현 가이드

### 1단계: Docker 배포 패턴

Node.js 다단계 Dockerfile:
```dockerfile
FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs && adduser --system --uid 1001 appuser
COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
USER appuser
EXPOSE 3000
CMD ["node", "dist/main.js"]
```

Python 프로덕션 Dockerfile:
```dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
RUN pip install --no-cache-dir poetry
COPY pyproject.toml poetry.lock ./
RUN poetry config virtualenvs.create false && poetry install --no-dev

FROM python:3.12-slim AS runner
WORKDIR /app
ENV PYTHONUNBUFFERED=1
COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY . .
RUN useradd --create-home appuser
USER appuser
EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

Go 프로덕션 Dockerfile:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates && adduser -D appuser
WORKDIR /app
COPY --from=builder /app/main .
USER appuser
EXPOSE 8080
CMD ["./main"]
```

### 2단계: Nixpacks 설정

Node.js nixpacks.toml:
```toml
[phases.setup]
nixPkgs = ["nodejs-20_x", "pnpm"]

[phases.install]
cmds = ["pnpm install --frozen-lockfile"]

[phases.build]
cmds = ["pnpm build"]

[start]
cmd = "pnpm start"
```

Python nixpacks.toml:
```toml
[phases.setup]
nixPkgs = ["python312", "poetry"]
aptPkgs = ["libpq-dev"]

[phases.install]
cmds = ["poetry install --no-dev"]

[start]
cmd = "poetry run gunicorn app:application --bind 0.0.0.0:$PORT"
```

### 3단계: Railway 설정

Nixpacks 기반 railway.toml:
```toml
[build]
builder = "NIXPACKS"
buildCommand = "npm run build"
watchPatterns = ["src/**", "package.json"]

[deploy]
startCommand = "npm start"
healthcheckPath = "/health"
healthcheckTimeout = 100
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 5
numReplicas = 2

[deploy.resources]
memory = "512Mi"
cpu = "0.5"
```

Docker 기반 railway.toml:
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "./Dockerfile"

[deploy]
healthcheckPath = "/api/health"
healthcheckTimeout = 60
restartPolicyType = "ALWAYS"
numReplicas = 3

[deploy.resources]
memory = "1Gi"
cpu = "1"
```

### 4단계: 다중 서비스 아키텍처

모노레포 다중 서비스 (각 서비스별 railway.toml):

모노레포의 각 서비스는 해당 디렉토리에 자체 railway.toml 필요

API 서비스 (apps/api/railway.toml):
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "./Dockerfile"

[deploy]
startCommand = "node dist/main.js"
healthcheckPath = "/health"
numReplicas = 3
```

워커 서비스 (apps/worker/railway.toml):
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "./Dockerfile"

[deploy]
startCommand = "node dist/worker.js"
numReplicas = 2
```

스케줄러 서비스 (apps/scheduler/railway.toml):
```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "./Dockerfile"

[deploy]
startCommand = "node dist/scheduler.js"
cronSchedule = "*/5 * * * *"
```

서비스 변수 참조:
서비스 설정에서 Railway의 변수 참조 구문을 사용하여 서비스 간 통신 구현. 형식: ${{ServiceName.VARIABLE_NAME}}

프라이빗 네트워킹:
```typescript
const getInternalUrl = (service: string, port = 3000): string => {
  const domain = process.env[`${service.toUpperCase()}_RAILWAY_PRIVATE_DOMAIN`]
  return domain ? `http://${domain}:${port}` : `http://localhost:${port}`
}

app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: process.env.RAILWAY_SERVICE_NAME,
    replica: process.env.RAILWAY_REPLICA_ID
  })
})
```

### 5단계: 지속적 볼륨

볼륨 설정:
```toml
[deploy]
startCommand = "npm start"

[[volumes]]
mountPath = "/app/data"
name = "app-data"
size = "10Gi"

[[volumes]]
mountPath = "/app/uploads"
name = "user-uploads"
size = "50Gi"
```

지속적 스토리지 패턴:
```typescript
import { join } from 'path'
import { existsSync, mkdirSync, writeFileSync, readFileSync } from 'fs'

const VOLUME_PATH = process.env.RAILWAY_VOLUME_MOUNT_PATH || '/app/data'

class PersistentStorage {
  constructor() {
    if (!existsSync(VOLUME_PATH)) mkdirSync(VOLUME_PATH, { recursive: true })
  }
  write(file: string, data: Buffer | string) { writeFileSync(join(VOLUME_PATH, file), data) }
  read(file: string): Buffer { return readFileSync(join(VOLUME_PATH, file)) }
}
```

### 6단계: 자동 확장

리소스 기반 확장:
```toml
[deploy.scaling]
minReplicas = 2
maxReplicas = 10
targetCPUUtilization = 70
targetMemoryUtilization = 80
```

요청 기반 확장:
```toml
[deploy.scaling]
minReplicas = 1
maxReplicas = 20
targetRequestsPerSecond = 100
scaleDownDelaySeconds = 300
```

애플리케이션 메트릭:
```typescript
import { register, Counter, Histogram } from 'prom-client'

const httpRequests = new Counter({
  name: 'http_requests_total',
  help: 'Total HTTP requests',
  labelNames: ['method', 'status']
})

app.get('/metrics', async (req, res) => {
  res.set('Content-Type', register.contentType)
  res.end(await register.metrics())
})
```

---

## 다중 지역 배포

사용 가능 지역: us-west1 (Oregon), us-east4 (Virginia), europe-west4 (Netherlands), asia-southeast1 (Singapore)

```bash
railway up --region us-west1
railway up --region europe-west4
```

지역 설정:
```toml
[[deploy.regions]]
name = "us-west1"
replicas = 3

[[deploy.regions]]
name = "europe-west4"
replicas = 2
```

데이터베이스 읽기 복제본:
```typescript
const primaryPool = new Pool({ connectionString: process.env.DATABASE_URL })
const replicaPool = new Pool({ connectionString: process.env.DATABASE_REPLICA_URL })

async function query(sql: string, params?: any[]) {
  const isRead = sql.trim().toLowerCase().startsWith('select')
  return (isRead ? replicaPool : primaryPool).query(sql, params)
}
```

---

## CI/CD 통합

GitHub Actions:
```yaml
name: Railway Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm i -g @railway/cli
      - run: railway up --detach
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

롤백 명령어:
```bash
railway deployments list
railway rollback <deployment-id>
railway rollback --previous
```

---

## 연관 스킬

- `do-platform-vercel` - 프론트엔드 애플리케이션용 엣지 배포
- `do-domain-backend` - 백엔드 서비스 아키텍처 패턴
- `do-lang-python` - Python FastAPI 배포 설정
- `do-lang-typescript` - TypeScript Node.js 배포 패턴

---

상태: Production Ready | 버전: 1.1.0 | 업데이트: 2025-12-30
