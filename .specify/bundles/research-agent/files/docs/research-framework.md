# Speckit Research Framework

**Version**: 1.0.0  
**Created**: 2025-12-28  
**Purpose**: Generic, project-agnostic research framework for the speckit.research agent

---

## Overview

The Speckit Research Framework enables structured, repeatable research for any software project. It provides:

✅ **Project-agnostic methodology** - Works with any tech stack or domain  
✅ **Configurable decision criteria** - Adapt to your project's priorities  
✅ **Constitution integration** - Align research with project principles  
✅ **Seamless agent integration** - Works harmoniously with other speckit agents  
✅ **Standardized outputs** - Consistent, traceable research artifacts

---

## Quick Start

### 1. Initialize Research for Your Project

```bash
# Create configuration file
mkdir -p .specify/config
cp examples/research-configs/default.yaml .specify/config/research.yaml

# Edit to match your project
vim .specify/config/research.yaml
```

### 2. Set Up Research Infrastructure

```bash
# Create research directory structure
mkdir -p {YOUR_RESEARCH_ROOT}

# Initialize from templates
cp .specify/templates/research-plan.yaml {YOUR_RESEARCH_ROOT}/
cp .specify/templates/research-template.yaml {YOUR_RESEARCH_ROOT}/
cp .specify/templates/decision-log.yaml {YOUR_RESEARCH_ROOT}/
```

### 3. Create Your First Research Area

```bash
# Copy template
mkdir -p {YOUR_RESEARCH_ROOT}/R001-your-topic
cp {YOUR_RESEARCH_ROOT}/research-template.yaml \
   {YOUR_RESEARCH_ROOT}/R001-your-topic/research.yaml

# Edit to define your research questions
vim {YOUR_RESEARCH_ROOT}/R001-your-topic/research.yaml
```

### 4. Run Research

```bash
# Invoke the research agent
/speckit.research R001

# Or auto-select next priority area
/speckit.research next
```

---

## Configuration Schema

The research agent is configured via `.specify/config/research.yaml`:

### Required Fields

```yaml
metadata:
  project_name: "YourProject"      # Project name
  project_type: "Web App"          # Project type (helps with defaults)
  schema_version: "1.0.0"          # Config schema version

paths:
  research_root: "path/to/research"      # Where research files live
  research_plan: "path/to/plan.yaml"     # Master research plan
  decision_log: "path/to/decisions.yaml" # Decision log
  feature_specs: "path/to/specs"         # Feature specs (for integration)
  architecture_docs: "path/to/arch"      # Architecture docs

decision_criteria:
  weights:
    # Define your criteria (must sum to 100)
    criterion_1: 30
    criterion_2: 25
    criterion_3: 20
    criterion_4: 15
    criterion_5: 10
```

### Optional Fields

```yaml
# Project constitution (if exists)
paths:
  constitution: "path/to/constitution.md"

# Constitution principles
constitution_principles:
  - id: "PRINCIPLE_ID"
    name: "Principle Name"
    description: "What this principle means"
    weight: "critical | high | medium | low"

# Performance gates (if quantitative targets exist)
performance_gates:
  metric_name: "<threshold"
  another_metric: ">minimum"

# Research methodology preferences
methodology:
  confidence_levels: ["LOW", "MEDIUM", "HIGH", "CERTAIN"]
  source_ranking:
    tier_1: ["Official docs", "Source code"]
    tier_2: ["Tech blogs", "Conference talks"]
    tier_3: ["Forums", "Community content"]

# Integration with other speckit agents
integration:
  handoff_agents:
    - agent: "speckit.plan"
      when: "Research informs technical approach"
  update_targets:
    - type: "feature_spec"
      pattern: "{feature_specs}/F{id}.yaml"
      sections: ["implementation"]
```

---

## Research Area Schema

Each research area uses a standard schema (`research.yaml`):

### Metadata

