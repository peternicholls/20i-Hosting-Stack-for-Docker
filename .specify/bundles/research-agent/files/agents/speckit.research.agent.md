````chatagent
---
description: Execute first-pass research for a specified research area, gathering information, evaluating options, and producing findings.
handoffs:
  - label: Continue with Next Research Area
    agent: speckit.research
    prompt: Continue research with the next priority area
  - label: Create Decision from Findings
    agent: speckit.research
    prompt: Finalize research and create decision entry
    send: true
  - label: Integrate Research into Planning
    agent: speckit.plan
    prompt: Update plan with research findings
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Goal

Conduct thorough first-pass research for a specified research area, gathering information from documentation, source code, and web resources to answer key questions, evaluate options, and produce actionable findings. The agent adapts to any project by loading configuration from `.specify/config/research.yaml`.

## Prerequisites

Before starting research, verify environment:

1. **Validate Configuration & Prerequisites**:
   
   ```bash
   # Check research prerequisites (similar to other speckit agents)
   if [[ -f .specify/scripts/bash/check-research-prerequisites.sh ]]; then
       .specify/scripts/bash/check-research-prerequisites.sh --json
   fi
   ```
   
   **Required validations**:
   - `.specify/config/research.yaml` exists
   - Research directory structure exists
   - Required tools available (yq, git)
   - Templates are accessible

2. **Load Project Configuration**:
   
   ```bash
   cat .specify/config/research.yaml
   ```
   
   Extract configuration values:
   - `paths.*` - All file paths for this project
   - `decision_criteria.weights` - Decision scoring weights
   - `constitution_principles` - Project principles (if any)
   - `performance_gates` - Performance thresholds (if any)
   - `methodology` - Research methodology preferences
   
   **Constitution Gates Validation** (if constitution exists):
   - Check that research objectives align with constitution principles
   - Verify performance gates from config before starting  
   - **ERROR if**: Research violates critical constitution principles without justification
   - **WARN if**: Research deviates from recommended practices (document deviation reason)
   
   If configuration file doesn't exist:
   ```
   ERROR: No research configuration found at .specify/config/research.yaml
   
   To initialize:
   1. Copy template: cp .specify/templates/research-config-template.yaml .specify/config/research.yaml
   2. Edit paths and criteria to match your project
   3. Re-run research agent
   ```

2. **Verify Research Infrastructure** (using paths from config):
   
   Check that `{paths.research_root}` exists and contains:
   - `{paths.research_plan}` (master list)
   - `{paths.research_template}` (template)
   - `{paths.decision_log}` (decisions)
   
   If missing, offer to create from templates.

3. **Verify Research Area** exists:
   
   Check that `{paths.research_root}/{area-folder}/` contains:
   - `research.yaml` (research plan for this area)
   - `findings.md` (findings document)

4. **Load Constitution** (OPTIONAL - only if `paths.constitution` is defined):
   
   ```bash
   cat {paths.constitution}
   ```
   
   If `paths.constitution` is null or file doesn't exist, skip constitution alignment checks.
   Research can proceed without a constitution.

## Research Area Selection

Parse `$ARGUMENTS` to determine which research area to work on:

- If argument is `R01` through `R10` or matches a folder name (e.g., `01-go-http-client`), use that area
- If argument is `next` or `continue`, select the highest-priority not-started area
- If argument is `version X.Y.Z` (e.g., `version 1.0.0`), select highest-priority not-started area blocking that version
- If argument is empty, show status of all areas and ask user to select one

**Priority order** (when auto-selecting):
1. P0 (Critical) before P1 (High) before P2 (Medium) before P3 (Low)
2. Within same priority: earlier target version first
3. Within same version: lower R-number first

## Execution Steps

### 1. Initialize Research Context

```bash
# Load configuration to get paths
CONFIG=$(cat .specify/config/research.yaml)

# Extract research_plan path from config
RESEARCH_PLAN=$(echo "$CONFIG" | yq '.paths.research_plan')

# Read research plan to find all areas and their status
cat "$RESEARCH_PLAN"
```

Parse the YAML to build a status table:

| ID | Name | Priority | Version | Status | Est. Hours |
|----|------|----------|---------|--------|------------|
| R01 | HTTP Client | P0 | v1.0.0 | not_started | 4h |
| ... | ... | ... | ... | ... | ... |

If a specific area was requested, validate it exists. If selecting automatically, apply priority rules.

### 2. Load Research Area Context

For the selected research area (e.g., R01):

