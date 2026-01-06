---
name: do-domain-frontend
description: Frontend development specialist covering React 19, Next.js 16, Vue 3.5, and modern UI/UX patterns with component architecture
version: 1.0.0
category: domain
allowed-tools: Read, Write, Edit, Grep, Glob
tags:
  - frontend
  - react
  - nextjs
  - vue
  - ui
  - components
updated: 2025-12-30
status: active
author: Do Team
---

# Frontend 개발 전문가

## 빠른 참조

React 19, Next.js 16, Vue 3.5 및 현대적 UI/UX 아키텍처를 포함한 종합 frontend 패턴 제공

핵심 능력:
- React 19: Server components, concurrent features, 최적화 패턴
- Next.js 16: App router, server actions, 고급 최적화
- Vue 3.5: Composition API, TypeScript 통합, reactivity
- Component Architecture: Design systems, component libraries
- Responsive Design: Mobile-first, 접근성, 성능 최적화

사용 시기:
- 현대적 웹 애플리케이션 개발
- Component library 생성 및 관리
- Frontend 성능 최적화
- 접근성을 고려한 UI/UX 구현

---

## 구현 가이드

### React 19 Server Components

Server Component와 Client Component 분리:
```tsx
// Server Component - 비동기 데이터 페칭
export default async function UserProfile({ userId }) {
  const user = await getUser(userId)
  return (
    <div>
      <h2>{user.name}</h2>
      <ClientActions userId={userId} />
    </div>
  )
}

// Client Component - 상호작용 처리
'use client'
function ClientActions({ userId }) {
  const [following, setFollowing] = useState(false)
  return <button onClick={() => setFollowing(!following)}>Follow</button>
}
```

Concurrent Features:
```tsx
<ErrorBoundary fallback={<Error />}>
  <Suspense fallback={<Loading />}>
    <UserProfile userId="123" />
  </Suspense>
</ErrorBoundary>
```

### Next.js 16 App Router

Server Actions:
```tsx
// app/actions/users.ts
'use server'
export async function createUser(formData: FormData) {
  const user = await db.user.create({
    data: { name: formData.get('name'), email: formData.get('email') }
  })
  revalidatePath('/users')
  redirect('/users/' + user.id)
}

// app/users/page.tsx
export default function UsersPage() {
  return (
    <form action={createUser}>
      <input name="name" required />
      <input name="email" type="email" required />
      <button type="submit">Create</button>
    </form>
  )
}
```

Dynamic Routes:
```tsx
// app/[category]/[slug]/page.tsx
export default async function Page({ params, searchParams }) {
  const data = await getData(params.category, params.slug)
  return <Content data={data} />
}

export async function generateStaticParams() {
  const posts = await getAllPosts()
  return posts.map(post => ({ category: post.category, slug: post.slug }))
}
```

### Vue 3.5 Composition API

```vue
<script setup lang="ts">
const props = defineProps<{ userId: string }>()
const emit = defineEmits<{ userLoaded: [user: User] }>()

const user = ref<User | null>(null)
const loading = ref(true)

const initials = computed(() =>
  user.value?.name.split(' ').map(n => n[0]).join('').toUpperCase() ?? '??'
)

const fetchUser = async () => {
  loading.value = true
  user.value = await getUser(props.userId)
  emit('userLoaded', user.value)
  loading.value = false
}

watchEffect(fetchUser)
</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="user">
    <span>{{ initials }}</span>
    <h3>{{ user.name }}</h3>
  </div>
</template>
```

### Component Architecture

CVA 기반 Button 컴포넌트:
```tsx
const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md font-medium transition-colors',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground',
        destructive: 'bg-destructive text-destructive-foreground',
        outline: 'border border-input hover:bg-accent',
        ghost: 'hover:bg-accent',
      },
      size: { default: 'h-10 px-4', sm: 'h-9 px-3', lg: 'h-11 px-8' },
    },
    defaultVariants: { variant: 'default', size: 'default' },
  }
)

const Button = forwardRef(({ className, variant, size, ...props }, ref) => (
  <button className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
))
```

Compound Components:
```tsx
const CardContext = createContext({ variant: 'default' })

function Card({ variant = 'default', children }) {
  return (
    <CardContext.Provider value={{ variant }}>
      <div className={`card card--${variant}`}>{children}</div>
    </CardContext.Provider>
  )
}
function CardHeader({ children }) { return <div className="card__header">{children}</div> }
function CardContent({ children }) { return <div className="card__content">{children}</div> }

// Usage: <Card><CardHeader>Title</CardHeader><CardContent>...</CardContent></Card>
```

---

## 고급 패턴

### 성능 최적화

