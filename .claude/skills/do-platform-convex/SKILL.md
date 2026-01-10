---
name: do-platform-convex
description: Convex 실시간 백엔드 전문가 - TypeScript 우선 reactive 패턴, optimistic updates, server functions. 실시간 협업 앱 구축 시 사용.
version: 1.0.0
category: platform
tags: convex, realtime, reactive, typescript, optimistic-updates
context7-libraries: /get-convex/convex
related-skills: do-platform-supabase, do-lang-typescript
allowed-tools: Read, Write, Bash, Grep, Glob
user-invocable: false
---

# do-platform-convex: Convex 실시간 백엔드 전문가

## 빠른 참조

Convex는 TypeScript 우선 설계, 자동 캐싱, optimistic updates를 갖춘 실시간 reactive 백엔드 플랫폼.

### 사용 시기

- 실시간 협업 애플리케이션 (문서, 화이트보드, 채팅)
- 수동 refetch 없이 즉각적인 UI 업데이트가 필요한 앱
- end-to-end 타입 안전성이 필요한 TypeScript 우선 프로젝트
- 복잡한 optimistic update 요구사항이 있는 애플리케이션

### 핵심 개념

- Server Functions: queries (읽기), mutations (쓰기), actions (외부 API)
- Reactive Queries: 기본 데이터 변경 시 자동 재실행
- Optimistic Updates: 서버 확인 전 즉각적인 UI 업데이트
- Automatic Caching: 지능형 무효화가 포함된 빌트인 캐싱

### 빠른 시작

```bash
npm create convex@latest
npx convex dev
```

### Context7 Library: /get-convex/convex

---

## 구현 가이드

### 프로젝트 구조

```
my-app/
  convex/
    _generated/         # 자동 생성 타입 및 API
    schema.ts           # 스키마 정의
    functions/          # 도메인별 서버 함수
    http.ts             # HTTP 엔드포인트
    crons.ts            # 예약 작업
  src/
    ConvexProvider.tsx  # 클라이언트 설정
```

### 스키마 정의

```typescript
// convex/schema.ts
import { defineSchema, defineTable } from 'convex/server'
import { v } from 'convex/values'

export default defineSchema({
  documents: defineTable({
    title: v.string(),
    content: v.string(),
    ownerId: v.string(),
    isPublic: v.boolean(),
    createdAt: v.number(),
    updatedAt: v.number()
  })
    .index('by_owner', ['ownerId'])
    .index('by_public', ['isPublic', 'createdAt'])
    .searchIndex('search_content', {
      searchField: 'content',
      filterFields: ['ownerId', 'isPublic']
    }),
  collaborators: defineTable({
    documentId: v.id('documents'),
    userId: v.string(),
    permission: v.union(v.literal('read'), v.literal('write'))
  })
    .index('by_document', ['documentId'])
    .index('by_user', ['userId'])
})
```

### Validators (v module)

```typescript
import { v } from 'convex/values'
// 기본형: v.string(), v.number(), v.boolean(), v.null(), v.int64(), v.bytes()
// 복합형: v.array(v.string()), v.object({...}), v.union(...), v.optional(...)
// 참조: v.id('tableName')
```

### Query 함수 (Reactive)

```typescript
import { query } from '../_generated/server'
import { v } from 'convex/values'

export const list = query({
  args: { ownerId: v.string() },
  handler: async (ctx, args) => {
    return await ctx.db
      .query('documents')
      .withIndex('by_owner', (q) => q.eq('ownerId', args.ownerId))
      .order('desc')
      .collect()
  }
})

export const getById = query({
  args: { id: v.id('documents') },
  handler: async (ctx, args) => {
    const doc = await ctx.db.get(args.id)
    if (!doc) throw new Error('Document not found')
    return doc
  }
})

export const searchContent = query({
  args: { searchQuery: v.string(), limit: v.optional(v.number()) },
  handler: async (ctx, args) => {
    return await ctx.db
      .query('documents')
      .withSearchIndex('search_content', (q) =>
        q.search('content', args.searchQuery).eq('isPublic', true)
      )
      .take(args.limit ?? 10)
  }
})
```

