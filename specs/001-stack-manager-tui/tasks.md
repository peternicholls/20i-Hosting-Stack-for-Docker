# Tasks: 20i Stack Manager TUI

**Feature Branch**: `001-stack-manager-tui`  
**Input**: Design documents from `/specs/001-stack-manager-tui/`  
**Prerequisites**: plan.md, spec.md (5 user stories), research.md, data-model.md, contracts/ (docker-api.md, ui-events.md)

**Tests**: INCLUDED - Unit tests, integration tests, and UI tests integrated throughout implementation phases using Go testing package, table-driven tests, mock Docker client, and Bubble Tea test utilities.

**Testing Strategy**: TDD-lite approach - write tests alongside implementation (not strictly before), validate each component independently, regression test suite for CI/CD.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

**📚 Agent Guidance**: Before starting, review:
- `/runbooks/research/QUICK-REFERENCE.md` - Keep this open while coding (cheat sheet)
- `/runbooks/research/INDEX.md` - Find which detailed guide you need
- `/runbooks/research/bubbletea-component-guide.md` - Component architecture patterns
- `/runbooks/research/lipgloss-styling-reference.md` - Styling patterns and color palette
- The TUI directory is `tui/` at the repository root

**Documentation Expectations**: While implementing tasks, add Go doc comments for any exported types, functions, and packages you touch, and add brief inline comments only where logic is non-obvious.

**Go File Header Template**: Use this standard header at the top of each Go source file (update fields as needed).
```go
// Project: 20i Stack Manager TUI
// File: <filename.go>
// Purpose: <short purpose>
// Version: <semver or revision>
// Updated: <YYYY-MM-DD>
```
---

## Format: `- [ ] [ID] [P?] [Story?] Description`

- **Checkbox**: `- [ ]` for tracking completion
- **[ID]**: Sequential task number (T001, T002, T003...)
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4, US5)
- **File paths**: Exact paths from tui/ directory structure in plan.md

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Initialize Go project structure and dependencies

- [X] T001 Create tui/ directory at repository root
- [X] T002 Initialize Go module with `go mod init github.com/peternicholls/20i-stack/tui`
- [X] T003 [P] Add Bubble Tea dependency v1.3.10+ in go.mod
- [X] T004 [P] Add Bubbles dependency v1.0.0+ in go.mod
- [X] T005 [P] Add Lipgloss dependency v1.0.0+ in go.mod
- [X] T006 [P] Add Docker SDK dependency v27.0.0+ in go.mod
- [X] T007 Run `go mod tidy` to generate go.sum
- [X] T008 Create directory structure: internal/app, internal/views/dashboard, internal/views/help, internal/views/projects, internal/docker, internal/ui
- [X] T009 [P] Create placeholder tui/main.go with basic Bubble Tea hello world
  > **📖 Reference**: See `/runbooks/research/QUICK-REFERENCE.md` - "Bubble Tea Basics" section for Init/Update/View pattern
- [X] T010 [P] Create Makefile with build, install, clean, test, test-coverage targets
- [X] T011 Create tests/ directory structure: tests/unit/, tests/integration/, tests/mocks/
- [X] T012 Verify `go run main.go` works (shows hello world, press 'q' to quit)
- [X] T013 [P] Create tests/mocks/docker_mock.go with MockDockerClient interface matching docker-api.md contract

**Checkpoint**: Go project initializes, dependencies resolve, basic TUI runs, test infrastructure ready

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T012 [P] Create internal/ui/styles.go with Lipgloss color palette (ColorRunning, ColorStopped, ColorError, ColorAccent, ColorBorder)
  > **🎨 CRITICAL**: Use EXACT color palette from `/runbooks/research/lipgloss-styling-reference.md` - "Standard Color Palette" section
  > ColorRunning = lipgloss.Color("10"), ColorStopped = lipgloss.Color("8"), ColorError = lipgloss.Color("9"), etc.
  > **⚠️ NEVER use raw ANSI codes** - Always use Lipgloss. See "Common Gotchas" in styling reference.
- [X] T013 [P] Create internal/ui/components.go with StatusIcon function (maps ContainerStatus enum to emoji icons)
  - **Implemented**: `StatusIcon(status string) string` in `internal/ui/components.go` — maps status strings to glyphs (●, ○, ✗, ⚠) and applies color styles from `internal/ui/styles.go`. Unit tests added in `internal/ui/components_test.go`.
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md` - "Unicode Icons Reference" section
  > Use "●" for running, "○" for stopped, "✗" for error, "⚠" for warning
- [X] T014 [P] Create internal/ui/layout.go with panel sizing functions (calculatePanelWidths, calculatePanelHeights)
  > Use "●" for running, "○" for stopped, "✗" for error, "⚠" for warning
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md` - "Layout Patterns" section
  > Use lipgloss.Width() and Height() to measure before placing. See "3-Panel Layout" example.
- [X] T015 Create internal/docker/client.go with Client struct and NewClient() method per docker-api.md contract
- [X] T016 Implement Docker connection check in client.go (Ping method, error handling for daemon unreachable)
- [X] T017 [P] Create internal/app/messages.go with custom tea.Msg types: statsMsg, containerListMsg, logLineMsg, containerActionMsg, containerActionResultMsg per ui-events.md
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md` - "Component Pattern" section for message type examples
  > Messages drive ALL state changes. Never mutate model directly.
- [X] T018 Create internal/app/root.go with RootModel struct (activeView, dashboard, help, projects, dockerClient fields)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md` - "Root Model Pattern" section
  > Root model is a COORDINATOR - delegates messages to child components. See full example in guide.
- [X] T019 Implement RootModel.Init() method (initialize Docker client, return initial commands)
- [X] T020 Implement RootModel.Update() method with global shortcut routing (q=quit, ?=help, p=projects)
  > **📖 Reference**: See `/runbooks/research/QUICK-REFERENCE.md` - "Common Update Pattern" section for switch statement template
  > Handle tea.KeyMsg for shortcuts, tea.WindowSizeMsg for responsive layout. NEVER block in Update()!
