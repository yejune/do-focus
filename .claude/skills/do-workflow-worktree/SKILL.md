---
name: do-worktree
description: 병렬 SPEC 개발을 위한 Git worktree 관리 시스템으로 고립된 워크스페이스, 자동 등록, 그리고 seamless Do 통합을 제공
version: 1.2.0
category: workflow
updated: 2026-01-06
status: active
tags:
  - git
  - worktree
  - parallel
  - development
  - spec
  - isolation
allowed-tools: Read, Write, Grep, Glob
user-invocable: true
---

# Do Worktree 관리

병렬 SPEC 개발을 위한 Git worktree 관리 시스템. 각 SPEC은 context switching 오버헤드 없이 진정한 병렬 개발을 가능하게 하는 고립된 워크스페이스를 제공받음.

핵심 철학: 각 SPEC은 진정한 병렬 개발을 위해 자신만의 고립된 워크스페이스를 가져야 함.

---

## 30초 개요

Do Worktree 관리란: 각 SPEC에 대해 고립된 개발 환경을 생성하는 Git worktree 시스템으로, 충돌 없이 병렬 개발을 가능하게 함.

핵심 기능:
- 고립된 워크스페이스: 각 SPEC은 독립적인 Git state를 가진 자신만의 worktree 보유
- 자동 등록: Worktree registry가 모든 활성 워크스페이스 추적
- 병렬 개발: 여러 SPEC을 동시에 개발 가능
- Seamless 통합: /do:1-plan, /do:2-run, /do:3-sync workflow와 연동
- 스마트 동기화: 필요시 base branch와 자동 동기화
- 자동 정리: merged worktree의 자동 정리

사용 사례:
- 복수 SPEC의 병렬 개발
- 고립된 테스트 환경
- 기능 브랜치 격리
- 코드 리뷰 워크플로우
- 실험적 기능 개발

---

## 핵심 아키텍처 (5분)

### 시스템 구성요소

Worktree 시스템의 5개 핵심 컴포넌트:

1. Worktree Registry - 모든 worktree metadata를 JSON으로 추적하는 중앙 데이터베이스
2. Manager Layer - create, switch, remove, sync 핵심 작업 처리
3. CLI Interface - 사용자 친화적인 명령 인터페이스
4. Models - Worktree metadata를 위한 데이터 구조
5. Integration Layer - Do workflow 통합

### Registry 구조

Registry 파일은 worktree metadata를 JSON 형식으로 저장. 각 worktree 항목에는 다음 정보 포함:
- identifier: SPEC 식별자
- path: 파일 경로
- branch: 브랜치 이름
- created_at: 생성 타임스탬프
- last_sync: 마지막 동기화 시간
- status: 상태 (active 또는 merged)
- base_branch: 베이스 브랜치 참조

Config 섹션에는 worktree root 디렉토리, auto-sync 설정, merged 브랜치 정리 동작이 정의됨.

### 파일 시스템 구조

Worktree 시스템은 전용 디렉토리 구조를 생성:
- Worktree root (~/worktrees/ProjectName/)에 중앙 registry JSON 파일 위치
- 각 SPEC에 대한 개별 디렉토리 존재
- 각 SPEC 디렉토리에는 worktree metadata용 .git 파일과 모든 프로젝트 파일의 완전한 복사본 포함

---

## CLI 명령어 레퍼런스

### 생성 명령어

**do-worktree new** - 새 worktree 생성

SPEC 개발을 위한 새로운 고립된 Git worktree 생성.

인자:
- spec-id: SPEC 식별자 (예: SPEC-001, SPEC-AUTH-001)
- description: 선택적 설명

옵션:
- --branch: 자동 생성 대신 특정 브랜치 이름 사용
- --base: 새 worktree의 베이스 브랜치 (기본값: main)
- --template: 사전 정의된 템플릿 사용
- --shallow: 빠른 설정을 위한 shallow clone 생성
- --depth: shallow clone의 깊이
- --force: worktree가 존재해도 강제 생성

자동 생성 브랜치 패턴:
- 기본값: feature/SPEC-ID-description-kebab-case
- 예시: SPEC-001 -> feature/SPEC-001-user-authentication

### 탐색 명령어

**do-worktree list** - worktree 목록

등록된 모든 worktree를 상태 및 metadata와 함께 표시.

옵션:
- --format: 출력 형식 (table, json, csv)
- --status: 상태로 필터 (active, merged, stale)
- --sort: 필드로 정렬 (name, created, modified, status)
- --verbose: 상세 정보 표시

**do-worktree switch** - worktree로 전환

지정된 worktree로 현재 작업 디렉토리 변경.

옵션:
- --auto-sync: 전환 전 자동 동기화
- --force: 커밋되지 않은 변경사항이 있어도 강제 전환

**do-worktree go** - worktree 경로 가져오기