### Mutation 함수

```typescript
import { mutation } from '../_generated/server'
import { v } from 'convex/values'

export const create = mutation({
  args: { title: v.string(), content: v.string(), isPublic: v.boolean() },
  handler: async (ctx, args) => {
    const identity = await ctx.auth.getUserIdentity()
    if (!identity) throw new Error('Unauthorized')
    return await ctx.db.insert('documents', {
      ...args,
      ownerId: identity.subject,
      createdAt: Date.now(),
      updatedAt: Date.now()
    })
  }
})

export const update = mutation({
  args: { id: v.id('documents'), title: v.optional(v.string()), content: v.optional(v.string()) },
  handler: async (ctx, args) => {
    const { id, ...updates } = args
    const existing = await ctx.db.get(id)
    if (!existing) throw new Error('Document not found')
    const identity = await ctx.auth.getUserIdentity()
    if (existing.ownerId !== identity?.subject) throw new Error('Forbidden')
    await ctx.db.patch(id, { ...updates, updatedAt: Date.now() })
  }
})

export const remove = mutation({
  args: { id: v.id('documents') },
  handler: async (ctx, args) => await ctx.db.delete(args.id)
})
```

### Action 함수 (외부 API)

```typescript
import { action } from '../_generated/server'
import { internal } from '../_generated/api'
import { v } from 'convex/values'

export const generateSummary = action({
  args: { documentId: v.id('documents') },
  handler: async (ctx, args) => {
    const doc = await ctx.runQuery(internal.documents.getById, { id: args.documentId })
    const response = await fetch('https://api.openai.com/v1/completions', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${process.env.OPENAI_API_KEY}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: 'gpt-4', prompt: `Summarize: ${doc.content}`, max_tokens: 150 })
    })
    const result = await response.json()
    await ctx.runMutation(internal.documents.updateSummary, { id: args.documentId, summary: result.choices[0].text })
    return result.choices[0].text
  }
})
```

### React 클라이언트 설정

```typescript
import { ConvexProvider, ConvexReactClient } from 'convex/react'
import { ConvexProviderWithClerk } from 'convex/react-clerk'

const convex = new ConvexReactClient(import.meta.env.VITE_CONVEX_URL)

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ConvexProviderWithClerk client={convex} useAuth={useAuth}>
      {children}
    </ConvexProviderWithClerk>
  )
}
```

### React Hooks 사용법

```typescript
import { useQuery, useMutation } from 'convex/react'
import { api } from '../../convex/_generated/api'

export function DocumentList({ userId }: { userId: string }) {
  const documents = useQuery(api.functions.documents.list, { ownerId: userId })
  const createDocument = useMutation(api.functions.documents.create)
  if (documents === undefined) return <Loading />
  return (
    <div>
      <button onClick={() => createDocument({ title: 'New', content: '', isPublic: false })}>
        새 문서
      </button>
      {documents.map((doc) => <DocumentCard key={doc._id} document={doc} />)}
    </div>
  )
}
```

---

## 고급 패턴

### Optimistic Updates

```typescript
import { useMutation } from 'convex/react'
import { api } from '../../convex/_generated/api'

export function useOptimisticUpdate() {
  return useMutation(api.functions.documents.update)
    .withOptimisticUpdate((localStore, args) => {
      const { id, ...updates } = args
      const existing = localStore.getQuery(api.functions.documents.getById, { id })
      if (existing) {
        localStore.setQuery(api.functions.documents.getById, { id }, { ...existing, ...updates, updatedAt: Date.now() })
      }
    })
}
```

### 파일 저장소

