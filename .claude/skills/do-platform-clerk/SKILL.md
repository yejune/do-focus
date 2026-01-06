---
name: do-platform-clerk
description: Clerk 현대 인증 플랫폼 전문가. WebAuthn, 패스키, 비밀번호 없는 인증, UI 컴포넌트 제공. 현대적 인증 UX 구현 시 사용.
version: 2.0.0
category: platform
tags: clerk, webauthn, passkeys, passwordless, authentication
context7-libraries: /clerk/clerk-docs
related-skills: do-platform-auth0, do-lang-typescript
updated: 2025-12-30
status: active
allowed-tools: Read, Write, Bash, Grep, Glob
---

# Clerk 현대 인증 플랫폼 전문가

WebAuthn, 패스키, 비밀번호 없는 인증 플로우, 사전 구축된 UI 컴포넌트, 멀티테넌트 조직 지원을 제공하는 현대 인증 플랫폼.

SDK 버전 (2025년 12월 기준):
- @clerk/nextjs: 6.x (Core 2, Next.js 13.0.4+, React 18+ 필수)
- @clerk/clerk-react: 5.x (Core 2, React 18+ 필수)
- @clerk/express: 1.x
- Node.js: 18.17.0+ 필수

---

## 빠른 참조

환경 변수:

```bash
# .env.local
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in
NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up
```

ClerkProvider 설정 (app/layout.tsx):

```tsx
import { ClerkProvider } from '@clerk/nextjs'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body>{children}</body>
      </html>
    </ClerkProvider>
  )
}
```

기본 미들웨어 (middleware.ts):

```typescript
import { clerkMiddleware } from '@clerk/nextjs/server'

export default clerkMiddleware()

export const config = {
  matcher: [
    '/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)',
    '/(api|trpc)(.*)',
  ],
}
```

Context7 접근:
- 라이브러리: /clerk/clerk-docs
- 해결: resolve-library-id에 "clerk" 사용 후 get-library-docs 호출

---

## 구현 가이드

### 인증 컴포넌트를 포함한 ClerkProvider

로그인/로그아웃 컨트롤을 포함한 전체 레이아웃:

```tsx
// app/layout.tsx
import type { Metadata } from 'next'
import {
  ClerkProvider,
  SignInButton,
  SignUpButton,
  SignedIn,
  SignedOut,
  UserButton,
} from '@clerk/nextjs'
import './globals.css'

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <ClerkProvider>
      <html lang="en">
        <body>
          <header className="flex justify-end items-center p-4 gap-4 h-16">
            <SignedOut>
              <SignInButton />
              <SignUpButton />
            </SignedOut>
            <SignedIn>
              <UserButton />
            </SignedIn>
          </header>
          {children}
        </body>
      </html>
    </ClerkProvider>
  )
}
```

### 미들웨어로 라우트 보호

createRouteMatcher를 사용한 라우트 보호:

```typescript
// middleware.ts
import { clerkMiddleware, createRouteMatcher } from '@clerk/nextjs/server'

const isProtectedRoute = createRouteMatcher([
  '/dashboard(.*)',
  '/forum(.*)',
  '/api/private(.*)',
])

export default clerkMiddleware(async (auth, req) => {
  if (isProtectedRoute(req)) {
    await auth.protect()
  }
})

export const config = {
  matcher: [
    '/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)',
    '/(api|trpc)(.*)',
  ],
}
```

공개 라우트를 제외한 모든 라우트 보호:

```typescript
// middleware.ts
import { clerkMiddleware, createRouteMatcher } from '@clerk/nextjs/server'

const isPublicRoute = createRouteMatcher([
  '/sign-in(.*)',
  '/sign-up(.*)',
  '/',
  '/about',
])

export default clerkMiddleware(async (auth, req) => {
  if (!isPublicRoute(req)) {
    await auth.protect()
  }
})
```

### useAuth 훅

인증 상태 및 토큰 접근:

```tsx
'use client'
import { useAuth } from '@clerk/nextjs'

export default function ExternalDataPage() {
  const { userId, sessionId, getToken, isLoaded, isSignedIn } = useAuth()

  const fetchExternalData = async () => {
    const token = await getToken()
    const response = await fetch('https://api.example.com/data', {
      headers: { Authorization: `Bearer ${token}` },
    })
    return response.json()
  }

  if (!isLoaded) return <div>Loading...</div>
  if (!isSignedIn) return <div>Sign in to view this page</div>

  return (
    <div>
      <p>User ID: {userId}</p>
      <p>Session ID: {sessionId}</p>
      <button onClick={fetchExternalData}>Fetch Data</button>
    </div>
  )
}
```

### useUser 훅

사용자 프로필 데이터 접근:

