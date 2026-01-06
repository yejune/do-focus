---
name: do-foundation-context
aliases: [do-foundation-context]
description: Enterprise context와 session 관리, token budget 최적화, state 지속성
version: 3.0.0
modularized: false
category: foundation
allowed-tools: Read, Grep, Glob
replaces: do-core-context-budget, do-core-session-state
---

## Quick Reference (30초)

# Enterprise Context & Session 관리

Claude Code를 위한 unified context 최적화 및 session state 관리. 200K token budget 관리, session 지속성, multi-agent handoff protocol 지원.

Core Capabilities:
- 200K token budget 할당 및 모니터링
- Session state tracking 및 지속성
- Context-aware token 최적화
- Multi-agent handoff protocol
- Progressive disclosure 및 memory 관리
- Session forking (parallel exploration)

사용 시기:
- Session 초기화 및 정리
- Long-running workflow (10분 초과)
- Multi-agent orchestration
- Context window 제한 접근 (150K tokens 초과)
- Model 전환 (Haiku ↔ Sonnet)
- Workflow phase 전환

Key Principles (2025):
1. 마지막 20% 회피 - context 마지막 구간에서 성능 저하 발생
2. Aggressive Clearing - SPEC workflow에서 1-3회 메시지마다 /clear 실행
3. Lean Memory Files - 각 파일 500줄 미만 유지
4. Unused MCP 비활성화 - tool definition overhead 최소화
5. Quality 우선 - 10% 관련 context가 90% noise보다 효과적

---

## Implementation Guide (5분)

### Features

- Claude Code session을 위한 intelligent context window 관리
- Priority-based caching과 progressive file loading
- Token budget tracking 및 optimization alert
- /clear 경계를 넘어선 selective context 보존
- MCP integration context 지속성

### 사용 시기

- 150K token limit 초과하는 large codebase 관리
- Long-running development session에서 token 사용 최적화
- Session reset 후에도 critical context 보존
- Shared context를 사용하는 multi-agent workflow 조율
- Claude Code에서 context 관련 issue 디버깅

### Core Patterns

Pattern 1: Progressive File Loading

우선순위 tier로 파일 로드:
- Tier 1: CLAUDE.md, config.json (항상 로드)
- Tier 2: 현재 SPEC 및 implementation files
- Tier 3: 관련 module 및 dependencies
- Tier 4: Reference documentation (on-demand)

Pattern 2: Context Checkpointing

1. Token 사용 모니터링: 150K에서 경고, 180K에서 critical
2. 보존해야 할 essential context 식별
3. /clear 실행하여 session reset
4. Tier 1 및 Tier 2 파일 자동 reload
5. 보존된 context로 작업 재개

Pattern 3: MCP Context Continuity

/clear 후에도 MCP agent context 보존:
- agent_id를 통해 context 식별
- /clear 이후 fresh MCP agent 초기화로 context 복원

---

## 5 Core Patterns (각 5-10분)

### Pattern 1: Token Budget 관리 (200K Context)

Concept: 200K token context window의 전략적 할당 및 모니터링.

Allocation Strategy:

200K Token Budget 분배:
- System Prompt 및 Instructions: 약 15K tokens (7.5%)
  - CLAUDE.md: 약 8K
  - Command definitions: 약 4K
  - Skill metadata: 약 3K
- Active Conversation: 약 80K tokens (40%)
  - Recent messages: 약 50K
  - Context cache: 약 20K
  - Active references: 약 10K
- Reference Context (Progressive Disclosure): 약 50K (25%)
  - Project structure: 약 15K
  - Related Skills: 약 20K
  - Tool definitions: 약 15K
- Reserve (Emergency Recovery): 약 55K tokens (27.5%)
  - Session state snapshot: 약 10K
  - TAGs 및 cross-references: 약 15K
  - Error recovery context: 약 20K
  - Free buffer: 약 10K

Monitoring Thresholds:

Token budget 실시간 모니터링:
- usage_percent가 85 초과: Critical - emergency compression 트리거 후 clear 실행
- usage_percent가 75 초과: Warning - non-critical context 지연 로드, 사용자에게 limit 접근 경고
- usage_percent가 60 초과: Monitor - context growth pattern 추적

Use Case: Long-running SPEC-First workflow에서 context overflow 방지.

---

### Pattern 2: Aggressive /clear Strategy

