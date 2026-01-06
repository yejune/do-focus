---
name: do-lang-typescript
description: TypeScript 5.9+ 개발 전문가 - React 19, Next.js 16 App Router, tRPC, Zod 검증 및 최신 TypeScript 패턴
version: 1.0.0
category: language
tags:
  - typescript
  - react
  - nextjs
  - frontend
  - fullstack
updated: 2025-12-07
status: active
allowed-tools:
  - Read
  - Grep
  - Glob
  - Edit
  - Write
  - Bash
---

## 빠른 참조

TypeScript 5.9+ 개발 전문가 - React 19, Next.js 16, 타입 안전 API 패턴

자동 트리거: `.ts`, `.tsx`, `.mts`, `.cts` 파일, TypeScript 설정, React/Next.js 프로젝트

핵심 스택:
- TypeScript 5.9: deferred module evaluation, decorators, satisfies 연산자
- React 19: Server Components, use() hook, Actions
- Next.js 16: App Router, Server Actions, ISR/SSG/SSR
- 타입 안전 API: tRPC 11, Zod 3.23, tanstack-query

빠른 명령어:
```bash
npx create-next-app@latest --typescript --tailwind --app
npm install @trpc/server @trpc/client @trpc/react-query zod @tanstack/react-query
npm install -D vitest @testing-library/react @playwright/test
```

---

## TypeScript 5.9 핵심 기능

### satisfies 연산자

```typescript
type Colors = "red" | "green" | "blue";
const palette = {
  red: [255, 0, 0],
  green: "#00ff00",
} satisfies Record<Colors, string | number[]>;

palette.red.map((n) => n * 2); // red는 number[]
palette.green.toUpperCase();   // green은 string
```

### Deferred Module Evaluation

```typescript
import defer * as analytics from "./heavy-analytics";
function trackEvent(name: string) {
  analytics.track(name); // 첫 사용 시 모듈 로드
}
```

### 고급 타입 패턴

```typescript
// 조건부 타입
type Awaited<T> = T extends Promise<infer U> ? U : T;

// 매핑 타입과 키 변환
type Getters<T> = {
  [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};

// 템플릿 리터럴 타입
type EventHandler = `on${Capitalize<"click" | "focus">}`;
// "onClick" | "onFocus"
```

---

## React 19 패턴

### Server Components

```typescript
// app/users/[id]/page.tsx
interface PageProps { params: Promise<{ id: string }>; }

export default async function UserPage({ params }: PageProps) {
  const { id } = await params;
  const user = await db.user.findUnique({ where: { id } });
  if (!user) notFound();
  return <h1>{user.name}</h1>;
}
```

### use() Hook

```typescript
"use client";
import { use } from "react";

function UserProfile({ userPromise }: { userPromise: Promise<User> }) {
  const user = use(userPromise);
  return <div>{user.name}</div>;
}
```

### Server Actions

```typescript
"use server";
import { z } from "zod";
import { revalidatePath } from "next/cache";

const CreateUserSchema = z.object({
  name: z.string().min(2),
  email: z.string().email(),
});

export async function createUser(formData: FormData) {
  const validated = CreateUserSchema.parse(Object.fromEntries(formData));
  await db.user.create({ data: validated });
  revalidatePath("/users");
}
```

### useActionState

```typescript
"use client";
import { useActionState } from "react";

export function CreateUserForm() {
  const [state, action, isPending] = useActionState(createUser, null);
  return (
    <form action={action}>
      <input name="name" disabled={isPending} />
      <button disabled={isPending}>{isPending ? "생성 중..." : "생성"}</button>
    </form>
  );
}
```

### useOptimistic

```typescript
"use client";
import { useOptimistic, useTransition } from "react";

export function MessageList({ messages }: { messages: Message[] }) {
  const [isPending, startTransition] = useTransition();
  const [optimisticMessages, addOptimistic] = useOptimistic(
    messages,
    (state, newMessage: Message) => [...state, { ...newMessage, sending: true }]
  );

  async function sendMessage(formData: FormData) {
    addOptimistic({ id: crypto.randomUUID(), text: formData.get("text") });
    startTransition(async () => await submitMessage(formData));
  }
  // ...
}
```

---

## Next.js 16 App Router

### 라우트 구조

```
app/
  layout.tsx, page.tsx, loading.tsx, error.tsx
  api/route.ts
  users/page.tsx, users/[id]/page.tsx
  (marketing)/about/page.tsx
```

