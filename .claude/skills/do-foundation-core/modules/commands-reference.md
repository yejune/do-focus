# Commands Reference - Do Core Commands

Purpose: Complete reference for Do's 6 core commands used in SPEC-First TDD workflow.

Last Updated: 2025-11-25
Version: 2.0.0

---

## Quick Reference (30 seconds)

Do provides 6 core commands for SPEC-First TDD execution:

| Command            | Purpose                | Phase         |
| ------------------ | ---------------------- | ------------- |
| `/do:0-project`  | Project initialization | Setup         |
| `/do:1-plan`     | SPEC generation        | Planning      |
| `/do:2-run`      | TDD implementation     | Development   |
| `/do:3-sync`     | Documentation sync     | Documentation |
| `/do:9-feedback` | Feedback collection    | Improvement   |
| `/do:99-release` | Production deployment  | Release       |

Required Workflow:
```
1. /do:0-project # Initialize
2. /do:1-plan "description" # Generate SPEC
3. /clear # Clear context (REQUIRED)
4. /do:2-run SPEC-001 # Implement
5. /do:3-sync SPEC-001 # Document
6. /do:9-feedback # Improve
```

Critical Rule: Execute `/clear` after `/do:1-plan` (saves 45-50K tokens)

---

## Implementation Guide (5 minutes)

### `/do:0-project` - Project Initialization

Purpose: Initialize project structure and generate configuration

Agent Delegation: `workflow-project`

Usage:
```bash
/do:0-project
/do:0-project --with-git
```

What It Does:
1. Creates `.do/` directory structure
2. Generates `config.json` with default settings
3. Initializes Git repository (if `--with-git` flag provided)
4. Sets up Do workflows

Output:
- `.do/` directory
- `.do/config/config.yaml`
- `.do/memory/` (empty, ready for session state)
- `.do/logs/` (empty, ready for logging)

Next Step: Ready for SPEC generation via `/do:1-plan`

Example:
```
User: /do:0-project
Do: Project initialized successfully.
 - .do/config/config.yaml created
 - Git workflow set to 'manual' mode
 Ready for SPEC generation.
```

---

### `/do:1-plan` - SPEC Generation

Purpose: Generate SPEC document in EARS format

Agent Delegation: `workflow-spec`

Usage:
```bash
/do:1-plan "Implement user authentication endpoint (JWT)"
/do:1-plan "Add dark mode toggle to settings page"
```

What It Does:
1. Analyzes user request
2. Generates EARS format SPEC document
3. Creates `.do/specs/SPEC-XXX/` directory
4. Saves `spec.md` with requirements

EARS Format (5 sections):
- WHEN (trigger conditions)
- IF (preconditions)
- THE SYSTEM SHALL (functional requirements)
- WHERE (constraints)
- UBIQUITOUS (quality requirements)

Output:
- `.do/specs/SPEC-001/spec.md` (EARS document)
- SPEC ID assigned (auto-incremented)

CRITICAL: Execute `/clear` immediately after completion
- Saves 45-50K tokens
- Prepares clean context for implementation

Example:
```
User: /do:1-plan "Implement user authentication endpoint (JWT)"
Do: SPEC-001 generated successfully.
 Location: .do/specs/SPEC-001/spec.md

 IMPORTANT: Execute /clear now to free 45-50K tokens.
```

---

### `/do:2-run` - TDD Implementation

Purpose: Execute RED-GREEN-REFACTOR cycle

Agent Delegation: `workflow-tdd`

Usage:
```bash
/do:2-run SPEC-001
/do:2-run SPEC-002
```

What It Does:
1. Reads SPEC document
2. Executes TDD cycle in 3 phases:
 - RED: Write failing tests
 - GREEN: Implement minimal code to pass tests
 - REFACTOR: Optimize and clean up code
3. Validates TRUST 5 quality gates
4. Generates implementation report

