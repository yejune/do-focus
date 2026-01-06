---
name: do-library-nextra
description: Next.js 기반 Nextra 문서화 프레임워크. 문서 사이트, 지식 베이스, API 참고 문서 구축
version: 2.1.0
modularized: true
allowed-tools: Read, Write, Edit, Grep, Glob
aliases:
  - do-library-nextra
category: library
---

## 역할

Nextra + Next.js 기반 문서 사이트 구축 전문 스킬

Nextra 핵심 장점:
- Zero config MDX (Markdown + JSX 통합)
- 파일 시스템 라우팅 (자동 경로 생성)
- 성능 최적화 (코드 스플리팅, 프리페칭)
- 테마 시스템 (커스터마이징)
- 국제화(i18n) 내장

핵심 파일:
- `pages/` - MDX 문서 페이지
- `theme.config.tsx` - 사이트 설정
- `_meta.js` - 네비게이션 구조

지원 버전:
- Nextra 3.x: Next.js 13.x, 14.x (Pages Router)
- Nextra 4.x: Next.js 14.x, 15.x (App Router, Turbopack)

---

## 핵심 패턴

### 프로젝트 초기화

```bash
npx create-nextra-app@latest my-docs --template docs
```

프로젝트 구조:
```
my-docs/
  pages/
    _app.tsx
    _meta.json
    index.mdx
    docs/
      _meta.json
      getting-started.mdx
  theme.config.tsx
  next.config.js
```

### 테마 설정 (Nextra 3.x)

```typescript
// theme.config.tsx
import { DocsThemeConfig } from 'nextra-theme-docs';

const config: DocsThemeConfig = {
  logo: <span>My Documentation</span>,
  logoLink: '/',
  project: { link: 'https://github.com/username/project' },
  docsRepositoryBase: 'https://github.com/username/project/tree/main/docs',

  sidebar: {
    defaultMenuCollapseLevel: 1,
    toggleButton: true,
  },

  toc: { title: 'On This Page', float: true, backToTop: true },
  navigation: { prev: true, next: true },
  editLink: { text: 'Edit this page on GitHub' },
  footer: { text: `MIT ${new Date().getFullYear()} - My Project` },
  darkMode: true,

  useNextSeoProps() {
    return { titleTemplate: '%s - My Documentation' };
  },

  head: (
    <>
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta property="og:title" content="My Documentation" />
    </>
  ),

  search: { placeholder: 'Search documentation...' },
  i18n: [
    { locale: 'en', text: 'English' },
    { locale: 'ko', text: 'Korean' },
  ],
};

export default config;
```

### 테마 설정 (Nextra 4.x App Router)

```typescript
// app/layout.tsx
import { Layout } from 'nextra-theme-docs';
import { getPageMap } from 'nextra/page-map';

export default async function RootLayout({ children }) {
  const pageMap = await getPageMap();

  return (
    <html lang="en">
      <body>
        <Layout
          pageMap={pageMap}
          docsRepositoryBase="https://github.com/username/project/tree/main"
          darkMode={true}
          sidebar={{ defaultMenuCollapseLevel: 2, toggleButton: true }}
          toc={{ title: 'On This Page', float: true, backToTop: 'Back to top' }}
          i18n={[
            { locale: 'en', name: 'English' },
            { locale: 'ko', name: 'Korean' },
          ]}
        >
          {children}
        </Layout>
      </body>
    </html>
  );
}
```

### 네비게이션 구조 (_meta.json)

```javascript
// pages/_meta.json (루트)
{
  "index": { "title": "Home", "type": "page", "display": "hidden" },
  "docs": { "title": "Documentation", "type": "page" }
}

// pages/docs/_meta.json (섹션)
{
  "index": "Overview",
  "getting-started": "Getting Started",
  "---": { "type": "separator", "title": "Guide" },
  "installation": "Installation",
  "advanced": {
    "title": "Advanced Topics",
    "type": "menu",
    "items": { "performance": "Performance", "security": "Security" }
  },
  "github": {
    "title": "GitHub",
    "href": "https://github.com/myorg/repo",
    "newWindow": true
  }
}
```

### MDX & React 통합

```mdx
import { Callout, Tabs, Tab, Steps, Cards, Card } from 'nextra/components';

# 컴포넌트 라이브러리

<Callout type="info">인터랙티브 문서화 기능 데모</Callout>

## 설치

<Tabs items={['npm', 'yarn', 'pnpm']}>
  <Tab>
    ```bash
    npm install @myproject/components
    ```
  </Tab>
</Tabs>

## 빠른 시작

<Steps>
### 컴포넌트 임포트
```tsx
import { Button } from '@myproject/components';
```

### 앱에서 사용
```tsx
<Button variant="primary">Click me</Button>
```
</Steps>

<Cards>
  <Card title="Button" href="/docs/button">주요 인터랙션 컴포넌트</Card>
</Cards>
```

