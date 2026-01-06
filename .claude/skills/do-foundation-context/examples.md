# Context & Session 관리 예제

200K token budget 최적화, session persistence, multi-agent handoff, progressive disclosure를 다루는 실용적 context 관리 pattern 예제.

---

## Example 1: Token Budget 모니터링 및 관리

**Scenario**: 자동 action이 포함된 실시간 token budget 모니터링 구현.

**Input**:
```
// 필요 기능:
// - 200K limit 대비 사용량 추적
// - Critical threshold에서 경고
// - 자동 압축 트리거
// - Context overflow 방지
```

**Output**:
```javascript
// TokenBudgetState 구조
const TokenBudgetState = {
    total_budget: 200000,
    used_tokens: 0,
    system_reserved: 15000,  // System prompt, instructions
    emergency_reserve: 55000, // 마지막 27.5% recovery용

    get available() {
        return this.total_budget - this.used_tokens;
    },

    get usage_percent() {
        return (this.used_tokens / this.total_budget) * 100;
    },

    get is_critical() {
        return this.usage_percent >= 85;
    },

    get is_warning() {
        return this.usage_percent >= 75;
    }
};

// TokenBudgetManager
const TokenBudgetManager = {
    // Budget 할당 (권장 분배)
    ALLOCATION: {
        system: 15000,       // 7.5% - System prompt, instructions
        conversation: 80000, // 40% - Active conversation
        reference: 50000,    // 25% - Reference context
        reserve: 55000       // 27.5% - Emergency reserve
    },

    THRESHOLDS: {
        normal: 60,
        warning: 75,
        critical: 85,
        emergency: 95
    },

    state: { ...TokenBudgetState },
    callbacks: {},
    checkpoints: [],

    registerCallback(event, callback) {
        this.callbacks[event] = callback;
    },

    updateUsage(tokens_used, source = "unknown") {
        this.state.used_tokens = tokens_used;

        console.log(
            `Token update: ${tokens_used.toLocaleString()}/${this.state.total_budget.toLocaleString()} ` +
            `(${this.state.usage_percent.toFixed(1)}%) from ${source}`
        );

        // Threshold 확인 및 action 트리거
        if (this.state.usage_percent >= this.THRESHOLDS.emergency) {
            this.handleEmergency();
        } else if (this.state.usage_percent >= this.THRESHOLDS.critical) {
            this.handleCritical();
        } else if (this.state.usage_percent >= this.THRESHOLDS.warning) {
            this.handleWarning();
        }
    },

    handleWarning() {
        console.warn(`Token usage at ${this.state.usage_percent.toFixed(1)}% - Context optimization 시작`);
        if (this.callbacks.on_warning) {
            this.callbacks.on_warning(this.state);
        }
    },

    handleCritical() {
        console.error(`Token usage CRITICAL at ${this.state.usage_percent.toFixed(1)}% - Context compression 트리거`);
        this.createCheckpoint("pre_compression");
        if (this.callbacks.on_critical) {
            this.callbacks.on_critical(this.state);
        }
    },

    handleEmergency() {
        console.error(`Token usage EMERGENCY at ${this.state.usage_percent.toFixed(1)}% - Context clear 강제`);
        if (this.callbacks.on_emergency) {
            this.callbacks.on_emergency(this.state);
        }
    },

    createCheckpoint(name) {
        const checkpoint = {
            name: name,
            timestamp: new Date().toISOString(),
            tokens_used: this.state.used_tokens,
            usage_percent: this.state.usage_percent
        };
        this.checkpoints.push(checkpoint);
        console.log(`Checkpoint created: ${name}`);
        return checkpoint;
    },

    getOptimizationSuggestions() {
        const suggestions = [];

        if (this.state.usage_percent > 50) {
            suggestions.push("Context reset을 위해 /clear 사용 고려");
        }
        if (this.state.usage_percent > 60) {
            suggestions.push("Non-critical reference documentation 지연 로드");
        }
        if (this.state.usage_percent > 75) {
            suggestions.push("Conversation에 progressive summarization 적용");
        }
        if (this.state.usage_percent > 85) {
            suggestions.push("긴급: 즉시 context compression 실행 필요");
        }

        return suggestions;
    }
};

// 사용 예제
TokenBudgetManager.registerCallback("on_warning", (state) => {
    console.log(`Warning: ${state.available.toLocaleString()} tokens 남음`);
});

TokenBudgetManager.registerCallback("on_critical", (state) => {
    console.log("Critical: Context 압축 중...");
});

TokenBudgetManager.registerCallback("on_emergency", (state) => {
    console.log("Emergency: /clear command 강제 실행");
});

// 사용량 업데이트 시뮬레이션
TokenBudgetManager.updateUsage(120000, "conversation");  // 60% - normal
TokenBudgetManager.updateUsage(155000, "file_loading");  // 77.5% - warning 트리거
TokenBudgetManager.updateUsage(175000, "agent_output");  // 87.5% - critical 트리거
```

