---
name: do-domain-uiux
aliases: [do-foundation-uiux]
category: domain
description: 엔터프라이즈 디자인 시스템, 컴포넌트 아키텍처, 접근성, 아이콘, 테마 통합 전문가
version: 2.0.0
modularized: true
tags:
  - domain
  - uiux
  - design-systems
  - accessibility
  - components
  - icons
  - theming
updated: 2025-11-30
status: active
---

## 빠른 참조

엔터프라이즈급 UI/UX 기반으로 디자인 시스템(W3C DTCG 2025.10), 컴포넌트 아키텍처(React 19, Vue 3.5), 접근성(WCAG 2.2), 아이콘 라이브러리(200K+ 아이콘), 테마 시스템을 통합

통합 기능:
- 디자인 시스템: W3C DTCG 2025.10 토큰, Style Dictionary 4.0, Figma MCP 워크플로우
- 컴포넌트 아키텍처: Atomic Design, React 19, Vue 3.5, shadcn/ui, Radix UI primitives
- 접근성: WCAG 2.2 AA/AAA 준수, 키보드 네비게이션, 스크린 리더 최적화
- 아이콘 라이브러리: 10+ 에코시스템 (Lucide, React Icons 35K+, Tabler 5900+, Iconify 200K+)
- 테마: CSS variables, light/dark 모드, 테마 프로바이더, 브랜드 커스터마이징

사용 시기:
- 디자인 시스템 기반의 모던 UI 컴포넌트 라이브러리 구축
- 접근 가능한 엔터프라이즈급 사용자 인터페이스 구현
- 멀티 플랫폼 프로젝트를 위한 디자인 토큰 아키텍처 설정
- 최적의 번들 크기로 포괄적인 아이콘 시스템 통합
- 다크 모드 지원으로 커스터마이징 가능한 테마 시스템 생성

모듈 구성:
- 컴포넌트: [컴포넌트 아키텍처](modules/component-architecture.md) - Atomic Design, 컴포넌트 패턴, Props API
- 디자인 시스템: [디자인 시스템 토큰](modules/design-system-tokens.md) - DTCG 토큰, Style Dictionary, Figma MCP
- 접근성: [접근성 WCAG](modules/accessibility-wcag.md) - WCAG 2.2 준수, 테스트, 네비게이션
- 아이콘: [아이콘 라이브러리](modules/icon-libraries.md) - 10+ 라이브러리, 선택 가이드, 성능 최적화
- 테마: [테마 시스템](modules/theming-system.md) - 테마 시스템, CSS variables, 브랜드 커스터마이징
- 예제: [예제](examples.md) - 실용적인 구현 예제
- 참고: [참고](reference.md) - 외부 문서 링크

---

## 구현 가이드

### 기반 스택 (2025년 11월)

핵심 기술:
- React 19 (Server Components, Concurrent Rendering)
- TypeScript 5.5 (완전한 타입 안전성, 개선된 추론)
- Tailwind CSS 3.4 (JIT 컴파일, CSS variables, 다크 모드)
- Radix UI (스타일 없는 접근 가능한 primitives)
- W3C DTCG 2025.10 (디자인 토큰 명세)
- Style Dictionary 4.0 (토큰 변환)
- Figma MCP (디자인-코드 자동화)
- Storybook 8.x (컴포넌트 문서화)

빠른 결정 매트릭스:
- 디자인 토큰 필요 시: Design System Tokens 모듈 참고 - DTCG 2025.10, Style Dictionary 4.0
- 컴포넌트 패턴 필요 시: Component Architecture 모듈 참고 - Atomic Design, React 19, shadcn/ui
- 접근성 필요 시: Accessibility WCAG 모듈 참고 - WCAG 2.2, jest-axe, 키보드 네비게이션
- 아이콘 필요 시: Icon Libraries 모듈 참고 - Lucide, React Icons, Tabler, Iconify
- 테마 필요 시: Theming System 모듈 참고 - CSS variables, Theme Provider
- 예제 필요 시: Examples 모듈 참고 - React/Vue 구현

---

## 빠른 시작 워크플로우

### 1. 디자인 시스템 설정

Step 1: 디자인 토큰 초기화