쉘 통합을 위한 cd 명령어 출력.

사용법:
- 표준 사용 (eval): eval $(do-worktree go SPEC-001)
- 절대 경로 출력: do-worktree go SPEC-001 --absolute
- 상대 경로: do-worktree go SPEC-001 --relative

### 관리 명령어

**do-worktree sync** - worktree 동기화

Worktree를 베이스 브랜치와 동기화.

옵션:
- --auto-resolve: 간단한 충돌 자동 해결
- --interactive: 대화형 충돌 해결
- --dry-run: 실제 동기화 없이 미리보기
- --force: 커밋되지 않은 변경사항이 있어도 강제 동기화
- --include: 특정 파일만 포함
- --exclude: 특정 파일 제외

**do-worktree remove** - worktree 제거

Worktree를 제거하고 등록 정리.

옵션:
- --force: 확인 없이 강제 제거
- --keep-branch: worktree 제거 후 브랜치 유지
- --backup: 제거 전 백업 생성
- --dry-run: 실제 제거 없이 미리보기

**do-worktree clean** - worktree 정리

Merged 브랜치 또는 stale worktree 제거.

옵션:
- --merged-only: merged 브랜치 worktree만 제거
- --stale: 지정된 일수 동안 업데이트되지 않은 worktree 제거
- --days: stale 임계값 일수 (기본값: 30)
- --interactive: 제거할 worktree 대화형 선택
- --dry-run: 실제 정리 없이 미리보기
- --force: 확인 프롬프트 건너뛰기

### 상태 및 설정

**do-worktree status** - 상태 표시

Worktree에 대한 상세 상태 정보 표시.

옵션:
- --all: 모든 worktree 상태 표시
- --sync-check: 동기화 필요 여부 확인
- --detailed: 상세 Git 상태 표시

**do-worktree config** - 설정 관리

설정 액션:
- get: 설정 값 가져오기
- set: 설정 값 설정
- list: 모든 설정 나열
- reset: 기본값으로 재설정

설정 키:
- worktree_root: worktree 루트 디렉토리
- auto_sync: 자동 동기화 활성화
- cleanup_merged: merged worktree 자동 정리
- default_base: 기본 베이스 브랜치
- sync_strategy: 동기화 전략 (merge, rebase, squash)

---

## 병렬 개발 워크플로우

### 고립 모델

각 worktree는 완전한 고립 레이어 제공:

1. Git State 고립: 각 worktree가 독립적인 브랜치 상태, 커밋, 히스토리 보유
2. 파일 시스템 고립: 독립적인 수정이 가능한 완전한 프로젝트 복사본
3. 의존성 고립: 별도의 node_modules, .venv, 빌드 아티팩트
4. 설정 고립: Worktree별 .env, IDE 설정, 도구 구성
5. 프로세스 고립: 독립적인 개발 서버, 테스트 러너, 빌드 프로세스

### 패턴 1: 독립 SPEC 개발

관련 없는 여러 기능을 동시에 개발하는 경우:

1. 환경 설정: 각 기능 영역에 대해 worktree 생성
2. 환경 초기화: 각 worktree에서 필요한 의존성 설치 및 서버 시작
3. 병렬 개발: 여러 터미널에서 각 SPEC 동시 작업
4. 통합: 모든 worktree 동기화 및 병렬 sync 실행

### 패턴 2: 순차 기능 개발

의존성이 있는 기능을 순차적으로 개발하되 준비 단계가 겹치는 경우:

1. 기반 작업 시작 (SPEC-FOUND-001)
2. 기반 개발 중 종속 기능 준비 (SPEC-API-001)
3. 기반 완료 후 develop에 병합
4. 종속 기능에서 기반 변경사항 동기화 후 개발 진행
5. 다음 레이어 준비 (예: UI)

### 패턴 3: 실험-프로덕션 병렬

안정적인 프로덕션 작업과 함께 실험적 기능 개발:

- 프로덕션 worktree: main 기반, 안정적인 버그 수정
- 실험 worktree: develop 기반, 프로덕션에 영향 없는 실험적 개발
- 리뷰용 빠른 컨텍스트 전환을 위한 별칭 설정

---

## Do Workflow 통합

### 계획 단계 (/do:1-plan)

SPEC 생성 후 자동 worktree 설정:
- SPEC 생성 완료 시 자동으로 worktree new 호출
- 자동 브랜칭 규칙에 따른 브랜치 생성
- 선택적으로 spec-development 템플릿 적용

출력 안내:
1. 생성된 worktree로 전환: do-worktree switch SPEC-ID
2. 쉘 통합: eval $(do-worktree go SPEC-ID)
3. 개발 시작: /do:2-run SPEC-ID

### 개발 단계 (/do:2-run)

