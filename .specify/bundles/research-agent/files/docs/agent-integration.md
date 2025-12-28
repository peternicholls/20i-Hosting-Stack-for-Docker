# Speckit Agent Integration Guide

**Version**: 1.0.0  
**Created**: 2025-12-28  
**Purpose**: How speckit agents work together as a cohesive system

---

## Overview

Speckit agents form an integrated workflow automation system for software development. Each agent has a specific role, and they pass data to each other through standardized interfaces.

```
┌─────────────────────────────────────────────────────────────┐
│                     Speckit Agent Family                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  User Request                                               │
│       ↓                                                     │
│  constitution  →  Define project principles                │
│       ↓                                                     │
│  specify       →  Create feature specification             │
│       ↓                                                     │
│  research      →  Answer technical questions               │
│       ↓                                                     │
│  clarify       →  Resolve ambiguities                      │
│       ↓                                                     │
│  plan          →  Create technical implementation plan     │
│       ↓                                                     │
│  tasks         →  Break plan into actionable tasks         │
│       ↓                                                     │
│  implement     →  Execute implementation                   │
│       ↓                                                     │
│  analyze       →  Validate consistency                     │
│       ↓                                                     │
│  checklist     →  Generate QA checklists                   │
│       ↓                                                     │
│  taskstoissues →  Convert to GitHub issues                 │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Agent Reference

### Core Workflow Agents

| Agent | Role | Input | Output | Next Agent |
|-------|------|-------|--------|------------|
| **constitution** | Define principles | User values | constitution.md | specify |
| **specify** | Create spec | Feature description | spec.md | research/clarify |
| **research** | Answer unknowns | Research questions | findings.md, decisions | plan |
| **clarify** | Resolve ambiguity | Spec with [NEEDS CLARIFICATION] | Clarified spec | plan |
| **plan** | Technical design | Feature spec | plan.md | tasks |
| **tasks** | Break into tasks | Implementation plan | tasks.md | implement |
| **implement** | Execute code | Task list | Working code | analyze |

### Supporting Agents

| Agent | Role | When Used | Input | Output |
|-------|------|-----------|-------|--------|
| **analyze** | Validate consistency | After spec/plan/tasks | Specs, plans, tasks | Analysis report |
| **checklist** | QA checklist | Before implementation | Feature context | Checklist file |
| **taskstoissues** | Create GitHub issues | Ready to track work | tasks.md | GitHub issues |

---

## Data Flow

### 1. Constitution → Specify

**Flow**: Project principles inform feature requirements

```yaml
# constitution.md (input)
Principle I: Minimal Dependencies
  - Prefer stdlib over external packages
  - Each dependency must be justified

# spec.md (output)
## Constraints
- Must use standard library where possible
- External dependencies require approval
```

**Data Passed**:
- Principle definitions
- Constraints
- Performance gates
- Values and priorities

### 2. Specify → Research

**Flow**: Spec identifies unknowns needing investigation

```yaml
# spec.md (input)
## Technical Notes
[NEEDS RESEARCH: HTTP client library selection]
Should we use third-party library or stdlib?

# research/R001/research.yaml (output)
metadata:
  id: "R001"
  name: "HTTP Client Library Selection"
key_questions:
  critical:
    - question: "Use library X or stdlib?"
```

**Data Passed**:
- Research questions
- Feature context
- Constitution constraints
- Success criteria

### 3. Research → Plan

**Flow**: Research findings inform technical decisions

```yaml
# decision-log.yaml (input from research)
decisions:
  - id: "D001"
    title: "HTTP Client Selection"
    selection: "stdlib net/http"
    rationale: "Minimal dependencies, full control"

# plan.md (output)
## Technology Stack
- HTTP Client: stdlib net/http
  - Decision: D001
  - Rationale: Aligns with constitution (minimal deps)
