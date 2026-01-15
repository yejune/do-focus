---
description: 업무일지 작성 - git log + DB observations 통합
allowed-tools: Bash, Read, Write, Glob, WebFetch
argument-hint: [날짜 YYYY-MM-DD 또는 빈칸=오늘] [--team]
---

# /do:report - 업무일지 작성

## 실행 단계

### Step 1: 날짜 결정

- `$ARGUMENTS`가 있으면 해당 날짜 사용 (YYYY-MM-DD 형식)
- 없으면 오늘 날짜 사용
- `--team` 플래그가 있으면 팀 활동도 포함

### Step 2: Git 로그 수집

```bash
# 해당 날짜의 커밋 조회
git log --since="{날짜} 00:00" --until="{날짜} 23:59" --format="%h %s" --reverse

# 각 커밋의 상세 diff
git show {commit_hash} --stat
```

### Step 3: DB에서 observations 조회

Worker API로 오늘의 작업 기록 조회:

```bash
curl -s "http://127.0.0.1:3778/api/observations?user=$DO_USER_NAME&date={날짜}"
```

응답 예시:
```json
{
  "observations": [
    {"type": "feature", "content": "로그인 API 구현", "created_at": 1705312800},
    {"type": "decision", "content": "JWT 대신 세션 기반 인증 선택", "reason": "보안 강화"},
    {"type": "delegation", "content": "API 구조 분석", "agent_name": "expert-backend"},
    {"type": "blocker", "content": "외부 API 응답 지연", "resolved": true}
  ]
}
```

**API 호출 실패 시**: git log만으로 업무일지 생성 (graceful fallback)

### Step 4: 팀 활동 조회 (--team 플래그 시)

```bash
curl -s "http://127.0.0.1:3778/api/observations/team?date={날짜}"
```

응답 예시:
```json
{
  "team_activities": [
    {"user": "kim", "summary": "프론트엔드 컴포넌트 작업", "commits": 5},
    {"user": "lee", "summary": "DB 스키마 수정", "commits": 3}
  ]
}
```

### Step 5: 업무일지 작성

파일 경로: `.do/reports/{YYYY}/{MM}/{DD}.md`

```markdown
# 업무일지 - {YYYY}-{MM}-{DD}

## 요약
- 총 {N}개 커밋
- 주요 작업: {핵심 내용 요약}

## 작업 내역

### 커밋
| 해시 | 메시지 | 변경 |
|------|--------|------|
| abc1234 | feat: Add user authentication | +150 -20 |
| def5678 | fix: Resolve login bug | +10 -5 |

### 의사결정
- JWT 대신 세션 기반 인증 선택 (보안 강화)
- API 응답 형식 통일 (일관성)

### 구현
- 로그인 API 구현
- 사용자 검증 로직 추가
- 에러 핸들링 개선

### 에이전트 활용
| 에이전트 | 작업 내용 |
|----------|----------|
| expert-backend | API 구조 분석 |
| expert-security | 보안 취약점 검토 |
| expert-frontend | UI 컴포넌트 검토 |

### 이슈/블로커
- [해결] 외부 API 응답 지연 -> 캐싱 적용
- [미해결] 테스트 커버리지 80% 미달

## 팀 활동
| 멤버 | 작업 내용 | 커밋 |
|------|----------|------|
| @kim | 프론트엔드 컴포넌트 작업 | 5 |
| @lee | DB 스키마 수정 | 3 |

## 내일 할 일
- {있으면 기록, 없으면 생략}

---
작성자: @{DO_USER_NAME}
생성일시: {YYYY-MM-DD HH:mm}
```

### Step 6: 디렉토리 생성 및 저장

```bash
mkdir -p .do/reports/{YYYY}/{MM}
```

파일 저장 후 경로 출력.

---

## 예시

```
/do:report
→ .do/reports/2026/01/07.md 생성 (개인 활동만)

/do:report 2026-01-06
→ .do/reports/2026/01/06.md 생성 (지정 날짜)

/do:report --team
→ .do/reports/2026/01/07.md 생성 (팀 활동 포함)

/do:report 2026-01-06 --team
→ .do/reports/2026/01/06.md 생성 (지정 날짜 + 팀 활동)
```

---

## Observation 타입

| 타입 | 설명 | 예시 |
|------|------|------|
| `feature` | 구현한 기능 | 로그인 API 구현 |
| `decision` | 기술적 결정 | JWT 대신 세션 선택 |
| `delegation` | 에이전트 활용 | expert-backend에 분석 위임 |
| `blocker` | 이슈/장애물 | 외부 API 응답 지연 |
| `learning` | 새로 배운 것 | React 19 use() hook |
| `todo` | 다음 할 일 | 테스트 커버리지 개선 |

---

## 주의사항

- git 로그가 상세해야 좋은 업무일지 생성 가능
- 커밋이 없는 날짜는 "작업 내역 없음"으로 기록
- Worker API (127.0.0.1:3778) 실행 중이어야 DB 조회 가능
- API 실패 시 git log만으로 fallback 생성
- `$DO_USER_NAME` 환경변수 필요 (`.claude/settings.local.json`에서 설정)