**설명**: 이 pattern은 threshold 기반 callback, checkpointing, 최적화 제안이 포함된 종합적인 token budget 모니터링 구현.

---

## Example 2: Session State Persistence

**Scenario**: 재개 가능한 workflow를 위한 session state persistence 구현.

**Input**:
```
// 필요 기능:
// - 중단 후에도 state 유지
// - Workflow 재개 지원
// - SPEC 진행 추적
// - Model 전환 처리
```

**Output**:
```javascript
// SessionPhase 열거
const SessionPhase = {
    INIT: "initialization",
    PLANNING: "planning",
    IMPLEMENTATION: "implementation",
    TESTING: "testing",
    DOCUMENTATION: "documentation",
    REVIEW: "review",
    COMPLETED: "completed"
};

// WorkState 구조
const createWorkState = () => ({
    current_spec: null,
    phase: SessionPhase.INIT,
    completed_steps: [],
    pending_steps: [],
    artifacts: {}
});

// ContextWindow 구조
const createContextWindow = () => ({
    total: 200000,
    used: 0,
    available: 200000,
    position_percent: 0.0
});

// SessionState 구조
const createSessionState = (session_id, model) => ({
    session_id: session_id,
    model: model,
    created_at: new Date().toISOString(),
    last_updated: new Date().toISOString(),
    context_window: createContextWindow(),
    work_state: createWorkState(),
    user_context: {},
    persistence: {
        auto_save: true,
        save_interval_seconds: 60,
        context_preservation: "critical_only"
    }
});

// SessionManager
const SessionManager = {
    storage_path: ".do/sessions",
    current_session: null,

    createSession(model = "claude-sonnet-4-5-20250929", user_context = null) {
        const session_id = `sess_${generateUUID().substring(0, 12)}`;

        const session = createSessionState(session_id, model);
        session.user_context = user_context || {};

        this.current_session = session;
        this.saveSession(session);
        return session;
    },

    loadSession(session_id) {
        const file_path = `${this.storage_path}/${session_id}.json`;

        // 파일에서 로드 (파일 시스템 접근)
        const data = readJsonFile(file_path);
        if (!data) return null;

        // Session 재구성
        const session = createSessionState(data.session_id, data.model);
        session.created_at = data.created_at;
        session.last_updated = data.last_updated;
        session.context_window = data.context_window;
        session.work_state = data.work_state;
        session.user_context = data.user_context || {};
        session.persistence = data.persistence || {};

        this.current_session = session;
        return session;
    },

    updateWorkState(options = {}) {
        if (!this.current_session) {
            throw new Error("No active session");
        }

        const work = this.current_session.work_state;

        if (options.spec_id) work.current_spec = options.spec_id;
        if (options.phase) work.phase = options.phase;
        if (options.completed_step) {
            work.completed_steps.push(options.completed_step);
            const idx = work.pending_steps.indexOf(options.completed_step);
            if (idx > -1) work.pending_steps.splice(idx, 1);
        }
        if (options.artifact) {
            work.artifacts[options.artifact.key] = options.artifact.value;
        }

        this.current_session.last_updated = new Date().toISOString();
        this.saveSession(this.current_session);
    },

    updateContextUsage(tokens_used) {
        if (!this.current_session) return;

        const ctx = this.current_session.context_window;
        ctx.used = tokens_used;
        ctx.available = ctx.total - tokens_used;
        ctx.position_percent = (tokens_used / ctx.total) * 100;

        this.current_session.last_updated = new Date().toISOString();
        this.saveSession(this.current_session);
    },

    getResumptionContext() {
        if (!this.current_session) return {};

        return {
            spec_id: this.current_session.work_state.current_spec,
            phase: this.current_session.work_state.phase,
            completed: this.current_session.work_state.completed_steps,
            pending: this.current_session.work_state.pending_steps,
            last_update: this.current_session.last_updated,
            context_usage: this.current_session.context_window.position_percent
        };
    },

    prepareForClear() {
        if (!this.current_session) return {};

        // 현재 state 저장
        this.saveSession(this.current_session);

        // Clear 후 reload할 minimal context 반환
        return {
            session_id: this.current_session.session_id,
            spec_id: this.current_session.work_state.current_spec,
            phase: this.current_session.work_state.phase,
            reload_files: [
                "CLAUDE.md",
                `.do/specs/${this.current_session.work_state.current_spec}.md`
            ]
        };
    },

    saveSession(session) {
        const file_path = `${this.storage_path}/${session.session_id}.json`;
        writeJsonFile(file_path, session);
    }
};

// 사용 예제
const session = SessionManager.createSession(
    "claude-sonnet-4-5-20250929",
    { language: "ko", user_name: "Developer" }
);
console.log(`Session 생성: ${session.session_id}`);

// 작업 진행 업데이트
SessionManager.updateWorkState({ spec_id: "SPEC-001", phase: SessionPhase.IMPLEMENTATION });
SessionManager.updateWorkState({ completed_step: "spec_complete" });
SessionManager.updateWorkState({ completed_step: "architecture_defined" });
SessionManager.updateWorkState({ artifact: { key: "api_schema", value: "schemas/user.json" } });

// /clear 전
const clear_context = SessionManager.prepareForClear();
console.log(`Clear 후 reload할 context: ${JSON.stringify(clear_context)}`);
```