```

**Data Passed**:
- Approved decisions
- Library choices
- Pattern recommendations
- Trade-offs

### 4. Research → Specify

**Flow**: Research updates spec with technical details

```yaml
# spec.md (before research)
## Technical Notes
[NEEDS RESEARCH: HTTP client library selection]

# spec.md (after research)
## Technical Notes
- HTTP Client: stdlib net/http (Decision D001)
- Resume support via Range headers
- Connection pooling for parallelism
```

**Data Passed**:
- Implementation constraints
- Technical requirements
- Dependency decisions

### 5. Specify/Clarify → Plan

**Flow**: Spec provides requirements for technical design

```yaml
# spec.md (input)
## Functional Requirements
1. User can resume interrupted downloads
2. System supports parallel downloads
3. Progress displayed in real-time

# plan.md (output)
## Architecture
- Download Manager: Handles resume logic
- Connection Pool: Manages parallel requests
- Progress Tracker: Emits progress events
```

**Data Passed**:
- Functional requirements
- Success criteria
- Constraints
- User scenarios

### 6. Plan → Tasks

**Flow**: Plan broken into actionable development tasks

```yaml
# plan.md (input)
## Phase 1: Core Download Manager
- Implement Range header support
- Create connection pool
- Add progress tracking

# tasks.md (output)
- [ ] T001: Implement Range header requests (4h)
- [ ] T002: Create connection pool (6h)
- [ ] T003: Add progress callback interface (3h)
```

**Data Passed**:
- Implementation phases
- Component designs
- Dependencies
- Estimates

### 7. Tasks → Implement

**Flow**: Tasks executed in sequence

```yaml
# tasks.md (input)
- [ ] T001: Implement Range header requests

# Code (output)
// download/manager.go
func (dm *Manager) Resume(url string, offset int64) error {
    req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
    // ...
}
```

**Data Passed**:
- Task descriptions
- Acceptance criteria
- Dependencies
- Test requirements

### 8. Implement → Analyze

**Flow**: Validate implementation against spec

```yaml
# spec.md (input)
## Success Criteria
- Downloads can be resumed after interruption

# analyze output
✓ Success criterion met
  - Test: TestDownloadResume passes
  - Code: download/manager.go:142-156
```

**Data Passed**:
- Requirements
- Success criteria
- Implementation code
- Test coverage

### 9. Tasks → TasksToIssues

**Flow**: Convert tasks to GitHub issues

```yaml
# tasks.md (input)
- [ ] T001: Implement Range header requests (4h)
  Acceptance: Resume download from byte offset

# GitHub Issue (output)
Title: T001: Implement Range header requests
Labels: feature, size:M
Estimate: 4h
Body: [Task description + acceptance criteria]
```

**Data Passed**:
- Task list
- Estimates
- Dependencies
- Acceptance criteria

---

## Integration Points

### Metadata Tracing

All artifacts share common metadata for traceability:

```yaml
metadata:
  agent: "speckit.research"       # Which agent created this
  version: "1.0.0"                # Schema version
  created: "2025-12-28"           # Creation date
  updated: "2025-12-28"           # Last update

trace:
  created_by: "speckit.specify"   # Originating agent
  updated_by: "speckit.research"  # Last modifier
  related_artifacts:
    - type: "spec"
      path: "specs/001-feature.md"
      relationship: "informs"
    - type: "decision"
      id: "D001"
      relationship: "produces"
```

### Status Synchronization

Agents maintain consistent status across artifacts:

```yaml
# spec.md
status: "in_progress"

# research/R001/research.yaml
status: "completed"

# plan.md
status: "not_started"  # Blocked by research

# When research completes:
# 1. research.yaml: status → "completed"
# 2. plan.md: status → "ready"
# 3. spec.md: updated with findings
```

### File References

Agents use standardized path references:

```yaml
# In any agent output
references:
  - type: "constitution"
    path: "{paths.constitution}"
  
  - type: "spec"
    path: "{feature_specs}/F{id}.yaml"
  
  - type: "decision"
    path: "{decision_log}"
    section: "decisions[id='D001']"
  
  - type: "plan"
    path: "{feature_dir}/plan.md"