Worktree 인식 TDD 구현:
- 시스템이 worktree 환경을 자동 감지
- Worktree 경로에서 현재 디렉토리 확인
- worktrees 디렉토리 내에서 SPEC 부분 식별
- TDD 실행 후 worktree metadata 업데이트

### 동기화 단계 (/do:3-sync)

자동화된 worktree 동기화:
- Worktree 존재 여부 확인
- 베이스 브랜치와 worktree 동기화
- 동기화 성공 시 일반 문서 동기화 진행
- Worktree에서 문서 업데이트 추출
- 필요시 PR 생성

### 정리 통합

성공적인 통합 후 worktree 정리:
- SPEC이 성공적으로 병합되면 정리 옵션 제공
- 권장: worktree 제거
- 대안: 참조용 유지 또는 아카이브
- Registry에 완료 상태 기록

---

## 문제 해결

### 생성 실패

Worktree 이미 존재 오류:
- 원인: 이전 생성 중단, 디렉토리는 있지만 registry 항목 누락, Git metadata 손상
- 해결: 디렉토리 확인, .git 파일 검사, git worktree prune 실행, 재생성

브랜치 이미 체크아웃 오류:
- 원인: 동일 브랜치를 여러 worktree에서 사용 시도
- 해결: worktree 목록 확인, git worktree prune, 다른 브랜치 이름 지정

### Registry 문제

Registry 파일 손상:
- 원인: 동시 쓰기, 디스크 공간 부족, 수동 편집 오류
- 해결: 백업 생성, JSON 구문 오류 수정, 기존 디렉토리에서 registry 재구성

고아 Registry 항목:
- 원인: 디렉토리 수동 삭제, 파일 시스템 오류
- 해결: status --all로 확인, registry prune으로 정리

### Git 상태 문제

분리된 HEAD 상태:
- 원인: 특정 커밋 체크아웃, 중단된 rebase/merge
- 해결: 커밋되지 않은 변경사항 저장, 새 브랜치 생성, 의도된 브랜치 체크아웃

동기화 중 병합 충돌:
- 원인: 베이스 브랜치 변경과 충돌, 장기 실행 worktree 발산
- 해결: 충돌 파일 식별, 해결 전략 선택 (worktree 버전/베이스 버전/수동 병합), git add로 스테이지, 커밋

### 통합 문제

Do 명령어 조정 실패:
- 원인: SPEC ID 형식 불일치, 수동 생성 worktree, 설정 불일치
- 해결: SPEC ID 형식 확인, 올바른 등록 확인, 설정 점검

자동 감지 불가:
- 원인: 예상 위치에 registry 파일 부재, 디렉토리 명명 규칙 불일치
- 해결: parent 디렉토리에 registry 존재 확인, 명명 패턴 확인, 쉘 프로필에 통합 함수 로드

### 진단 명령어

상태 확인:
- do-worktree status --all: 모든 worktree 동기화 상태
- git worktree list: 기본 Git worktree 목록
- do-worktree status SPEC-ID --detailed: 특정 worktree 상세 상태

정리 및 복구:
- do-worktree clean --dry-run: 변경 없이 정리 미리보기
- git worktree prune: stale Git worktree metadata 제거
- do-worktree registry rebuild: 디렉토리에서 registry 재구성

---

## 모범 사례

### 정기 유지보수

주간:
- do-worktree status --all로 문제 확인
- 활성 worktree를 베이스 브랜치와 동기화하여 충돌 최소화
- Merged worktree 정리

월간:
- Stale Git worktree metadata 정리
- Registry 무결성 검증
- 오래된 worktree 검토 및 아카이브

### 안전한 작업 패턴

생성:
- 수동 Git 명령어 대신 항상 do-worktree new 사용
- 생성 전 브랜치 이름 고유성 확인
- 쉬운 식별을 위한 설명적 SPEC ID 사용

수정:
- Worktree 전환 전 변경사항 커밋
- 중요 변경 전 sync 명령어 사용
- 충돌 발생 시 즉시 해결

제거:
- 수동 삭제 대신 do-worktree remove 사용
- 필요시 keep-branch 옵션으로 브랜치 보존
- 커밋되지 않은 작업이 있는 worktree는 백업 생성

---

## 함께 사용하기

Commands:
- do:1-plan - 자동 worktree 설정 포함 SPEC 생성
- do:2-run - 고립된 worktree 환경에서 개발
- do:3-sync - 자동 worktree 동기화와 통합

Skills:
- do-foundation-core - 병렬 개발 패턴
- do-workflow-project - 프로젝트 관리 통합
- do-workflow-spec - SPEC 기반 개발

Tools:
- Git worktree - 기본 Git worktree 기능
- Rich CLI - 형식화된 터미널 출력

---

Version: 1.2.0
Last Updated: 2026-01-06
Status: Active (통합 한국어 문서, 핵심 기능 완전 포함)
