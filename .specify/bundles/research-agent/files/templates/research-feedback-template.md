# Research Feedback Integration: [RESEARCH AREA NAME]

**Research ID**: [R00]  
**Created**: [DATE]  
**Status**: Integration Pending  
**Research Owner**: [NAME]

<!--
  RESEARCH FEEDBACK INTEGRATION TEMPLATE
  
  This template guides the integration of research findings back into project artifacts.
  Use this after completing a research area to ensure findings are properly applied.
  
  Based on .specify/config/research.yaml configuration and SpecKit integration workflows.
-->

## Prerequisites

Before starting integration, verify:

- [ ] Research area status is "completed" in research.yaml
- [ ] All critical questions have answers with confidence >= MEDIUM
- [ ] At least one decision is documented
- [ ] Findings.md is complete with executive summary
- [ ] All sources are cited with URLs
- [ ] Constitution alignment checked (if applicable)

## Research Summary

### Executive Overview

[Brief 2-3 sentence summary of research findings]

### Key Decisions Made

| Decision ID | Title | Chosen Option | Confidence | Impact |
|-------------|-------|---------------|------------|--------|
| D001 | [Decision Title] | [Option Selected] | HIGH/MEDIUM/LOW | [Brief impact] |

### Key Findings

1. **[Finding 1]**: [Brief description and significance]
2. **[Finding 2]**: [Brief description and significance]
3. **[Finding 3]**: [Brief description and significance]

## Integration Checklist

### 1. Decision Log Updates

**Action**: Add approved decisions to decision-log.yaml

**Decisions to integrate**:
- [ ] D001: [Decision Title]
  - Status: APPROVED / PENDING
  - Date: [YYYY-MM-DD]
  - Integrated: No → Yes

**Process**:
```bash
# Open decision log
vim {decision_log_path}

# Add decision entry following schema
# Update metadata counters
```

---

### 2. Feature Specification Updates

**Features blocked by this research**:

#### Feature [001]: [Feature Name]

**File**: `{feature_specs}/[feature-file].yaml` or `.md`

**Sections requiring updates**:

- [ ] **Implementation → Libraries/Dependencies**
  - Current: [What it says now / blank]
  - Update to: [New library/package based on research]
  - Decision ref: D001

- [ ] **Implementation → Patterns**
  - Current: [What it says now / blank]
  - Update to: [Architecture pattern chosen]
  - Decision ref: D002

- [ ] **Implementation → Constraints**
  - Add: [Technical constraint discovered]
  - Rationale: [Why this constraint exists]
  - Decision ref: D003

- [ ] **Technical Context**
  - Replace "NEEDS CLARIFICATION" with: [Resolved information]
  - Decision ref: D001

- [ ] **Requirements → Non-Functional**
  - Add/Update: [Performance gate or quality requirement]
  - Based on: [Research benchmark/finding]

- [ ] **Add Research Reference**
  ```yaml
  references:
    research:
      - id: R00
        decision: D001
        summary: "[Brief summary of research impact]"
  ```

**Validation**:
- [ ] Feature spec still aligns with constitution principles
- [ ] All "NEEDS CLARIFICATION" resolved
- [ ] Story points re-estimated if complexity changed

---

### 3. Architecture Documentation Updates

**Documents requiring updates**:

#### [architecture-doc.md]

**File**: `{architecture_docs}/[doc-name].md`

**Updates needed**:

- [ ] **Technology Stack Section**
  - Add: [New library/framework]
  - Version: [X.Y.Z]
  - Rationale: [Link to decision D001]
  - Alternatives considered: [List]

- [ ] **Design Patterns Section**
  - Add/Update: [Pattern name]
  - Description: [How it applies]
  - Code example: [Link or inline snippet]
  - Based on: [Research finding F001]

- [ ] **Diagrams**
  - [ ] Update [diagram name] to reflect [change]
  - [ ] Add new diagram for [component interaction]

- [ ] **Trade-offs Section**
  - Add: [Trade-off accepted]
  - Why acceptable: [Rationale]
  - Decision ref: D001

**Validation**:
- [ ] Architecture doc consistent with feature specs
- [ ] All technical decisions referenced
- [ ] Diagrams up to date

---

### 4. Sprint/Roadmap Updates

**Sprint**: [sprint-id]

**Tasks previously blocked**:

- [ ] **Task [T001]**: [Task Description]
  - Remove blocker: "pending_research: R00"
  - Re-estimate: [Old estimate] → [New estimate]
  - Reason for change: [What research revealed]
  - Status: blocked → ready

