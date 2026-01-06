---
name: do-library-shadcn
description: Do Lib Shadcn Ui - 전문 구현 가이드
version: 2.0.0
allowed-tools:
  - Read
  - Glob
  - Grep
---

# Enterprise shadcn/ui 컴포넌트 라이브러리

## 개요

shadcn/ui는 Radix UI와 Tailwind CSS로 구축된 재사용 가능한 컴포넌트 모음. npm 패키지가 아닌 프로젝트에 직접 복사하는 방식으로, 완전한 제어와 소유권을 보장함.

핵심 이점:
- 컴포넌트에 대한 완전한 제어와 소유권
- 의존성 없음 (Radix UI 프리미티브만 사용)
- Tailwind CSS로 완전한 커스터마이제이션
- TypeScript 우선으로 뛰어난 타입 안전성
- WCAG 2.1 AA 준수의 내장 접근성

---

## 설치 및 설정

### 1단계: shadcn/ui 초기화

```bash
npx shadcn-ui@latest init
```

### 2단계: components.json 구성

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "app/globals.css",
    "baseColor": "slate",
    "cssVariables": true,
    "prefix": ""
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui"
  }
}
```

### 3단계: 컴포넌트 추가

```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add form
npx shadcn-ui@latest add dialog
```

---

## 컴포넌트 카테고리

### 폼 컴포넌트
- Input, Select, Checkbox, Radio, Textarea
- react-hook-form + Zod 검증 통합
- 적절한 ARIA 레이블로 접근성 보장

### 디스플레이 컴포넌트
- Card, Dialog, Sheet, Drawer, Popover
- 반응형 디자인 패턴
- 다크 모드 지원

### 네비게이션 컴포넌트
- Navigation Menu, Breadcrumb, Tabs, Pagination
- 키보드 네비게이션 지원
- 포커스 관리

### 데이터 컴포넌트
- Table, Calendar, DatePicker, Charts
- 대규모 데이터세트용 가상 스크롤링
- TanStack Table 통합

### 피드백 컴포넌트
- Alert, Toast, Progress, Badge, Avatar
- 로딩 상태와 스켈레톤
- 에러 바운더리

---

## 핵심 구현 표준

[HARD] 색상 값에 CSS 변수만 사용
- 동적 테마 활성화
- 다크 모드 전환 지원
- 모든 컴포넌트에서 디자인 시스템 일관성 유지
- CSS 변수 없이는 테마 변경에 코드 수정 필요

[HARD] 모든 대화형 요소에 접근성 속성 포함
- WCAG 2.1 AA 준수 보장
- 스크린 리더 호환성
- 장애가 있는 사용자를 위한 포용적 사용자 경험

[HARD] 모든 대화형 컴포넌트에 키보드 네비게이션 구현
- 키보드 사용자를 위한 필수 네비게이션 방법 제공
- 보조 기술 지원
- 전반적인 사용자 경험 효율성 향상

[SOFT] 비동기 작업에 로딩 상태 제공
- 사용자에게 작업 진행 상황 전달
- 체감 지연 시간 감소
- 앱 응답성에 대한 사용자 신뢰도 향상

[HARD] 컴포넌트 트리 주위에 에러 바운더리 구현
- 격리된 컴포넌트 실패로 전체 앱 충돌 방지
- 우아한 오류 복구 가능
- 앱 안정성 유지

[HARD] 인라인 스타일 대신 Tailwind CSS 클래스 적용
- 디자인 시스템과 일관성 유지
- JIT 컴파일 이점 활용
- 반응형 디자인 변형 지원
- 번들 크기 최적화 개선

[SOFT] 모든 컴포넌트에 다크 모드 지원 구현
- 사용자 선호도 존중
- 저조도 환경에서 눈 피로 감소

---

## 테마 시스템

### 테마 프로바이더

```typescript
import { createContext, useContext, useEffect, useState } from "react";

type Theme = "dark" | "light" | "system";

interface ThemeProviderState {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

const ThemeProviderContext = createContext<ThemeProviderState | undefined>(undefined);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "ui-theme",
}) {
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window !== "undefined") {
      return (localStorage.getItem(storageKey) as Theme) || defaultTheme;
    }
    return defaultTheme;
  });

  useEffect(() => {
    const root = window.document.documentElement;
    root.classList.remove("light", "dark");

    if (theme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches ? "dark" : "light";
      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(theme);
  }, [theme]);

  return (
    <ThemeProviderContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);
  if (!context) throw new Error("useTheme must be used within ThemeProvider");
  return context;
};
```

### CSS 변수 설정 (globals.css)

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 240 10% 3.9%;
    --primary: 240 9% 10%;
    --primary-foreground: 0 0% 98%;
    --secondary: 240 4.8% 95.9%;
    --muted: 240 4.8% 95.9%;
    --accent: 240 4.8% 95.9%;
    --destructive: 0 72.22% 50.59%;
    --border: 240 5.9% 90%;
    --radius: 0.5rem;
  }

  .dark {
    --background: 240 10% 3.9%;
    --foreground: 0 0% 98%;
    --primary: 0 0% 98%;
    --primary-foreground: 240 9% 10%;
    --secondary: 240 3.7% 15.9%;
    --muted: 240 3.7% 15.9%;
    --accent: 240 3.7% 15.9%;
    --destructive: 0 62.8% 30.6%;
    --border: 240 3.7% 15.9%;
  }

  * { @apply border-border; }
  body { @apply bg-background text-foreground; }
}
```