```tsx
'use client'
import { useUser } from '@clerk/nextjs'

export default function ProfilePage() {
  const { isSignedIn, user, isLoaded } = useUser()

  if (!isLoaded) return <div>Loading...</div>
  if (!isSignedIn) return <div>Sign in to view your profile</div>

  return (
    <div>
      <h1>Welcome, {user.firstName}!</h1>
      <p>Email: {user.primaryEmailAddress?.emailAddress}</p>
      <img src={user.imageUrl} alt="Profile" width={100} height={100} />
    </div>
  )
}
```

### 로그인 및 회원가입 페이지

전용 인증 페이지:

```tsx
// app/sign-in/[[...sign-in]]/page.tsx
import { SignIn } from '@clerk/nextjs'

export default function SignInPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <SignIn />
    </div>
  )
}
```

```tsx
// app/sign-up/[[...sign-up]]/page.tsx
import { SignUp } from '@clerk/nextjs'

export default function SignUpPage() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <SignUp />
    </div>
  )
}
```

### 서버 측 인증

App Router 서버 컴포넌트:

```tsx
// app/dashboard/page.tsx
import { auth, currentUser } from '@clerk/nextjs/server'
import { redirect } from 'next/navigation'

export default async function DashboardPage() {
  const { userId } = await auth()
  if (!userId) redirect('/sign-in')

  const user = await currentUser()

  return (
    <div>
      <h1>Dashboard</h1>
      <p>Welcome, {user?.firstName}!</p>
    </div>
  )
}
```

Route Handler 인증:

```typescript
// app/api/user/route.ts
import { auth } from '@clerk/nextjs/server'
import { NextResponse } from 'next/server'

export async function GET() {
  const { userId } = await auth()
  if (!userId) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 })
  }
  return NextResponse.json({ userId })
}
```

### 조직 관리

OrganizationSwitcher 컴포넌트:

```tsx
// app/dashboard/layout.tsx
import { OrganizationSwitcher } from '@clerk/nextjs'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div>
      <nav className="flex items-center gap-4 p-4">
        <OrganizationSwitcher />
      </nav>
      {children}
    </div>
  )
}
```

useOrganizationList를 사용한 커스텀 조직 전환:

```tsx
'use client'
import { useOrganizationList } from '@clerk/nextjs'

export function CustomOrganizationSwitcher() {
  const { isLoaded, setActive, userMemberships } = useOrganizationList({
    userMemberships: { infinite: true },
  })

  if (!isLoaded) return <p>Loading...</p>

  return (
    <div>
      <h2>Your Organizations</h2>
      <ul>
        {userMemberships.data?.map((membership) => (
          <li key={membership.id}>
            <span>{membership.organization.name}</span>
            <button
              onClick={() => setActive({ organization: membership.organization.id })}
            >
              Select
            </button>
          </li>
        ))}
      </ul>
      {userMemberships.hasNextPage && (
        <button onClick={() => userMemberships.fetchNext()}>Load more</button>
      )}
    </div>
  )
}
```

---

## 고급 패턴

### Core 2 마이그레이션

환경 변수 변경:

```bash
# Core 1 (지원 종료)
CLERK_FRONTEND_API=clerk.xxx.lcl.dev
CLERK_API_KEY=sk_xxx

# Core 2 (현재)
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_xxx
CLERK_SECRET_KEY=sk_test_xxx
```

authMiddleware에서 clerkMiddleware로 마이그레이션:

```typescript
// Core 1 (지원 종료) - 사용 금지
import { authMiddleware } from '@clerk/nextjs'
export default authMiddleware()

// Core 2 (현재) - 이것을 사용
import { clerkMiddleware, createRouteMatcher } from '@clerk/nextjs/server'

const isPublicRoute = createRouteMatcher(['/sign-in(.*)', '/sign-up(.*)'])

export default clerkMiddleware(async (auth, request) => {
  if (!isPublicRoute(request)) {
    await auth.protect()
  }
})
```

서버 임포트 경로 변경:

```typescript
// Core 1 (지원 종료)
import { auth } from '@clerk/nextjs'

// Core 2 (현재)
import { auth } from '@clerk/nextjs/server'
```

마이그레이션 도구:

```bash
npx @clerk/upgrade --from=core-1
```

### 역할 기반 접근 제어

권한으로 라우트 보호:

```typescript
// middleware.ts
import { clerkMiddleware, createRouteMatcher } from '@clerk/nextjs/server'

const isAdminRoute = createRouteMatcher(['/admin(.*)'])
const isMemberRoute = createRouteMatcher(['/dashboard(.*)'])

export default clerkMiddleware(async (auth, req) => {
  if (isAdminRoute(req)) {
    await auth.protect((has) => has({ permission: 'org:admin:access' }))
  }

  if (isMemberRoute(req)) {
    await auth.protect()
  }
})
```

컴포넌트에서 권한 확인:

```tsx
'use client'
import { useAuth } from '@clerk/nextjs'

export function AdminPanel() {
  const { has, isLoaded } = useAuth()

  if (!isLoaded) {
    return <div>Loading...</div>
  }

  const isAdmin = has?.({ permission: 'org:admin:access' })

  if (!isAdmin) {
    return <div>Access denied</div>
  }

  return <div>Admin Panel Content</div>
}
```

### 웹훅 통합

웹훅 핸들러:

```typescript
// app/api/webhooks/clerk/route.ts
import { Webhook } from 'svix'
import { headers } from 'next/headers'
import { WebhookEvent } from '@clerk/nextjs/server'

export async function POST(req: Request) {
  const WEBHOOK_SECRET = process.env.CLERK_WEBHOOK_SECRET

  if (!WEBHOOK_SECRET) {
    throw new Error('Missing CLERK_WEBHOOK_SECRET')
  }

  const headerPayload = await headers()
  const svix_id = headerPayload.get('svix-id')
  const svix_timestamp = headerPayload.get('svix-timestamp')
  const svix_signature = headerPayload.get('svix-signature')

  if (!svix_id || !svix_timestamp || !svix_signature) {
    return new Response('Missing svix headers', { status: 400 })
  }

  const payload = await req.json()
  const body = JSON.stringify(payload)

  const wh = new Webhook(WEBHOOK_SECRET)
  let evt: WebhookEvent

  try {
    evt = wh.verify(body, {
      'svix-id': svix_id,
      'svix-timestamp': svix_timestamp,
      'svix-signature': svix_signature,
    }) as WebhookEvent
  } catch (err) {
    return new Response('Invalid signature', { status: 400 })
  }

  const eventType = evt.type

  if (eventType === 'user.created') {
    const { id, email_addresses, first_name, last_name } = evt.data
    // 데이터베이스에 사용자 동기화
  }

  if (eventType === 'user.updated') {
    const { id, first_name, last_name } = evt.data
    // 데이터베이스의 사용자 업데이트
  }

  if (eventType === 'user.deleted') {
    const { id } = evt.data
    // 사용자 삭제 처리
  }

  return new Response('Webhook received', { status: 200 })
}
```

### 커스텀 인증 플로우

useSignIn을 사용한 커스텀 로그인:

```tsx
'use client'
import { useSignIn } from '@clerk/nextjs'
import { useState } from 'react'

export function CustomSignIn() {
  const { signIn, isLoaded, setActive } = useSignIn()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  if (!isLoaded) {
    return <div>Loading...</div>
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')

    try {
      const result = await signIn.create({
        identifier: email,
        password,
      })

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId })
      }
    } catch (err: any) {
      setError(err.errors?.[0]?.message || 'Sign in failed')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
      />
      {error && <p className="text-red-500">{error}</p>}
      <button type="submit">Sign In</button>
    </form>
  )
}
```

### 외부 서비스 연동

Clerk와 Supabase 연동:

```typescript
// lib/supabase.ts
import { createClient } from '@supabase/supabase-js'
import { auth } from '@clerk/nextjs/server'

export async function createClerkSupabaseClient() {
  const { getToken } = await auth()
  const supabaseToken = await getToken({ template: 'supabase' })

  return createClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      global: {
        headers: {
          Authorization: `Bearer ${supabaseToken}`,
        },
      },
    }
  )
}
```

Clerk와 Convex 연동:

```tsx
// app/providers.tsx
'use client'
import { ClerkProvider, useAuth } from '@clerk/nextjs'
import { ConvexProviderWithClerk } from 'convex/react-clerk'
import { ConvexReactClient } from 'convex/react'

const convex = new ConvexReactClient(process.env.NEXT_PUBLIC_CONVEX_URL!)

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ClerkProvider>
      <ConvexProviderWithClerk client={convex} useAuth={useAuth}>
        {children}
      </ConvexProviderWithClerk>
    </ClerkProvider>
  )
}
```

---

## 리소스

공식 문서:
- 빠른 시작: https://clerk.com/docs/quickstarts/nextjs
- SDK 레퍼런스: https://clerk.com/docs/reference/nextjs/overview
- Core 2 마이그레이션: https://clerk.com/docs/guides/development/upgrading/upgrade-guides/core-2/nextjs
- 웹훅: https://clerk.com/docs/integrations/webhooks

연관 스킬:
- do-platform-auth0: 엔터프라이즈 SSO 솔루션
- do-platform-supabase: Supabase 인증 연동
- do-platform-vercel: Vercel 배포와 Clerk 연동
- do-lang-typescript: TypeScript 개발 패턴
- do-domain-frontend: React 및 Next.js 연동

---

상태: 프로덕션 준비 완료
버전: 2.0.0
최종 업데이트: 2025-12-30
SDK 버전: @clerk/nextjs 6.x (Core 2)
