# Tasks: Release Versioning and Workflow

**Input**: Design documents from `/specs/004-release-workflow/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Not explicitly requested in spec - test validation covered by workflow validation steps.

**Organization**: Tasks grouped by user story (P1→P4) for independent implementation.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Exact file paths included in descriptions

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Create foundational files and directory structure for release workflow

- [ ] T001 Create VERSION file with initial version `1.0.0` at `/VERSION`
- [ ] T002 [P] Create release-please config at `/release-please-config.json`
- [ ] T003 [P] Create release-please manifest at `/.release-please-manifest.json`
- [ ] T004 [P] Create directory structure for scripts at `/scripts/release/`
- [ ] T005 [P] Create directory structure for workflows at `/.github/workflows/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core release scripts that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Implement version.sh script at `/scripts/release/version.sh` (get, bump, set commands)
- [ ] T007 [P] Implement validate.sh script at `/scripts/release/validate.sh` (--version, --changelog, --all options)
- [ ] T008 Update existing CHANGELOG.md to comply with Keep a Changelog format at `/CHANGELOG.md`
- [ ] T009 [P] Update .gitignore to exclude dist/ artifacts at `/.gitignore`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Create Versioned Release (Priority: P1) 🎯 MVP

**Goal**: Enable maintainers to create versioned releases via GitHub Actions workflow

**Independent Test**: Trigger release workflow, verify git tag created and GitHub Release published

### Implementation for User Story 1

- [ ] T010 [US1] Create main release workflow at `/.github/workflows/release.yml` with release-please-action@v4
- [ ] T011 [US1] Add release-please job to workflow with outputs (release_created, tag_name, version)
- [ ] T012 [US1] Configure workflow permissions (contents: write, pull-requests: write)
- [ ] T013 [US1] Add concurrency control to prevent simultaneous releases in workflow
- [ ] T014 [P] [US1] Add version badge to README.md at `/README.md`
- [ ] T015 [P] [US1] Update PR template with conventional commit guide at `/.github/PULL_REQUEST_TEMPLATE.md`

**Checkpoint**: User Story 1 complete - releases can be created via release-please PR workflow

---

## Phase 4: User Story 2 - Automated CHANGELOG Generation (Priority: P2)

**Goal**: Auto-generate CHANGELOG entries from conventional commits

**Independent Test**: Make feat:/fix:/docs: commits, merge Release PR, verify CHANGELOG categorized correctly

### Implementation for User Story 2

- [ ] T016 [US2] Configure changelog-sections in release-please-config.json for commit type mapping
- [ ] T017 [US2] Create PR validation workflow at `/.github/workflows/validate-pr.yml`
- [ ] T018 [US2] Add conventional commit validation using amannn/action-semantic-pull-request@v5
- [ ] T019 [P] [US2] Add ShellCheck linting step to validate-pr.yml for scripts/
- [ ] T020 [P] [US2] Add Docker Compose validation step to validate-pr.yml
- [ ] T021 [US2] Create changelog-preview.sh script at `/scripts/release/changelog-preview.sh`
- [ ] T022 [US2] Create changelog preview workflow at `/.github/workflows/changelog-preview.yml`

**Checkpoint**: User Story 2 complete - CHANGELOG auto-generated from conventional commits

---

## Phase 5: User Story 3 - Release Validation Before Publishing (Priority: P3)

**Goal**: Validate prerequisites before allowing release publication

**Independent Test**: Attempt release with invalid VERSION format, verify workflow fails with clear error

### Implementation for User Story 3

- [ ] T023 [US3] Add pre-release validation step to release.yml using scripts/release/validate.sh --all
- [ ] T024 [US3] Implement VERSION file validation in validate.sh (semver format check)
- [ ] T025 [US3] Implement CHANGELOG validation in validate.sh (entry exists for version)
- [ ] T026 [US3] Implement required files check in validate.sh (VERSION, CHANGELOG.md, docker-compose.yml)
- [ ] T027 [P] [US3] Add emoji status indicators to validate.sh output (✅ pass, ❌ fail)
- [ ] T028 [US3] Add validation failure handling to release.yml with clear error messages

**Checkpoint**: User Story 3 complete - invalid releases are blocked with clear feedback

---

## Phase 6: User Story 4 - Attach Release Artifacts (Priority: P4)

**Goal**: Build and attach downloadable artifacts to GitHub Releases

**Independent Test**: Create release, verify tar.gz archive, install.sh, and checksums.sha256 are attached

### Implementation for User Story 4