**설명**: 이 pattern은 persistence, 작업 진행 추적, 중단 또는 /clear 후 seamless 재개를 위한 종합적인 session state 관리 제공.

---

## Example 3: Multi-Agent Handoff Protocol

**Scenario**: Agent 간 효율적인 context handoff 구현.

**Input**:
```
// 필요 기능:
// - Token overhead 최소화
// - Handoff integrity 검증
// - Context compression 지원
// - Agent 호환성 처리
```

**Output**:
```javascript
// AgentType 열거
const AgentType = {
    SPEC_BUILDER: "spec-builder",
    TDD_IMPLEMENTER: "tdd-implementer",
    BACKEND_EXPERT: "backend-expert",
    FRONTEND_EXPERT: "frontend-expert",
    DOCS_MANAGER: "docs-manager",
    QUALITY_GATE: "quality-gate"
};

// Agent 호환성 매트릭스
const AGENT_COMPATIBILITY = {
    [AgentType.SPEC_BUILDER]: [
        AgentType.TDD_IMPLEMENTER,
        AgentType.BACKEND_EXPERT,
        AgentType.FRONTEND_EXPERT
    ],
    [AgentType.TDD_IMPLEMENTER]: [
        AgentType.QUALITY_GATE,
        AgentType.DOCS_MANAGER
    ],
    [AgentType.BACKEND_EXPERT]: [
        AgentType.FRONTEND_EXPERT,
        AgentType.QUALITY_GATE,
        AgentType.DOCS_MANAGER
    ],
    [AgentType.FRONTEND_EXPERT]: [
        AgentType.QUALITY_GATE,
        AgentType.DOCS_MANAGER
    ],
    [AgentType.QUALITY_GATE]: [
        AgentType.DOCS_MANAGER
    ]
};

// HandoffPackage 생성
const createHandoffPackage = (from_agent, to_agent, session_ctx, task_ctx, recovery_info) => {
    const handoff_id = `hoff_${generateUUID().substring(0, 8)}`;
    const created_at = new Date().toISOString();

    const pkg = {
        handoff_id: handoff_id,
        from_agent: from_agent,
        to_agent: to_agent,
        session_context: session_ctx,
        task_context: task_ctx,
        recovery_info: recovery_info,
        created_at: created_at,
        checksum: ""
    };

    pkg.checksum = calculateChecksum(pkg);
    return pkg;
};

const calculateChecksum = (pkg) => {
    const content = `${pkg.handoff_id}${pkg.from_agent}${pkg.to_agent}${pkg.task_context.spec_id}${pkg.created_at}`;
    return hashString(content).substring(0, 16);
};

// HandoffManager
const HandoffManager = {
    MINIMUM_SAFE_TOKENS: 30000,
    handoff_history: [],

    canAgentsCooperate(from_agent, to_agent) {
        const compatible = AGENT_COMPATIBILITY[from_agent] || [];
        return compatible.includes(to_agent);
    },

    createHandoff(from_agent, to_agent, session_ctx, task_ctx, recovery_info) {
        // 호환성 검증
        if (!this.canAgentsCooperate(from_agent, to_agent)) {
            throw new Error(`Agent ${from_agent}가 ${to_agent}로 handoff 불가`);
        }

        // Token budget 검증
        if (session_ctx.available_tokens < this.MINIMUM_SAFE_TOKENS) {
            throw new Error(
                `Token 부족: ${session_ctx.available_tokens} < ${this.MINIMUM_SAFE_TOKENS} 필요`
            );
        }

        const handoff = createHandoffPackage(from_agent, to_agent, session_ctx, task_ctx, recovery_info);
        this.handoff_history.push(handoff);
        return handoff;
    },

    validateHandoff(pkg) {
        // Checksum 검증
        const expected_checksum = calculateChecksum(pkg);
        if (pkg.checksum !== expected_checksum) {
            throw new Error("Handoff checksum 불일치");
        }

        // Agent 호환성 검증
        if (!this.canAgentsCooperate(pkg.from_agent, pkg.to_agent)) {
            throw new Error("Agent 간 협력 불가");
        }

        // Token budget 검증
        if (pkg.session_context.available_tokens < this.MINIMUM_SAFE_TOKENS) {
            return this.triggerContextCompression(pkg);
        }

        return true;
    },

    triggerContextCompression(pkg) {
        console.log(`Handoff ${pkg.handoff_id}에 대해 context 압축 중`);
        // Progressive summarization 적용
        return true;
    },

    extractMinimalContext(full_context) {
        const priority_fields = [
            "spec_id",
            "current_phase",
            "next_step",
            "critical_decisions",
            "blocking_issues"
        ];

        const result = {};
        for (const key of priority_fields) {
            if (full_context[key] !== undefined) {
                result[key] = full_context[key];
            }
        }
        return result;
    },

    getHandoffSummary() {
        return {
            total_handoffs: this.handoff_history.length,
            handoffs: this.handoff_history.map(h => ({
                id: h.handoff_id,
                from: h.from_agent,
                to: h.to_agent,
                spec: h.task_context.spec_id,
                timestamp: h.created_at
            }))
        };
    }
};

// 사용 예제
const session_ctx = {
    session_id: "sess_abc123",
    model: "claude-sonnet-4-5-20250929",
    context_position: 42.5,
    available_tokens: 115000,
    user_language: "ko"
};

const task_ctx = {
    spec_id: "SPEC-001",
    current_phase: "planning_complete",
    completed_steps: ["requirement_analysis", "spec_creation", "architecture_design"],
    next_step: "write_tests",
    key_artifacts: {
        spec_document: ".do/specs/SPEC-001.md",
        architecture: ".do/architecture/SPEC-001.mermaid"
    }
};

const recovery = {
    last_checkpoint: new Date().toISOString(),
    recovery_tokens_reserved: 55000,
    session_fork_available: true
};

try {
    const handoff = HandoffManager.createHandoff(
        AgentType.SPEC_BUILDER,
        AgentType.TDD_IMPLEMENTER,
        session_ctx,
        task_ctx,
        recovery
    );

    console.log(`Handoff 생성: ${handoff.handoff_id}`);
    console.log(`Checksum: ${handoff.checksum}`);

    // 다음 agent로 전송 전 검증
    const is_valid = HandoffManager.validateHandoff(handoff);
    console.log(`Handoff valid: ${is_valid}`);

} catch (e) {
    console.log(`Handoff 실패: ${e.message}`);
}
```

