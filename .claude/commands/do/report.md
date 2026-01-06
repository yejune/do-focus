---
description: 업무일지 작성 - git log/diff 기반
allowed-tools: Bash, Read, Write, Glob
argument-hint: [날짜 YYYY-MM-DD 또는 빈칸=오늘]
---

# /do:report - 업무일지 작성

## 실행 단계

### Step 1: 날짜 결정

- `$ARGUMENTS`가 있으면 해당 날짜 사용 (YYYY-MM-DD 형식)
- 없으면 오늘 날짜 사용

### Step 2: Git 로그 수집

```bash
# 해당 날짜의 커밋 조회
git log --since="{날짜} 00:00" --until="{날짜} 23:59" --format="%h %s" --reverse

# 각 커밋의 상세 diff
git show {commit_hash} --stat
```

### Step 3: 업무일지 작성

파일 경로: `./daily_report/{YYYY}/{MM}/{DD}.md`

```markdown
# {YYYY}-{MM}-{DD} 업무일지

## 요약
- 총 {N}개 커밋
- 주요 작업: {핵심 내용 요약}

## 상세 작업 내역

### 1. {커밋 메시지 제목}
- 커밋: `{hash}`
- 변경 파일:
  - {파일1} (+{추가} -{삭제})
  - {파일2} (+{추가} -{삭제})
- 내용: {커밋 메시지 본문 또는 diff 요약}

### 2. {다음 커밋}
...

## 기술적 결정사항
- {있으면 기록}

## 내일 할 일
- {있으면 기록, 없으면 생략}
```

### Step 4: 디렉토리 생성 및 저장

```bash
mkdir -p ./daily_report/{YYYY}/{MM}
```

파일 저장 후 경로 출력.

---

## 예시

```
/do:report
→ ./daily_report/2026/01/07.md 생성

/do:report 2026-01-06
→ ./daily_report/2026/01/06.md 생성
```

---

## 주의사항

- git 로그가 상세해야 좋은 업무일지 생성 가능
- 커밋이 없는 날짜는 "작업 내역 없음" 으로 기록