- [ ] T029 [US4] Implement artifacts.sh script at `/scripts/release/artifacts.sh`
- [ ] T030 [US4] Add archive creation logic to artifacts.sh (tar.gz with version in name)
- [ ] T031 [US4] Define archive contents in artifacts.sh per contracts/release-workflow.md
- [ ] T032 [US4] Add SHA256 checksum generation to artifacts.sh (checksums.sha256)
- [ ] T033 [P] [US4] Create standalone install.sh script at `/scripts/release/install.sh`
- [ ] T034 [US4] Add build-artifacts job to release.yml (conditional on release_created)
- [ ] T035 [US4] Configure softprops/action-gh-release@v2 for artifact upload in release.yml
- [ ] T036 [US4] Add artifact upload with glob patterns (dist/*.tar.gz, dist/install.sh, dist/checksums.sha256)

**Checkpoint**: User Story 4 complete - releases include downloadable artifacts with checksums

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, edge cases, and final validation

- [ ] T037 [P] Document release workflow in CONTRIBUTING.md at `/CONTRIBUTING.md`
- [ ] T038 [P] Add release-request issue template at `/.github/ISSUE_TEMPLATE/release-request.yml`
- [ ] T039 Add pre-release version support to release-please-config.json (alpha, beta, rc)
- [ ] T040 [P] Add workflow_dispatch manual trigger option to release.yml for exceptional cases
- [ ] T041 Run full release workflow validation using quickstart.md test scenarios
- [ ] T042 Update copilot-instructions.md with release workflow commands at `/.github/agents/copilot-instructions.md`

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1: Setup
    ↓
Phase 2: Foundational (BLOCKS all user stories)
    ↓
┌───────────────────────────────────────────────────┐
│  User Stories can proceed in priority order:      │
│  Phase 3 (US1) → Phase 4 (US2) → Phase 5 (US3)   │
│                      → Phase 6 (US4)              │
└───────────────────────────────────────────────────┘
    ↓
Phase 7: Polish
```

### User Story Dependencies

| Story | Depends On | Independent Test |
|-------|------------|------------------|
| US1 (P1) | Foundational | Trigger workflow → tag + release created |
| US2 (P2) | US1 (release.yml exists) | Conventional commits → CHANGELOG categorized |
| US3 (P3) | US1 (release.yml exists) | Invalid input → workflow fails with message |
| US4 (P4) | US1 (release.yml exists) | Release → artifacts attached |

### Within Each Phase

- Tasks marked [P] can run in parallel
- Sequential tasks depend on previous tasks in that phase
- Complete phase before moving to next

### Critical Path

```
T001 → T006 → T010 → T011 → T012 → T013 (MVP release workflow)
```

---

## Parallel Opportunities

### Phase 1 (Setup)
```bash
# Run in parallel:
T002: Create release-please-config.json
T003: Create .release-please-manifest.json
T004: Create scripts/release/ directory
T005: Create .github/workflows/ directory
```

### Phase 3 (US1)
```bash
# Run in parallel after T013:
T014: Add version badge to README.md
T015: Update PR template
```

### Phase 4 (US2)
```bash
# Run in parallel:
T019: Add ShellCheck linting
T020: Add Docker Compose validation
```

### Phase 6 (US4)
```bash
# After T032:
T033: Create install.sh (parallel with T034)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. ✅ Complete Phase 1: Setup (T001-T005)
2. ✅ Complete Phase 2: Foundational (T006-T009)
3. ✅ Complete Phase 3: User Story 1 (T010-T015)
4. **VALIDATE**: Create a test release, verify tag and GitHub Release
5. **Deploy**: MVP is functional - maintainers can create releases

### Incremental Delivery

| Increment | Stories | Value Delivered |
|-----------|---------|-----------------|
| MVP | US1 | Basic release workflow with tags and GitHub Releases |
| +1 | US1+US2 | Auto-generated CHANGELOG from conventional commits |
| +2 | US1+US2+US3 | Validation prevents bad releases |
| Complete | All | Full release with artifacts and checksums |

### Task Count Summary

| Phase | Tasks | Parallel |
|-------|-------|----------|
| Setup | 5 | 4 |
| Foundational | 4 | 2 |
| US1 (P1) | 6 | 2 |
| US2 (P2) | 7 | 2 |
| US3 (P3) | 6 | 1 |
| US4 (P4) | 8 | 1 |
| Polish | 6 | 3 |
| **Total** | **42** | **15** |

---

## Files Created/Modified Summary

### New Files (26)
- `/VERSION`
- `/release-please-config.json`
- `/.release-please-manifest.json`
- `/.github/workflows/release.yml`
- `/.github/workflows/validate-pr.yml`
- `/.github/workflows/changelog-preview.yml`
- `/.github/ISSUE_TEMPLATE/release-request.yml`
- `/scripts/release/version.sh`
- `/scripts/release/validate.sh`
- `/scripts/release/artifacts.sh`
- `/scripts/release/install.sh`
- `/scripts/release/changelog-preview.sh`

### Modified Files (5)
- `/CHANGELOG.md` - Update format
- `/README.md` - Add version badge
- `/.gitignore` - Add dist/
- `/.github/PULL_REQUEST_TEMPLATE.md` - Add conventional commit guide
- `/CONTRIBUTING.md` - Add release documentation
- `/.github/agents/copilot-instructions.md` - Add release commands

---

## Notes

- Tasks reference exact file paths from plan.md project structure
- [P] tasks operate on different files with no dependencies
- [US#] labels map tasks to user stories for traceability
- Each user story can be tested independently after completion
- Conventional commits are enforced by validate-pr.yml (US2)
- All scripts will be ShellCheck compliant per constitution
