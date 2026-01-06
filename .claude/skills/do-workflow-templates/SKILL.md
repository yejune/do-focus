---
name: do-workflow-templates
description: 엔터프라이즈 템플릿 관리 (코드 보일러플레이트, 피드백 템플릿, 프로젝트 최적화 워크플로우 포함)
version: 3.0.0
category: workflow
modularized: true
replaces: do-core-code-templates, do-core-feedback-templates, do-project-template-optimizer
allowed-tools: Read, Write, Edit, Grep, Glob
---

# 엔터프라이즈 템플릿 관리

코드 보일러플레이트, 피드백 템플릿, 프로젝트 최적화 워크플로우를 통합한 템플릿 시스템

## 빠른 참조

핵심 기능:
- 코드 템플릿 라이브러리 (FastAPI, React, Vue, Next.js)
- GitHub 이슈 피드백 템플릿 (6가지 유형)
- 프로젝트 템플릿 최적화 및 스마트 병합
- 템플릿 버전 관리 및 이력
- 백업 발견 및 복원

사용 시기:
- 새 프로젝트 또는 기능 스캐폴딩
- `/do:9-feedback`으로 GitHub 이슈 생성
- Do 업데이트 후 템플릿 구조 최적화
- 프로젝트 백업에서 복원
- 템플릿 버전 및 커스터마이징 관리

모듈 참조:
- 코드 템플릿 - modules/code-templates.md
- 피드백 템플릿 - modules/feedback-templates.md
- 템플릿 옵티마이저 - modules/template-optimizer.md

---

## 구현 가이드

### 기능

- 공통 아키텍처용 프로젝트 템플릿
- 모범 사례를 적용한 보일러플레이트 코드 생성
- 구성 가능한 템플릿 변수 및 커스터마이징
- 다중 프레임워크 지원 (React, FastAPI, Spring 등)
- 통합 테스트 및 CI/CD 구성

### 언제 사용하나

- 검증된 아키텍처 패턴으로 새 프로젝트 부트스트랩
- 조직 내 여러 프로젝트에서 일관성 보장
- 적절한 구조로 새 기능 빠르게 프로토타이핑
- 표준화된 프로젝트 레이아웃으로 신규 개발자 온보딩

### 핵심 패턴

패턴 1 - 템플릿 구조:

템플릿 디렉토리 구조는 프레임워크별로 구성:
- fastapi-backend: template.json (변수 정의), src, main.py, models, tests
- nextjs-frontend: template.json, app, components
- fullstack: backend, frontend 하위 디렉토리

패턴 2 - 템플릿 변수:

변수 정의:
- PROJECT_NAME: 프로젝트 식별자 (필수)
- AUTHOR: 작성자 이름 (기본값: Do Team)
- LICENSE: 라이선스 유형 (MIT, Apache-2.0)
- PYTHON_VERSION: 백엔드용 (3.11, 3.12, 3.13)
- DATABASE: 데이터베이스 유형 (postgresql, mysql, mongodb)
- AUTH_TYPE: 인증 유형 (jwt, oauth, session)

파일 처리 액션:
- substitute: 변수 치환 수행
- copy: 그대로 복사
- ignore: 무시

패턴 3 - 템플릿 생성:

프로젝트 생성 단계:
1. 템플릿 디렉토리 로드
2. 마킹된 파일에서 변수 치환
3. 정적 파일 그대로 복사
4. 생성 후 훅 실행 (의존성 설치, git 초기화)
5. 생성된 프로젝트 구조 검증

---

## 5가지 핵심 패턴

### 패턴 1: 코드 템플릿 스캐폴딩

개념: 프로덕션 준비 완료된 보일러플레이트로 빠르게 프로젝트 스캐폴드

사용 예제 (슈도 코드):

```
// FastAPI 프로젝트 구조 생성
template = loadTemplate("backend/fastapi")
project = template.scaffold({
    name: "my-api",
    features: ["auth", "database", "celery"],
    customizations: { db: "postgresql" }
})
```

FastAPI 템플릿 생성 결과:
- src: main.py, config.py, database.py
- src/auth: router.py, schemas.py, models.py, jwt.py
- src/models: user.py
- src/tasks: celery_app.py
- tests: conftest.py, test_auth.py, test_main.py
- alembic: versions, env.py
- docker: Dockerfile, docker-compose.yml
- 루트: pyproject.toml, .env.example, README.md