**설명**: 이 pattern은 호환성 검사, token budget 검증, integrity 확인, 효율적인 agent 조율을 위한 context compression이 포함된 robust한 multi-agent handoff 구현.

---

## Common Patterns

### Pattern 1: Aggressive /clear Strategy

전략적 checkpoint에서 /clear 실행:

```javascript
const ClearStrategy = {
    CLEAR_TRIGGERS: {
        post_spec_creation: true,
        token_threshold_150k: true,
        message_count_50: true,
        phase_transition: true,
        model_switch: true
    },

    shouldClear(context) {
        // Token threshold 확인
        if ((context.token_usage || 0) > 150000) {
            return { clear: true, reason: "Token threshold 초과 (>150K)" };
        }

        // Message count 확인
        if ((context.message_count || 0) > 50) {
            return { clear: true, reason: "Message count 초과 (>50)" };
        }

        // Phase transition 확인
        if (context.phase_changed) {
            return { clear: true, reason: "Phase transition 감지" };
        }

        // Post-SPEC creation 확인
        if (context.spec_just_created) {
            return { clear: true, reason: "SPEC 생성 완료" };
        }

        return { clear: false, reason: "Clear 불필요" };
    },

    prepareClearContext(session) {
        return {
            session_id: session.session_id,
            spec_id: session.work_state.current_spec,
            phase: session.work_state.phase,
            reload_sequence: [
                "CLAUDE.md",
                `.do/specs/${session.work_state.current_spec}.md`,
                "src/main.php"  // 현재 작업 파일
            ],
            preserved_decisions: session.work_state.artifacts.decisions || []
        };
    }
};
```