React 최적화:
```tsx
const ExpensiveList = memo(({ items, onItemClick }) => {
  const total = useMemo(() => items.reduce((sum, item) => sum + item.value, 0), [items])
  const handleClick = useCallback((item) => onItemClick(item), [onItemClick])

  return (
    <div>
      <p>Total: {total}</p>
      {items.map(item => <Item key={item.id} onClick={() => handleClick(item)} />)}
    </div>
  )
})
```

Next.js 최적화:
```tsx
// Dynamic import for code splitting
const DynamicChart = dynamic(() => import('@/components/Chart'), {
  loading: () => <div>Loading...</div>,
  ssr: false
})

// Image optimization
<Image src={src} alt={alt} placeholder="blur"
  sizes="(max-width: 768px) 100vw, 50vw" />
```

### 상태 관리 (Zustand)

```tsx
const useUserStore = create(
  devtools(persist(
    (set) => ({
      user: null,
      isLoading: false,
      login: async (email, password) => {
        set({ isLoading: true })
        const user = await authService.login(email, password)
        set({ user, isLoading: false })
      },
      logout: () => set({ user: null }),
    }),
    { name: 'user-storage', partialize: (state) => ({ user: state.user }) }
  ))
)
```

### Custom Hook 패턴

```tsx
function useQuery(queryFn, { enabled = true, onSuccess, onError } = {}) {
  const [data, setData] = useState(null)
  const [isLoading, setIsLoading] = useState(enabled)
  const [error, setError] = useState(null)

  const fetchData = useCallback(async () => {
    setIsLoading(true)
    try {
      const result = await queryFn()
      setData(result)
      onSuccess?.(result)
    } catch (err) {
      setError(err)
      onError?.(err)
    } finally {
      setIsLoading(false)
    }
  }, [queryFn, onSuccess, onError])

  useEffect(() => { if (enabled) fetchData() }, [enabled, fetchData])

  return { data, isLoading, error, refetch: fetchData }
}
```

---

## 문제 해결

### React 이슈

Hydration mismatch:
- 서버/클라이언트 동일 컨텐츠 렌더링 보장
- 동적 컨텐츠에 suppressHydrationWarning 사용
- 브라우저 전용 코드는 useEffect 내부에서 실행
- 클라이언트 전용 컴포넌트에 ssr: false 사용

무한 리렌더링:
- useEffect dependencies에서 object/array reference 확인
- useMemo/useCallback으로 reference 안정성 확보
- effect에서 조건 없이 상태 설정하지 않기

메모리 누수:
- fetch에 AbortController 사용
- useEffect return에서 구독 정리
- unmount 시 async 작업 취소

### Next.js 이슈

빌드 실패:
- 순환 의존성 확인
- 서버 전용 코드가 클라이언트 컴포넌트에 import 되지 않았는지 확인
- 조건부 import에 next/dynamic 사용

성능 문제:
- @next/bundle-analyzer로 번들 분석
- 대형 컴포넌트에 dynamic import 사용
- PPR(Partial Prerendering) 활성화

---

## 안티패턴

### Prop Drilling

잘못된 방법:
```tsx
function App() {
  const [user, setUser] = useState(null)
  return <Layout user={user} setUser={setUser} />  // 여러 레벨 통과
}
```

올바른 방법: Context 사용
```tsx
const UserContext = createContext(null)
function UserProvider({ children }) {
  const [user, setUser] = useState(null)
  return <UserContext.Provider value={{ user, setUser }}>{children}</UserContext.Provider>
}
function useUser() { return useContext(UserContext) }
```

### 인라인 함수

잘못된 방법: 매 렌더링마다 새 함수 생성
```tsx
{items.map(item => <Item onClick={() => handleClick(item.id)} />)}
```

올바른 방법: useCallback 사용
```tsx
const handleClick = useCallback((id) => { /* ... */ }, [])
{items.map(item => <Item id={item.id} onClick={handleClick} />)}
```

### useEffect에서 Fetching

잘못된 방법: 클라이언트 측 초기 데이터 페칭
```tsx
'use client'
useEffect(() => { fetch('/api/users').then(setUser) }, [])
```

올바른 방법: Server Component 사용
```tsx
async function Page({ params }) {
  const user = await getUser(params.id)
  return <Profile user={user} />
}
```

---

## 기술 스택

Frameworks: React 19, Next.js 16, Vue 3.5, Nuxt 3
Languages: TypeScript 5.9+, JavaScript ES2024
Styling: Tailwind CSS 3.4+, CSS Modules
State: Zustand, Redux Toolkit, Pinia
Testing: Vitest, Testing Library, Playwright
Build: Vite 5, Turbopack, SWC
Components: shadcn/ui, Radix UI, Headless UI

---

## 협력 기술

- do-domain-backend: Full-stack 개발
- do-library-shadcn: Component library 통합
- do-domain-uiux: UI/UX 디자인 원칙
- do-lang-typescript: TypeScript 개발

---

Status: Production Ready
Last Updated: 2025-12-30
Maintained by: Do Frontend Team