- [X] T021 Implement RootModel.View() method with modal overlay logic
- [X] T022 Update tui/main.go to create and run RootModel instead of hello world
- [X] T023 [P] [TEST] Create internal/docker/client_test.go with table-driven tests for Docker connection (success, daemon unreachable, permission denied)
- [X] T024 [P] [TEST] Create internal/app/root_test.go with unit tests for RootModel Init/Update/View (test global shortcuts, view routing)
- [X] T024a [TEST] Add RootModel.lastError clearing test - verify error clears on success/view change
- [X] T024b [TEST] Add stubbed client path test - verify app handles Docker unavailable gracefully (soft-fail vs hard-fail startup)
- [X] T025 [TEST] Run `make test` to verify foundational tests pass
- [X] T025a [TEST] Update main.go to write errors to stderr instead of stdout

**Checkpoint**: Foundation ready - Docker client connects, root model routes messages, test coverage >80%

### Possible Next Steps (Pre-Phase 3)

Before proceeding to Phase 3 implementation, consider these architectural decisions:

- [X] T025b [DECISION] Decide startup resilience strategy for Docker outages:
  - **Decision**: Hard-fail approach chosen for MVP (documented in STARTUP-RESILIENCE-DECISION.md)
  - **Rationale**: Simpler implementation, clear error messaging, fail-fast principle
  - **Future**: Soft-fail with degraded UI could be added in v1.1+ if needed
- [X] T025c [ENHANCEMENT] Implement RootModel.lastError clearing:
  - Clear lastError on successful operation completion (SuccessMsg handler)
  - Clear lastError on view change (?, p, esc key handlers)
  - Test coverage added in root_test.go (T024a)
- [X] T025d [FIX] Update main.go error handling:
  - Changed error output to stderr using `fmt.Fprintf(os.Stderr, ...)`
  - Non-zero exit codes (exit 1) on all fatal errors
  - Implementation verified in main.go

**Note**: These tasks are OPTIONAL before Phase 3. They improve robustness but are not blockers for MVP.

---

## Phase 3: User Story 2 - Container Lifecycle (Priority: P0 - Core) 🎯 MVP FIRST

**Goal**: Start, stop, restart individual services or entire stack - REPLICATE 20i-gui core functionality

**Independent Test**: Run TUI, see service list, press 's' on stopped apache, verify it starts; press 'S' to stop all, verify stack stops

**Rationale**: This is the PRIMARY use case - get stacks running and verified. Dashboard monitoring is secondary.

**📚 CRITICAL - Read Before Starting Phase 3**:

1. **[PHASE3-IMPLEMENTATION-NOTES.md](PHASE3-IMPLEMENTATION-NOTES.md)** - Architectural decisions and implementation patterns
   - Section 1: Entity Design Strategy (minimal Container schema - 6 fields only)
   - Section 2: Dashboard Layout Strategy (2-panel layout, NOT 3-panel)
   - Section 3: Service List Rendering (simple icons + names, NO stats yet)
   - Section 4: Message Type Design (string-based actions, not enums)
   - Section 5: Command Function Pattern (generic containerActionCmd)
   - Section 6: Compose Operations Strategy (Stop/Restart/Down only, NO Up)
   - Section 7: Error Message Formatting (centralized formatDockerError)
   - Section 8: Testing Strategy (3-layer approach)

2. **[PHASE3-ROADMAP.md](PHASE3-ROADMAP.md)** - Step-by-step execution plan
   - Block 1-4: Docker Client Layer (entities, list, lifecycle, compose)
   - Block 5: Message Types Enhancement
   - Block 6-8: Dashboard Foundation (model, rendering, navigation)
   - Block 9-11: Actions & Polish (container actions, stack actions, status messages)
   - Block 12-13: Optional polish + Integration testing
   - Time estimates: 25-29 hours solo, 14-18 hours with 3 developers

3. **[PHASE3-ADR.md](PHASE3-ADR.md)** - Architecture Decision Records
   - ADR-001: Minimal Container schema (6 fields → extend to 9 in Phase 5)
   - ADR-002: 2-panel layout in Phase 3 (expand to 3-panel in Phase 5)
   - ADR-003: String-based actions (not typed enums)
   - ADR-004: Generic containerActionCmd (not separate functions)
   - ADR-005: No ComposeUp implementation (focus on management, not setup)
   - ADR-006: Centralized error formatter (consistent UX)

**⚠️ Implementation Order**: Follow critical path in PHASE3-IMPLEMENTATION-NOTES.md Section "Implementation Order"
- Start: T026-T027 (entities) → T028-T029 (state mapping) → T030-T031 (list containers)
- Then: T044-T046 (dashboard model) → T047-T049 (rendering) → T050 (wire to root)
- See PHASE3-ROADMAP.md for detailed block-by-block breakdown

### Implementation for User Story 2 - LIFECYCLE FIRST

- [X] T026 [P] [US2] Create Container entity struct in internal/docker/client.go (ID, Name, Service, Image, Status, State fields - minimal for lifecycle; will be extended later with Ports, CreatedAt, StartedAt)
  > **📋 Reference**: See PHASE3-ADR.md ADR-001 for minimal schema decision (6 fields only)
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 1 for entity design strategy
- [X] T027 [P] [US2] Create ContainerStatus enum in internal/docker/client.go (Running, Stopped, Restarting, Error)
- [X] T028 [P] [TEST] Create internal/docker/client_test.go unit tests for mapDockerState() with table-driven tests (all Docker states → ContainerStatus)
- [X] T029 [US2] Implement mapDockerState() helper function in internal/docker/client.go to map Docker states to ContainerStatus enum
- [X] T030 [US2] Implement ListContainers(projectName string) method in internal/docker/client.go per docker-api.md contract
- [X] T031 [P] [TEST] Add unit tests to client_test.go for ListContainers (with mock client, test project filtering, empty results, errors)
- [X] T032 [P] [US2] Implement StartContainer(containerID string) method in internal/docker/client.go per docker-api.md contract
- [X] T033 [P] [US2] Implement StopContainer(containerID string, timeout int) method in internal/docker/client.go
- [X] T034 [P] [US2] Implement RestartContainer(containerID string, timeout int) method in internal/docker/client.go
- [X] T035 [P] [TEST] Add unit tests for Start/Stop/Restart methods (mock client, test success, container not found, timeout errors)
- [X] T036 [P] [US2] Implement ComposeStop(projectPath string) method in internal/docker/client.go per docker-api.md
  > **📋 Reference**: See PHASE3-ADR.md ADR-005 for rationale (NO ComposeUp - management not setup)
