---
name: do-docs-generation
description: Sphinx, MkDocs, TypeDoc, Nextra 등 실제 도구를 사용한 기술 명세서, API 문서, 사용자 가이드 생성 패턴. 코드에서 문서 생성, 문서 사이트 구축, 문서화 워크플로우 자동화에 활용.
version: 2.0.0
category: workflow
allowed-tools: Read, Write, Edit, Grep, Glob, Bash
updated: 2025-01-06
status: active
---

# 문서 생성 패턴

## 빠른 참조

목적: 검증된 도구와 프레임워크로 전문적인 문서 생성

핵심 도구:
- Python: Sphinx (autodoc), MkDocs (Material 테마), pydoc
- TypeScript/JavaScript: TypeDoc, JSDoc, TSDoc
- API 문서: OpenAPI/Swagger (FastAPI/Express), Redoc, Stoplight
- 정적 사이트: Nextra (Next.js), Docusaurus (React), VitePress (Vue)
- 범용: Markdown, MDX, reStructuredText

스킬 활용 시점:
- 코드 주석에서 API 문서 자동 생성
- 검색 및 네비게이션 기능이 있는 문서 사이트 구축
- 사용자 가이드 및 기술 명세서 작성
- CI/CD 파이프라인에서 문서 자동 업데이트
- 문서 포맷 간 변환

---

## 구현 가이드

### Python 문서화: Sphinx

Sphinx 설정 및 구성:

pip install sphinx sphinx-autodoc-typehints sphinx-rtd-theme myst-parser로 설치

sphinx-quickstart docs 실행으로 프로젝트 초기화

conf.py 핵심 설정:
- extensions에 autodoc, napoleon, typehints, myst_parser 추가
- html_theme을 sphinx_rtd_theme으로 설정 (전문적 외관)
- autodoc_typehints를 description으로 설정 (인라인 타입 힌트)

sphinx-apidoc으로 API 문서 생성 후, docs 디렉토리에서 make html 실행

### Python 문서화: MkDocs

MkDocs Material 설정:

pip install mkdocs mkdocs-material mkdocstrings mkdocstrings-python으로 설치

mkdocs.yml 구성:
- site_name, site_url 설정
- theme에서 name을 material로, 원하는 색상 팔레트 지정
- plugins에 search, mkdocstrings 추가
- nav 구조로 섹션 및 페이지 정의

Markdown 파일에서 ::: module.path 구문으로 docstring에서 API 문서 자동 생성

mkdocs serve (로컬), mkdocs build (빌드), mkdocs gh-deploy (배포)

### TypeScript 문서화: TypeDoc

TypeDoc 설정:

npm install typedoc --save-dev로 설치

package.json scripts에 추가: typedoc --out docs/api src/index.ts

typedoc.json 구성:
- entryPoints: 소스 파일 지정
- out: docs/api 출력 경로
- includeVersion, categorizeByGroup 활성화
- theme: 기본 또는 커스텀 테마

npm run docs:generate로 문서 생성

### JavaScript 문서화: JSDoc

JSDoc 설정:

npm install jsdoc --save-dev로 설치

jsdoc.json 구성:
- source include 경로 및 패턴
- templates 및 출력 경로
- markdown 플러그인 활성화

함수 문서화 태그:
- @param: 매개변수 타입 및 설명
- @returns: 반환 값 문서화
- @example: 사용 예시
- @throws: 에러 문서화

### OpenAPI/Swagger 문서화

FastAPI 자동 문서:

FastAPI는 자동 OpenAPI 문서 제공. /docs에서 Swagger UI, /redoc에서 ReDoc 접근

문서 개선 방법:
- 라우트 핸들러에 docstring 추가
- response_model로 타입 응답 정의
- Pydantic 모델 Config 클래스에 examples 정의
- tags로 엔드포인트 그룹화
- 라우트 데코레이터에 상세 설명 추가

app.openapi()로 OpenAPI 스펙 프로그래밍 방식 추출 후 openapi.json 저장

Express + Swagger:

swagger-jsdoc, swagger-ui-express 설치

OpenAPI 정의 및 API 파일 경로로 swagger-jsdoc 구성

라우트 핸들러에 @openapi 주석으로 경로, 매개변수, 응답 문서화

/api-docs 엔드포인트에서 Swagger UI 제공

### 정적 문서 사이트

Nextra (Next.js):

Skill("do-library-nextra") 참조 - Nextra 패턴 상세 정보

핵심 장점: MDX 지원, 파일 시스템 라우팅, 내장 검색, 테마 커스터마이징

npx create-nextra-app으로 생성, theme.config.tsx 구성, pages 디렉토리에 페이지 정리

Docusaurus (React):

