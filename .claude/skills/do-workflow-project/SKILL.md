---
name: do-workflow-project
description: 문서 생성, 언어 초기화, 템플릿 최적화를 통합한 프로젝트 관리 시스템
version: 2.0.0
modularized: true
updated: 2026-01-06
status: active
category: workflow
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
user-invocable: true
---

# Do Workflow Project - 프로젝트 관리 시스템

목적: 문서 생성, 다국어 지원, 템플릿 최적화를 통합한 종합 프로젝트 관리 시스템

범위: 프로젝트 초기화부터 유지보수까지 전체 라이프사이클 지원

대상: 프로젝트 설정, 문서 생성, 다국어 지원, 성능 최적화를 수행하는 Claude Code 에이전트

---

## 빠른 참조

핵심 역량:
- 문서 관리: 다국어 지원 템플릿 기반 문서 생성
- 언어 초기화: 언어 감지, 설정, 지역화 관리
- 템플릿 최적화: 고급 템플릿 분석 및 성능 최적화
- 통합 인터페이스: 모든 기능을 통합하는 단일 진입점

주요 기능:
- 프로젝트 타입 자동 감지 및 템플릿 선택
- 다국어 문서 생성 (영어, 한국어, 일본어, 중국어)
- 성능 벤치마킹을 통한 지능형 템플릿 최적화
- SPEC 기반 문서 업데이트
- 다중 형식 내보내기 (Markdown, HTML, PDF)

지원 프로젝트 타입:
- 웹 애플리케이션, 모바일 애플리케이션, CLI 도구, 라이브러리, ML 프로젝트

---

## 구현 가이드

### 모듈 아키텍처

문서 관리 모듈:
- 템플릿 기반 문서 생성, 프로젝트 타입 감지, SPEC 데이터 통합, 다중 형식 내보내기

언어 초기화 모듈:
- 자동 언어 감지, 언어 설정 관리, 에이전트 프롬프트 지역화, 토큰 비용 최적화

템플릿 최적화 모듈:
- 복잡도 메트릭 분석, 크기 축소, 백업/복구 시스템, 벤치마킹

### 핵심 워크플로우

완전한 프로젝트 초기화:

1단계: 프로젝트 디렉토리 경로를 지정하여 시스템 초기화

2단계: 설정 매개변수 지정
- 언어 설정: 기본 언어 코드 (예: "en", "ko")
- 사용자 이름: 개인화를 위한 개발자 또는 팀 이름
- 도메인: 프로젝트 도메인 목록 (예: backend, frontend)
- 프로젝트 타입: web_application, mobile_app 등
- 최적화 활성화: 초기화 중 템플릿 최적화 여부

3단계: 결과 검토 - 언어 설정, 문서 구조, 템플릿 분석 보고서

SPEC 기반 문서 생성:

1단계: SPEC 데이터 준비 (id, title, description, requirements, status, priority, apiEndpoints)

2단계: generateDocumentationFromSpec 호출로 문서 생성

3단계: 기능 문서, API 문서, 다국어 버전 검토

템플릿 성능 최적화:

1단계: analyzeProjectTemplates로 현재 템플릿 분석

2단계: 최적화 옵션 설정 (backupFirst, applySizeOptimizations, preserveFunctionality)

3단계: optimizeProjectTemplates 실행 및 결과 검토

### 언어 및 지역화

자동 언어 감지 프로세스:
- 파일 콘텐츠 분석 (주석, 문자열)
- 설정 파일 검사 (로케일 설정)
- 시스템 로케일 감지
- 디렉토리 구조 패턴

다국어 문서 구조:
- 언어별 디렉토리 (docs/ko, docs/en)
- 언어 협상 설정, 자동 리디렉션 설정

에이전트 프롬프트 지역화:
- 언어별 지침, 문화적 맥락 적응, 토큰 비용 최적화 권장

---

## 사용 예제

### 예제 1: 완전한 프로젝트 초기화

시나리오: 새로운 웹 애플리케이션을 다국어 지원으로 초기화

슈도 코드:

project = new DoMenuProject("./my-ecommerce-app")

result = project.initializeCompleteProject(
    language: "ko",
    userName: "김개발",
    domains: ["backend", "frontend", "mobile"],
    projectType: "web_application",
    optimizationEnabled: true
)

결과 구조:
- initialization_status: 초기화 성공 여부
- project_metadata: 이름, 타입, 초기화 시간
- language_config: 대화 언어, 에이전트 프롬프트 언어, 토큰 비용 분석
- documentation_structure: 생성된 디렉토리 및 파일 목록
- template_analysis: 분석된 템플릿 수, 최적화 기회

### 예제 2: SPEC 기반 문서 자동 생성

시나리오: 사용자 인증 SPEC에서 API 문서와 기능 문서를 자동 생성

슈도 코드:

specData = {
    id: "SPEC-001",
    title: "User Authentication System",
    description: "Implement secure authentication with JWT",
    requirements: [
        "User registration with email verification",
        "JWT token generation and validation",
        "Password reset functionality"
    ],
    status: "Planned",
    priority: "High",
    apiEndpoints: [
        {
            path: "/api/auth/register",
            method: "POST",
            description: "User registration endpoint",
            requestBody: { email: "string", password: "string" },
            response: { userId: "string", verificationSent: "boolean" }
        },
        {
            path: "/api/auth/login",
            method: "POST",
            description: "User login endpoint"
        }
    ]
}

docsResult = project.generateDocumentationFromSpec(specData)

### 예제 3: 템플릿 성능 최적화

시나리오: 기존 프로젝트의 템플릿을 분석하고 최적화

슈도 코드:

analysis = project.templateOptimizer.analyzeProjectTemplates()

optimizationOptions = {
    backupFirst: true,
    applySizeOptimizations: true,
    applyPerformanceOptimizations: true,
    preserveFunctionality: true
}

optimizationResult = project.optimizeProjectTemplates(optimizationOptions)

결과: 크기 축소 백분율, 성능 개선 메트릭, 백업 경로, 권장사항

---

## 공통 패턴

### 패턴 1: 언어 자동 감지

슈도 코드:

language = project.languageInitializer.detectProjectLanguage()

결과: detected_language, confidence, indicators, recommendation

### 패턴 2: 다국어 문서 구조 생성

슈도 코드:

multilingual = project.languageInitializer.createMultilingualDocumentationStructure("ko")

생성 구조: docs/ko/ (기본), docs/en/ (폴백), _meta.json (언어 협상)

### 패턴 3: 에이전트 프롬프트 지역화

슈도 코드:

localized = project.languageInitializer.localizeAgentPrompts(
    basePrompt: "Generate user authentication system with JWT",
    language: "ko"
)

비용 최적화 전략:
- full_english: 비용 0%, 사용자 경험 낮음
- full_localized: 비용 +20%, 사용자 경험 높음
- hybrid: 비용 +10%, 사용자 경험 중간

---

## 안티 패턴 (피해야 할 패턴)

### 안티 패턴 1: 수동 문서 구조 생성

문제: 각 언어별 문서 폴더를 수동으로 생성하고 관리
해결: DoMenuProject로 자동 생성하여 일관된 구조와 언어 협상 설정 보장

### 안티 패턴 2: 최적화 없는 템플릿 사용

문제: 기본 템플릿을 그대로 사용하여 성능 저하
해결: 템플릿 분석 및 최적화 적용 후 사용

### 안티 패턴 3: SPEC 데이터 불완전

문제: 필수 필드 없이 SPEC에서 문서 생성 시도 (불완전한 문서, 빈 섹션)
해결: id, title, description, requirements, status, priority, api_endpoints 완전히 제공

### 안티 패턴 4: 백업 없는 최적화

문제: 백업 없이 템플릿 최적화 실행 (복구 불가)
해결: backupFirst: true, preserveFunctionality: true 설정

---

## 고급 패턴

### 커스텀 템플릿 개발

특정 프로젝트 타입에 대한 커스텀 템플릿 생성:
- project_type: 대상 프로젝트 타입
- language: 기본 언어 코드
- sections: 문서 섹션 정의 (mission, metrics, frameworks)

### 템플릿 캐싱

1단계: 최적화 결과를 위한 캐시 스토리지 초기화
2단계: 캐시 키로 중복 최적화 방지
3단계: 모든 템플릿 작업에 캐시된 최적화 함수 사용

### 배치 처리

1단계: 처리할 템플릿 파일 경로 목록 수집
2단계: 각 템플릿 최적화 시도 (성공/실패 기록)
3단계: 완전한 결과 컬렉션 반환

---

## 설정 참조

프로젝트 설정:
- project.name: 표시 이름
- project.type: web_application, mobile_app, cli_tool, library, ml_project
- language.conversation_language: en, ko, ja, zh, es, fr, de
- language.agent_prompt_language: english (비용 최적화) 또는 localized

최적화 옵션:
- backup_first: 백업 생성 여부
- apply_size_optimizations: 크기 축소
- apply_performance_optimizations: 성능 개선
- preserve_functionality: 기능 유지

언어별 토큰 영향:
- 영어: 0% (기준), 한국어: +20%, 일본어: +25%, 중국어: +15%, 유럽어: +5%

---

## 연계 기술

- do-foundation-core: 핵심 실행 패턴 및 SPEC 기반 개발 워크플로우
- do-foundation-claude: Claude Code 통합 및 설정
- do-workflow-docs: 통합 문서 관리
- do-workflow-templates: 템플릿 최적화 전략
- do-library-nextra: 고급 문서 아키텍처

---

Version: 2.0.0
Last Updated: 2026-01-06
Integration Status: Complete