```typescript
// 서버 측
export const generateUploadUrl = mutation({
  handler: async (ctx) => await ctx.storage.generateUploadUrl()
})

export const saveFile = mutation({
  args: { storageId: v.id('_storage'), fileName: v.string() },
  handler: async (ctx, args) => await ctx.db.insert('files', { ...args, uploadedAt: Date.now() })
})

export const getFileUrl = query({
  args: { storageId: v.id('_storage') },
  handler: async (ctx, args) => await ctx.storage.getUrl(args.storageId)
})
```

```typescript
// 클라이언트 측 업로드
export function useFileUpload() {
  const generateUploadUrl = useMutation(api.functions.files.generateUploadUrl)
  const saveFile = useMutation(api.functions.files.saveFile)
  return async (file: File) => {
    const uploadUrl = await generateUploadUrl()
    const response = await fetch(uploadUrl, { method: 'POST', headers: { 'Content-Type': file.type }, body: file })
    const { storageId } = await response.json()
    await saveFile({ storageId, fileName: file.name })
    return storageId
  }
}
```

### 예약 함수 (Crons)

```typescript
import { cronJobs } from 'convex/server'
import { internal } from './_generated/api'

const crons = cronJobs()
crons.interval('cleanup old drafts', { hours: 24 }, internal.documents.cleanupOldDrafts)
crons.cron('daily analytics', '0 0 * * *', internal.analytics.generateDailyReport)
export default crons
```

### HTTP 엔드포인트

```typescript
import { httpRouter } from 'convex/server'
import { httpAction } from './_generated/server'

const http = httpRouter()
http.route({
  path: '/webhook/stripe',
  method: 'POST',
  handler: httpAction(async (ctx, request) => {
    const body = await request.text()
    await ctx.runMutation(internal.payments.processWebhook, { body, signature: request.headers.get('stripe-signature') })
    return new Response('OK', { status: 200 })
  })
})
export default http
```

### 인증 (Clerk 연동)

```typescript
export const current = query({
  handler: async (ctx) => {
    const identity = await ctx.auth.getUserIdentity()
    if (!identity) return null
    return await ctx.db.query('users').withIndex('by_token', (q) => q.eq('tokenIdentifier', identity.tokenIdentifier)).first()
  }
})

export const ensureUser = mutation({
  handler: async (ctx) => {
    const identity = await ctx.auth.getUserIdentity()
    if (!identity) throw new Error('Unauthorized')
    const existing = await ctx.db.query('users').withIndex('by_token', (q) => q.eq('tokenIdentifier', identity.tokenIdentifier)).first()
    if (existing) return existing._id
    return await ctx.db.insert('users', { tokenIdentifier: identity.tokenIdentifier, email: identity.email, name: identity.name, createdAt: Date.now() })
  }
})
```

### 에러 처리

```typescript
import { ConvexError } from 'convex/values'

export const secureOperation = mutation({
  args: { id: v.id('documents') },
  handler: async (ctx, args) => {
    const identity = await ctx.auth.getUserIdentity()
    if (!identity) throw new ConvexError('UNAUTHORIZED')
    const doc = await ctx.db.get(args.id)
    if (!doc) throw new ConvexError({ code: 'NOT_FOUND', message: 'Document not found' })
  }
})
```

---

## 모범 사례

Query 최적화:
- 모든 필터링된 쿼리에 인덱스 사용
- 대용량 데이터셋에는 페이지네이션 쿼리 선호
- 전문 검색에는 search index 사용

Mutation 설계:
- mutation을 집중적이고 원자적으로 유지
- 다단계 작업에는 internal mutation 사용
- v module로 모든 입력 검증

---

## 관련 Skills

- do-platform-supabase - 대체 PostgreSQL 기반 백엔드
- do-lang-typescript - TypeScript 패턴 및 모범 사례
- do-domain-frontend - React 통합 패턴
- do-foundation-quality - 인증 및 권한 패턴

---

Status: Production Ready
Last Updated: 2026-01-06
Platform: Convex Real-time Backend