- [X] T037 [P] [US2] Implement ComposeRestart(projectPath string) method in internal/docker/client.go
- [X] T038 [P] [US2] Implement ComposeDown(projectPath string, removeVolumes bool) method in internal/docker/client.go
  > **⚠️ WARNING**: removeVolumes=true DESTROYS ALL DATA - wire to 'D' key with confirmation in Phase 4
- [X] T039 [P] [TEST] Add unit tests for Compose operations (mock exec, test success, invalid path, permission errors)
- [X] T040 [US2] Add ContainerAction enum to internal/app/messages.go (Start, Stop, Restart)
  > **📋 Reference**: See PHASE3-ADR.md ADR-003 for string-based action decision (use comments, NOT typed enums)
- [X] T041 [US2] Add ComposeAction enum to internal/app/messages.go (StopAll, RestartAll, Destroy)
  > **📋 Reference**: Document valid values in comments per ADR-003
- [X] T042 [US2] Add composeActionMsg and composeActionResultMsg types to internal/app/messages.go
- [X] T043 [P] [TEST] Create internal/app/messages_test.go with tests for message type creation and field validation
- [ ] T044 [P] [US2] Create internal/views/dashboard/dashboard.go with DashboardModel struct (serviceList, containers, selectedIndex fields - NO stats yet)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md` - "Component Structure" section
  > Every component must implement tea.Model interface (Init/Update/View). See full pattern example.
  > **📋 Reference**: See PHASE3-ADR.md ADR-002 for 2-panel layout decision (NOT 3-panel)
- [ ] T045 [US2] Implement DashboardModel.Init() method to load container list
- [ ] T046 [US2] Implement containerListMsg handling in DashboardModel.Update()
- [ ] T047 [US2] Create internal/views/dashboard/service_list.go with simple list rendering (status icon + name only)
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md` - "List Row Styles" section
  > AND `/runbooks/research/QUICK-REFERENCE.md` - "Service List" pattern (copy-paste example)
  > Use StatusIcon from T013 + RowStyle/SelectedRowStyle. DEFINE STYLES ONCE at package level, not in View()!
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 3 for simple rendering pattern (NO stats)
- [ ] T048 [P] [TEST] Create internal/views/dashboard/dashboard_test.go with Bubble Tea test program (test Init returns correct cmd, Update handles messages, View renders)
- [ ] T049 [US2] Implement DashboardModel.View() with simple 2-panel layout (service list 30% | status messages 70% | footer)
  > **📋 Reference**: See PHASE3-ADR.md ADR-002 for 2-panel layout justification and Phase 5 migration plan
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 2 for layout ASCII diagram
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md` - "3-Panel Layout (Dashboard)" section
  > Use lipgloss.JoinHorizontal() for side-by-side, JoinVertical() for stacking. Measure panels first with lipgloss.Width().
  > **⚠️ Anti-Pattern**: Don't hardcode widths! Use m.width from tea.WindowSizeMsg to calculate panel sizes.
- [ ] T050 [US2] Wire DashboardModel into RootModel in internal/app/root.go
- [ ] T051 [US2] Implement navigation keys (↑/k=up, ↓/j=down) in DashboardModel.Update()
  > **📖 Reference**: See `/runbooks/research/QUICK-REFERENCE.md` - "Common Update Pattern" for key handling switch statement
  > Handle both arrow keys ("up"/"down") AND vim bindings ("k"/"j") for better UX
- [ ] T052 [P] [TEST] Add tests to dashboard_test.go for navigation (test up/down key messages update selectedIndex correctly)
- [ ] T053 [US2] Implement 's' key handler to toggle start/stop for selected container
- [ ] T054 [US2] Implement 'r' key handler to restart selected container
- [ ] T055 [US2] Implement 'S' key handler to stop all stack containers (with simple confirmation)
- [ ] T056 [US2] Implement 'R' key handler to restart entire stack
- [ ] T057 [P] [TEST] Add tests for key handlers (test 's' sends correct containerActionMsg, 'S' sends composeActionMsg)
- [ ] T058 [US2] Create containerActionCmd() function in dashboard.go to launch async Docker operation
  > **📋 CRITICAL**: See PHASE3-ADR.md ADR-004 - Implement ONE generic function for all actions (NOT separate functions)
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 5 for generic command pattern code example
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Component Pattern" section, fetchServicesCmd() example
  > Commands return tea.Msg when complete. Use goroutine inside tea.Cmd to avoid blocking UI.
  > **⚠️ Anti-Pattern**: NEVER call Docker API directly in Update() - this blocks the UI! Always use tea.Cmd.
- [ ] T059 [US2] MERGED INTO T058 - Use generic containerActionCmd() per ADR-004
- [ ] T060 [US2] MERGED INTO T058 - Use generic containerActionCmd() per ADR-004
- [ ] T061 [US2] Implement containerActionResultMsg handler with success/error feedback
- [ ] T062 [P] [TEST] Add tests for command functions (use mock client, verify commands send correct result messages)
- [ ] T063 [US2] Add status message panel to show "✅ Container started" or "❌ Failed: error"
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md" - "Error/Warning Message Styles" section
  > Use ErrorStyle/WarningStyle/InfoStyle with Unicode icons (✓, ✗, ⚠, ℹ) from "Unicode Icons Reference"
- [ ] T064 [US2] Trigger containerListMsg refresh after successful action
- [ ] T065 [US2] Implement error message formatting per docker-api.md (user-friendly port conflicts, timeouts)
  > **📋 Reference**: See PHASE3-ADR.md ADR-006 for centralized formatDockerError() function pattern
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 7 for complete code example with regex
- [ ] T066 [P] [TEST] Add tests for error message formatting (test port conflict → user-friendly message, timeout → actionable message)
  > **📋 Reference**: See ADR-006 for table-driven test structure
- [ ] T067 [US2] Add footer with basic shortcuts: "s:start/stop  r:restart  S:stop-all  R:restart-all  D:destroy  q:quit"
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md" - "Footer/Help Style" section
  > Use FooterStyle with ColorMuted foreground. Keep footer visible at bottom of every view.
- [ ] T068 [US2] Implement Enter key handler to show detail panel for selected container (basic info: image, status, uptime)
- [ ] T069 [US2] Implement Tab key to cycle focus between service list and status message panel
- [ ] T070 [P] [TEST] Create tests/integration/lifecycle_test.go - integration test with mock Docker client simulating full lifecycle workflow
  > **📋 Reference**: See PHASE3-ROADMAP.md Block 13 for test scenarios and acceptance criteria
- [ ] T071 [TEST] Manual test per US2 acceptance scenarios 1-5 in spec.md: start stopped container, verify status changes; stop running container, verify status changes
  > **📋 Reference**: See PHASE3-ROADMAP.md "Test Scenarios" section for all 6 acceptance tests
- [ ] T072 [TEST] Run `make test` to verify all Phase 3 tests pass (unit + integration), achieve >85% coverage
  > **📋 Reference**: See PHASE3-IMPLEMENTATION-NOTES.md Section 8 for testing strategy and coverage targets

**Checkpoint**: Core 20i-gui functionality working - can start/stop/restart containers, verify status, >85% test coverage

**📊 Phase 3 Success Criteria** (from PHASE3-ROADMAP.md):
- ✅ All 47 tasks (T026-T072) checked off
- ✅ `make test` passes with >85% coverage
- ✅ All 6 manual acceptance scenarios verified
- ✅ No blocking bugs or crashes
- ✅ Error messages are user-friendly
- ✅ Code follows Go best practices (gofmt, golint)

**🎯 Ready for Phase 4**: Dashboard layout established, message handling patterns proven, Docker client tested and reliable

---

## Phase 4: User Story 4 - Destroy Stack (Priority: P0 - Core)

**Goal**: Destroy stack (stop containers + remove volumes) - COMPLETE 20i-gui baseline

**Independent Test**: Press 'D', confirm with 'yes', verify stack destroyed and volumes removed

### Implementation for User Story 4

- [ ] T073 [P] [US4] Create ConfirmationModal component in internal/ui/components.go with text input and warning styling
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md" - "Centered Modal Dialog" section for layout
  > AND `/runbooks/research/bubbletea-component-guide.md" - "Text Input Component" for Bubbles textinput.Model
  > Use lipgloss.Place() to center modal, ErrorStyle/WarningStyle for borders/background