상세: modules/code-templates.md 참조

---

### 패턴 2: GitHub 피드백 템플릿

개념: 일관된 GitHub 이슈 생성을 위한 구조화된 템플릿

6가지 템플릿 유형:
- Bug Report: 버그 설명, 재현 단계, 예상/실제 동작, 환경 정보
- Feature Request: 기능 설명, 사용 시나리오, 예상 효과, 우선순위
- Improvement: 현재 상태, 개선된 상태, 성능/품질 영향
- Refactor: 리팩토링 범위, 현재/개선 구조, 영향 분석
- Documentation: 문서 내용, 대상 독자, 문서 구조
- Question: 배경, 질문 또는 제안, 옵션, 결정 기준

통합: `/do:9-feedback` 명령으로 자동 트리거

피드백 처리 흐름:
1. 사용자가 `/do:9-feedback "설명"` 실행
2. 스킬이 적절한 템플릿 유형 자동 선택
3. 사용자 입력으로 템플릿 채우기
4. GitHub 이슈 자동 생성

상세: modules/feedback-templates.md 참조

---

### 패턴 3: 템플릿 최적화 및 스마트 병합

개념: 사용자 커스터마이징을 보존하면서 템플릿 업데이트 지능적으로 병합

스마트 병합 알고리즘 (슈도 코드):

```
function smartMerge(backup, template, current) {
    // 백업에서 사용자 커스터마이징 추출
    userContent = extractUserCustomizations(backup)

    // 최신 템플릿 기본값 가져오기
    templateDefaults = getCurrentTemplates()

    // 우선순위로 병합
    merged = {
        template_structure: templateDefaults,  // 항상 최신
        user_config: userContent,              // 보존
        custom_content: userContent            // 추출
    }
    return merged
}
```

6단계 최적화 워크플로우:
- 1단계: 백업 발견 및 분석 (.do-backups 디렉토리 스캔)
- 2단계: 템플릿 비교 (해시 기반 파일 비교, 커스터마이징 감지)
- 3단계: 스마트 병합 (사용자 콘텐츠 추출, 템플릿 업데이트 적용)
- 4단계: 템플릿 기본값 감지 (플레이스홀더 패턴 식별)
- 5단계: 버전 관리 (HISTORY 섹션 업데이트)
- 6단계: 구성 업데이트 (최적화 플래그 설정)

상세: modules/template-optimizer.md 참조

---

### 패턴 4: 백업 발견 및 복원

개념: 지능형 복원을 갖춘 자동 백업 관리

백업 메타데이터 구조:
- backup_id: 백업 식별자 (예: backup-2025-11-24-v0.28.2)
- created_at: 생성 시각 (ISO 8601 형식)
- template_version: 템플릿 버전
- project_state: 프로젝트 상태 (이름, SPEC 목록, 백업된 파일 수)
- customizations: 커스터마이징 (언어, 팀 설정, 도메인)

복원 프로세스 (슈도 코드):

```
function restoreFromBackup(backupId) {
    backup = loadBackup(backupId)

    if (!validateBackupIntegrity(backup)) {
        throw new BackupIntegrityError("백업 손상됨")
    }

    customizations = extractCustomizations(backup)
    applyCustomizations(customizations)

    updateConfig({
        restored_from: backupId,
        restored_at: now()
    })
}
```

상세: modules/template-optimizer.md 참조

---

### 패턴 5: 템플릿 버전 관리

개념: 템플릿 버전 추적 및 업데이트 이력 유지

template_optimization 구성:
- last_optimized: 마지막 최적화 시각
- backup_version: 백업 참조 버전
- template_version: 현재 템플릿 버전
- customizations_preserved: 보존된 커스터마이징 목록

version_history:
- 버전별 변경 내역 배열

migration_log:
- 이전 버전에서 새 버전으로의 마이그레이션 기록 (날짜, 상태, 충돌 수)

상세: modules/template-optimizer.md 참조

---

## 모듈 참조

핵심 모듈:
- modules/code-templates.md - 보일러플레이트 라이브러리, 스캐폴드 패턴
- modules/feedback-templates.md - 6가지 GitHub 이슈 유형
- modules/template-optimizer.md - 스마트 병합, 백업 복원, 버전 관리