- [ ] **Task [T002]**: [Task Description]
  - Remove blocker: "pending_research: R00"
  - No estimate change needed
  - Status: blocked → ready

**New tasks created from research**:
- [ ] [New task identified from findings]
- [ ] [Another task needed based on decision]

**Roadmap impact**:
- [ ] No change to version sequencing
- [ ] OR: Version [X.Y.Z] should move [earlier/later] because: [Reason]

---

### 5. Constitution Updates (if applicable)

**Constitution conflicts or gaps identified**:

- [ ] **Conflict/Gap**: [Description]
  - Current principle: [PRINCIPLE_ID]
  - Research finding: [What research revealed]
  - Proposed resolution: [How to address]
  - Action: Create issue for constitution amendment
  - Issue #: [Link when created]

---

### 6. Documentation Updates

**User-facing documentation**:

- [ ] **README.md**
  - Section: [Section name]
  - Change: [What to add/update]
  - Reason: [Why this matters to users]

- [ ] **CHANGELOG.md** (if user-facing impact)
  - Version: [Next version]
  - Entry: "Added support for [feature based on research]"
  - Breaking change: Yes / No

- [ ] **Installation/Setup Docs**
  - New dependency: [Package name]
  - Installation: [Command/steps]
  - Configuration: [New settings]

---

## Constitution Alignment Verification

**Constitution principles addressed**:

| Principle | Complies | Notes |
|-----------|----------|-------|
| [PRINCIPLE_1] | ✅ / ⚠️ / ❌ | [Explanation] |
| [PRINCIPLE_2] | ✅ / ⚠️ / ❌ | [Explanation] |

**Performance gates validated**:

| Gate | Requirement | Research Finding | Status |
|------|-------------|------------------|--------|
| [gate_name] | [threshold] | [actual result] | ✅ PASS / ❌ FAIL |

**Action if gates failed**: [Mitigation strategy or principle exception justification]

---

## Risk Assessment

**Risks introduced by research decisions**:

| Risk | Probability | Impact | Mitigation | Owner |
|------|-------------|--------|------------|-------|
| [Risk 1] | LOW/MED/HIGH | LOW/MED/HIGH | [How to address] | [Name] |

---

## Communication Plan

**Stakeholders to notify**:

- [ ] **[Team/Person]**: [What they need to know and why]
- [ ] **[Team/Person]**: [What they need to know and why]

**Communication method**:
- [ ] Team meeting
- [ ] Slack/Discord announcement
- [ ] Email summary
- [ ] GitHub discussion
- [ ] Documentation update (self-service)

---

## Completion Verification

### Final Checklist

- [ ] All decision log entries added
- [ ] All feature specifications updated
- [ ] All architecture docs updated
- [ ] All sprint/roadmap blockers removed
- [ ] Constitution alignment verified
- [ ] User documentation updated
- [ ] Stakeholders notified
- [ ] Research.yaml marked as integrated
- [ ] Integration commit SHA recorded
- [ ] This feedback document committed to repo

### Integration Metadata

**Integration completed by**: [Name]  
**Integration date**: [YYYY-MM-DD]  
**Integration commit**: [SHA]  
**Review required**: Yes / No  
**Reviewed by**: [Name]  
**Review date**: [YYYY-MM-DD]

---

## Open Questions / Follow-up

**Questions that remain unanswered**:
1. [Question that needs follow-up research]
2. [Question that needs stakeholder decision]

**Future research needed**:
- [ ] [Follow-up research area]
- [ ] [Related investigation]

**Technical debt acknowledged**:
- [Technical debt accepted as trade-off]
- [When to revisit]: [Timeline or trigger]

---

## Notes

[Any additional notes, lessons learned, or context for future reference]

---

## Appendix

### Related Documents

- Research YAML: `{research_root}/[folder]/research.yaml`
- Findings Document: `{research_root}/[folder]/findings.md`
- Decision Log: `{decision_log_path}`
- Feature Specs: `{feature_specs}/`
- Architecture Docs: `{architecture_docs}/`

### Decision Cross-Reference

| Decision ID | Research Question | Chosen Option | Affected Artifacts |
|-------------|-------------------|---------------|-------------------|
| D001 | [Question] | [Option] | [List of specs/docs/tasks] |
| D002 | [Question] | [Option] | [List of specs/docs/tasks] |

---

**Template Version**: 1.0.0  
**Schema Compatibility**: research.yaml v1.0.0  
**Last Updated**: 2025-12-28