- [ ] T074 [P] [TEST] Create internal/ui/components_test.go with tests for ConfirmationModal (test input validation, yes/no/esc handling)
- [ ] T075 [US4] Add confirmationModal field to DashboardModel in dashboard.go
- [ ] T076 [US4] Implement 'D' key handler in DashboardModel.Update() to show destroy confirmation modal
- [ ] T077 [US4] Render confirmation modal overlay with warning "⚠️ This will REMOVE ALL VOLUMES and data. Type 'yes' to confirm:"
- [ ] T078 [US4] Add text input to confirmation modal using Bubbles textinput.Model component
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Text Input Component" section
  > Create with textinput.New(), set Placeholder, Focus(). Update in modal's Update method.
- [ ] T079 [US4] Implement confirmation modal input handling (type "yes", press Enter to confirm)
- [ ] T080 [US4] Implement Esc handler in confirmation modal to cancel without destroying
- [ ] T081 [P] [TEST] Add tests to dashboard_test.go for destroy confirmation flow (test 'D' key, typing "yes", Esc cancel)
- [ ] T082 [US4] Create composeDownCmd() function in dashboard.go to call ComposeDown with removeVolumes=true
- [ ] T083 [US4] Implement composeActionResultMsg handler for Destroy action
- [ ] T084 [US4] Show success message "✅ Stack destroyed" after ComposeDown completes
- [ ] T085 [US4] Refresh container list after destroy (should show empty or no running containers)
- [ ] T086 [US4] Update footer to show "D:destroy" shortcut
- [ ] T087 [P] [TEST] Create tests/integration/destroy_test.go - integration test for full destroy workflow with mock client
- [ ] T088 [TEST] Manual test per US4 acceptance: Press 'D', type "yes", verify stack destroyed and volumes removed
- [ ] T089 [TEST] Run regression test suite (all previous tests + US4 tests) to ensure baseline functionality intact

**Checkpoint**: ✅ BASELINE COMPLETE - All 20i-gui core functions replicated, >85% test coverage, regression suite passing

---

## Phase 5: User Story 1 - Dashboard Monitoring (Priority: P1 - Enhancement)

**Goal**: Add live CPU/memory monitoring and detailed container info - ENHANCES baseline beyond 20i-gui

**Independent Test**: After lifecycle working, verify CPU% and memory bars update every 2s, detail panel shows ports/image/uptime

**Note**: NFRs (NFR-001 to NFR-008) for performance, stats refresh rates, and memory limits apply to this enhanced phase, not the baseline MVP.

### Implementation for User Story 1 - Dashboard Enhancement

- [ ] T090 [P] [US1] Create Stats entity struct in internal/docker/client.go (CPUPercent, MemoryUsed, MemoryLimit, MemoryPercent, NetworkRxBytes, NetworkTxBytes, Timestamp)
- [ ] T091 [P] [US1] Create PortMapping entity struct in internal/docker/client.go (ContainerPort, HostPort, Protocol)
- [ ] T092 [P] [TEST] Add unit tests for Stats and PortMapping struct validation
- [ ] T093 [US1] Extend Container entity with Ports, CreatedAt, StartedAt fields (builds on minimal schema from earlier)
- [ ] T094 [US1] Implement WatchStats() method in internal/docker/stats.go per docker-api.md
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - search for "Command factory" pattern
  > Return channel-based tea.Cmd that sends statsMsg periodically. Use goroutine, don't block Update().
