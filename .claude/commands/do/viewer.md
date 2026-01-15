# /do:viewer

Do Memory 웹뷰어 실행. 세션, 관찰, 플랜, 업무일지를 시각적으로 확인.

## 사용법

```
/do:viewer          # 웹뷰어 시작 (localhost:3777)
/do:viewer stop     # 웹뷰어 중지
/do:viewer status   # 상태 확인
```

## allowed-tools
Bash

## Step 1: 상태 확인

현재 뷰어 실행 상태 확인:
```bash
lsof -i :3777 2>/dev/null | grep LISTEN
```

## Step 2: 시작/중지

### 시작 (기본)
```bash
cd .do/viewer && npm run dev &
echo "✓ 웹뷰어 시작: http://localhost:3777"
```

### 중지 (stop 인자)
```bash
pkill -f "vite.*3777"
echo "✓ 웹뷰어 중지됨"
```

### 상태 (status 인자)
```bash
if lsof -i :3777 | grep -q LISTEN; then
  echo "✓ 웹뷰어 실행 중: http://localhost:3777"
else
  echo "✗ 웹뷰어 중지됨"
fi
```

## Step 3: 브라우저 열기 (macOS)

시작 시 자동으로 브라우저 열기:
```bash
open http://localhost:3777
```

## 페이지 구성

| 경로 | 설명 |
|------|------|
| `/` | 대시보드 (통계, 최근 활동) |
| `/sessions` | 세션 목록 및 상세 |
| `/observations` | 관찰 검색 및 필터 |
| `/plans` | 플랜 목록 및 상태 |
| `/reports` | 업무일지 생성/조회 |

## 요구사항
- Go Worker 실행 중 (localhost:3778)
- Node.js 18+