DTCG 포맷 JSON 구조:
- schema: "https://tr.designtokens.org/format/"
- color 토큰: type을 "color"로 설정하고 primary.500 값 정의
- spacing 토큰: type을 "dimension"으로 설정하고 md 값을 "1rem"으로 정의

Step 2: Style Dictionary로 토큰 변환
- style-dictionary 패키지를 dev 의존성으로 설치
- npx style-dictionary build 명령으로 빌드

Step 3: 컴포넌트에 통합
- colors, spacing을 tokens 경로에서 import

자세한 토큰 아키텍처: [디자인 시스템 토큰](modules/design-system-tokens.md) 참고

---

### 2. 컴포넌트 라이브러리 설정

Step 1: shadcn/ui 초기화
- npx shadcn-ui init 명령으로 초기화
- npx shadcn-ui add button form dialog 명령으로 컴포넌트 추가

Step 2: Atomic Design 구조 설정

components 디렉토리 구조:
- atoms/ (Button, Input, Label)
- molecules/ (FormGroup, Card)
- organisms/ (DataTable, Modal)

Step 3: 접근성과 함께 구현
- Button에 aria-label="Submit form" 속성 추가
- variant="primary" 속성으로 스타일 변형 적용

자세한 패턴과 예제: [컴포넌트 아키텍처](modules/component-architecture.md) 참고

---

### 3. 아이콘 시스템 통합

Step 1: 아이콘 라이브러리 선택
- 일반 목적: lucide-react
- 최대 다양성: @iconify/react
- 대시보드 최적화: @tabler/icons-react

Step 2: 타입 안전 아이콘 구현
- lucide-react에서 Heart, Search를 개별 import
- className으로 크기와 색상 스타일링 적용

자세한 라이브러리 비교 및 최적화: [아이콘 라이브러리](modules/icon-libraries.md) 참고

---

### 4. 테마 시스템 설정

Step 1: CSS variables 구성

root 선택자에서 정의:
- --primary: 222.2 47.4% 11.2%;
- --background: 0 0% 100%;

dark 클래스에서 정의:
- --primary: 210 40% 98%;
- --background: 222.2 84% 4.9%;

Step 2: Theme Provider 구현
- ThemeProvider 컴포넌트로 App 래핑
- attribute="class", defaultTheme="system" 설정

자세한 테마 시스템: [테마 시스템](modules/theming-system.md) 참고

---

## 핵심 원칙

1. Design Token 우선:
- 디자인 결정의 단일 정보원
- 시맨틱 네이밍 (`color.primary.500` not `blue-500`)
- 멀티 테마 지원 (light/dark)
- 플랫폼 독립적인 변환

2. 기본 접근성:
- WCAG 2.2 AA 최소 (4.5:1 텍스트 대비)
- 모든 인터랙티브 요소의 키보드 네비게이션
- 스크린 리더를 위한 ARIA 속성
- Focus 관리 및 시각적 표시기

3. 컴포넌트 구성:
- Atomic Design 계층 (Atoms -> Molecules -> Organisms)
- 재사용을 위한 Props API
- 변형 기반 스타일링 (별도 컴포넌트가 아님)
- TypeScript로 타입 안전성

4. 성능 최적화:
- 아이콘 Tree-shaking (특정 항목만 import, * 사용 금지)
- 대형 컴포넌트의 지연 로딩
- 비용이 큰 렌더링을 위한 React.memo
- 번들 크기 모니터링

---

## 도구 에코시스템

디자인 토큰 도구:
- W3C DTCG 2025.10: 토큰 명세
- Style Dictionary 4.0+: 토큰 변환

컴포넌트 도구:
- React 19: UI 프레임워크
- shadcn/ui Latest: 컴포넌트 라이브러리
- Radix UI Latest: 접근 가능한 primitives

아이콘 도구:
- Lucide Latest: 1000+ 모던 아이콘
- React Icons Latest: 35K+ 멀티 라이브러리
- Iconify Latest: 200K+ 유니버설

테마 도구:
- Tailwind CSS 3.4: Utility-first CSS
- CSS Variables Native: 테마 토큰

접근성 도구:
- axe DevTools Latest: 접근성 테스트
- jest-axe Latest: 자동화된 a11y 테스트

문서화 도구:
- Storybook 8.x: 컴포넌트 문서
- Figma MCP Latest: 디자인-코드