---

## 고급 패턴

### 컴포넌트 조합

```typescript
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export function DashboardCard({ title, children }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {children}
      </CardContent>
    </Card>
  );
}
```

### 폼 검증 (Zod + React Hook Form)

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

type FormValues = z.infer<typeof formSchema>;

export function LoginForm() {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
  });

  return (
    <form onSubmit={form.handleSubmit(onSubmit)}>
      {/* 폼 필드 */}
    </form>
  );
}
```

### 테마 인식 컴포넌트

```typescript
export function StatusBadge({ status }: { status: "success" | "error" | "warning" }) {
  const statusConfig = {
    success: "bg-green-500/20 text-green-700 dark:text-green-400",
    error: "bg-red-500/20 text-red-700 dark:text-red-400",
    warning: "bg-yellow-500/20 text-yellow-700 dark:text-yellow-400",
  };

  return (
    <span className={`px-3 py-1 rounded-md text-sm font-medium ${statusConfig[status]}`}>
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  );
}
```

### 반응형 디자인

```typescript
export function ResponsiveCard() {
  return (
    <Card className="w-full sm:max-w-sm md:max-w-md lg:max-w-lg">
      <CardHeader className="p-4 sm:p-6">
        <CardTitle className="text-lg sm:text-xl md:text-2xl">
          반응형 카드
        </CardTitle>
      </CardHeader>
      <CardContent className="p-4 sm:p-6">
        <p className="text-sm sm:text-base md:text-lg">
          화면 크기에 맞춤
        </p>
      </CardContent>
    </Card>
  );
}
```

---

## 성능 최적화

### 번들 크기 최적화

```typescript
// 좋음: 특정 컴포넌트만 임포트
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

// 피하기: 전체 라이브러리 임포트
import * as UI from "@/components/ui";
```

### 컴포넌트 지연 로딩

```typescript
import React from "react";

const DataTableComponent = React.lazy(() =>
  import("@/components/data-table").then(mod => ({ default: mod.DataTable }))
);

export function Dashboard() {
  return (
    <React.Suspense fallback={<p>로딩 중...</p>}>
      <DataTableComponent />
    </React.Suspense>
  );
}
```

### React.memo로 메모이제이션

```typescript
import React from "react";
import { Card, CardContent } from "@/components/ui/card";

interface UserCardProps {
  user: { name: string; email: string };
}

export const UserCard = React.memo(function UserCard({ user }: UserCardProps) {
  return (
    <Card>
      <CardContent>
        <p>{user.name}</p>
        <p className="text-gray-500">{user.email}</p>
      </CardContent>
    </Card>
  );
});
```

### 대규모 목록용 가상 스크롤링

```typescript
import { FixedSizeList } from "react-window";
import { Card } from "@/components/ui/card";

export function VirtualList({ items }: { items: { id: string; name: string }[] }) {
  const Row = ({ index, style }) => (
    <div style={style}>
      <Card className="p-2">
        <p>{items[index].name}</p>
      </Card>
    </div>
  );

  return (
    <FixedSizeList
      height={600}
      itemCount={items.length}
      itemSize={80}
      width="100%"
    >
      {Row}
    </FixedSizeList>
  );
}
```

### 성능 모범 사례 요약

1. 컴포넌트 지연 로딩 - 비필수 컴포넌트에 React.lazy 사용
2. 선택적 메모이제이션 - 과도한 메모이제이션 금지, 먼저 프로파일링
3. 트리 쉐이킹 임포트 - 전체 라이브러리가 아닌 특정 컴포넌트 임포트
4. CSS 최적화 - Tailwind content 속성 올바르게 구성
5. 대규모 목록 가상 스크롤 - 1000개 이상 항목에 react-window 사용
6. 계산 캐싱 - useMemo와 useCallback 적절히 사용
7. 번들 모니터링 - 번들 분석기로 병목 현상 식별

---

## 기술 스택 (2025년 11월 기준)

기반 기술:
- React 19: 서버 컴포넌트 지원, 동시 렌더링
- TypeScript 5.5: 완전한 타입 안전성
- Tailwind CSS 3.4: JIT 컴파일, CSS 변수 지원
- Radix UI: 스타일 없는 접근 가능한 프리미티브

통합 스택:
- React Hook Form: 폼 상태 관리
- Zod: 스키마 검증
- class-variance-authority: 변형 관리
- Framer Motion: 애니메이션 라이브러리
- Lucide React: 아이콘 라이브러리

---

## 관련 스킬

- [shadcn 컴포넌트](modules/shadcn-components.md) - 고급 컴포넌트 패턴
- [shadcn 테마](modules/shadcn-theming.md) - 테마 시스템과 커스터마이제이션
- `do-domain-uiux` - 디자인 시스템 아키텍처
- `do-lang-typescript` - TypeScript 모범 사례
- `do-domain-frontend` - 프론트엔드 개발 패턴

---

## 공식 리소스

- [shadcn/ui 문서](https://ui.shadcn.com/docs)
- [API 레퍼런스](https://ui.shadcn.com/docs/components)
- [Radix UI 문서](https://www.radix-ui.com/)
- [Tailwind CSS 문서](https://tailwindcss.com/)

---

Last Updated: 2025-01-06
Status: Production Ready
