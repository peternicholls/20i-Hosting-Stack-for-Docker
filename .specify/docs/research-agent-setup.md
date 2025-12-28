# Research Agent Setup Guide

This guide explains how to set up and use the research agent in any project.

## Quick Start

### For New Projects

1. **Initialize the research agent**:
   ```bash
   .specify/scripts/bash/research.sh init
   ```

2. **Follow the interactive prompts** to create `.specify/config/research.yaml`

3. **Review and customize** the generated configuration

4. **Start researching**:
   ```bash
   @speckit.research "Evaluate React vs Vue for frontend"
   ```

### For Existing Projects

If you're migrating from an older version or setting up on an existing project:

1. **Choose a template** from `.specify/examples/research-configs/`:
   - `default.yaml` - Generic balanced configuration
   - `web-app.yaml` - Web applications
   - `library.yaml` - Libraries and SDKs
   - `infrastructure.yaml` - Infrastructure/DevOps

2. **Copy and customize**:
   ```bash
   cp .specify/examples/research-configs/web-app.yaml .specify/config/research.yaml
   # Edit to match your project
   ```

3. **Create research directory structure**:
   ```bash
   mkdir -p dev/research  # Or your configured path
   ```

## Configuration

### Essential Settings

The minimum required configuration in `.specify/config/research.yaml`:

```yaml
metadata:
  project_name: "YourProject"
  project_type: "web-app"  # or library, infrastructure, cli, etc.
  schema_version: "1.0.0"

paths:
  research_root: "research"  # Where research files live
  research_plan: "research/research-plan.yaml"
  research_template: ".specify/templates/research-template.yaml"
  research_feedback_template: ".specify/templates/research-feedback-template.md"
  decision_log: "research/decision-log.yaml"

decision_criteria:
  weights:
    quality: 30
    performance: 25
    maintainability: 20
    # ... add criteria relevant to your project
```

### Optional Settings

```yaml
paths:
  constitution: "docs/principles.md"  # If you have guiding principles
  feature_specs: "specs"
  architecture_docs: "architecture"

performance_gates:
  page_load: "<2s"        # Web app example
  bundle_size: "<500kb"   # Library example
  deployment_time: "<5min" # Infrastructure example

constitution_principles:
  - id: "SECURITY"
    name: "Security First"
    description: "Security is non-negotiable"
    weight: "critical"
```

## Project Types

### Web Application

Priorities: Security, Scalability, User Experience

```yaml
decision_criteria:
  weights:
    security: 30
    scalability: 25
    user_experience: 20
    performance: 15
    maintainability: 10

performance_gates:
  page_load: "<2s"
  time_to_interactive: "<3s"
  lighthouse_score: ">90"
```

**Use cases**:
- Frontend framework selection
- API design patterns
- Authentication approaches
- State management solutions

### Library/SDK

Priorities: API Design, Backward Compatibility, Documentation

```yaml
decision_criteria:
  weights:
    api_design: 30
    backward_compatibility: 25
    documentation: 20
    performance: 15
    bundle_size: 10

performance_gates:
  bundle_size: "<50kb"
  tree_shakeable: true
```

**Use cases**:
- Public API design
- Breaking change assessment
- Dependency analysis
- Build tool selection

### Infrastructure/DevOps

Priorities: Reliability, Operational Simplicity, Cost

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

**Use cases**:
- Platform selection (AWS vs GCP vs Azure)
- Monitoring solutions
- CI/CD pipeline design
- Container orchestration

### CLI Tool

Priorities: Performance, User Experience, Minimal Dependencies

```yaml
decision_criteria:
  weights:
    constitution_alignment: 25  # If you have guiding principles
    performance: 20
    maintainability: 20
    developer_experience: 15
    dependency_footprint: 10

performance_gates:
  startup_time: "<500ms"
  response_time: "<100ms"
```

**Use cases**:
- Library selection
- TUI framework choice
- Build optimization
- Cross-platform compatibility

## Usage Patterns

### Research Workflow

1. **Define research area**:
   ```yaml
   # In research-plan.yaml or as a new research.yaml
   id: "R01"
   name: "State Management Research"
   description: "Evaluate state management solutions"
   type: "library_evaluation"
   ```

2. **Execute research**:
   ```bash
   @speckit.research R01
   # or
   @speckit.research "Compare Redux vs MobX vs Zustand"
   ```

3. **Review findings**:
   - Agent creates `research/<area>/research.yaml`
   - Detailed findings in `research/<area>/findings.md`
   - Decision recorded in `decision-log.yaml`

4. **Apply decision**:
   - Update feature specs
   - Update architecture docs
   - Create implementation tasks

### Common Research Types

#### Library Evaluation

Research Question: "Which library should we use for X?"

Agent will:
- Fetch repository metadata (stars, license, activity)
- Review documentation quality
- Check dependency tree
- Analyze maintenance status
- Score against your criteria
- Provide recommendation

Example:
```bash
@speckit.research "Evaluate date libraries: date-fns vs dayjs vs luxon"
```

#### Pattern Research

Research Question: "What's the best approach for X?"

Agent will:
- Search for authoritative sources
- Study reference implementations
- Extract code examples
- Document trade-offs
- Recommend pattern

Example:
```bash
@speckit.research "Research error handling patterns in React"
```

#### Architecture Design

Research Question: "How should we architect X?"

Agent will:
- Identify requirements and constraints
- Research existing solutions
- Compare architectural patterns
- Document trade-offs
- Recommend approach

Example:
```bash
@speckit.research "Design microservices communication architecture"
```

### Helper Scripts