- [ ] T095 [P] [TEST] Create internal/docker/stats_test.go with tests for WatchStats (mock Docker stats stream, test channel output)
- [ ] T096 [US1] Add statsMsg type to internal/app/messages.go per ui-events.md
- [ ] T097 [US1] Implement statsSubscriptionCmd() in dashboard.go to subscribe to stats stream
- [ ] T098 [US1] Implement statsMsg handler in DashboardModel.Update()
- [ ] T099 [P] [TEST] Add tests for stats subscription and message handling (mock stats channel, verify UI updates)
- [ ] T100 [US1] Enhance service_list.go rendering to show CPU/Memory usage per container (inline sparklines)
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md" - "Progress Bar" section for bar rendering
  > Use RenderProgressBar() pattern with filled/empty blocks (▓/░) or vertical bars (▁▂▃▄▅▆▇█) for sparklines
- [ ] T101 [US1] Add basic port mapping display in detail panel ("8080:80/tcp")
- [ ] T102 [US1] Update detail panel to show CreatedAt and Uptime
- [ ] T103 [P] [TEST] Create internal/views/dashboard/service_list_test.go with tests for enhanced rendering (test CPU/mem formatting, port display)
- [ ] T104 [US1] Implement auto-refresh for stats every 2 seconds using tea.Tick
- [ ] T105 [US1] Add visual loading indicator while waiting for stats
- [ ] T106 [US1] Add sparkline component for CPU/Memory history (last 60 data points)
- [ ] T107 [P] [TEST] Add UI snapshot tests for dashboard with stats data (verify layout, sparklines render correctly)
- [ ] T108 [US1] Optimize rendering to update only changed stats (not entire container list)
- [ ] T109 [TEST] Performance test: Verify stats refresh <200ms with 10 containers, <50ms panel switching
- [ ] T110 [TEST] Manual test per US1 acceptance: Open dashboard, verify stats update every 2s, CPU/mem sparklines visible

**Checkpoint**: Dashboard enhanced with real-time monitoring, performance targets met, tests passing

---

## Phase 6: User Story 3 - Log Viewer (Priority: P2 - Enhancement)

**Goal**: View live container logs with follow mode - ADDS debugging capability beyond 20i-gui

**Independent Test**: Press 'l' on running container, verify logs show; press 'f' for follow mode, make web request, see new log line

### Implementation for User Story 3

- [ ] T111 [P] [US3] Create LogStream entity struct in internal/docker/logs.go (Timestamp, Source [stdout/stderr], Message)
- [ ] T112 [P] [TEST] Create internal/docker/logs_test.go with tests for log parsing and filtering
- [ ] T113 [P] [US3] Create internal/views/logs/logs.go with LogsModel struct (viewport, containerID, logLines, isStreaming fields)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Viewport Component" section
  > Use viewport.New(width, height) from Bubbles. Perfect for scrollable logs with auto-scroll.
  > See "Real-World Examples → Log Viewer with Auto-Scroll" for complete pattern.
- [ ] T114 [US3] Implement StreamLogs(containerID string, tail int) method in internal/docker/logs.go per docker-api.md
- [ ] T115 [P] [TEST] Add tests for StreamLogs (mock container logs API, test channel output, error handling)
- [ ] T116 [US3] Add logLineMsg type to internal/app/messages.go per ui-events.md
- [ ] T117 [US3] Implement LogsModel.Init() to subscribe to log stream
- [ ] T118 [US3] Implement logLineMsg handler in LogsModel.Update() to append to viewport
- [ ] T119 [P] [TEST] Create internal/views/logs/logs_test.go with tests for log message handling (test append, scrolling, stream start/stop)
- [ ] T120 [US3] Implement LogsModel.View() with Bubbles viewport.Model (auto-scroll to bottom)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Viewport Component" usage
  > Use m.viewport.GotoBottom() for auto-scroll in follow mode. Handle scroll keys (↑↓ PgUp PgDn) via viewport.Update()
- [ ] T121 [US3] Wire LogsModel into RootModel, bind 'l' key in DashboardModel to switch to logs view
- [ ] T122 [US3] Add stdout/stderr color coding (green/red) to log rendering
- [ ] T123 [P] [TEST] Add tests for color coding logic (verify ANSI codes applied correctly)
- [ ] T124 [US3] Implement scroll controls (↑/↓/PgUp/PgDn) for viewport navigation
- [ ] T125 [US3] Add 'f' key to toggle follow mode (auto-scroll vs. manual scroll)
- [ ] T126 [US3] Implement timestamp formatting (HH:MM:SS.mmm) for log lines
- [ ] T127 [P] [TEST] Add tests for follow mode toggle and scroll controls
- [ ] T128 [US3] Add log filtering UI: '/' key to show search input, filter logs by text
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Text Input Component" for search box
  > Toggle visibility with '/' key. Filter log lines before setting viewport content. Use strings.Contains() or regex.
- [ ] T129 [US3] Implement search input handler using Bubbles textinput.Model
- [ ] T130 [US3] Add regex support to log filter (toggle with 'Alt+r')
- [ ] T131 [P] [TEST] Add tests for log filtering (text search, regex patterns, error handling)
- [ ] T132 [US3] Add 'c' key to clear log buffer
- [ ] T133 [US3] Add footer with shortcuts: "f:follow  /:search  c:clear  Esc:back  q:quit"
- [ ] T134 [US3] Implement graceful cleanup: stop log stream when switching views or quitting
- [ ] T135 [P] [TEST] Create tests/integration/logs_test.go - integration test for full log streaming workflow
- [ ] T136 [TEST] Performance test: Verify log rendering <100ms for 1000 lines, memory stable under continuous streaming
- [ ] T137 [TEST] Manual test per US3 acceptance: Open logs, verify streaming, test follow toggle, search, clear

**Checkpoint**: Log viewer functional, streaming works, filtering implemented, tests passing

---

## Phase 7: User Story 5 - Project Switcher (Priority: P3 - Enhancement)

**Goal**: Multi-project support - OPTIONAL enhancement for power users

**Independent Test**: Press 'p', see project list, select different project, verify dashboard switches context

### Implementation for User Story 5