코드 템플릿:
- FastAPI REST API 템플릿
- React 컴포넌트 템플릿
- Docker 및 CI/CD 템플릿

피드백 템플릿:
- Bug Report, Feature Request
- Improvement, Refactor, Documentation, Question

템플릿 옵티마이저:
- 6단계 최적화 워크플로우
- 스마트 병합, 백업 복원, 버전 추적

---

## 모범 사례

### 품질 표준

[HARD] 변경 적용 전 모든 템플릿 기본값 수정을 문서화
왜: 템플릿 기본값은 모든 프로젝트의 기준선이며, 문서화되지 않은 변경은 팀 간 혼란 초래

[HARD] 템플릿 최적화 워크플로우 실행 전 백업 생성
왜: 템플릿 최적화는 되돌리기 어려운 구조적 변경 포함

[HARD] 템플릿 업데이트 워크플로우 중 모든 병합 충돌 해결
왜: 미해결 충돌은 적절한 템플릿 기능을 방해

[HARD] 모든 템플릿 업데이트에서 완전한 커스터마이징 이력 보존
왜: 커스터마이징 이력은 감사 추적 제공 및 롤백 가능

[HARD] 프로덕션 배포 전 테스트를 통해 템플릿 기능 검증
왜: 테스트되지 않은 템플릿은 프로덕션에서만 나타나는 오류 포함 가능

[HARD] 내장 버전 관리 시스템을 사용하여 템플릿 버전 추적
왜: 버전 추적은 템플릿 진화 이해 및 호환성 확인 가능

[SOFT] 프로젝트 전체에서 일관된 템플릿 패턴 사용 유지

[SOFT] 유지보수를 위해 합리적인 복잡도 제한 내에서 템플릿 설계

---

## 함께 잘 작동하는 것들

에이전트:
- workflow-project: 프로젝트 초기화
- core-planner: 템플릿 계획
- workflow-spec: SPEC 템플릿 생성

스킬:
- do-project-config-manager: 구성 관리 및 검증
- do-foundation-specs: SPEC 템플릿 생성
- do-docs-generation: 문서 템플릿 스캐폴딩

명령어:
- /do:0-project: 템플릿으로 프로젝트 초기화
- /do:9-feedback: 피드백 템플릿 선택 및 이슈 생성

---

## 워크플로우 통합

프로젝트 초기화:
1. 코드 템플릿 선택 (패턴 1)
2. 프로젝트 구조 스캐폴드
3. 커스터마이징 적용
4. 버전 추적 초기화 (패턴 5)

피드백 제출:
1. /do:9-feedback 실행
2. 이슈 유형 선택 (패턴 2)
3. 템플릿 필드 채우기
4. GitHub 이슈 자동 생성

템플릿 업데이트:
1. 템플릿 버전 변경 감지
2. 백업 생성 (패턴 4)
3. 스마트 병합 실행 (패턴 3)
4. 버전 이력 업데이트 (패턴 5)

---

## 성공 지표

- 스캐폴드 시간: 새 프로젝트 2분 (수동 30분 대비)
- 템플릿 채택률: 95% 프로젝트가 템플릿 사용
- 커스터마이징 보존: 업데이트 중 100% 사용자 콘텐츠 유지
- 피드백 완전성: 95% GitHub 이슈가 완전한 정보 포함
- 병합 성공률: 99% 충돌 자동 해결

---

## 파일 위치

- 템플릿 저장: .claude/skills/do-workflow-templates/templates/
- 백업: .do/backups/templates/
- 병합 충돌: .do/merge-conflicts/
- 버전 이력: .do/config/template-versions.json

---

## 변경 이력

- v3.0.0 (2025-11-24): do-core-code-templates, do-core-feedback-templates, do-project-template-optimizer를 5가지 핵심 패턴의 단일 스킬로 통합
- v2.0.0 (2025-11-22): 원본 개별 스킬

---

상태: 프로덕션 준비 완료 (엔터프라이즈)
모듈 아키텍처: SKILL.md + 3 핵심 모듈
통합: Plan-Run-Sync 워크플로우 최적화
생성자: Do Skill Factory