```

Paths are resolved from `.specify/config/research.yaml` (or equivalent config).

---

## Handoff Mechanisms

### Explicit Handoffs

Agents declare which agents they can hand off to:

```yaml
# In agent definition
handoffs:
  - label: "Continue to Planning"
    agent: "speckit.plan"
    prompt: "Create implementation plan for this feature"
    send: true  # Send context automatically
  
  - label: "Resolve Ambiguities"
    agent: "speckit.clarify"
    prompt: "Clarify [NEEDS CLARIFICATION] markers"
    send: false  # Ask user first
```

### Implicit Handoffs

Agents can invoke others based on state:

```python
# Pseudo-code in specify agent
if spec_has_research_needs():
    trigger_agent("speckit.research", questions=extract_research_needs())
elif spec_has_clarifications():
    trigger_agent("speckit.clarify", markers=extract_clarifications())
else:
    trigger_agent("speckit.plan", spec=spec_path)
```

### Context Passing

When handing off, agents pass relevant context:

```yaml
# specify → research handoff
context:
  spec_path: "specs/001-feature.md"
  research_needs:
    - topic: "HTTP client selection"
      question: "Library X vs stdlib?"
      context: "Need minimal dependencies per constitution"
  constitution_constraints:
    - "Minimal dependencies"
    - "Performance <500ms"
```

---

## Agent Interfaces

### Standard Input Format

All agents accept:

```yaml
# Via $ARGUMENTS
arguments:
  primary: "Main command or description"
  flags:
    - "--auto"     # Skip confirmations
    - "--json"     # JSON output
    - "--verbose"  # Detailed logging
  context:
    feature_id: "001"
    version: "1.0.0"
    related_files: []
```

### Standard Output Format

All agents produce:

```yaml
# Metadata
metadata:
  agent: "speckit.{name}"
  status: "success | error"
  duration: "3.5s"
  timestamp: "2025-12-28T10:30:00Z"

# Results
results:
  files_created: []
  files_updated: []
  files_read: []
  decisions_made: []

# Next steps
next_steps:
  suggested_agent: "speckit.plan"
  actions:
    - "Review decision D001"
    - "Update spec with findings"
```

### Error Handling

Standard error format:

```yaml
status: "error"
error:
  type: "PrerequisiteError | ValidationError | ExecutionError"
  message: "Human-readable error"
  details:
    missing_file: "specs/001-feature.md"
    suggestion: "Run /speckit.specify first"
  recovery:
    - "Create missing file"
    - "Re-run agent"
```

---

## Workflow Scenarios

### Scenario 1: New Feature (Full Workflow)

```
1. User: "Add user authentication"
   ↓
2. /speckit.specify "Add user authentication"
   → Creates specs/001-auth/spec.md
   → Adds [NEEDS RESEARCH: Auth library selection]
   ↓
3. Auto-triggers /speckit.research R001
   → Researches auth libraries
   → Produces decision D001: Use library X
   → Updates spec.md with choice
   ↓
4. /speckit.clarify
   → Resolves [NEEDS CLARIFICATION: Session vs token?]
   → User chooses: Token-based (JWT)
   → Updates spec.md
   ↓
5. /speckit.plan
   → Creates plan.md with JWT implementation
   → References decision D001
   ↓
6. /speckit.tasks
   → Breaks into T001-T010 tasks
   → Creates tasks.md
   ↓
7. /speckit.implement T001
   → Implements first task
   → Writes code
   ↓
8. /speckit.analyze
   → Validates T001 meets spec requirements
   ↓
9. /speckit.taskstoissues
   → Creates GitHub issues for remaining tasks
```

### Scenario 2: Research Before Spec

```
1. User: "Research state management libraries"
   ↓