- [ ] T138 [P] [US5] Create Project entity struct in internal/docker/client.go (Name, Path, IsActive)
- [ ] T139 [P] [TEST] Create internal/docker/project_test.go with tests for project detection and validation
- [ ] T140 [US5] Implement DiscoverProjects(basePath string) method in internal/docker/client.go to scan for docker-compose.yml
- [ ] T141 [P] [TEST] Add tests for DiscoverProjects (mock filesystem, test single project, multiple projects, no compose file)
- [ ] T142 [P] [US5] Create internal/views/projects/projects.go with ProjectListModel struct (projects list.Model, basePath field)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "List Component" section
  > Use Bubbles list.Model for scrollable, filterable project list. Define list.Item interface for Project entity.
  > See "List Component" example for item implementation (Title/Description/FilterValue methods)
- [ ] T143 [US5] Implement ProjectListModel.Init() to discover projects
- [ ] T144 [US5] Implement ProjectListModel.View() with Bubbles list.Model
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "List Component" View() section
  > Simply return m.list.View() - Bubbles handles all rendering, navigation, filtering. Update list in Update() method.
- [ ] T145 [P] [TEST] Create internal/views/projects/projects_test.go with tests for project list rendering and selection
- [ ] T146 [US5] Wire ProjectListModel into RootModel, bind 'p' key in DashboardModel to switch to projects view
- [ ] T147 [US5] Implement project selection: Enter key to switch active project
- [ ] T148 [US5] Implement projectSwitchMsg to notify RootModel of project change
- [ ] T149 [P] [TEST] Add tests for project switching workflow (test message flow, verify dashboard updates)
- [ ] T150 [US5] Update DashboardModel to reload containers when projectSwitchMsg received
- [ ] T151 [US5] Add visual indicator in dashboard header showing current project name
- [ ] T152 [US5] Add footer in projects view: "Enter:select  Esc:back  q:quit"
- [ ] T153 [P] [TEST] Create tests/integration/projects_test.go - integration test for multi-project workflow
- [ ] T154 [TEST] Manual test per US5 acceptance: Discover multiple projects, switch between them, verify container list updates

**Checkpoint**: Multi-project support working, can switch between stacks, tests passing

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Production-ready improvements for release

- [ ] T155 [P] Apply consistent Lipgloss styling across all views (colors, borders, padding)
  > **📖 Reference**: Review ALL styles against `/runbooks/research/lipgloss-styling-reference.md` - "Common Style Patterns"
  > Ensure HeaderStyle, FooterStyle, PanelStyle, RowStyle are CONSISTENT across dashboard/logs/projects views.
  > Use SAME color palette everywhere (ColorRunning, ColorStopped, ColorError, ColorBorder from T012)
- [ ] T156 [P] Implement responsive layout: adjust panel widths based on terminal size
  > **📖 Reference**: See `/runbooks/research/lipgloss-styling-reference.md" - "Common Gotchas → Hardcoding Dimensions"
  > ALWAYS handle tea.WindowSizeMsg in Update(). Calculate panel widths from m.width, NOT hardcoded 80.
  > Use calculated widths in View() to set style.Width(). Test with various terminal sizes (minimum 80x24).
- [ ] T157 [P] Add global help modal ('?' key) with all keyboard shortcuts organized by view
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Help Component" section
  > Use Bubbles help.Model with key.Binding structs. Renders ShortHelp (footer) and FullHelp (modal) automatically.
  > See keyMap struct example - define all shortcuts with descriptions
- [ ] T158 [P] [TEST] Create internal/ui/styles_test.go with tests for style consistency (verify color scheme, spacing)
- [ ] T159 Improve error messages: add "What to do next" suggestions for common failures
- [ ] T160 Add loading spinners for long operations (compose down, stats initialization)
  > **📖 Reference**: See `/runbooks/research/bubbletea-component-guide.md" - "Spinner Component" section
  > Use Bubbles spinner.Model. Create with spinner.New(), set style (spinner.Dot, spinner.Line, etc.)
  > Update spinner in Update(), render in View() while isLoading flag is true
- [ ] T161 Implement graceful degradation: fallback UI if Docker daemon unreachable
- [ ] T162 [TEST] Add accessibility tests: verify keyboard-only navigation, screen reader compatibility hints
- [ ] T163 Add startup banner with version number and Docker status
- [ ] T164 Implement config file support: ~/.20istackman/settings.json for default project path
- [ ] T165 Add theme toggle: cycle built-in palettes, persist selection in ~/.20istackman/settings.json
- [ ] T166 Add '--project' CLI flag to override default project
- [ ] T167 Add '--version' and '--help' CLI flags
- [ ] T168 [TEST] Create tests/e2e/ directory with end-to-end test suite using Bubble Tea test utilities
- [ ] T169 [TEST] Create e2e test: Launch app, navigate dashboard, perform lifecycle action, verify result
- [ ] T170 [TEST] Create e2e test: Open logs, verify streaming, apply filter, verify results
- [ ] T171 [TEST] Create e2e test: Switch projects, verify containers reload
- [ ] T172 [TEST] Create regression test suite: Run all unit + integration + e2e tests in sequence
- [ ] T173 [TEST] Performance regression: Benchmark startup time <2s, memory usage <30MB, stats refresh <200ms
- [ ] T174 [TEST] Cross-platform test: Run test suite on macOS and Linux (GitHub Actions)
- [ ] T175 Update README.md with screenshots, installation instructions, keyboard shortcuts reference
- [ ] T176 Create TESTING.md with guide for running unit/integration/e2e tests, adding new tests
- [ ] T177 Add inline code comments for complex Bubble Tea message flows
- [ ] T178 Generate API documentation: `go doc` for all exported types
- [ ] T179 Create ARCHITECTURE.md diagram showing Model-Update-View flow
- [ ] T180 [TEST] Documentation test: Verify all CLI flags documented, all shortcuts in help modal
- [ ] T181 [TEST] Final acceptance test: Run all 5 user story acceptance scenarios end-to-end
- [ ] T182 [TEST] Test coverage report: Generate HTML coverage report, verify >85% total coverage
- [ ] T183 [TEST] Run `make test-all` (unit + integration + e2e + regression), all tests must pass