TDD Process:
```
Phase 1 (RED):
 - Write failing tests for each requirement
 - Run tests → ALL FAIL (expected)

Phase 2 (GREEN):
 - Implement minimal code to pass tests
 - Run tests → ALL PASS

Phase 3 (REFACTOR):
 - Optimize code structure
 - Improve readability
 - Remove duplication
 - Run tests → ALL PASS (maintained)
```

Output:
- Implemented code (in source directories)
- Test files (in test directories)
- Quality report (TRUST 5 validation)

Requirement: Test coverage ≥ 85% (TRUST 5)

Example:
```
User: /do:2-run SPEC-001
Do: TDD implementation cycle started for SPEC-001.

 RED: 12 failing tests written
 GREEN: Implementation complete, all tests passing
 REFACTOR: Code optimized

 Test Coverage: 92% ( meets 85% threshold)
 TRUST 5: All gates passed
```

---

### `/do:3-sync` - Documentation Synchronization

Purpose: Auto-generate API documentation and project artifacts

Agent Delegation: `workflow-docs`

Usage:
```bash
/do:3-sync SPEC-001
/do:3-sync SPEC-002
```

What It Does:
1. Reads implemented code
2. Generates API documentation (OpenAPI format)
3. Creates architecture diagrams
4. Produces project completion report

Output:
- API documentation (OpenAPI/Swagger format)
- Architecture diagrams (Mermaid)
- `.do/docs/SPEC-001/` directory
- Project report

Example:
```
User: /do:3-sync SPEC-001
Do: Documentation synchronized for SPEC-001.

 Generated:
 - API documentation: .do/docs/SPEC-001/api.yaml
 - Architecture diagram: .do/docs/SPEC-001/architecture.md
 - Completion report: .do/docs/SPEC-001/report.md
```

---

### `/do:9-feedback` - Improvement Feedback Collection

Purpose: Error analysis and improvement suggestions

Agent Delegation: `core-quality`

Usage:
```bash
/do:9-feedback
/do:9-feedback --analyze SPEC-001
```

What It Does:
1. Analyzes errors encountered during workflow
2. Collects improvement suggestions
3. Reports to Do development team
4. Proposes error recovery strategies

Use Cases:
- Errors: When errors occur during any workflow phase
- Improvements: When Do enhancements are identified
- Analysis: Post-implementation review

Example:
```
User: /do:9-feedback
Do: Collecting feedback for recent session.

 Errors: 2 permission issues detected
 Improvements: 1 token optimization suggestion

 Feedback submitted to Do development team.
```

---

### `/do:99-release` - Production Deployment

Purpose: Production deployment workflow

Agent Delegation: `infra-devops`

Usage:
```bash
/do:99-release
```

What It Does:
1. Validates all TRUST 5 quality gates
2. Runs full test suite
3. Builds production artifacts
4. Deploys to production environment

Note: This command is local-only and NOT synchronized to the package template. It's for local development and testing.

---

## Advanced Implementation (10+ minutes)

### Context Initialization Rules

Rule 1: Execute `/clear` AFTER `/do:1-plan` (mandatory)
- SPEC generation uses 45-50K tokens
- `/clear` frees this context for implementation phase
- Prevents context overflow

Rule 2: Execute `/clear` when context > 150K tokens
- Monitor context usage via `/context` command
- Prevents token limit exceeded errors

Rule 3: Execute `/clear` after 50+ conversation messages
- Accumulated context from conversation history
- Reset for fresh context

Why `/clear` is critical:
```
Without /clear:
 SPEC generation: 50K tokens
 Implementation: 100K tokens
 Total: 150K tokens (approaching 200K limit)

With /clear:
 SPEC generation: 50K tokens
 /clear: 0K tokens (reset)
 Implementation: 100K tokens
 Total: 100K tokens (50K budget remaining)
```

---

### Command Delegation Patterns

Each command delegates to a specific agent:

| Command            | Agent              | Agent Type              |
| ------------------ | ------------------ | ----------------------- |
| `/do:0-project`  | `workflow-project` | Tier 1 (Always Active)  |
| `/do:1-plan`     | `workflow-spec`    | Tier 1 (Always Active)  |
| `/do:2-run`      | `workflow-tdd`     | Tier 1 (Always Active)  |
| `/do:3-sync`     | `workflow-docs`    | Tier 1 (Always Active)  |
| `/do:9-feedback` | `core-quality`     | Tier 2 (Auto-triggered) |
| `/do:99-release` | `infra-devops`     | Tier 3 (Lazy-loaded)    |

Delegation Flow:
```
User executes command
 ↓
Do receives command
 ↓
Command processor agent invoked
 ↓
Agent executes workflow
 ↓
Results reported to user
```

---

### Token Budget by Command

| Command        | Average Tokens | Phase Budget                          |
| -------------- | -------------- | ------------------------------------- |
| `/do:1-plan` | 45-50K         | Planning Phase (30K allocated)        |
| `/do:2-run`  | 80-100K        | Implementation Phase (180K allocated) |
| `/do:3-sync` | 20-25K         | Documentation Phase (40K allocated)   |
| Total          | 145-175K       | 250K per feature                      |

Optimization:
- Use Haiku 4.5 for `/do:2-run` (fast, cost-effective)
- Use Sonnet 4.5 for `/do:1-plan` (high-quality SPEC)
- Execute `/clear` between phases (critical)

---

### Error Handling

Common Errors:

| Error                     | Command                | Solution                                    |
| ------------------------- | ---------------------- | ------------------------------------------- |
| "Project not initialized" | `/do:1-plan`         | Run `/do:0-project` first                 |
| "SPEC not found"          | `/do:2-run SPEC-999` | Verify SPEC ID exists                       |
| "Token limit exceeded"    | Any                    | Execute `/clear` immediately                |
| "Test coverage < 85%"     | `/do:2-run`          | `core-quality` auto-generates missing tests |

Recovery Pattern:
```bash
# Error: Token limit exceeded
1. /clear # Reset context
2. /do:2-run SPEC-001 # Retry with clean context
```

---

### Workflow Variations

Standard Workflow (Full SPEC):
```
/do:0-project → /do:1-plan → /clear → /do:2-run → /do:3-sync
```

Quick Workflow (No SPEC for simple tasks):
```
/do:0-project → Direct implementation (for 1-2 file changes)
```

Iterative Workflow (Multiple SPECs):
```
/do:1-plan "Feature A" → /clear → /do:2-run SPEC-001 → /do:3-sync SPEC-001
/do:1-plan "Feature B" → /clear → /do:2-run SPEC-002 → /do:3-sync SPEC-002
```

---

### Integration with Git Workflow

Commands automatically integrate with Git based on `config.json` settings:

Manual Mode (Local Git):
- `/do:1-plan`: Prompts for branch creation
- `/do:2-run`: Auto-commits to local branch
- No auto-push

Personal Mode (GitHub Individual):
- `/do:1-plan`: Auto-creates feature branch + auto-push
- `/do:2-run`: Auto-commits + auto-push
- `/do:3-sync`: Suggests PR creation (user choice)

Team Mode (GitHub Team):
- `/do:1-plan`: Auto-creates feature branch + Draft PR
- `/do:2-run`: Auto-commits + auto-push
- `/do:3-sync`: Prepares PR for team review

---

## Works Well With

Skills:
- [do-foundation-core](../SKILL.md) - Parent skill
- [do-foundation-context](../../do-foundation-context/SKILL.md) - Token budget management

Other Modules:
- [spec-first-tdd.md](spec-first-tdd.md) - Detailed SPEC-First TDD process
- [token-optimization.md](token-optimization.md) - /clear execution strategies
- [agents-reference.md](agents-reference.md) - Agent catalog

Agents:
- [workflow-project](agents-reference.md#tier-1-command-processors) - `/do:0-project`
- [workflow-spec](agents-reference.md#tier-1-command-processors) - `/do:1-plan`
- [workflow-tdd](agents-reference.md#tier-1-command-processors) - `/do:2-run`
- [workflow-docs](agents-reference.md#tier-1-command-processors) - `/do:3-sync`

---

Maintained by: Do Team
Status: Production Ready