### Pattern 2: Progressive Summarization

Key 정보 보존하면서 context 압축:

```javascript
const ProgressiveSummarizer = {
    summarizeConversation(messages, target_ratio = 0.3) {
        // Key 정보 추출
        const decisions = this.extractDecisions(messages);
        const code_changes = this.extractCodeChanges(messages);
        const issues = this.extractIssues(messages);

        // 압축된 summary 생성
        const summary = `
## Conversation Summary

### Key Decisions
${this.formatList(decisions)}

### Code Changes
${this.formatList(code_changes)}

### Open Issues
${this.formatList(issues)}

### Reference
Original conversation: ${messages.length} messages
Compression ratio: ${(target_ratio * 100).toFixed(0)}%
`;
        return summary;
    },

    extractDecisions(messages) {
        const decisions = [];
        const decision_markers = ["decided", "agreed", "will use", "chosen"];

        for (const msg of messages) {
            const content = (msg.content || "").toLowerCase();
            if (decision_markers.some(marker => content.includes(marker))) {
                decisions.push(this.extractSentence(msg.content));
            }
        }

        return decisions.slice(0, 5);  // Top 5 decisions
    },

    extractCodeChanges(messages) {
        const changes = [];
        for (const msg of messages) {
            if ((msg.content || "").includes("```")) {
                changes.push(`Modified: ${msg.file || "unknown"}`);
            }
        }
        return changes;
    },

    extractIssues(messages) {
        const issues = [];
        const issue_markers = ["todo", "fixme", "issue", "problem", "bug"];

        for (const msg of messages) {
            const content = (msg.content || "").toLowerCase();
            if (issue_markers.some(marker => content.includes(marker))) {
                issues.push(this.extractSentence(msg.content));
            }
        }
        return issues;
    },

    extractSentence(text) {
        const sentences = text.split(".");
        return sentences[0] ? sentences[0].substring(0, 100) : text.substring(0, 100);
    },

    formatList(items) {
        if (!items.length) return "- None";
        return items.map(item => `- ${item}`).join("\n");
    }
};
```

### Pattern 3: Context Tag References

Inline content 대신 효율적인 reference 사용:

```javascript
const ContextTagManager = {
    tags: {},

    registerTag(tag_id, content) {
        this.tags[tag_id] = content;
        return `@${tag_id}`;
    },

    resolveTag(tag_ref) {
        if (tag_ref.startsWith("@")) {
            const tag_id = tag_ref.substring(1);
            return this.tags[tag_id] || null;
        }
        return null;
    },

    createMinimalReference(full_context) {
        const references = {};

        for (const [key, value] of Object.entries(full_context)) {
            if (typeof value === "string" && value.length > 200) {
                // Full content 저장, reference 반환
                const tag_id = `${key.toUpperCase()}-001`;
                this.registerTag(tag_id, value);
                references[key] = `@${tag_id}`;
            } else {
                references[key] = value;
            }
        }

        return references;
    }
};