**Checkpoint**: Production-ready - all features polished, comprehensive test suite passing, documentation complete, >85% coverage

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - START HERE
- **Foundational (Phase 2)**: Depends on Setup (T001-T011) - BLOCKS all user stories
- **User Story 2 - Lifecycle (Phase 3)**: Depends on Foundational - **THIS IS THE MVP** 🎯
- **User Story 4 - Destroy (Phase 4)**: Depends on US2 dashboard - COMPLETES 20i-gui baseline
- **User Story 1 - Dashboard (Phase 5)**: OPTIONAL enhancement - adds monitoring to baseline
- **User Story 3 - Logs (Phase 6)**: OPTIONAL enhancement - adds debugging
- **User Story 5 - Projects (Phase 7)**: OPTIONAL enhancement - adds multi-project
- **Polish (Phase 8)**: Depends on at minimum US2+US4 (baseline), or all stories (full feature set)

### User Story Dependencies & Priority

**BASELINE (20i-gui parity)**:
1. **User Story 2 (Lifecycle - P0)**: START HERE after Foundational - Creates simple dashboard with lifecycle ops
2. **User Story 4 (Destroy - P0)**: Depends on US2 dashboard (T036-T041) - Completes baseline

**ENHANCEMENTS (beyond 20i-gui)**:
3. **User Story 1 (Dashboard - P1)**: Depends on US2 container model (T023-T026) - Adds monitoring
4. **User Story 3 (Logs - P2)**: Depends on US2 dashboard layout - Adds log viewer
5. **User Story 5 (Projects - P3)**: Independent after Foundational - Adds multi-project

### Recommended Execution Order

**MINIMUM VIABLE PRODUCT (20i-gui replacement)**:
1. Phase 1: Setup (T001-T013) - Initialize project with test infrastructure
2. Phase 2: Foundational (T014-T025) - Docker client + root model with tests
3. Phase 3: User Story 2 - Lifecycle (T026-T072) - MVP with comprehensive tests
4. Phase 4: User Story 4 - Destroy (T073-T089) - Baseline complete with regression tests
5. Phase 8: Polish (T155-T182) - Production-ready with full test suite

**STOP HERE FOR v1.0 RELEASE** - You now have full 20i-gui functionality in TUI with >85% test coverage

**ENHANCED VERSION (add monitoring & debugging)**:
6. Phase 5: User Story 1 - Dashboard monitoring (T090-T110) - Live stats with performance tests
7. Phase 6: User Story 3 - Log viewer (T111-T137) - Streaming logs with integration tests
8. Phase 7: User Story 5 - Project switcher (T138-T154) - Multi-project with e2e tests

**Total Tasks**: 182 (includes ~60 test tasks across all phases)
7. Phase 6: User Story 3 - Logs (T084-T103) ← Adds log viewer
8. Phase 7: User Story 5 - Projects (T104-T117) ← Adds multi-project

**Parallel (multiple developers)**:
- Team: Complete Phase 1 + Phase 2 together (T001-T022)
- After Phase 2:
  - Dev A: User Story 2 Lifecycle (T023-T055) ← PRIORITY, BLOCKS OTHERS
  - Once US2 complete:
    - Dev B: User Story 4 Destroy (T058-T069) ← Completes baseline
    - Dev C: User Story 1 Dashboard enhancement (T070-T083) ← Can run parallel
    - Dev D: User Story 3 Logs (T084-T103) ← Can run parallel  
    - Dev E: User Story 5 Projects (T104-T117) ← Can run parallel
- Team: Complete Phase 8 together (T118-T140)

### Within Each User Story

- Models/entities before services (e.g., T023-T026 before T027)
- Services before UI components (e.g., T027 before T029)
- Core implementation before integration (e.g., T036 before T038)

### Parallel Opportunities

**Setup Phase (T001-T011)**:
- T003-T006 can run in parallel (different dependencies)
- T009-T010 can run in parallel (different files)

**Foundational Phase (T012-T022)**:
- T012-T014 can run in parallel (different files in ui/)
- T017 can run in parallel with T015-T016 (different files)

**User Story 2 - Lifecycle (T023-T055)** - MVP:
- T023-T024 can run in parallel (entity structs)
- T027-T032 can run in parallel (Docker client methods)
- T033-T035 can run in parallel (message types)
- T047-T049 can run in parallel (command functions)

**User Story 4 - Destroy (T058-T069)**:
- T058 can run in parallel with T059 (different files)

**User Story 1 - Dashboard Enhancement (T070-T083)**:
- T070-T071 can run in parallel (different entity structs)
- T073-T074 can run in parallel (stats.go methods)
- T078 can run in parallel with stats work

**User Story 3 - Logs (T084-T103)**:
- T084-T086 can run in parallel (LogStream methods)
- T088 can run in parallel with T084-T086

**User Story 5 - Projects (T104-T117)**:
- T104-T106 can run in parallel (entity and methods)
- T110 can run in parallel with other tasks

**Polish Phase (T125-T150)**:
- T125-T128 can run in parallel (different documentation files)
- T130-T131 can run in parallel (different components)
- T137-T139 can run in parallel (different files)

---

## Parallel Example: User Story 2 (Lifecycle - MVP)

```bash
# After Foundational complete, launch these in parallel:
Task T023: Create Container struct (minimal fields)
Task T024: Create ContainerStatus enum

# After entities, launch Docker client methods in parallel:
Task T027: Implement StartContainer
Task T028: Implement StopContainer
Task T029: Implement RestartContainer
Task T030: Implement ComposeStop
Task T031: Implement ComposeRestart
Task T032: Implement ComposeDown

# While Docker methods being built, create message types in parallel:
Task T033: Add ContainerAction enum
Task T034: Add ComposeAction enum
Task T035: Add action message types

# After all Docker methods ready, build UI commands in parallel:
Task T047: Create startContainerCmd
Task T048: Create stopContainerCmd
Task T049: Create restartContainerCmd
```

---

## Implementation Strategy

### MVP First (Lifecycle + Destroy = 20i-gui Baseline) ✅ RECOMMENDED

1. Complete Phase 1: Setup (T001-T011)
2. Complete Phase 2: Foundational (T012-T022) ← FOUNDATION READY
3. Complete Phase 3: User Story 2 - Lifecycle (T023-T055) ← **MVP CHECKPOINT** 🎯
4. Complete Phase 4: User Story 4 - Destroy (T056-T067) ← **BASELINE COMPLETE** ✅
5. Complete Phase 8: Polish (T118-T140) ← Production-ready
6. **VALIDATE**: Start stack, verify running, stop stack, verify stopped, destroy stack
7. **RELEASE v1.0**: Full 20i-gui replacement in TUI

