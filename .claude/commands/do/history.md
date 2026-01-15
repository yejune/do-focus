# /do:history

작업 히스토리 검색. 키워드로 과거 작업 내용 조회.

## 사용법

```
/do:history 인증           # "인증" 관련 모든 작업
/do:history bugfix        # 버그 수정 타입만
/do:history @kim          # 팀원 kim의 작업
/do:history --team        # 팀 전체 히스토리
```

## allowed-tools
Bash, Read

## Step 1: 검색 쿼리 파싱

- 키워드: 일반 텍스트 → FTS5 검색
- 타입 필터: decision, bugfix, feature, refactor, delegation
- 팀 필터: @username 또는 --team

## Step 2: Worker API 호출

```bash
# 개인 히스토리
curl "http://127.0.0.1:3778/api/observations/search?q={query}&user=$DO_USER_NAME"

# 팀 히스토리
curl "http://127.0.0.1:3778/api/team/observations?project_path=$(pwd)"
```

## Step 3: 결과 포맷팅

```markdown
## 검색 결과: "{query}"

### 2024-01-15
- [decision] JWT 대신 세션 기반 인증 선택 (@max)
- [feature] 로그인 API 구현 (@max)

### 2024-01-14
- [bugfix] 동시 요청 race condition 수정 (@kim)
- [delegation → expert-backend] API 구조 분석 (@max)

총 4건
```

## 옵션

| 옵션 | 설명 |
|------|------|
| `--team` | 팀 전체 히스토리 |
| `--days N` | 최근 N일 (기본: 7) |
| `--type TYPE` | 특정 타입만 |
