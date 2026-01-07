---
description: Do 프레임워크 실행 모드 전환 (do/focus/auto)
allowed-tools: Read, Write, Edit, AskUserQuestion
---

# /do:mode - 실행 모드 전환

## 실행

### Step 1: 현재 설정 확인

`.claude/settings.local.json` 파일을 읽어서 `env.DO_MODE` 값 확인.
- 파일이 없거나 `DO_MODE`가 없으면 기본값 "do" 사용

### Step 2: 모드 표시 또는 전환

**인자가 없으면**: 현재 모드와 사용 가능한 모드 표시

**인자가 있으면** ($ARGUMENTS):
1. 인자 검증: "do", "focus", "auto" 중 하나인지 확인
2. 유효하지 않으면 오류 메시지 표시하고 종료
3. 유효하면 Step 3으로 진행

### Step 3: 모드 전환 (인자가 있을 때만)

`.claude/settings.local.json` 파일 업데이트:

```json
{
  "env": {
    "DO_MODE": "{선택한 모드}",
    ... (기존 env 필드 유지)
  }
  ... (기존 다른 필드 유지)
}
```

기존 설정 내용을 모두 유지하면서 `env.DO_MODE`만 업데이트.

### Step 4: 결과 표시

**현재 모드 표시 형식**:

```
현재 모드: {모드명} ({모드 설명})
{모드 특성}

사용 가능한 모드:
┌─ do: 복잡한 멀티 도메인 작업
├─ focus: 간단한 1-3 파일 변경
└─ auto: 자동 선택 (권장)

전환: /do:mode <do|focus|auto>
```

**모드 전환 후 표시 형식**:

```
모드 변경 완료!

{모드명} ({모드 설명})
{모드 특성}
```

## 모드 정보

### do (Strategic Orchestrator)
- 모든 작업을 전문 에이전트에게 위임
- 병렬 실행 우선
- 복잡한 멀티 도메인 작업에 최적
- 높은 품질과 일관성 보장

### focus (Direct Executor)
- 간단한 1-3 파일 변경을 직접 처리
- 에이전트 위임 없이 빠른 실행
- 단순 수정, 작은 버그 수정에 적합
- 빠른 응답 시간

### auto (Automatic Selection)
- 작업 복잡도에 따라 자동으로 모드 선택
- 간단한 작업은 focus, 복잡한 작업은 do
- 가장 효율적인 실행 방식 자동 결정
- 권장 모드