// 사용
// Inline content 대신 (high token cost):
// "이전 20개 메시지의 user configuration..."

// 효율적인 reference 사용 (low token cost):
ContextTagManager.registerTag("CONFIG-001", full_config_content);
const reference = "@CONFIG-001";  // 10 tokens vs 500+ tokens
```

---

## Anti-Patterns (피해야 할 Pattern)

### Anti-Pattern 1: Token Warning 무시

**문제점**: Token warning 처리 없이 작업 계속.

```javascript
// 잘못된 접근
if (token_usage > 150000) {
    console.log("Warning: High token usage");
    // 그냥 계속 작업 - context overflow 유발
    continueWork();
}
```

**해결**: Warning에 즉시 action 취하기.

```javascript
// 올바른 접근
if (token_usage > 150000) {
    console.warn("Token warning triggered");
    // Checkpoint 생성 후 clear
    const checkpoint = saveCurrentState();
    executeClear();
    restoreEssentialContext(checkpoint);
}
```

### Anti-Pattern 2: Handoff에 Full Context 포함

**문제점**: Agent 간 complete context 전달은 token 낭비.

```javascript
// 잘못된 접근
const handoff = {
    full_conversation: all_messages,  // 50K tokens
    all_files_content: file_contents, // 100K tokens
    complete_history: history         // 30K tokens
};
```

**해결**: Critical context만 전달.

```javascript
// 올바른 접근
const handoff = {
    spec_id: "SPEC-001",
    current_phase: "implementation",
    next_step: "write_tests",
    key_decisions: ["Use JWT", "PostgreSQL"],
    file_references: ["@API-SCHEMA", "@DB-MODEL"]
};
```

### Anti-Pattern 3: Session Persistence 없음

**문제점**: 중단 시 작업 진행 손실.

```javascript
// 잘못된 접근
// State 저장 안함 - /clear 또는 중단 시 모든 진행 손실
const work_in_progress = processSpec(spec_id);
// 연결 끊김 - 작업 손실
```

**해결**: 지속적으로 state 저장.

```javascript
// 올바른 접근
SessionManager.createCheckpoint("pre_processing");

const work_in_progress = processSpec(spec_id);
SessionManager.updateWorkState({ completed_step: "processing_done" });
SessionManager.saveSession(session);  // State 보존

// 중단 후
const resumed = SessionManager.loadSession(session_id);
// Checkpoint에서 계속
```

---

## Workflow Integration

### Context 관리와 SPEC-First Workflow

```javascript
// Context 최적화가 포함된 complete workflow

// Phase 1: Planning (약 40K tokens 사용)
const analysis = Task({ subagent_type: "spec-builder", prompt: "Analyze: user auth" });
SessionManager.updateWorkState({ phase: SessionPhase.PLANNING });

// Planning 후 mandatory /clear (45-50K tokens 절약)
const clear_context = SessionManager.prepareForClear();
executeClear();
restoreFromCheckpoint(clear_context);

// Phase 2: Implementation (fresh 200K budget)
const implementation = Task({
    subagent_type: "tdd-implementer",
    prompt: `Implement: ${clear_context.spec_id}`
});
SessionManager.updateWorkState({ phase: SessionPhase.IMPLEMENTATION });

// 필요시 모니터링 및 clear
if (token_usage > 150000) {
    clearAndResume();
}

// Phase 3: Documentation
const docs = Task({ subagent_type: "docs-manager", prompt: "Generate docs" });
SessionManager.updateWorkState({ phase: SessionPhase.DOCUMENTATION });
```

---

Version: 3.0.0
Last Updated: 2025-12-06