```yaml
metadata:
  id: "R001"                           # Unique ID
  name: "Research Topic"               # Human-readable name
  folder: "R001-research-topic"        # Folder name
  research_type: "library_evaluation"  # Type of research
  target_version: "1.0.0"              # Version this supports
  priority: "P0"                       # P0 (critical) → P3 (low)
  status: "not_started"                # not_started | in_progress | completed | blocked
  estimated_hours: 4                   # Estimate
  actual_hours: 0                      # Actual (update as you go)
  
  # What this research blocks
  blocks:
    - id: "001"
      name: "Feature Name"
  
  # Dependencies on other research
  depends_on: ["R002", "R003"]
```

### Objectives

```yaml
objectives:
  primary:
    - goal: "Primary objective"
      description: "What this achieves"
      success_criteria: "How we know it's complete"
      status: "not_started"
  
  secondary:
    - goal: "Optional objective"
      description: "Nice to have"
      success_criteria: "Completion criteria"
      status: "not_started"
```

### Key Questions

```yaml
key_questions:
  critical:
    - question: "Critical question to answer?"
      context: "Why this matters"
      answer: null              # Fill after research
      confidence: null          # HIGH | MEDIUM | LOW | CERTAIN
      sources: []               # URLs
  
  important:
    - question: "Important question?"
      context: "Background"
      answer: null
      confidence: null
      sources: []
  
  exploratory:
    - question: "Nice to know?"
      context: "Why useful"
      answer: null
      confidence: null
      sources: []
```

### Sources

```yaml
sources:
  libraries:
    - name: "library-name"
      url: "https://github.com/org/repo"
      license: "MIT"
      relevance: "Why evaluate this"
      evaluation_status: "not_evaluated"
      recommendation: null  # use | consider | avoid
  
  documentation:
    - name: "Doc name"
      url: "https://example.com/docs"
      type: "official"
      relevance: "What we'll learn"
      read_status: "not_read"
  
  projects:
    - name: "Reference project"
      url: "https://github.com/org/project"
      relevance: "Patterns to learn"
      study_status: "not_studied"
```

### Constitution Alignment (Optional)

```yaml
constitution_alignment:
  principles:
    - id: "PRINCIPLE_ID"
      name: "Principle Name"
      relevance: "How research relates"
      compliance_notes: null  # Fill after research
  
  gates:
    - gate: "Performance"
      requirement: "<500ms"
      research_relevance: "How research informs this"
      findings: null
```

### Tasks

```yaml
tasks:
  - id: "T001"
    description: "Research task"
    estimated_hours: 1
    actual_hours: 0
    status: "not_started"
    completed_date: null
```

### Findings

```yaml
findings:
  discoveries:
    - id: "D001"
      date: "2025-12-28"
      title: "Finding title"
      description: "What was discovered"
      significance: "high | medium | low"
      source: "Where from"
      implications: "What this means"
  
  comparisons:
    - category: "Comparison category"
      options:
        - name: "Option A"
          pros: ["Pro 1", "Pro 2"]
          cons: ["Con 1", "Con 2"]
          score: 85
        - name: "Option B"
          pros: ["Pro 1"]
          cons: ["Con 1", "Con 2"]
          score: 65
```

### Decisions

```yaml
decisions:
  - id: "DEC001"
    title: "Decision title"
    question: "What question does this answer?"
    options_considered:
      - name: "Option A"
      - name: "Option B"
    chosen: "Option A"
    rationale: "Why chosen"
    confidence: "high"
    impacts:
      - document: "Feature spec"
        path: "specs/F001.yaml"
        change: "What to update"
```

---

## Research Workflow

### 1. Selection

The agent selects research based on:

```
Priority: P0 > P1 > P2 > P3
    ↓
Version: Earlier version first
    ↓
ID: Lower R-number first
```

**Manual selection**:
```bash
/speckit.research R001    # Specific area
```

**Auto-selection**:
```bash
/speckit.research next                # Next priority
/speckit.research version 1.0.0       # Next for version
/speckit.research type library_evaluation  # Next of type
```

### 2. Research Brief

Agent presents a brief before starting:

```markdown
## Research Brief: [Topic]

**ID**: R001
**Priority**: P0 (Critical)
**Estimated**: 4 hours
**Blocks**: Feature 001

### Critical Questions
- [ ] Q1: Question 1?
- [ ] Q2: Question 2?

### Sources
- Library A, Library B
- Documentation X, Y

Ready to proceed? (yes/skip)
```

### 3. Execution

Agent performs research tasks:

- **Library evaluation**: GitHub metadata, docs, dependencies
- **Pattern research**: Search sources, extract examples
- **Documentation review**: Read and synthesize
- **Benchmarks**: Compare performance

### 4. Recording

Agent updates:
- **research.yaml**: Structured answers
- **findings.md**: Detailed notes
- **decision-log.yaml**: Decisions

### 5. Completion

Agent reports:

```markdown
## Complete: R001

**Time**: 3.5h (estimated: 4h)
**Questions**: 4/4 answered
**Confidence**: HIGH
**Decision**: D001

### Next Steps
1. Review decision D001
2. Update spec F001
3. Continue to R002
```

---

## Decision Framework

### Scoring Options

1. **Define criteria** in `.specify/config/research.yaml`:

```yaml
decision_criteria:
  weights:
    security: 30
    performance: 25
    maintainability: 20
    cost: 15
    usability: 10
```

2. **Score each option** (0-100 scale):

| Option | Security | Performance | Maintain | Cost | Usability | **Total** |
|--------|----------|-------------|----------|------|-----------|-----------|
| A      | 90       | 80          | 70       | 60   | 90        | **79**    |
| B      | 70       | 90          | 80       | 80   | 70        | **78**    |

3. **Apply thresholds**:
- Pass: ≥60
- Recommend: ≥75

### Recording Decisions

```yaml
# In decision-log.yaml
decisions:
  - id: "D001"
    title: "Decision Title"
    category: "LIBRARY_SELECTION"
    research_id: "R001"
    status: "PROPOSED"
    
    options_considered:
      - name: "Option A"
        score: 79
      - name: "Option B"
        score: 78
    
    decision:
      selection: "Option A"
      rationale: "Slightly better security, meets all criteria"
      trade_offs_accepted:
        - "Slightly slower performance"
      fallback_plan: "Can switch to Option B if performance issues arise"
```

---

## Integration with Speckit Agents

### Agent Flow

```
User request
    ↓
speckit.specify  → Create spec
    ↓
[Needs research?]
    ↓ YES
speckit.research → Answer questions
    ↓
[Update spec with findings]
    ↓
speckit.clarify  → Resolve ambiguities
    ↓
speckit.plan     → Technical plan
    ↓
speckit.tasks    → Task breakdown
    ↓
speckit.implement → Execute
```

### Data Flow

**specify → research**:
- Spec contains `[NEEDS RESEARCH: topic]`
- Research agent creates research areas

**research → plan**:
- Decision log has approved decisions
- Plan reads decisions for tech choices

**research → specify**:
- Research updates spec's `implementation` section
- Adds library/dependency details

**research → analyze**:
- Analyze validates all questions answered
- Checks constitution alignment

### Shared Metadata

All agents use consistent traceability:

```yaml
metadata:
  agent: "speckit.research"
  version: "1.0.0"
  created: "2025-12-28"
  
trace:
  created_by: "speckit.specify"
  updated_by: "speckit.research"
  related_artifacts:
    - type: "spec"
      path: "specs/001-feature.md"
    - type: "decision"
      id: "D001"
```

---

## Adaptation Examples

### Web Application

```yaml
# .specify/config/research.yaml

decision_criteria:
  weights:
    security: 30           # Critical for web
    scalability: 25        # Handle growth
    developer_experience: 20
    cost: 15               # Cloud costs
    maintainability: 10

performance_gates:
  page_load: "<2s"
  api_response: "<200ms"
  time_to_interactive: "<3s"
```

### Library/SDK

```yaml
decision_criteria:
  weights:
    api_design: 30                # Most important
    backward_compatibility: 25
    documentation: 20
    dependency_footprint: 15
    performance: 10

research_types:
  - api_design_research
  - ecosystem_compatibility
  - migration_path_analysis
```

### Infrastructure/DevOps