Concept: 효율성 유지를 위해 전략적 checkpoint에서 proactive context clearing.

Clear 실행 규칙:

MANDATORY /clear Points:
- /do:1-plan 완료 후 (45-50K tokens 절약)
- Context 150K tokens 초과 (overflow 방지)
- Conversation 50개 메시지 초과 (stale history 제거)
- Major phase transition 전 (clean slate)
- Model 전환 (Haiku ↔ Sonnet handoff)

Implementation:

/clear 실행 여부 결정:
- spec_created가 true인 경우: clear 필요
- token_usage가 150000 초과인 경우: clear 필요
- message_count가 50 초과인 경우: clear 필요
- phase_changed가 true인 경우: clear 필요
- 위 조건 중 하나라도 true이면 /clear 실행

Use Case: SPEC-Run-Sync cycle 전반에서 token 효율성 극대화.

---

### Pattern 3: Session State Persistence

Concept: State snapshot으로 중단 이후에도 session 연속성 유지.

Session State Architecture:

Session State Layers:
- L1: Context-Aware Layer (Claude 4.5+ feature)
  - Token budget tracking
  - Context window position
  - Auto-summarization trigger
  - Model-specific optimization
- L2: Active Context (current task, variables, scope)
- L3: Session History (recent actions, decisions)
- L4: Project State (SPEC progress, milestones)
- L5: User Context (preferences, language, expertise)
- L6: System State (tools, permissions, environment)

State Snapshot Structure:

session_state = {
  session_id: "sess_uuid_v4",
  model: "claude-sonnet-4-5-20250929",
  created_at: "2025-11-24T10:30:00Z",
  context_window: {
    total: 200000,
    used: 85000,
    available: 115000,
    position_percent: 42.5
  },
  persistence: {
    auto_load_history: true,
    context_preservation: "critical_only",
    cache_enabled: true
  },
  work_state: {
    current_spec: "SPEC-001",
    phase: "implementation",
    completed_steps: ["spec_complete", "architecture_defined"]
  }
}

Use Case: 중단 후에도 context 손실 없이 long-running task 재개.

---

### Pattern 4: Multi-Agent Handoff Protocol

Concept: Minimal token overhead로 agent 간 seamless context transfer.

Handoff Package:

handoff_package = {
  handoff_id: "uuid-v4",
  from_agent: "spec-builder",
  to_agent: "tdd-implementer",
  session_context: {
    session_id: "sess_uuid",
    model: "claude-sonnet-4-5-20250929",
    context_position: 42.5,
    available_tokens: 115000,
    user_language: "ko"
  },
  task_context: {
    spec_id: "SPEC-001",
    current_phase: "implementation",
    completed_steps: ["spec_complete", "architecture_defined"],
    next_step: "write_tests"
  },
  recovery_info: {
    last_checkpoint: "2025-11-24T10:25:00Z",
    recovery_tokens_reserved: 55000,
    session_fork_available: true
  }
}

Handoff Validation:

Handoff package integrity 검증:
- available_tokens가 30000 미만이면: context compression 트리거
- from_agent와 to_agent 호환성 검증
- 호환되지 않으면: AgentCompatibilityError 발생

Use Case: Plan → Run → Sync workflow의 효율적 실행.

---

### Pattern 5: Progressive Disclosure & Memory 최적화

Concept: 관련성과 필요에 따라 context를 점진적으로 로드.

Progressive Summarization:

Context 압축 단계:
1. Key sentence 추출 (50K → 15K)
2. 원본 content에 대한 pointer 추가
3. 복구를 위해 원본 context를 session_archive에 저장
결과: 약 35K tokens 절약

Context Tagging:

Bad (high token cost):
"이전 20개 메시지의 user configuration..."

Good (efficient reference):
"User preferences는 @CONFIG-001 참조"

Use Case: Token overhead 최소화하면서 context 연속성 유지.

---

## Best Practices

### DO

- SPEC 생성 직후 /clear 실행
- Token 사용 모니터링 및 계획 수립
- Context-aware token budget tracking 사용
- Major operation 전 checkpoint 생성
- Long workflow에 progressive summarization 적용
- Recovery를 위한 session persistence 활성화
- Parallel exploration을 위한 session forking 사용
- Memory files 500줄 미만 유지
- Unused MCP server 비활성화로 overhead 감소

### REQUIREMENTS