---

## 모듈 상호 참조

### 컴포넌트 아키텍처
포커스: 컴포넌트 아키텍처 및 구현 패턴

핵심 주제:
- Atomic Design (Atoms, Molecules, Organisms)
- React 19 + Server Components
- Vue 3.5 + Composition API
- shadcn/ui 컴포넌트 패턴
- Props API 설계
- Storybook 통합

사용 시기: UI 컴포넌트 라이브러리 구축 또는 아키텍처 설계

---

### 디자인 시스템 토큰
포커스: 디자인 토큰 아키텍처 및 도구

핵심 주제:
- W3C DTCG 2025.10 토큰 구조
- Style Dictionary 구성
- 멀티 테마 지원
- Figma MCP 워크플로우
- 시맨틱 네이밍 규칙

사용 시기: 디자인 시스템 기반 설정

---

### 접근성 WCAG
포커스: WCAG 2.2 준수 및 접근성 테스트

핵심 주제:
- 색상 대비 검증 (4.5:1 AA, 7:1 AAA)
- 키보드 네비게이션 패턴
- 스크린 리더 최적화 (ARIA)
- Focus 관리
- 자동화된 테스트 (jest-axe)

사용 시기: 접근성 준수 보장

---

### 아이콘 라이브러리
포커스: 아이콘 라이브러리 선택 및 통합

핵심 주제:
- 10+ 라이브러리 비교 (Lucide, React Icons, Tabler, Iconify)
- 번들 크기 최적화
- Tree-shaking 전략
- 타입 안전 아이콘 컴포넌트
- 성능 패턴

사용 시기: 최적의 번들 크기로 아이콘 시스템 통합

---

### 테마 시스템
포커스: 테마 시스템 구현

핵심 주제:
- CSS variable 아키텍처
- Light/dark 모드 전환
- 시스템 기본 설정 감지
- 브랜드 커스터마이징
- Tailwind CSS 통합

사용 시기: 커스터마이징 가능한 테마 구현

---

### 예제
포커스: 실용적인 코드 예제

핵심 주제:
- Button 컴포넌트 (React, Vue)
- 폼 검증 (Zod + React Hook Form)
- 데이터 테이블 (TanStack Table)
- 모달 다이얼로그 (focus trap)
- 테마 프로바이더
- 아이콘 사용 패턴

사용 시기: 참조 구현

---

### 참고
포커스: 외부 문서 링크

핵심 주제:
- 공식 문서 (DTCG, WCAG, Figma, Storybook)
- 라이브러리 참고 (React, Tailwind, Radix UI)
- 도구 문서 (Style Dictionary, jest-axe)
- 모범 사례 가이드

사용 시기: 공식 리소스 검색

---

## 모범 사례

해야 할 일:
- 시맨틱 디자인 토큰 사용 (`color.primary.500` not `blue-500`)
- Atomic Design 계층 준수 (Atoms -> Molecules -> Organisms)
- 모든 텍스트에 4.5:1 대비 비율 확인 (WCAG AA)
- 모든 인터랙티브 요소에 키보드 네비게이션 구현
- 아이콘 Tree-shaking (특정 항목 import, `import *` 방지)
- 테마 커스터마이징을 위해 CSS variables 사용
- TypeScript 타입으로 모든 props 문서화
- jest-axe로 접근성 컴포넌트 테스트

필수 사례:

[HARD] 모든 색상, 간격, 타이포그래피 값에 디자인 토큰만 사용
WHY: 디자인 토큰은 단일 정보원을 제공하여 일관된 테마, 멀티 플랫폼 지원, 확장 가능한 디자인 시스템 활성화
IMPACT: 하드코딩된 값은 유지보수 부채 생성, 테마 전환 방해, 디자인 시스템 원칙 위반

[HARD] 아이콘만 사용하는 모든 인터랙티브 요소에 ARIA 레이블 포함
WHY: 스크린 리더는 텍스트 대안 없이 시각적 아이콘을 해석할 수 없어 시각 장애 사용자가 접근 불가
IMPACT: 누락된 ARIA 레이블은 WCAG 2.2 AA 준수 위반 및 보조 기술 의존 사용자 제외