npx create-docusaurus@latest my-docs classic으로 초기화

docusaurus.config.js 구성:
- siteMetadata에 title, tagline, url 설정
- presets에 docs, blog 설정
- themeConfig에 navbar, footer 추가
- algolia 플러그인으로 검색 활성화

docs 폴더에 문서 정리, category.json으로 사이드바 구조 설정

VitePress (Vue):

npm init vitepress로 초기화

.vitepress/config.js 구성:
- title, description, base path 설정
- themeConfig에 nav, sidebar 정의
- 검색 및 소셜 링크 설정

Vue 컴포넌트, 코드 하이라이팅, frontmatter와 함께 Markdown 사용

---

## 고급 패턴

### SPEC 파일에서 문서 생성

Do SPEC 파일에서 문서 생성 패턴:

SPEC 파일 내용 읽기 및 핵심 섹션 추출: id, title, description, requirements, api_endpoints

구조화된 Markdown 문서 생성:
- description에서 개요 섹션 생성
- requirements를 기능 항목으로 나열
- 각 API 엔드포인트에 method, path, description 문서화
- 엔드포인트 정의 기반 사용 예시 추가

docs 디렉토리 적절한 위치에 생성 문서 저장

### CI/CD 문서화 파이프라인

GitHub Actions 워크플로우:

main 브랜치 푸시 시 src 또는 docs 경로 변경에 트리거되는 .github/workflows/docs.yml 생성

워크플로우 단계:
- 리포지토리 체크아웃
- 언어 런타임 설정 (Python, Node.js)
- 문서화 의존성 설치
- 적절한 도구로 문서 생성
- GitHub Pages, Netlify, Vercel에 배포

Python/Sphinx 예시:
- pip install sphinx sphinx-rtd-theme로 설치
- sphinx-build -b html docs/source docs/build로 생성
- actions-gh-pages 액션으로 배포

TypeScript/TypeDoc 예시:
- npm ci로 설치
- npm run docs:generate로 생성
- Pages에 배포

### 문서 유효성 검사

링크 검사:

linkchecker로 HTML 출력 로컬 링크 유효성 검사

Markdown은 pre-commit 훅에서 markdown-link-check 사용

맞춤법 검사:

pyspelling + Aspell로 자동 맞춤법 검사

.pyspelling.yml에 다양한 파일 타입 매트릭스 항목 구성

문서 커버리지:

Python은 interrogate로 docstring 커버리지 확인

pyproject.toml에 최소 커버리지 임계값 설정

임계값 미달 시 CI 빌드 실패 처리

### 다국어 문서

Nextra 국제화:

next.config.js에서 locales 배열과 defaultLocale로 i18n 구성

pages/[locale] 디렉토리에 로케일별 페이지 생성

next-intl 또는 유사 도구로 번역 처리

Docusaurus 국제화:

docusaurus.config.js에서 defaultLocale, locales로 i18n 구성

docusaurus write-translations로 번역 파일 생성

i18n/[locale] 디렉토리 구조로 번역 정리

---

## 연관 기술

스킬:
- do-library-nextra: Nextra 문서 프레임워크 패턴
- do-lang-python: Python docstring 규칙 및 타이핑
- do-lang-typescript: TypeScript/JSDoc 문서화 패턴
- do-domain-backend: 백엔드 서비스 API 문서
- do-workflow-project: 프로젝트 문서 통합

에이전트:
- manager-docs: 문서 워크플로우 오케스트레이션
- expert-backend: API 엔드포인트 문서
- expert-frontend: 컴포넌트 문서

커맨드:
- /do:3-sync: 코드 변경과 문서 동기화

---

## 도구 참조

Python 문서화:
- Sphinx: https://www.sphinx-doc.org/
- MkDocs: https://www.mkdocs.org/
- MkDocs Material: https://squidfunk.github.io/mkdocs-material/
- mkdocstrings: https://mkdocstrings.github.io/

JavaScript/TypeScript 문서화:
- TypeDoc: https://typedoc.org/
- JSDoc: https://jsdoc.app/
- TSDoc: https://tsdoc.org/

API 문서화:
- OpenAPI Specification: https://spec.openapis.org/
- Swagger UI: https://swagger.io/tools/swagger-ui/
- Redoc: https://redocly.github.io/redoc/
- Stoplight: https://stoplight.io/

정적 사이트 생성기:
- Nextra: https://nextra.site/
- Docusaurus: https://docusaurus.io/
- VitePress: https://vitepress.dev/

스타일 가이드:
- Google Developer Documentation Style Guide: https://developers.google.com/style
- Microsoft Writing Style Guide: https://learn.microsoft.com/style-guide/

---

Version: 2.0.0
Last Updated: 2025-01-06