### next.config.js 설정

Nextra 3.x:
```javascript
const withNextra = require('nextra')({
  theme: 'nextra-theme-docs',
  themeConfig: './theme.config.tsx',
  staticImage: true,
  flexsearch: { codeblocks: true },
  defaultShowCopyCode: true,
});

module.exports = withNextra({ reactStrictMode: true });
```

Nextra 4.x:
```javascript
import nextra from 'nextra';
const withNextra = nextra({});
export default withNextra({ experimental: { turbopack: true } });
```

---

## 내장 MDX 컴포넌트

- Callout: 알림 박스 (info, warning, error, default)
- Tabs/Tab: 탭 콘텐츠
- Cards/Card: 카드 그리드
- Steps: 단계별 가이드
- FileTree: 파일 구조 시각화

---

## 국제화 (i18n)

```javascript
// next.config.js
module.exports = withNextra({
  i18n: { locales: ['en', 'ko', 'ja'], defaultLocale: 'en' },
});
```

파일 구조:
```
pages/
  index.mdx           # 영어 (기본)
  index.ko.mdx        # 한국어
  docs/
    guide.mdx         # 영어
    guide.ko.mdx      # 한국어
```

---

## 검색 및 SEO

```typescript
const config: DocsThemeConfig = {
  search: { placeholder: 'Search...', emptyResult: <span>No results</span> },

  head: function useHead() {
    const { title } = useConfig();
    const { route } = useRouter();
    return (
      <>
        <meta property="og:title" content={title} />
        <meta property="og:url" content={`https://docs.myproject.com${route}`} />
      </>
    );
  },

  useNextSeoProps() {
    return { titleTemplate: '%s - My Docs' };
  },
};
```

---

## 배포

### Vercel (권장)

```bash
npm install -g vercel
vercel --prod
```

### 정적 내보내기

```javascript
module.exports = withNextra({
  output: 'export',
  images: { unoptimized: true },
  trailingSlash: true,
});
```

### GitHub Actions

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: npm ci && npm run build
      - run: vercel --prod --token ${{ secrets.VERCEL_TOKEN }}
```

---

## 테마 커스터마이징

```css
:root {
  --nextra-primary-hue: 212deg;
  --nextra-content-width: 90rem;
}

.dark { --nextra-bg: 17 17 17; }
```

```javascript
// tailwind.config.js
module.exports = {
  content: ['./pages/**/*.{js,ts,jsx,tsx,mdx}', './theme.config.tsx'],
};
```

---

## 문제 해결

빌드 에러:
1. `.next` 캐시 삭제: `rm -rf .next`
2. 의존성 재설치: `rm -rf node_modules && npm install`
3. Next.js/Nextra 버전 호환성 확인

MDX 파싱 에러:
1. JSX 문법 확인
2. 컴포넌트 임포트 검증
3. 닫히지 않은 태그 확인

검색 작동 안함:
1. FlexSearch 활성화 확인
2. 검색 인덱스 재빌드: `npm run build`

---

## 피해야 할 패턴

- _meta.json 누락: 알파벳 순 정렬 문제
- 5단계 이상 중첩: UX 저하
- 최적화되지 않은 이미지: Next.js Image 사용
- 하드코딩된 링크: 상대 경로 또는 next/link 사용

---

## 추가 학습

- [modules/configuration.md](modules/configuration.md) - theme.config 완전 참고
- [modules/mdx-components.md](modules/mdx-components.md) - MDX 컴포넌트 라이브러리
- [modules/i18n-setup.md](modules/i18n-setup.md) - 국제화 가이드
- [modules/deployment.md](modules/deployment.md) - 호스팅 및 배포

---

## 관련 스킬

- Skill("do-docs-generation") - 코드로부터 자동 문서 생성
- Skill("do-library-mermaid") - 다이어그램 통합
- Skill("do-domain-frontend") - 프론트엔드 개발 패턴

---

## 버전 이력

2.1.0 (2025-12-30): Nextra 4.x App Router, Turbopack 지원
2.0.0 (2025-11-23): Progressive Disclosure 리팩토링
1.0.0 (2025-11-12): 초기 릴리스

---

Domain: 문서화 아키텍처
Generated with: Do Skill Factory