[HARD] 네임스페이스 imports 대신 개별 아이콘 import
WHY: 네임스페이스 imports는 전체 라이브러리를 번들링하여 tree-shaking 최적화 무효화
IMPACT: 번들 크기가 아이콘 라이브러리당 500KB-2MB 증가하여 로드 성능 및 사용자 경험 저하

[HARD] light 및 dark 모드에서 모든 컴포넌트 테스트
WHY: 테마 전환은 색상 대비, 가독성, 모든 UI 상태의 접근성 준수에 영향
IMPACT: 테스트되지 않은 다크 모드 구현은 WCAG 대비 요구 사항 실패 가능

[HARD] 모든 인터랙티브 컴포넌트에 키보드 네비게이션 구현
WHY: 키보드 전용 사용자는 Tab, Enter, Escape, Arrow 키 지원 필요
IMPACT: 누락된 키보드 지원은 WCAG 2.2 AA 위반 및 포인팅 장치 사용 불가 사용자 제외

[HARD] 모든 포커스 가능 요소에 시각적 focus 표시기 제공
WHY: Focus 표시기는 현재 키보드 위치를 전달하여 네비게이션 및 접근성 필수
IMPACT: 보이지 않는 focus 상태는 혼란을 일으키고 WCAG 2.2 AA 위반

[SOFT] 인라인 스타일 대신 Tailwind utility 클래스 사용
WHY: Tailwind는 일관된 간격 스케일, 반응형 디자인, 최적의 번들 크기를 위한 자동 정리 제공
IMPACT: 인라인 스타일은 디자인 시스템 제약 우회, 불일치한 간격 생성, CSS 번들 크기 증가

[SOFT] 모든 비동기 작업에 로딩 상태 포함
WHY: 로딩 상태는 데이터 페칭 중 피드백을 제공하여 사용자 불확실성 및 중복 작업 방지
IMPACT: 누락된 로딩 상태는 불명확한 인터페이스 상태로 사용자 경험 저하

---

## 함께 사용하면 좋은 것들

Skills:
- `do-lang-typescript` - TypeScript와 JavaScript 모범 사례
- `do-foundation-core` - TRUST 5 품질 검증
- `do-library-nextra` - 문서 생성
- `do-library-shadcn` - shadcn/ui 특화 패턴 (보완)

Agents:
- `code-frontend` - 프론트엔드 컴포넌트 구현
- `design-uiux` - 디자인 시스템 아키텍처
- `mcp-figma` - Figma 통합 워크플로우
- `core-quality` - 접근성 및 품질 검증

Commands:
- `/do:2-run` - TDD 구현 사이클
- `/do:3-sync` - 문서 생성

---

## 레거시 Skills에서 마이그레이션

이 skill은 4개의 이전 skills를 통합:

do-component-designer -> 컴포넌트 아키텍처 모듈
- Atomic Design 패턴
- React 19 / Vue 3.5 예제
- 컴포넌트 아키텍처

do-design-systems -> 디자인 시스템 토큰 모듈 + 접근성 WCAG 모듈
- DTCG 토큰 아키텍처
- Figma MCP 워크플로우
- WCAG 2.2 준수

do-icons-vector -> 아이콘 라이브러리 모듈
- 아이콘 라이브러리 비교
- 성능 최적화
- 통합 패턴

do-library-shadcn (부분) -> 컴포넌트 아키텍처 모듈 + 테마 시스템 모듈
- shadcn/ui 패턴
- 테마 시스템
- 컴포넌트 구성

참고: `do-library-shadcn`은 shadcn/ui 전용 고급 패턴을 위한 보완 skill로 유지

---

## 공식 리소스

- W3C DTCG: https://designtokens.org
- WCAG 2.2: https://www.w3.org/WAI/WCAG22/quickref/
- React 19: https://react.dev
- Tailwind CSS: https://tailwindcss.com
- Radix UI: https://www.radix-ui.com
- shadcn/ui: https://ui.shadcn.com
- Storybook: https://storybook.js.org
- Figma MCP: https://help.figma.com/hc/en-us/articles/32132100833559
- Style Dictionary: https://styledictionary.com
- Lucide Icons: https://lucide.dev
- Iconify: https://iconify.design

---

Last Updated: 2025-11-26
Status: Production Ready
Version: 2.0.0