2. /speckit.research R005 "State management"
   → Evaluates Redux, Zustand, Jotai
   → Produces decision D005: Use Zustand
   ↓
3. User: "Create feature using Zustand"
   ↓
4. /speckit.specify "Global state management"
   → Creates spec referencing D005
   → No [NEEDS RESEARCH] markers (already done)
   ↓
5. /speckit.plan
   → Plans implementation using Zustand
```

### Scenario 3: Plan Update from New Research

```
1. Existing plan.md uses library X
   ↓
2. /speckit.research R010 "Better alternatives to X"
   → Finds library Y is superior
   → Creates decision D010: Switch to Y
   ↓
3. /speckit.plan --update
   → Reads D010
   → Updates plan.md to use library Y
   → Marks old tasks as obsolete
   ↓
4. /speckit.tasks --regenerate
   → Creates new tasks for library Y
```

---

## Configuration Integration

### Per-Agent Configuration

Each agent can have its own config:

```
.specify/config/
├── research.yaml       # Research agent config
├── plan.yaml          # Plan agent config
├── tasks.yaml         # Tasks agent config
└── implement.yaml     # Implement agent config
```

### Shared Configuration

Common settings in base config:

```yaml
# .specify/config/common.yaml

paths:
  specs: "specs"
  architecture: "docs/architecture"
  constitution: ".specify/memory/constitution.md"

agents:
  auto_handoff: true          # Automatically trigger next agent
  require_approval: ["implement"]  # Require approval for these agents
  parallel_execution: false   # Run agents sequentially

traceability:
  link_artifacts: true        # Cross-link all outputs
  versioning: "git"           # Use git for version tracking
  changelog: "CHANGELOG.md"   # Update changelog
```

---

## Best Practices

### Agent Design

✅ **Single responsibility** - Each agent does one thing well  
✅ **Clear inputs/outputs** - Documented interfaces  
✅ **Idempotent** - Can run multiple times safely  
✅ **Traceable** - Logs all actions and decisions

### Integration

✅ **Validate before handoff** - Ensure output is valid  
✅ **Pass sufficient context** - Next agent has what it needs  
✅ **Handle errors gracefully** - Clear error messages  
✅ **Document dependencies** - What must exist first?

### Data Flow

✅ **One source of truth** - Don't duplicate data  
✅ **Update atomically** - All-or-nothing changes  
✅ **Link artifacts** - Trace requirements → code  
✅ **Version everything** - Track changes over time

---

## Troubleshooting

### "Agent can't find prerequisite file"

Check agent's expected paths match actual structure:

```bash
# In agent config
cat .specify/config/research.yaml | grep paths

# Compare to actual
ls -la specs/
```

### "Handoff fails with missing context"

Ensure previous agent completed successfully:

```bash
# Check status
cat specs/001-feature/spec.md | grep status

# Should be: status: "complete" or "ready"
```

### "Circular dependency between agents"

Review handoff configuration:

```yaml
# specify → research → specify (circular!)
# FIX: specify → research → plan
```

---

## Reference

### Agent Directory

```
.github/agents/
├── speckit.constitution.agent.md
├── speckit.specify.agent.md
├── speckit.research.agent.md
├── speckit.clarify.agent.md
├── speckit.plan.agent.md
├── speckit.tasks.agent.md
├── speckit.implement.agent.md
├── speckit.analyze.agent.md
├── speckit.checklist.agent.md
└── speckit.taskstoissues.agent.md
```

### Configuration Directory

```
.specify/config/
├── common.yaml        # Shared settings
├── research.yaml      # Research agent
├── plan.yaml          # Plan agent
└── tasks.yaml         # Tasks agent
```

### Documentation

- [Research Framework](research-framework.md)
- [Agent Development Guide](agent-development.md)
- [Configuration Reference](configuration-reference.md)

---

**Version History**:
- v1.0.0 (2025-12-28): Initial integration guide