**Stop here for first release** - You have all essential 20i-gui functionality

### Enhanced Version (Add Monitoring & Debugging)

8. Add Phase 5: User Story 1 - Dashboard monitoring (T068-T083) ← Live stats
9. Add Phase 6: User Story 3 - Log viewer (T084-T103) ← Debugging
10. Add Phase 7: User Story 5 - Projects (T104-T117) ← Multi-project
11. **RELEASE v1.1**: Enhanced TUI with monitoring, logs, multi-project

### Incremental Delivery Strategy

Each phase delivers testable value:

- **After Phase 2 (Foundational)**: Docker client connects, can list containers
- **After Phase 3 (US2 Lifecycle)**: ← **MVP!** Can start/stop/restart services, verify status ✅
- **After Phase 4 (US4 Destroy)**: ← **BASELINE COMPLETE** Full 20i-gui parity ✅
- **After Phase 5 (US1 Dashboard)**: Live CPU/memory monitoring added
- **After Phase 6 (US3 Logs)**: Log viewer with follow mode added
- **After Phase 7 (US5 Projects)**: Multi-project switching added
- **After Phase 8 (Polish)**: Production-ready with docs and error handling

---

## Task Summary

- **Total Tasks**: 187 (includes comprehensive testing throughout all phases)
- **Setup Phase**: 13 tasks (T001-T013) - includes test infrastructure
- **Foundational Phase**: 17 tasks (T014-T025d) - includes foundation tests + architectural decisions
- **User Story 2 (Lifecycle) - MVP**: 47 tasks (T026-T072) 🎯 - includes unit + integration tests
- **User Story 4 (Destroy) - Baseline**: 17 tasks (T073-T089) ✅ - includes regression tests
- **User Story 1 (Dashboard Enhancement)**: 21 tasks (T090-T110) - includes performance tests
- **User Story 3 (Logs)**: 27 tasks (T111-T137) - includes streaming tests
- **User Story 5 (Projects)**: 17 tasks (T138-T154) - includes e2e tests
- **Polish Phase**: 28 tasks (T155-T182) - includes comprehensive test suite

**Test Coverage Breakdown**:
- Unit tests: ~30 tasks (Docker client, UI components, message handlers, entities)
- Integration tests: ~15 tasks (lifecycle workflow, destroy workflow, logs streaming, projects)
- Performance tests: ~5 tasks (stats refresh, log rendering, startup time, memory usage)
- E2E tests: ~5 tasks (full user journeys using Bubble Tea test utilities)
- Manual acceptance tests: ~5 tasks (per user story acceptance criteria)
- Total test tasks: ~60 (33% of all tasks)

**Parallel Opportunities**:
- Setup: 4 tasks can run in parallel (Go modules, deps, test mocks)
- Foundational: 5 tasks can run in parallel (Docker client + tests, root model + tests)
- User Story 2: 12 tasks can run in parallel (Docker methods + tests, messages + tests, UI + tests)
- User Story 4: 3 tasks can run in parallel (modal component + tests)
- User Story 1: 8 tasks can run in parallel (entity/view/stats work + tests)
- User Story 3: 6 tasks can run in parallel (log streaming + tests)
- User Story 5: 5 tasks can run in parallel (project detection + tests)
- Polish: 8 tasks can run in parallel (styling, e2e tests, docs)

**Estimated Time with Testing** (US2+US4 baseline):
- Single developer: 25-29 hours (includes writing tests + architectural decisions)
- Team of 3: 12-14 hours (parallel execution after Foundational)

**Test Execution Strategy**:
- Run unit tests after each component implementation (`go test ./internal/...`)
- Run integration tests at end of each user story phase
- Run regression suite after each phase to ensure no breakage
- Run full test suite (`make test-all`) before merging to main
- Target: >85% code coverage for production readiness

**Independent Test Criteria**:
- **US1**: Dashboard displays 4 services with status icons and live CPU/memory bars
- **US2**: Pressing 's' on apache stops it (gray), pressing again starts (green)
- **US3**: Pressing 'l' opens logs, 'f' enables follow, new requests appear in real-time
- **US4**: Pressing 'D' shows warning, typing "yes" removes all volumes and containers
- **US5**: Pressing 'p' shows project list, selecting switches context, dashboard reloads

**📚 Remember**: Keep `/runbooks/research/QUICK-REFERENCE.md` open while implementing ALL tasks!

**Format Validation**: ✅ All 182 tasks follow strict checklist format:
- ✅ Checkbox: `- [ ]` prefix on every task
- ✅ Task ID: T001-T140 sequential (T141-T142 reserved for future)
- ✅ [P] marker: Only on parallelizable tasks (different files, no blockers)
- ✅ [Story] label: US1-US5 on all user story phase tasks
- ✅ File paths: Exact paths included in descriptions (e.g., internal/docker/client.go)
- ✅ No story label on Setup, Foundational, Polish phases (correct)

---

## Notes

- **[P] tasks**: Different files, no dependencies within their phase - can run in parallel
- **[Story] labels**: Map tasks to user stories for traceability and independent testing
- **Checkpoints**: Stop after Phase 2, Phase 3, and Phase 6 to validate progress
- **MVP scope**: Phases 1-4 replicate all 20i-gui CORE functionality (US2 Lifecycle + US4 Destroy) ✅
- **Enhancements**: Phases 5-7 add monitoring (US1), logs (US3), multi-project (US5) beyond 20i-gui
- **Production-ready**: Phase 8 adds documentation, error handling, polish
- **Priority**: GET STACKS RUNNING FIRST (lifecycle), then add nice-to-haves (monitoring/logs)
- **Tests**: Not included per spec (no TDD requirement, manual testing via quickstart.md)
- **File organization**: All tasks reference exact file paths from plan.md structure
- **Bubble Tea patterns**: Tasks follow Elm Architecture (Model-Update-View) per research findings
- **Docker SDK**: All Docker operations use SDK (no shell commands) per contract
- **Lipgloss**: All styling uses Lipgloss (no raw ANSI codes) per NFR requirements