### Metadata API

```typescript
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: { default: "My App", template: "%s | My App" },
};

export async function generateMetadata({ params }): Promise<Metadata> {
  const { id } = await params;
  return { title: (await getUser(id)).name };
}
```

### 렌더링 전략

```typescript
export const dynamic = "force-static";  // 정적 생성
export const dynamic = "force-dynamic"; // 매 요청마다
export const revalidate = 60;           // 60초마다 재검증

// 캐싱
import { unstable_cache, revalidatePath, revalidateTag } from "next/cache";

const getCachedUser = unstable_cache(
  async (id: string) => db.user.findUnique({ where: { id } }),
  ["user"],
  { revalidate: 3600, tags: ["user"] }
);
```

---

## tRPC 타입 안전 API

### 서버 설정

```typescript
import { initTRPC, TRPCError } from "@trpc/server";

const t = initTRPC.context<Context>().create();

export const router = t.router;
export const publicProcedure = t.procedure;
export const protectedProcedure = t.procedure.use(async ({ ctx, next }) => {
  if (!ctx.session?.user) throw new TRPCError({ code: "UNAUTHORIZED" });
  return next({ ctx: { ...ctx, user: ctx.session.user } });
});
```

### 라우터 정의

```typescript
import { z } from "zod";

export const userRouter = router({
  getById: publicProcedure
    .input(z.object({ id: z.string().uuid() }))
    .query(({ input, ctx }) => ctx.db.user.findUnique({ where: { id: input.id } })),

  create: protectedProcedure
    .input(z.object({ name: z.string().min(2), email: z.string().email() }))
    .mutation(({ input, ctx }) => ctx.db.user.create({ data: input })),
});
```

### 클라이언트

```typescript
"use client";
export function UserList() {
  const { data, isLoading } = trpc.user.list.useQuery({ page: 1 });
  if (isLoading) return <div>로딩 중...</div>;
  return <ul>{data?.map((u) => <li key={u.id}>{u.name}</li>)}</ul>;
}
```

---

## Zod 스키마 패턴

### 복합 검증

```typescript
const UserSchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(2).max(100),
  email: z.string().email(),
  role: z.enum(["admin", "user"]),
}).strict();

type User = z.infer<typeof UserSchema>;

const CreateUserSchema = UserSchema.omit({ id: true })
  .extend({ password: z.string().min(8), confirmPassword: z.string() })
  .refine((d) => d.password === d.confirmPassword, {
    message: "비밀번호 불일치", path: ["confirmPassword"],
  });
```

### 고급 패턴

```typescript
// Discriminated unions
const EventSchema = z.discriminatedUnion("type", [
  z.object({ type: z.literal("click"), x: z.number(), y: z.number() }),
  z.object({ type: z.literal("keypress"), key: z.string() }),
]);

// Branded types
const UserId = z.string().uuid().brand<"UserId">();
```

---

## 상태 관리

### Zustand

```typescript
import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";

interface AuthState {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  devtools(persist((set) => ({
    user: null,
    login: (user) => set({ user }),
    logout: () => set({ user: null }),
  }), { name: "auth-storage" }))
);
```

### Jotai

```typescript
import { atom } from "jotai";
import { atomWithStorage } from "jotai/utils";

const countAtom = atom(0);
const themeAtom = atomWithStorage<"light" | "dark">("theme", "light");
```

---

## 테스팅

```typescript
// vitest.config.ts
import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: ["./src/test/setup.ts"],
  },
});

// 컴포넌트 테스트
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";

describe("PostCard", () => {
  it("제목 렌더링", () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText("Test Post")).toBeInTheDocument();
  });
});
```

---

## 문제 해결

```bash
npx tsc --noEmit                    # 타입 체크
npm run build                       # 빌드 오류 확인
rm -rf .next && npm run dev         # 캐시 클리어
```

타입 안전성:
```typescript
function assertNever(x: never): never { throw new Error(`Unexpected: ${x}`); }

function isUser(v: unknown): v is User {
  return typeof v === "object" && v !== null && "id" in v;
}
```

---

## 관련 Skills

- `do-domain-frontend` - UI 컴포넌트, 스타일링 패턴
- `do-domain-backend` - API 설계, 데이터베이스 통합
- `do-library-shadcn` - 컴포넌트 라이브러리
- `do-workflow-testing` - 테스팅 전략

---

Version: 1.0.0
Last Updated: 2025-12-07
Status: Production Ready
