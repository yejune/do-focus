---
name: do-platform-vercel
description: Vercel edge deployment 전문가 - Edge Functions, Next.js 최적화, preview 배포, ISR 커버. Next.js 또는 edge-first 애플리케이션 배포 시 사용.
version: 1.0.0
category: platform
tags: vercel, edge, nextjs, isr, preview, cdn
context7-libraries: /vercel/next.js, /vercel/vercel
related-skills: do-platform-railway, do-lang-typescript
updated: 2025-01-06
status: active
allowed-tools: Read, Write, Bash, Grep, Glob
---

# Vercel Edge 배포 전문가

## 역할

Vercel 플랫폼에서 Edge-first 배포를 전문으로 담당. 글로벌 CDN, Next.js 최적화 런타임, 개발자 중심 preview 워크플로우를 활용한 고성능 웹 애플리케이션 배포 지원.

---

## 핵심 기능

### Edge Functions

- 30개 이상의 글로벌 edge 위치에서 저지연 연산
- 50ms 이하 cold start로 최적의 사용자 경험 제공
- 지역 기반 라우팅 및 개인화
- Edge Middleware로 요청/응답 변환

### Next.js 최적화 런타임

- Next.js에 대한 1급 지원과 자동 최적화
- Server Components 및 App Router 통합
- Streaming SSR로 TTFB 개선
- next/image 내장 이미지 최적화

### Preview 배포

- PR 기반 자동 preview URL 생성
- 브랜치별 환경 변수 관리
- PR 리뷰용 댓글 통합
- 즉시 롤백 기능

### ISR (Incremental Static Regeneration)

- On-demand revalidation으로 동적 콘텐츠 관리
- Stale-while-revalidate 캐싱 전략
- 태그 기반 캐시 무효화
- 사용자 영향 없는 백그라운드 재생성

---

## 선택 기준

Vercel 선택이 적합한 경우:
- Next.js가 주 프레임워크인 프로젝트
- Edge 성능이 핵심 요구사항
- 팀 협업을 위한 preview 배포 필요
- Web Vitals 모니터링이 우선순위

---

## 구현 가이드

### Phase 1: Edge Functions 아키텍처

Edge Runtime 설정:

```typescript
// app/api/edge-handler/route.ts
export const runtime = 'edge'
export const preferredRegion = ['iad1', 'sfo1', 'fra1']

export async function GET(request: Request) {
  const { geo, ip } = request

  return Response.json({
    country: geo?.country ?? 'Unknown',
    city: geo?.city ?? 'Unknown',
    region: geo?.region ?? 'Unknown',
    ip: ip ?? 'Unknown',
    timestamp: new Date().toISOString()
  })
}
```

지역 기반 콘텐츠 전달:

```typescript
// app/api/localized/route.ts
export const runtime = 'edge'

const CONTENT_BY_REGION: Record<string, { currency: string; locale: string }> = {
  US: { currency: 'USD', locale: 'en-US' },
  DE: { currency: 'EUR', locale: 'de-DE' },
  JP: { currency: 'JPY', locale: 'ja-JP' },
  KR: { currency: 'KRW', locale: 'ko-KR' }
}

export async function GET(request: Request) {
  const country = request.geo?.country ?? 'US'
  const config = CONTENT_BY_REGION[country] ?? CONTENT_BY_REGION.US

  return Response.json(config, {
    headers: {
      'Cache-Control': 'public, s-maxage=3600',
      'CDN-Cache-Control': 'public, max-age=86400'
    }
  })
}
```

Edge에서 A/B 테스트:

```typescript
// middleware.ts
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const existingBucket = request.cookies.get('ab-bucket')?.value

  if (!existingBucket) {
    const bucket = Math.random() < 0.5 ? 'control' : 'variant'
    const response = NextResponse.next()
    response.cookies.set('ab-bucket', bucket, {
      maxAge: 60 * 60 * 24 * 30, // 30일
      httpOnly: true,
      sameSite: 'lax'
    })
    return response
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)']
}
```

### Phase 2: Next.js 최적화 패턴

ISR 구현:

```typescript
// app/products/[id]/page.tsx
import { notFound } from 'next/navigation'

// 60초마다 재검증
export const revalidate = 60

// 상위 제품에 대한 정적 파라미터 생성
export async function generateStaticParams() {
  const products = await fetch('https://api.example.com/products/top').then(r => r.json())
  return products.map((p: { id: string }) => ({ id: p.id }))
}

// SEO를 위한 동적 메타데이터
export async function generateMetadata({ params }: { params: { id: string } }) {
  const product = await fetchProduct(params.id)
  if (!product) return {}

  return {
    title: product.name,
    description: product.description,
    openGraph: {
      images: [product.imageUrl]
    }
  }
}

export default async function ProductPage({ params }: { params: { id: string } }) {
  const product = await fetchProduct(params.id)
  if (!product) notFound()

  return <ProductDetail product={product} />
}

async function fetchProduct(id: string) {
  const res = await fetch(`https://api.example.com/products/${id}`, {
    next: { tags: [`product-${id}`] }
  })
  if (!res.ok) return null
  return res.json()
}
```

On-Demand 재검증:

```typescript
// app/api/revalidate/route.ts
import { revalidateTag, revalidatePath } from 'next/cache'
import { NextRequest } from 'next/server'

export async function POST(request: NextRequest) {
  const { tag, path, secret } = await request.json()

  // 웹훅 시크릿 검증
  if (secret !== process.env.REVALIDATION_SECRET) {
    return Response.json({ error: 'Invalid secret' }, { status: 401 })
  }

  try {
    if (tag) {
      revalidateTag(tag)
      return Response.json({ revalidated: true, tag })
    }

    if (path) {
      revalidatePath(path)
      return Response.json({ revalidated: true, path })
    }

    return Response.json({ error: 'Missing tag or path' }, { status: 400 })
  } catch (error) {
    return Response.json({ error: 'Revalidation failed' }, { status: 500 })
  }
}
```

Suspense를 활용한 Streaming:

```typescript
// app/dashboard/page.tsx
import { Suspense } from 'react'

export default function DashboardPage() {
  return (
    <div className="dashboard">
      <h1>Dashboard</h1>

      <Suspense fallback={<MetricsSkeleton />}>
        <Metrics />
      </Suspense>

      <Suspense fallback={<ChartSkeleton />}>
        <AnalyticsChart />
      </Suspense>

      <Suspense fallback={<TableSkeleton />}>
        <RecentOrders />
      </Suspense>
    </div>
  )
}

async function Metrics() {
  const data = await fetch('https://api.example.com/metrics', {
    next: { revalidate: 30 }
  }).then(r => r.json())

  return <MetricsDisplay data={data} />
}
```

### Phase 3: Vercel 설정

vercel.json 설정:

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "framework": "nextjs",
  "buildCommand": "pnpm build",
  "installCommand": "pnpm install",
  "outputDirectory": ".next",
  "regions": ["iad1", "sfo1", "fra1", "hnd1"],
  "functions": {
    "app/api/**/*.ts": {
      "memory": 1024,
      "maxDuration": 30
    },
    "app/api/heavy/**/*.ts": {
      "memory": 3008,
      "maxDuration": 60
    }
  },
  "crons": [
    {
      "path": "/api/cron/cleanup",
      "schedule": "0 0 * * *"
    },
    {
      "path": "/api/cron/sync",
      "schedule": "*/15 * * * *"
    }
  ],
  "headers": [
    {
      "source": "/api/(.*)",
      "headers": [
        { "key": "Access-Control-Allow-Origin", "value": "*" },
        { "key": "Cache-Control", "value": "s-maxage=60, stale-while-revalidate" }
      ]
    },
    {
      "source": "/(.*)",
      "headers": [
        { "key": "X-Content-Type-Options", "value": "nosniff" },
        { "key": "X-Frame-Options", "value": "DENY" },
        { "key": "X-XSS-Protection", "value": "1; mode=block" }
      ]
    }
  ],
  "rewrites": [
    { "source": "/blog/:slug", "destination": "/posts/:slug" },
    { "source": "/api/v1/:path*", "destination": "/api/:path*" }
  ],
  "redirects": [
    { "source": "/old-page", "destination": "/new-page", "permanent": true }
  ]
}
```

환경 변수 관리:

```bash
# 프로덕션 환경
vercel env add DATABASE_URL production
vercel env add API_SECRET production
vercel env add NEXT_PUBLIC_API_URL production

# Preview 환경 (PR 배포)
vercel env add DATABASE_URL preview
vercel env add API_SECRET preview

# 개발 환경
vercel env add DATABASE_URL development

# 로컬 개발을 위한 환경 변수 가져오기
vercel env pull .env.local
```

### Phase 4: Preview 배포

GitHub 통합 설정:

```yaml
# .github/workflows/vercel-preview.yml
name: Vercel Preview Deployment
on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  deploy-preview:
    runs-on: ubuntu-latest
    environment:
      name: Preview
      url: ${{ steps.deploy.outputs.url }}
    steps:
      - uses: actions/checkout@v4

      - name: Install Vercel CLI
        run: npm i -g vercel@latest

      - name: Pull Vercel Environment
        run: vercel pull --yes --environment=preview --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project
        run: vercel build --token=${{ secrets.VERCEL_TOKEN }}

      - name: Deploy to Vercel
        id: deploy
        run: |
          url=$(vercel deploy --prebuilt --token=${{ secrets.VERCEL_TOKEN }})
          echo "url=$url" >> $GITHUB_OUTPUT

      - name: Comment PR with Preview URL
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: `Preview deployed: ${{ steps.deploy.outputs.url }}`
            })
```

프로덕션 배포:

```yaml
# .github/workflows/vercel-production.yml
name: Vercel Production Deployment
on:
  push:
    branches: [main]

jobs:
  deploy-production:
    runs-on: ubuntu-latest
    environment:
      name: Production
      url: https://your-domain.com
    steps:
      - uses: actions/checkout@v4

      - name: Install Vercel CLI
        run: npm i -g vercel@latest

      - name: Pull Vercel Environment
        run: vercel pull --yes --environment=production --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project
        run: vercel build --prod --token=${{ secrets.VERCEL_TOKEN }}

      - name: Deploy to Vercel
        run: vercel deploy --prebuilt --prod --token=${{ secrets.VERCEL_TOKEN }}
```

### Phase 5: Web Vitals 모니터링

Analytics 통합:

```typescript
// app/layout.tsx
import { Analytics } from '@vercel/analytics/react'
import { SpeedInsights } from '@vercel/speed-insights/next'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ko">
      <body>
        {children}
        <Analytics />
        <SpeedInsights />
      </body>
    </html>
  )
}
```

커스텀 Web Vitals 리포팅:

```typescript
// app/components/WebVitals.tsx
'use client'

import { useReportWebVitals } from 'next/web-vitals'

export function WebVitals() {
  useReportWebVitals((metric) => {
    const body = JSON.stringify({
      name: metric.name,
      value: metric.value,
      rating: metric.rating,
      delta: metric.delta,
      id: metric.id,
      navigationType: metric.navigationType
    })

    // 커스텀 분석 엔드포인트로 전송
    if (navigator.sendBeacon) {
      navigator.sendBeacon('/api/vitals', body)
    } else {
      fetch('/api/vitals', { body, method: 'POST', keepalive: true })
    }
  })

  return null
}
```

---

## 고급 패턴

### Turborepo를 활용한 Monorepo

```json
{
  "buildCommand": "cd ../.. && pnpm turbo build --filter=web",
  "installCommand": "cd ../.. && pnpm install",
  "framework": "nextjs"
}
```

### Blue-Green 배포

새 버전 배포 후 preview URL에서 스모크 테스트 실행, 이후 Vercel SDK의 aliases.assign() 메서드로 프로덕션 alias 전환. 무중단 릴리스 달성.

### Context7 통합

Edge Function 패턴에는 `/vercel/vercel`, App Router, ISR, Streaming 패턴에는 `/vercel/next.js` 사용. 적절한 토큰 할당 고려.

---

## 연관 스킬

- `do-platform-railway` - 컨테이너 기반 배포 대안
- `do-lang-typescript` - Next.js용 TypeScript 패턴
- `do-domain-frontend` - React 및 Next.js 컴포넌트 패턴
- `do-foundation-quality` - 배포 검증 및 테스트

---

Status: Production Ready | Version: 1.0.0 | Updated: 2025-01-06
