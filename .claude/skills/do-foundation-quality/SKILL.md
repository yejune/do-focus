---
name: do-foundation-quality
aliases:
  - do-foundation-quality
category: foundation
description: TRUST 5 검증, 사전 분석, 자동화된 모범 사례 적용을 포함한 엔터프라이즈 코드 품질 오케스트레이터
version: 2.0.0
modularized: true
updated: 2025-01-06
status: active
tags:
  - foundation
  - quality
  - testing
  - validation
  - trust-5
  - best-practices
  - code-review
user-invocable: false
---

# 엔터프라이즈 코드 품질 오케스트레이터

체계적인 코드 검토, 사전 개선 제안, 자동화된 모범 사례 적용을 결합한 엔터프라이즈급 코드 품질 관리 시스템. TRUST 5 프레임워크 검증과 Context7 통합을 통해 실시간 모범 사례로 포괄적인 품질 보증을 제공.

## 빠른 참조

핵심 기능:
- TRUST 5 검증: Testable, Readable, Unified, Secured, Trackable 품질 게이트
- 사전 분석: 자동화된 문제 감지 및 개선 제안
- 모범 사례: Context7 기반 실시간 표준 검증
- 다중 언어: 25개 이상 프로그래밍 언어 지원
- 엔터프라이즈: CI/CD 파이프라인, 품질 메트릭, 보고

핵심 패턴:
1. 품질 게이트 파이프라인 - 설정 가능한 임계값을 통한 자동화된 검증
2. 사전 스캐너 - 개선 권장 사항을 포함한 지속적 분석
3. 모범 사례 엔진 - Context7 기반 표준 적용
4. 품질 메트릭 대시보드 - 종합 보고 및 추세 분석

사용 시기:
- 코드 검토 자동화 및 품질 게이트 적용
- 사전 코드 품질 개선 및 기술 부채 감소
- 엔터프라이즈 코딩 표준 적용 및 준수 검증
- CI/CD 파이프라인 통합 및 자동 품질 검사

빠른 접근:
- TRUST 5 프레임워크 - modules/trust5-validation.md
- 사전 분석 - modules/proactive-analysis.md
- 모범 사례 - modules/best-practices.md
- 통합 패턴 - modules/integration-patterns.md

## 구현 가이드

### 시작하기

기본 품질 검증:

```javascript
// 품질 오케스트레이터 초기화
const qualityOrchestrator = new QualityOrchestrator({
    trust5Enabled: true,
    proactiveAnalysis: true,
    bestPracticesEnforcement: true,
    context7Integration: true
});

// 종합 품질 분석 실행
const result = await qualityOrchestrator.analyzeCodebase({
    path: "src/",
    languages: ["php", "javascript", "typescript"],
    qualityThreshold: 0.85
});

// TRUST 5 품질 게이트 검증
const qualityGate = new QualityGate();
const validationResult = await qualityGate.validateTrust5({
    codebasePath: "src/",
    testCoverageThreshold: 0.90,
    complexityThreshold: 10
});
```

사전 품질 분석:

```javascript
// 사전 스캐너 초기화
const proactiveScanner = new ProactiveQualityScanner({
    context7Client: context7Client,
    ruleEngine: new BestPracticesEngine()
});

// 개선 기회 스캔
const improvements = await proactiveScanner.scanCodebase({
    path: "src/",
    scanTypes: ["security", "performance", "maintainability", "testing"]
});

// 개선 권장 사항 생성
const recommendations = await proactiveScanner.generateRecommendations({
    issues: improvements,
    priority: "high",
    autoFix: true
});
```

### 핵심 컴포넌트

#### 1. 품질 오케스트레이션 엔진

```javascript
class QualityOrchestrator {
    // TRUST 5 프레임워크를 통한 엔터프라이즈 품질 오케스트레이션

    constructor(config) {
        this.trust5Validator = new TRUST5Validator();
        this.proactiveScanner = new ProactiveScanner();
        this.bestPracticesEngine = new BestPracticesEngine();
        this.context7Client = new Context7Client();
        this.metricsCollector = new QualityMetricsCollector();
    }

    async analyzeCodebase(request) {
        // 종합 코드베이스 품질 분석

        // 1단계: TRUST 5 검증
        const trust5Result = await this.trust5Validator.validate({
            codebase: request.path,
            thresholds: request.qualityThresholds
        });

        // 2단계: 사전 분석
        const proactiveResult = await this.proactiveScanner.scan({
            codebase: request.path,
            focusAreas: request.focusAreas
        });

        // 3단계: 모범 사례 확인
        const practicesResult = await this.bestPracticesEngine.validate({
            codebase: request.path,
            languages: request.languages,
            context7Docs: true
        });

        // 4단계: 메트릭 수집
        const metrics = await this.metricsCollector.collectComprehensiveMetrics({
            codebase: request.path,
            analysisResults: [trust5Result, proactiveResult, practicesResult]
        });

        return {
            trust5Validation: trust5Result,
            proactiveAnalysis: proactiveResult,
            bestPractices: practicesResult,
            metrics: metrics,
            overallScore: this.calculateOverallScore([
                trust5Result, proactiveResult, practicesResult
            ])
        };
    }
}
```

상세 구현:
- TRUST 5 검증기 구현 - modules/trust5-validation.md 참조
- 사전 스캐너 구현 - modules/proactive-analysis.md 참조
- 모범 사례 엔진 구현 - modules/best-practices.md 참조

### 구성 및 커스터마이징

품질 구성:

```yaml
quality_orchestration:
  trust5_framework:
    enabled: true
    thresholds:
      overall: 0.85
      testable: 0.90
      readable: 0.80
      unified: 0.85
      secured: 0.90
      trackable: 0.80

  proactive_analysis:
    enabled: true
    scan_frequency: "daily"
    focus_areas:
      - "performance"
      - "security"
      - "maintainability"
      - "technical_debt"

    auto_fix:
      enabled: true
      severity_threshold: "medium"
      confirmation_required: true

  best_practices:
    enabled: true
    context7_integration: true
    auto_update_standards: true
    compliance_target: 0.85

    language_rules:
      php:
        style_guide: "psr-12"
        formatter: "php-cs-fixer"
        linter: "phpstan"

      javascript:
        style_guide: "airbnb"
        formatter: "prettier"
        linter: "eslint"

      typescript:
        style_guide: "google"
        formatter: "prettier"
        linter: "eslint"

  reporting:
    enabled: true
    metrics_retention_days: 90
    trend_analysis: true
    executive_dashboard: true

    notifications:
      quality_degradation: true
      security_vulnerabilities: true
      technical_debt_increase: true
```

통합 예시:
- CI/CD 파이프라인 통합 - modules/integration-patterns.md 참조
- GitHub Actions 통합
- 품질-서비스 REST API
- 교차 프로젝트 벤치마킹

## 고급 패턴

### 1. 사용자 정의 품질 규칙

```javascript
class CustomQualityRule {
    // 사용자 정의 품질 검증 규칙 정의

    constructor(name, validator, severity = "medium") {
        this.name = name;
        this.validator = validator;
        this.severity = severity;
    }

    async validate(codebase) {
        // 사용자 정의 규칙 검증 실행
        try {
            const result = await this.validator(codebase);
            return {
                ruleName: this.name,
                passed: result.passed,
                severity: this.severity,
                details: result.details,
                recommendations: result.recommendations
            };
        } catch (error) {
            return {
                ruleName: this.name,
                passed: false,
                severity: "error",
                details: { error: error.message },
                recommendations: ["규칙 구현 수정 필요"]
            };
        }
    }
}
```

상세 내용 - modules/best-practices.md 참조

### 2. ML 기반 품질 예측

코드 특성 추출 및 예측 모델을 활용한 ML 기반 품질 문제 예측.

상세 구현 - modules/proactive-analysis.md 참조

### 3. 실시간 품질 모니터링

품질 저하 및 보안 취약점에 대한 자동 알림을 포함한 지속적 품질 모니터링.

상세 구현 - modules/proactive-analysis.md 참조

### 4. 교차 프로젝트 품질 벤치마킹

업계 유사 프로젝트 대비 프로젝트 품질 메트릭 비교.

상세 구현 - modules/integration-patterns.md 참조

## 모듈 참조

### 핵심 모듈

- modules/trust5-validation.md - 종합 품질 프레임워크 검증
- modules/proactive-analysis.md - 자동화된 문제 감지 및 개선
- modules/best-practices.md - Context7 기반 표준 적용
- modules/integration-patterns.md - CI/CD 및 엔터프라이즈 통합

### 모듈별 핵심 컴포넌트

TRUST 5 검증:
- `TRUST5Validator` - 5계층 품질 검증
- `TestableValidator` - 테스트 커버리지 및 품질
- `SecuredValidator` - 보안 및 OWASP 준수
- 품질 게이트 파이프라인 통합

사전 분석:
- `ProactiveQualityScanner` - 자동화된 문제 감지
- `QualityPredictionEngine` - ML 기반 예측
- `RealTimeQualityMonitor` - 지속적 모니터링
- 성능 및 유지보수성 분석

모범 사례:
- `BestPracticesEngine` - 표준 검증
- Context7 통합으로 최신 문서 적용
- 사용자 정의 품질 규칙
- 언어별 검증기

통합 패턴:
- CI/CD 파이프라인 통합
- GitHub Actions 워크플로우
- 품질-서비스 REST API
- 교차 프로젝트 벤치마킹

## Context7 라이브러리 매핑

품질 분석 도구 및 프레임워크를 위한 필수 라이브러리 매핑.

상세 목록 - modules/best-practices.md 참조

## 연관 도구

에이전트:
- core-planner - 품질 요구사항 계획
- workflow-tdd - TDD 구현 검증
- security-expert - 보안 취약점 분석
- code-backend - 백엔드 코드 품질
- code-frontend - 프론트엔드 코드 품질

스킬:
- do-foundation-core - TRUST 5 프레임워크 참조
- do-workflow-testing - TDD 워크플로우 검증
- do-domain-backend - 백엔드 품질 기준
- do-domain-frontend - 프론트엔드 품질 기준

명령어:
- `/do:2-run` - TDD 검증 통합
- `/do:3-sync` - 문서 품질 검사
- `/do:9-feedback` - 품질 개선 피드백

## 빠른 참조 요약

핵심 기능: TRUST 5 검증, 사전 스캔, Context7 기반 모범 사례, 다중 언어 지원, 엔터프라이즈 통합

핵심 클래스: `QualityOrchestrator`, `TRUST5Validator`, `ProactiveQualityScanner`, `BestPracticesEngine`, `QualityMetricsCollector`

필수 메서드: `analyzeCodebase()`, `validateTrust5()`, `scanForIssues()`, `validateBestPractices()`, `generateQualityReport()`

통합 준비 완료: CI/CD 파이프라인, GitHub Actions, REST API, 실시간 모니터링, 교차 프로젝트 벤치마킹

엔터프라이즈 기능: 사용자 정의 규칙, ML 예측, 실시간 모니터링, 벤치마킹, 종합 보고

품질 표준: OWASP 준수, TRUST 5 프레임워크, Context7 통합, 자동화된 개선 권장 사항