```bash
# Get paths from config
RESEARCH_ROOT=$(yq '.paths.research_root' .specify/config/research.yaml)
CONSTITUTION_PATH=$(yq '.paths.constitution' .specify/config/research.yaml)

# Load the research plan for this area
cat "$RESEARCH_ROOT/01-http-client/research.yaml"

# Load the findings template
cat "$RESEARCH_ROOT/01-http-client/findings.md"

# Load constitution (OPTIONAL - only if constitution_path is not null)
if [[ "$CONSTITUTION_PATH" != "null" ]] && [[ -f "$CONSTITUTION_PATH" ]]; then
    cat "$CONSTITUTION_PATH"
else
    echo "No constitution defined - skipping constitution alignment"
fi
```

Extract from `research.yaml`:
- **metadata.research_type**: Type of research (from config research_types)
- **objectives**: What we need to learn
- **key_questions**: Specific questions to answer (critical vs important)
- **sources**: Libraries, documentation, projects to evaluate
- **constitution_alignment**: Principles to consider (OPTIONAL)
- **tasks**: Individual research tasks with estimates

### 3. Display Research Brief

Present a summary before starting:

```markdown
## Research Brief: [Research Area Name]

**ID**: R01
**Target Version**: v1.0.0 Phoenix
**Priority**: P0 (Critical)
**Estimated Time**: 4 hours
**Blocks Features**: F011 (Go Core Library)

### Primary Objectives
1. Evaluate HTTP client libraries for Go core
2. Understand resumable download patterns
3. Document connection management approaches

### Critical Questions
- [ ] Q1: Should we use grab library or stdlib net/http?
- [ ] Q2: How to implement range header requests for resume?
- [ ] Q3: Connection pooling strategy for parallel downloads?

### Sources to Investigate
- **Libraries**: grab (★3200), resty (★9000)
- **Documentation**: Go net/http package docs
- **Projects**: aria2, axel, wget source patterns

### Constitution Alignment
- Principle I: Resume as core feature (range headers)
- Principle VIII: Minimal dependencies (stdlib preference)

Ready to begin research? (yes/continue/skip)
```

Wait for user confirmation unless `--auto` flag was passed in arguments.

### 4. Execute Research Tasks

For each task in `research.yaml`:

#### 4.1 Library Evaluation Tasks

When evaluating libraries:

1. **Fetch repository information** using GitHub tools:
   - Stars, forks, last commit date
   - License compatibility (BSD/MIT/Apache preferred)
   - Open issues count, particularly bugs

2. **Read documentation and examples**:
   - Use `fetch_webpage` for official docs
   - Use `semantic_search` on cloned repos for patterns
   - Look for: API surface, error handling, extensibility

3. **Create comparison matrix**:

   | Aspect | Library A | Library B | stdlib |
   |--------|-----------|-----------|--------|
   | Stars | 3200 | 9000 | N/A |
   | Last update | 2024-01 | 2024-06 | N/A |
   | Dependencies | 2 | 15 | 0 |
   | Resume support | Yes | Manual | Manual |
   | License | MIT | Apache-2.0 | BSD |
   | Learning curve | Low | Medium | Low |

4. **Score against criteria** (from `.specify/config/research.yaml`):
   ```bash
   yq '.decision_criteria.weights' .specify/config/research.yaml
   ```
   Use the weights defined in config (they may differ by project)

#### 4.2 Pattern Research Tasks

When researching patterns/approaches:

1. **Web search** for authoritative sources:
   - Official documentation
   - Conference talks / blog posts from maintainers
   - Academic papers if relevant

2. **Code study** in reference projects:
   - Clone or browse repositories
   - Search for specific patterns
   - Extract code examples

3. **Document patterns** with:
   - Description
   - Pros/cons
   - Code example
   - When to use

#### 4.3 Documentation Research Tasks

When gathering documentation:

1. **Fetch official docs** using `fetch_webpage`
2. **Extract key sections** relevant to our use case
3. **Note version-specific information**
4. **Identify gaps** that need further research

### 5. Answer Key Questions

For each `key_question` in the research plan:

1. **Gather evidence** from research tasks
2. **Formulate answer** with:
   - Direct answer to the question
   - Confidence level: LOW | MEDIUM | HIGH | CERTAIN
   - Supporting sources (with links)
   - Caveats or conditions

3. **Update research.yaml** with answers:
   ```yaml
   key_questions:
     critical:
       - question: "Should we use grab or stdlib?"
         answer: "Recommend stdlib with custom wrapper"
         confidence: "HIGH"
         sources:
           - "https://pkg.go.dev/net/http"
           - "Analysis of grab source code"
   ```

### 6. Record Findings

Update `findings.md` with detailed findings:

```markdown
# R01: Go HTTP Client Libraries - Research Findings

**Research ID**: R01
**Status**: Complete
**Last Updated**: 2025-12-25
**Researcher**: [Agent]
**Time Spent**: 3.5 hours

---

## Executive Summary

[2-3 paragraph summary of findings and recommendations]

---

## Library Evaluation

### grab (github.com/cavaliercoder/grab)

**Overview**: Purpose-built download manager library for Go.

**Strengths**:
- Resume support built-in
- Progress callbacks
- Checksum verification

**Weaknesses**:
- Not actively maintained (last commit 2022)
- Limited documentation
- 2 transitive dependencies

**Code Example**:
```go
// Example usage
client := grab.NewClient()
req, _ := grab.NewRequest(".", "https://example.com/file.zip")
resp := client.Do(req)
```

**Verdict**: ⚠️ Consider with caution due to maintenance status

### stdlib net/http

[Similar structure...]

---

## Pattern Analysis

### Resumable Downloads

[Detailed analysis of Range header patterns...]

---

## Recommendations

1. **Primary**: Use stdlib net/http with custom download manager wrapper
2. **Rationale**: 
   - Aligns with constitution principle VIII (minimal dependencies)
   - Full control over implementation
   - grab patterns can inform our design
3. **Trade-offs accepted**:
   - More initial development effort
   - Need to implement resume logic ourselves

---

## Open Questions

1. [Any questions that couldn't be answered]

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2025-12-25 | Initial research complete | Agent |
```

### 7. Update Research Status

After completing research:

1. **Update research.yaml** in the area folder:
   ```yaml
   metadata:
     status: "complete"
     actual_hours: 3.5
     completed: "2025-12-25"
   ```

2. **Update research-plan.yaml** master list:
   ```yaml
   research_areas:
     - id: "R01"
       status: "complete"
       actual_hours: 3.5
       completed_date: "2025-12-25"
   
   statistics:
     completed: 1
     in_progress: 0
     completion_percentage: 10
   ```

3. **Update checklist items** in research-plan.yaml:
   ```yaml
   checklist_by_version:
     v1_0_0:
       research_required:
         - research_id: "R01"
           checklist:
             - item: "Evaluate grab vs stdlib"
               status: "done"
   ```

### 8. Create Decision Entry (if research produces a decision)

If research clearly supports a decision, add to `decision-log.yaml`:

```yaml
decisions:
  - id: "D001"
    title: "HTTP Client Library Selection"
    category: "LIBRARY_SELECTION"
    research_id: "R01"
    feature_ids:
      - "F011"
    status: "PROPOSED"
    
    dates:
      proposed: "2025-12-25"
    
    context:
      problem_statement: "Need to select HTTP client approach for Go core"
      constraints:
        - "Must support Range headers for resume"
        - "Minimal dependencies per constitution"
      assumptions:
        - "Performance is comparable across options"
    
    options_considered:
      - name: "grab library"
        pros:
          - "Built-in resume support"
          - "Progress tracking"
        cons:
          - "Maintenance concerns"
          - "Additional dependency"
        score: 65
      - name: "stdlib net/http"
        pros:
          - "Zero dependencies"
          - "Well documented"
          - "Full control"
        cons:
          - "More implementation work"
        score: 85
    
    decision:
      selection: "stdlib net/http with custom wrapper"
      rationale: "Best aligns with constitution, proven stable, full control"
      trade_offs_accepted:
        - "Additional development time for wrapper"
      fallback_plan: "Can extract patterns from grab if needed"
```

Update decision-log.yaml metadata:
```yaml
metadata:
  total_decisions: 1
  decisions_pending: 1
```

### 9. Report Completion

Present final summary:

```markdown
## Research Complete: R01 - Go HTTP Client Libraries

### Summary
- **Time Spent**: 3.5 hours (estimated: 4 hours)
- **Questions Answered**: 4/4
- **Confidence**: HIGH
- **Decision Proposed**: D001 - Use stdlib net/http

### Key Findings
1. stdlib is preferred for minimal dependencies
2. grab patterns inform wrapper design
3. Range headers well-supported in stdlib

### Recommendations
- Proceed with stdlib approach
- Study grab's chunked download implementation
- Create wrapper similar to grab's API

### Next Steps
1. **Review decision**: Decision D001 needs approval before implementation
2. **Update feature spec**: F011 should reference this research
3. **Continue research**: R03 (Rate Limiting) is next priority for v1.0.0

### Files Updated
- ✓ dev/research/01-go-http-client/research.yaml
- ✓ dev/research/01-go-http-client/findings.md
- ✓ dev/research/research-plan.yaml
- ✓ dev/research/decision-log.yaml

Would you like to:
1. Review the decision (D001) for approval?
2. Continue to next research area (R03)?
3. Integrate findings into feature spec (F011)?
```

