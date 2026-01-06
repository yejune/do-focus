# do-foundation-context Reference

## API Reference

### Token Budget Monitoring API

Functions:
- monitor_token_budget(context_usage): 실시간 사용량 모니터링
- get_usage_percent(): 현재 사용량 백분율 반환 (0-100)
- trigger_emergency_compression(): Critical 시 context 압축
- defer_non_critical_context(): non-essential context를 cache로 이동

Thresholds:
- 60%: 모니터링 및 growth pattern 추적
- 75%: Warning - progressive disclosure 시작
- 85%: Critical - emergency compression 트리거

### Session State API

Session State Structure:
- session_id: Unique identifier (UUID v4)
- model: 현재 model identifier
- created_at: ISO 8601 timestamp
- context_window: Token 사용 통계
- persistence: Recovery 설정
- work_state: 현재 task state

State Management Functions:
- create_session_snapshot(): 현재 state 캡처
- restore_session_state(snapshot): Snapshot에서 복원
- validate_session_state(state): State integrity 검증

### Handoff Protocol API

Handoff Package Structure:
- handoff_id: Unique transfer identifier
- from_agent: Source agent type
- to_agent: Destination agent type
- session_context: Token 및 model 정보
- task_context: 현재 work state
- recovery_info: Checkpoint 및 fork data

Validation Functions:
- validate_handoff(package): Package integrity 검증
- can_agents_cooperate(from, to): 호환성 확인

---

## Configuration Options

### Token Budget Allocation

Default Allocation (200K total):
- System Prompt and Instructions: 15K tokens (7.5%)
- Active Conversation: 80K tokens (40%)
- Reference Context: 50K tokens (25%)
- Reserve (Emergency): 55K tokens (27.5%)

Customizable Settings:
- system_prompt_budget: System 할당 override
- conversation_budget: Conversation 할당 override
- reference_budget: Reference 할당 override
- reserve_budget: Emergency reserve override

### Clear Execution Settings

Mandatory Clear Points:
- /do:1-plan 완료 후
- Context 150K tokens 초과
- Conversation 50개 메시지 초과
- Major phase transition 전
- Model 전환 (Haiku → Sonnet)

Configuration Options:
- auto_clear_enabled: 자동 clearing 활성화
- clear_threshold_tokens: Auto-clear용 token threshold
- clear_threshold_messages: Message count threshold
- preserve_on_clear: Clear 시 보존할 context types

### Session Persistence Settings

State Layers Configuration:
- L1: Context-Aware Layer (model features)
- L2: Active Context (current task)
- L3: Session History (recent actions)
- L4: Project State (SPEC progress)
- L5: User Context (preferences)
- L6: System State (tools, permissions)

Persistence Options:
- auto_load_history: 이전 context 복원
- context_preservation: 보존 level
- cache_enabled: Context caching 활성화

---

## Integration Patterns

### Plan-Run-Sync Workflow Integration

Workflow Sequence:
1. /do:1-plan 실행
2. /clear (mandatory - 45-50K tokens 절약)
3. /do:2-run SPEC-XXX
4. Multi-agent handoff
5. /do:3-sync SPEC-XXX
6. Session state persistence

Token Savings:
- Post-plan clear: 45-50K tokens 절약
- Progressive disclosure: 30-40% 감소
- Handoff optimization: Transfer당 15-20K 절약

### Multi-Agent Coordination

Handoff Workflow:
1. Source agent가 task phase 완료
2. Minimal context로 handoff package 생성
3. Handoff integrity 검증
4. Target agent가 수신 및 검증
5. Target agent가 workflow 계속

Context Minimization:
- SPEC ID 및 key requirements만 포함
- Architecture summary 200자 제한
- Background 및 reasoning 제외
- Critical state만 transfer

### Progressive Disclosure Integration

Loading Tiers:
- Tier 1: CLAUDE.md, config.json (항상 로드)
- Tier 2: 현재 SPEC 및 implementation files
- Tier 3: 관련 modules 및 dependencies
- Tier 4: Reference documentation (on-demand)

Disclosure Triggers:
- 명시적 user 요청
- Error recovery 필요
- Complex implementation 필요
- Documentation reference 필요

---

## Troubleshooting

### Context Overflow Issues

Symptoms: 성능 저하, 불완전한 response

Solutions:
1. /clear 즉시 실행
2. 로드된 context tier 감소
3. Progressive summarization 적용
4. 여러 session으로 task 분할

Prevention:
- 60% threshold에서 모니터링
- Major milestone 후 clear
- Aggressive clearing strategy 사용

### Session Recovery Failures

Symptoms: 중단 후 state 손실

Solutions:
1. Session ID가 persist 되었는지 확인
2. Snapshot integrity 확인
3. 가장 최근 checkpoint에서 복원
4. Project files에서 state 재구성

Prevention:
- Operation 전 checkpoint 생성
- Clearing 전 session ID persist
- State snapshot auto-save 활성화

### Handoff Validation Errors

Symptoms: Agent transition 실패

Solutions:
1. Available tokens가 30K 초과하는지 확인
2. Agent 호환성 확인
3. Handoff package 크기 감소
4. Transfer 전 context compression 트리거

Prevention:
- Package 생성 전 검증
- Critical context만 포함
- Handoff overhead용 token 예약

### Token Budget Exhaustion

Symptoms: 강제 중단, emergency behavior

Solutions:
1. 즉시 /clear 실행
2. Tier 1 context만으로 재개
3. 추가 context를 점진적으로 로드
4. 남은 작업을 여러 session으로 분할

Prevention:
- 55K emergency reserve 유지
- 85% threshold에서 clear 실행
- Progressive disclosure 일관되게 적용

---

## External Resources

### Related Documentation

- Token Management Best Practices
- Session State Architecture Guide
- Multi-Agent Coordination Patterns
- Context Optimization Strategies

### Performance Metrics

Target Metrics:
- Token Efficiency: Clearing으로 60-70% 감소
- Context Overhead: System metadata에 15K 미만
- Handoff Success Rate: 95% 초과
- Session Recovery: 5초 미만
- Memory Files: 파일당 500줄 미만

### Related Skills

- do-foundation-claude: Claude Code authoring 및 configuration
- do-foundation-core: Core execution patterns 및 SPEC workflow
- do-workflow-project: Project management 및 documentation
- do-cc-memory: Memory management 및 persistence

---

Version: 3.0.0
Last Updated: 2025-12-06