```yaml
decision_criteria:
  weights:
    reliability: 30
    operational_complexity: 25
    cost: 20
    security: 15
    vendor_lock_in: 10

performance_gates:
  deployment_time: "<5min"
  recovery_time: "<15min"
  uptime: ">99.9%"
```

---

## Migration Guide

### From Hardcoded to Configurable

**Step 1**: Create configuration

```bash
# Generate config from current setup
cp .specify/templates/research-config-template.yaml \
   .specify/config/research.yaml

# Edit paths to match your structure
vim .specify/config/research.yaml
```

**Step 2**: Update agent

Agent now loads config at startup:

```markdown
## Prerequisites

1. Load configuration:
   ```bash
   cat .specify/config/research.yaml
   ```
2. Extract paths, criteria, gates
3. Proceed with research
```

**Step 3**: Update research areas

Replace hardcoded values:

```yaml
# Before (hardcoded)
blocks_features:
  - "F011"  # SafeDownload-specific

# After (generic)
blocks:
  - id: "011"
    name: "Core Download Engine"
```

**Step 4**: Test

```bash
# Validate configuration
speckit.research validate

# Test with one area
speckit.research R001
```

---

## Best Practices

### Configuration

✅ **Use relative paths** - Makes project portable  
✅ **Document criteria weights** - Explain why each weight chosen  
✅ **Version your config** - Track changes to methodology  
✅ **Keep constitution optional** - Not all projects have one

### Research Areas

✅ **One topic per area** - Keep focused  
✅ **Clear success criteria** - Know when done  
✅ **Cite all sources** - Include URLs and dates  
✅ **Record confidence** - Be honest about uncertainty

### Decisions

✅ **Document trade-offs** - What are you accepting?  
✅ **Include fallback** - What if this doesn't work?  
✅ **Get approval** - Don't implement without review  
✅ **Track impacts** - What needs updating?

### Integration

✅ **Update promptly** - Don't let research go stale  
✅ **Link artifacts** - Trace research → spec → plan → tasks  
✅ **Communicate** - Tell team what was learned  
✅ **Archive old research** - When decisions change

---

## Troubleshooting

### "Configuration not found"

```bash
# Create from template
cp .specify/templates/research-config-template.yaml \
   .specify/config/research.yaml
```

### "Research plan not found"

```bash
# Initialize research infrastructure
mkdir -p {YOUR_RESEARCH_ROOT}
cp .specify/templates/research-plan.yaml {YOUR_RESEARCH_ROOT}/
```

### "Can't select next area"

Check research-plan.yaml has areas defined:

```yaml
research_areas:
  - id: "R001"
    status: "not_started"
    priority: "P0"
```

### "Constitution alignment failing"

Make config paths.constitution optional:

```yaml
paths:
  constitution: null  # No constitution for this project
```

Agent will skip constitution checks.

---

## Reference

### File Structure

```
project/
├── .specify/
│   ├── config/
│   │   └── research.yaml       # Configuration
│   └── templates/
│       ├── research-plan.yaml
│       └── research-template.yaml
├── {research_root}/
│   ├── research-plan.yaml      # Master plan
│   ├── decision-log.yaml       # Decisions
│   └── R001-topic/
│       ├── research.yaml       # Research area
│       └── findings.md         # Findings
└── {feature_specs}/
    └── F001.yaml               # Integration target
```

### Configuration Reference

See [examples/research-configs/](../examples/research-configs/) for:
- `default.yaml` - Starter template
- `web-app.yaml` - Web application example
- `library.yaml` - Library/SDK example
- `infrastructure.yaml` - DevOps/infra example

### Agent Reference

- [speckit.research.agent.md](../../.github/agents/speckit.research.agent.md) - Full agent definition
- [speckit.research.prompt.md](../../.github/prompts/speckit.research.prompt.md) - Activation prompt

---

## Support

For issues or questions:
1. Check [Troubleshooting](#troubleshooting)
2. Review [examples/research-configs/](../examples/research-configs/)
3. Open an issue in the Speckit repository

---

**Version History**:
- v1.0.0 (2025-12-28): Initial generic framework