[HARD] Bounded context history 유지하며 regular clearing cycle 실행
WHY: 무제한 context 축적은 성능 저하와 token cost 기하급수 증가 야기
IMPACT: Context overflow 방지, 일관된 response quality 유지, 60-70% token 낭비 감소

[HARD] Token budget warning 발생 시 즉시 대응 (150K tokens 초과 시)
WHY: Context window 마지막 20%에서 운영 시 상당한 성능 저하 발생
IMPACT: 최적 model 성능 보장, context overflow 실패 방지, workflow 연속성 유지

[HARD] Session recovery 중 state validation check 실행
WHY: Invalid state는 multi-step process에서 workflow 실패 및 data loss 야기
IMPACT: Session integrity 보장, silent failure 방지, 95% 초과 복구 성공률

[HARD] Context clearing 전 session identifier persist
WHY: Session ID가 중단된 workflow 재개를 위한 유일한 신뢰 가능 메커니즘
IMPACT: Seamless workflow 재개, 작업 손실 방지, multi-agent 조율 지원

[SOFT] 여러 concurrent session 작업 시 명확한 session boundary 설정
WHY: Session mixing은 context 오염 및 예측 불가능한 agent 동작 야기
IMPACT: 디버깅 명확성 향상, cross-session 간섭 방지, clean audit trail 유지

[SOFT] Session 연속성 가정 전 checkpoint snapshot 생성
WHY: Network issue, timeout, system event로 context 손실 가능
IMPACT: Recovery point 제공, rework 시간 감소, system 신뢰성에 대한 user 신뢰 유지

[SOFT] Priority tier 사용하여 codebase component progressive load
WHY: 전체 codebase loading은 token budget 소진 및 irrelevant context 포함
IMPACT: 40-50% token 사용 최적화, response 관련성 향상, 더 큰 project 지원

[SOFT] Handoff package를 critical context만으로 제한
WHY: Non-critical context는 handoff overhead 증가 및 available working token 감소
IMPACT: Agent transition 30% 가속, 실제 작업을 위한 token budget 보존, transfer error 감소

[HARD] 사용량 85% threshold 도달 시 context compression 또는 clearing 실행
WHY: Context limit 접근 시 emergency behavior 트리거 및 model capability 감소
IMPACT: 55K token emergency reserve 유지, 강제 중단 방지, graceful degradation 보장

---

## Works Well With

- `do-cc-memory` - Memory management 및 context persistence
- `do-cc-configuration` - Session configuration 및 preferences
- `do-core-workflow` - Workflow state persistence 및 recovery
- `do-cc-agents` - Session 간 agent state management
- `do-foundation-trust` - Quality gate integration

---

## Workflow Integration

Session Initialization:

1. Token budget 초기화 (Pattern 1)
2. Session state 로드 (Pattern 3)
3. Progressive disclosure 설정 (Pattern 5)
4. Handoff protocol 구성 (Pattern 4)

SPEC-First Workflow:

1. /do:1-plan 실행
   ↓
2. /clear (mandatory - 45-50K tokens 절약)
   ↓
3. /do:2-run SPEC-XXX
   ↓
4. Multi-agent handoff (Pattern 4)
   ↓
5. /do:3-sync SPEC-XXX
   ↓
6. Session state persistence (Pattern 3)

Context Monitoring:

Continuous:
- Token 사용 추적 (Pattern 1)
- Progressive disclosure 적용 (Pattern 5)
- Threshold에서 /clear 실행 (Pattern 2)
- Handoff 검증 (Pattern 4)

---

## Success Metrics

- Token Efficiency: Aggressive clearing으로 60-70% 감소
- Context Overhead: System/skill metadata에 15K tokens 미만
- Handoff Success Rate: Validation으로 95% 초과
- Session Recovery: State persistence로 5초 미만
- Memory Optimization: Memory file당 500줄 미만

---

## Changelog

- v3.0.0 (2025-12-06): Do project용 한국어 버전 생성
- v2.0.0 (2025-11-24): do-core-context-budget 및 do-core-session-state를 5 core pattern으로 통합
- v1.0.0 (2025-11-22): 원본 개별 skill

---

Status: Production Ready (Enterprise)
Modular Architecture: SKILL.md + reference + examples
Integration: Plan-Run-Sync workflow 최적화
Generated with: Do-ADK Skill Factory