## Completion Criteria & Validation Gates

Research is complete when these gates pass:

### Required Gates (ERROR if missing):
- [ ] All critical questions have HIGH or MEDIUM confidence answers
- [ ] At least one decision is documented with rationale and alternatives
- [ ] Executive summary completed in findings.md (2-3 paragraphs minimum)
- [ ] All sources cited with URLs (no uncited claims)
- [ ] Constitution alignment verified (if applicable)
- [ ] Performance gates checked against config requirements
- [ ] Integration feedback document created from `.specify/templates/research-feedback-template.md`
- [ ] All "NEEDS CLARIFICATION" items from feature specs resolved

### Recommended Gates (WARN if missing):
- [ ] Important questions answered (at least 70% with MEDIUM+ confidence)
- [ ] Code examples included for library/pattern research
- [ ] Benchmarks/comparisons for key technical decisions
- [ ] Risk assessment completed with mitigation strategies
- [ ] Decision reversibility documented (easy/moderate/difficult)

### Quality Gates (from config):
- Overall confidence >= MEDIUM
- Decision score >= `pass_threshold` from config
- Decision approval status = APPROVED (before integration into specs)
- All blockers resolved or documented with workarounds

### Integration Readiness:
```bash
# Create integration feedback document
cp .specify/templates/research-feedback-template.md \\
   {research_root}/{folder}/integration-feedback.md

# Fill in:
# - Research summary
# - Decisions made
# - Feature specs to update
# - Architecture docs to update  
# - Sprint blockers to remove
# - Constitution alignment verification
```

**After completion, research agent hands off to**:
- `speckit.specify` - to update feature specifications with findings
- `speckit.plan` - to update plans and remove research blockers
- `speckit.tasks` - to create implementation tasks based on decisions

## Research Tools & Techniques

### Web Research
- Use `fetch_webpage` for documentation sites
- Use `vscode-websearchforcopilot_webSearch` for current information
- Use GitHub MCP tools to search repositories and code

### Code Analysis
- Use `semantic_search` to find patterns in codebases
- Use `grep_search` for specific API usage
- Use `read_file` to study implementation details

### Evaluation Criteria
Apply scoring from `.specify/config/research.yaml`:

```bash
# Load weights from config
yq '.decision_criteria.weights' .specify/config/research.yaml
```

Weights are project-specific (each project defines their own priorities).

### Constitution Alignment (Optional)

**ONLY if project has a constitution** (`paths.constitution` defined in config):

```bash
# Check if constitution exists
CONSTITUTION=$(yq '.paths.constitution' .specify/config/research.yaml)
if [[ "$CONSTITUTION" != "null" ]] && [[ -f "$CONSTITUTION" ]]; then
    cat "$CONSTITUTION"
    # Load principles from config
    yq '.constitution_principles' .specify/config/research.yaml
fi
```

**If no constitution**: Skip this section entirely. Score options based on other criteria only.

## Error Handling

### Research Area Not Found
```
ERROR: Research area 'R99' not found.

To see available areas:
    cat {paths.research_plan}  # Path from .specify/config/research.yaml
```

### Prerequisites Missing
```
ERROR: Research infrastructure not found.

Expected files (from .specify/config/research.yaml):
  - {paths.research_plan}
  - {paths.research_template}
  - {paths.decision_log}

Initialize research infrastructure or check config paths.
```

### Unable to Answer Question
If a question cannot be answered with available resources:
1. Document what was tried
2. Mark confidence as LOW
3. Add to "Open Questions" in findings
4. Suggest how to resolve (e.g., "need to test locally")

## Progress Tracking

The agent tracks progress through:

1. **Task-level**: Each task in research.yaml has status field
2. **Question-level**: Each key_question has answer and confidence
3. **Area-level**: research.yaml status field
4. **Master-level**: research-plan.yaml statistics

Always update all levels when completing work.

## Integration After Research

Once research is complete, it feeds back through paths defined in `.specify/config/research.yaml`:

1. **Feature specs** (`{feature_specs}/`): Update with chosen libraries/patterns
2. **Architecture docs** (`{architecture_docs}/`): Update with decisions
3. **Sprint planning** (`{sprint_plans}/`): Unblock dependent tasks (if defined)
4. **Decision log** (`{decision_log}`): Decisions await approval before implementation

Use the `integration.update_targets` section from config to determine which files to update.

````