```bash
# Show research status
.specify/scripts/bash/research.sh status

# Select a research area
.specify/scripts/bash/research.sh select

# Validate research configuration
.specify/scripts/bash/research.sh validate

# Initialize new project
.specify/scripts/bash/research.sh init
```

## Integration with SpecKit

The research agent is part of the 10-agent SpecKit family:

```
┌──────────────┐
│ constitution │  ← Defines project principles
└──────┬───────┘
       ↓
┌──────────────┐
│   specify    │  ← Creates feature specifications
└──────┬───────┘
       ↓
┌──────────────┐
│   clarify    │  ← Answers technical questions
└──────┬───────┘
       ↓
┌──────────────┐
│   research   │  ← YOU ARE HERE - Deep research on blockers
└──────┬───────┘
       ↓
┌──────────────┐
│     plan     │  ← Creates sprint plans
└──────┬───────┘
       ↓
┌──────────────┐
│    tasks     │  ← Breaks down into tasks
└──────┬───────┘
       ↓
┌──────────────┐
│  implement   │  ← Writes code
└──────┬───────┘
       ↓
┌──────────────┐
│   analyze    │  ← Reviews architecture
└──────┬───────┘
       ↓
┌──────────────┐
│  checklist   │  ← Validates compliance
└──────┬───────┘
       ↓
┌──────────────┐
│taskstoissues │  ← Creates GitHub issues
└──────────────┘
```

### Data Flow

**Research provides**:
- Technical feasibility answers
- Library recommendations
- Architecture patterns
- Implementation constraints
- Performance considerations

**Research consumes**:
- Constitution principles (optional)
- Feature specifications
- Sprint blockers
- Technical questions

## Troubleshooting

### Configuration Not Found

```
Error: .specify/config/research.yaml not found
```

**Solution**: Run `.specify/scripts/bash/research.sh init`

### Missing yq Dependency

```
Error: yq command not found
```

**Solution**: 
```bash
brew install yq  # macOS
# or
apt-get install yq  # Debian/Ubuntu
```

### Constitution File Not Found

If you specified a constitution path but file doesn't exist:

**Solution**: Either create the file or set `constitution: null` in config

### Research Directory Not Found

```
Error: Research directory not found
```

**Solution**: Create the directory specified in `paths.research_root`

## Best Practices

### 1. Define Clear Research Questions

❌ Bad: "Research state management"
✅ Good: "Compare Redux vs MobX for large-scale React app with focus on developer experience and bundle size"

### 2. Use Constitution Principles

If your project has guiding principles, document them and reference in config:

```yaml
constitution_principles:
  - id: "SIMPLICITY"
    name: "Simplicity First"
    description: "Prefer simple solutions over complex ones"
    weight: "critical"
```

### 3. Set Realistic Performance Gates

Don't set gates you can't measure or that don't matter:

❌ Bad: `api_response: "<1ms"` (unrealistic)
✅ Good: `api_p95: "<500ms"` (measurable, realistic)

### 4. Version Your Configuration

Track changes to your research config:

```bash
git add .specify/config/research.yaml
git commit -m "chore: update research criteria weights"
```

### 5. Document Decisions

Always record decisions in `decision-log.yaml`:

```yaml
decisions:
  - id: "D001"
    date: "2025-01-15"
    decision: "Use Zustand for state management"
    rationale: "Smallest bundle, simplest API, good TypeScript support"
    alternatives_considered: ["Redux", "MobX", "Jotai"]
    trade_offs: "Less ecosystem than Redux, fewer middleware options"
```

## Examples

### Example 1: Web Framework Selection

```bash
@speckit.research "Compare Next.js vs Remix vs SvelteKit for our marketing site"
```

Agent will:
1. Load your web-app config (security: 30%, scalability: 25%, UX: 20%)
2. Research each framework against criteria
3. Test performance gates (page_load, lighthouse_score)
4. Provide scored comparison
5. Recommend based on weights

### Example 2: Database Choice

```bash
@speckit.research "PostgreSQL vs MySQL vs MongoDB for user data with 10M+ records"
```

Agent will:
1. Research scalability patterns
2. Compare query performance
3. Evaluate operational complexity
4. Check cost implications
5. Recommend with rationale

### Example 3: Authentication Approach

```bash
@speckit.research "Authentication strategy: JWT vs session cookies vs OAuth2 only"
```

Agent will:
1. Research security implications (30% weight for web apps)
2. Compare UX impact
3. Evaluate implementation complexity
4. Check scalability
5. Recommend approach

## Advanced Topics

### Custom Research Types

Define project-specific research types in your config:

```yaml
research_types:
  - type: "compliance_research"
    description: "Research compliance requirements"
    deliverables:
      - "Compliance checklist"
      - "Required documentation"
      - "Implementation tasks"
```

### Multi-Criteria Optimization

For complex decisions, use multiple weighted criteria:

```yaml
decision_criteria:
  weights:
    security: 25
    performance: 20
    cost: 20
    developer_experience: 15
    vendor_lock_in: 10
    compliance: 10
```

### Integration with CI/CD

Validate research config in CI:

```yaml
# .github/workflows/validate.yml
- name: Validate Research Config
  run: |
    .specify/scripts/bash/research.sh validate
```

## Further Reading

- [Research Framework Documentation](.specify/docs/research-framework.md)
- [Agent Integration Guide](.specify/docs/agent-integration.md)
- [Configuration Examples](.specify/examples/research-configs/)
- [Generalization Implementation Plan](dev/research/SPECKIT_RESEARCH_GENERALIZATION.md)

## Support

If you encounter issues:

1. Check this setup guide
2. Review example configurations
3. Validate your config with `research.sh validate`
4. Check [agent integration docs](.specify/docs/agent-integration.md)
5. Review the generalization plan for implementation details
