# Tasks: 20i Stack Manager TUI

**Feature Branch**: `001-stack-manager-tui`  
**Input**: Design documents from `/specs/001-stack-manager-tui/`  
**Prerequisites**: plan.md, spec.md (5 user stories), research.md, data-model.md, contracts/ (docker-api.md, ui-events.md)

**Tests**: NOT REQUESTED - Tests omitted from task list (no TDD requirement in spec)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

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

- [ ] T001 Create tui/ directory at repository root
- [ ] T002 Initialize Go module with `go mod init github.com/peternicholls/20i-stack/tui`
- [ ] T003 [P] Add Bubble Tea dependency v1.3.10+ in go.mod
- [ ] T004 [P] Add Bubbles dependency v1.0.0+ in go.mod
- [ ] T005 [P] Add Lipgloss dependency v1.0.0+ in go.mod
- [ ] T006 [P] Add Docker SDK dependency v27.0.0+ in go.mod
- [ ] T007 Run `go mod tidy` to generate go.sum
- [ ] T008 Create directory structure: internal/app, internal/views/dashboard, internal/views/help, internal/views/projects, internal/docker, internal/ui
- [ ] T009 [P] Create placeholder tui/main.go with basic Bubble Tea hello world
- [ ] T010 [P] Create Makefile with build, install, clean targets
- [ ] T011 Verify `go run main.go` works (shows hello world, press 'q' to quit)

**Checkpoint**: Go project initializes, dependencies resolve, basic TUI runs

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T012 [P] Create internal/ui/styles.go with Lipgloss color palette (ColorRunning, ColorStopped, ColorError, ColorAccent, ColorBorder)
- [ ] T013 [P] Create internal/ui/components.go with StatusIcon function (maps ContainerStatus enum to emoji icons)
- [ ] T014 [P] Create internal/ui/layout.go with panel sizing functions (calculatePanelWidths, calculatePanelHeights)
- [ ] T015 Create internal/docker/client.go with Client struct and NewClient() method per docker-api.md contract
- [ ] T016 Implement Docker connection check in client.go (Ping method, error handling for daemon unreachable)
- [ ] T017 [P] Create internal/app/messages.go with custom tea.Msg types: statsMsg, containerListMsg, logLineMsg, containerActionMsg, containerActionResultMsg per ui-events.md
- [ ] T018 Create internal/app/root.go with RootModel struct (activeView, dashboard, help, projects, dockerClient fields)
- [ ] T019 Implement RootModel.Init() method (initialize Docker client, return initial commands)
- [ ] T020 Implement RootModel.Update() method with global shortcut routing (q=quit, ?=help, p=projects)
- [ ] T021 Implement RootModel.View() method with modal overlay logic
- [ ] T022 Update tui/main.go to create and run RootModel instead of hello world

**Checkpoint**: Foundation ready - Docker client connects, root model routes messages, can extend with views

---

## Phase 3: User Story 2 - Container Lifecycle (Priority: P0 - Core) 🎯 MVP FIRST

**Goal**: Start, stop, restart individual services or entire stack - REPLICATE 20i-gui core functionality

**Independent Test**: Run TUI, see service list, press 's' on stopped apache, verify it starts; press 'S' to stop all, verify stack stops

**Rationale**: This is the PRIMARY use case - get stacks running and verified. Dashboard monitoring is secondary.

### Implementation for User Story 2 - LIFECYCLE FIRST

- [ ] T023 [P] [US2] Create Container entity struct in internal/docker/client.go (ID, Name, Service, Image, Status, State fields - minimal for lifecycle)
- [ ] T024 [P] [US2] Create ContainerStatus enum in internal/docker/client.go (Running, Stopped, Restarting, Error)
- [ ] T025 [US2] Implement mapDockerState() helper function in internal/docker/client.go to map Docker states to ContainerStatus enum
- [ ] T026 [US2] Implement ListContainers(projectName string) method in internal/docker/client.go per docker-api.md contract
- [ ] T027 [P] [US2] Implement StartContainer(containerID string) method in internal/docker/client.go per docker-api.md contract
- [ ] T028 [P] [US2] Implement StopContainer(containerID string, timeout int) method in internal/docker/client.go
- [ ] T029 [P] [US2] Implement RestartContainer(containerID string, timeout int) method in internal/docker/client.go
- [ ] T030 [P] [US2] Implement ComposeStop(projectPath string) method in internal/docker/client.go per docker-api.md
- [ ] T031 [P] [US2] Implement ComposeRestart(projectPath string) method in internal/docker/client.go
- [ ] T032 [P] [US2] Implement ComposeDown(projectPath string, removeVolumes bool) method in internal/docker/client.go
- [ ] T033 [US2] Add ContainerAction enum to internal/app/messages.go (Start, Stop, Restart)
- [ ] T034 [US2] Add ComposeAction enum to internal/app/messages.go (StopAll, RestartAll, Destroy)
- [ ] T035 [US2] Add composeActionMsg and composeActionResultMsg types to internal/app/messages.go
- [ ] T036 [P] [US2] Create internal/views/dashboard/dashboard.go with DashboardModel struct (serviceList, containers, selectedIndex fields - NO stats yet)
- [ ] T037 [US2] Implement DashboardModel.Init() method to load container list
- [ ] T038 [US2] Implement containerListMsg handling in DashboardModel.Update()
- [ ] T039 [US2] Create internal/views/dashboard/service_list.go with simple list rendering (status icon + name only)
- [ ] T040 [US2] Implement DashboardModel.View() with simple 2-panel layout (service list 30% | status messages 70% | footer)
- [ ] T041 [US2] Wire DashboardModel into RootModel in internal/app/root.go
- [ ] T042 [US2] Implement navigation keys (↑/k=up, ↓/j=down) in DashboardModel.Update()
- [ ] T043 [US2] Implement 's' key handler to toggle start/stop for selected container
- [ ] T044 [US2] Implement 'r' key handler to restart selected container
- [ ] T045 [US2] Implement 'S' key handler to stop all stack containers (with simple confirmation)
- [ ] T046 [US2] Implement 'R' key handler to restart entire stack
- [ ] T047 [US2] Create startContainerCmd() function in dashboard.go to launch async Docker operation
- [ ] T048 [US2] Create stopContainerCmd() function in dashboard.go
- [ ] T049 [US2] Create restartContainerCmd() function in dashboard.go
- [ ] T050 [US2] Implement containerActionResultMsg handler with success/error feedback
- [ ] T051 [US2] Add status message panel to show "✅ Container started" or "❌ Failed: error"
- [ ] T052 [US2] Trigger containerListMsg refresh after successful action
- [ ] T053 [US2] Implement error message formatting per docker-api.md (user-friendly port conflicts, timeouts)
- [ ] T054 [US2] Add footer with basic shortcuts: "s:start/stop  r:restart  S:stop-all  R:restart-all  D:destroy  q:quit"
- [ ] T055 [US2] Test lifecycle: start stopped container, verify status changes; stop running container, verify status changes

**Checkpoint**: Core 20i-gui functionality working - can start/stop/restart containers and verify status

---

## Phase 4: User Story 4 - Destroy Stack (Priority: P0 - Core)

**Goal**: Destroy stack (stop containers + remove volumes) - COMPLETE 20i-gui baseline

**Independent Test**: Press 'D', confirm with 'yes', verify stack destroyed and volumes removed

### Implementation for User Story 4

- [ ] T056 [P] [US4] Create ConfirmationModal component in internal/ui/components.go with text input and warning styling
- [ ] T057 [US4] Add confirmationModal field to DashboardModel in dashboard.go
- [ ] T058 [US4] Implement 'D' key handler in DashboardModel.Update() to show destroy confirmation modal
- [ ] T059 [US4] Render confirmation modal overlay with warning "⚠️ This will REMOVE ALL VOLUMES and data. Type 'yes' to confirm:"
- [ ] T060 [US4] Add text input to confirmation modal using Bubbles textinput.Model component
- [ ] T061 [US4] Implement confirmation modal input handling (type "yes", press Enter to confirm)
- [ ] T062 [US4] Implement Esc handler in confirmation modal to cancel without destroying
- [ ] T063 [US4] Create composeDownCmd() function in dashboard.go to call ComposeDown with removeVolumes=true
- [ ] T064 [US4] Implement composeActionResultMsg handler for Destroy action
- [ ] T065 [US4] Show success message "✅ Stack destroyed" after ComposeDown completes
- [ ] T066 [US4] Refresh container list after destroy (should show empty or no running containers)
- [ ] T067 [US4] Update footer to show "D:destroy" shortcut

**Checkpoint**: ✅ BASELINE COMPLETE - All 20i-gui core functions replicated (start/stop/restart/destroy)

---

## Phase 5: User Story 1 - Dashboard Monitoring (Priority: P1 - Enhancement)

**Goal**: Add live CPU/memory monitoring and detailed container info - ENHANCES baseline beyond 20i-gui

**Independent Test**: After lifecycle working, verify CPU% and memory bars update every 2s, detail panel shows ports/image/uptime

### Implementation for User Story 1 - Dashboard Enhancement

- [ ] T068 [P] [US1] Create Stats entity struct in internal/docker/client.go (CPUPercent, MemoryUsed, MemoryLimit, MemoryPercent, NetworkRxBytes, NetworkTxBytes, Timestamp)
- [ ] T069 [P] [US1] Create PortMapping entity struct in internal/docker/client.go (ContainerPort, HostPort, Protocol)
- [ ] T070 [US1] Add Ports, CreatedAt, StartedAt fields to Container entity (expand from minimal US2 version)
- [ ] T071 [US1] Implement WatchStats() method in internal/docker/stats.go per docker-api.md
- [ ] T072 [US1] Implement calculateCPUPercent() and calculateMemoryPercent() helper functions in internal/docker/stats.go
- [ ] T073 [US1] Add stats field to DashboardModel (map[string]Stats)
- [ ] T074 [US1] Implement tickMsg handler in DashboardModel.Update() to trigger stats refresh every 2s
- [ ] T075 [US1] Implement statsMsg handler to update stats map
- [ ] T076 [US1] Update service_list.go to display CPU/memory bars from stats map
- [ ] T077 [US1] Add ProgressBar component to internal/ui/components.go for CPU/memory visualization
- [ ] T078 [P] [US1] Create internal/views/dashboard/detail.go with DetailPanel struct
- [ ] T079 [US1] Implement detail panel rendering (image, ports, uptime, container ID, volumes)
- [ ] T080 [US1] Expand DashboardModel.View() from 2-panel to 3-panel layout (service list 25% | detail 50% | status 25% | footer)
- [ ] T081 [US1] Implement Enter key handler to show detail panel for selected container
- [ ] T082 [US1] Implement Tab key to cycle focus between panels
- [ ] T083 [US1] Update footer to show "Enter:detail  Tab:panels" shortcuts

**Checkpoint**: Dashboard enhanced with live monitoring - CPU/memory stats, detailed info panel

---

## Phase 6: User Story 3 - Log Viewer (Priority: P2 - Enhancement)

**Goal**: View live container logs with follow mode - ADDS debugging capability beyond 20i-gui

**Independent Test**: Press 'l' on running container, verify logs show; press 'f' for follow mode, make web request, see new log line

### Implementation for User Story 3

- [ ] T084 [P] [US3] Create LogStream entity struct in internal/docker/client.go (ContainerID, Buffer, Following, FilterText, Head, Size)
- [ ] T085 [US3] Implement LogStream.Append() method for ring buffer management
- [ ] T086 [US3] Implement LogStream.GetFilteredLines() method with search filter support
- [ ] T087 [US3] Implement StreamLogs(containerID string, since time.Time, follow bool) method in internal/docker/client.go
- [ ] T088 [P] [US3] Create internal/views/dashboard/logs.go with LogPanel struct
- [ ] T089 [US3] Add logPanel and logVisible fields to DashboardModel
- [ ] T090 [US3] Implement 'l' key handler to toggle log panel visibility
- [ ] T091 [US3] Create toggleLogPanelMsg type in internal/app/messages.go
- [ ] T092 [US3] Implement log panel opening: create Bubbles viewport.Model, resize layout
- [ ] T093 [US3] Implement streamLogsCmd() to launch background log streaming
- [ ] T094 [US3] Implement logLineMsg handler to append lines to viewport
- [ ] T095 [US3] Implement renderLogs() method in logs.go
- [ ] T096 [US3] Update View() to show logs panel when visible (detail 30%, logs 70%)
- [ ] T097 [US3] Implement 'f' key handler to toggle follow mode
- [ ] T098 [US3] Implement auto-scroll when Following=true
- [ ] T099 [US3] Implement '/' key for search/filter mode
- [ ] T100 [US3] Implement scroll navigation (↑/↓/j/k/g/G)
- [ ] T101 [US3] Implement Esc/q to close logs panel
- [ ] T102 [US3] Update footer with log shortcuts when visible
- [ ] T103 [US3] Implement log buffer limit (10k lines)

**Checkpoint**: Log viewer functional with follow mode and search

---

## Phase 7: User Story 5 - Project Switcher (Priority: P3 - Enhancement)

**Goal**: Multi-project support - OPTIONAL enhancement for power users

**Independent Test**: Press 'p', see project list, select different project, verify dashboard switches context

### Implementation for User Story 5

- [ ] T104 [P] [US5] Create Project entity struct in internal/docker/client.go (Name, Path, ComposeFile, IsActive, ContainerCount, Is20iStack)
- [ ] T105 [US5] Implement is20iStack() detection function (check for .20i-local OR apache+mariadb+nginx)
- [ ] T106 [US5] Implement GetComposeProject(composeFilePath string) method
- [ ] T107 [US5] Create internal/docker/filters.go with ScanForProjects() function
- [ ] T108 [US5] Implement project detection filtering (20i stacks only)
- [ ] T109 [US5] Implement GetContainerCount() function in filters.go
- [ ] T110 [P] [US5] Create internal/views/projects/projects.go with ProjectListModel
- [ ] T111 [US5] Implement ProjectListModel.Init() to scan projects
- [ ] T112 [US5] Implement navigation and selection in ProjectListModel.Update()
- [ ] T113 [US5] Implement ProjectListModel.View() with project list
- [ ] T114 [US5] Add 'p' key handler in RootModel to open project switcher
- [ ] T115 [US5] Implement project switching logic
- [ ] T116 [US5] Clear state when switching projects
- [ ] T117 [US5] Update header to show current project name

**Checkpoint**: Project switcher working (optional feature)

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Production-ready improvements for release

- [ ] T118 [P] Create tui/README.md with installation instructions
- [ ] T119 [P] Document usage in tui/README.md (keyboard shortcuts, requirements)
- [ ] T120 [P] Add "TUI Interface" section to repository README.md
- [ ] T121 [P] Add CHANGELOG.md entry for v2.0.0
- [ ] T122 Update .gitignore to exclude tui/20i-tui binary
- [ ] T123 [P] Implement Docker daemon retry logic (auto-retry every 5s, show error screen)
- [ ] T124 [P] Create error screen component in internal/ui/components.go
- [ ] T125 Implement terminal size validation (minimum 80x24)
- [ ] T126 Implement SIGWINCH handler for terminal resize
- [ ] T127 [P] Create internal/views/help/help.go with HelpModel
- [ ] T128 Implement '?' key handler in RootModel for help modal
- [ ] T129 Implement help modal View() with keyboard shortcuts
- [ ] T130 [P] Add inline code comments for Bubble Tea patterns
- [ ] T131 [P] Add inline comments for Docker SDK usage
- [ ] T132 [P] Add docstrings to all public functions
- [ ] T133 Implement graceful shutdown on Ctrl-C
- [ ] T134 Test with 4 running containers to verify functionality
- [ ] T135 Test terminal resize behavior
- [ ] T136 Test all keyboard shortcuts (s, r, S, R, D, q, arrows, vim keys)
- [ ] T137 Verify error messages are user-friendly
- [ ] T138 Run lifecycle validation: start stack, verify running, stop stack, verify stopped, destroy stack
- [ ] T139 Add build target to Makefile for /usr/local/bin
- [ ] T140 Create symlink 'tui' → '20i-tui' in Makefile

**Checkpoint**: Production-ready TUI that replicates 20i-gui core functions

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
1. Phase 1: Setup (T001-T011)
2. Phase 2: Foundational (T012-T022) ← FOUNDATION COMPLETE
3. Phase 3: User Story 2 - Lifecycle (T023-T055) ← **MVP CHECKPOINT** ✅
4. Phase 4: User Story 4 - Destroy (T056-T067) ← **BASELINE COMPLETE** ✅
5. Phase 8: Polish (T118-T140) ← Production-ready

**STOP HERE FOR v1.0 RELEASE** - You now have full 20i-gui functionality in TUI

**ENHANCED VERSION (add monitoring & debugging)**:
6. Phase 5: User Story 1 - Dashboard (T068-T083) ← Adds live stats
7. Phase 6: User Story 3 - Logs (T084-T103) ← Adds log viewer
8. Phase 7: User Story 5 - Projects (T104-T117) ← Adds multi-project

**Parallel (multiple developers)**:
- Team: Complete Phase 1 + Phase 2 together (T001-T022)
- After Phase 2:
  - Dev A: User Story 2 Lifecycle (T023-T055) ← PRIORITY, BLOCKS OTHERS
  - Once US2 complete:
    - Dev B: User Story 4 Destroy (T056-T067) ← Completes baseline
    - Dev C: User Story 1 Dashboard enhancement (T068-T083) ← Can run parallel
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

**User Story 4 - Destroy (T056-T067)**:
- T056 can run in parallel with T057 (different files)

**User Story 1 - Dashboard Enhancement (T068-T083)**:
- T068-T069 can run in parallel (different entity structs)
- T071-T072 can run in parallel (stats.go methods)
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

- **Total Tasks**: 150
- **Setup Phase**: 11 tasks (T001-T011)
- **Foundational Phase**: 11 tasks (T012-T022)
- **User Story 1 (Dashboard)**: 24 tasks (T023-T046)
- **User Story 2 (Lifecycle)**: 21 tasks (T047-T067)
- **User Story 3 (Logs)**: 24 tasks (T068-T091)
- **User Story 4 (Destroy)**: 13 tasks (T092-T104)
- **User Story 5 (Projects)**: 20 tasks (T105-T124)
- **Polish Phase**: 26 tasks (T125-T150)

**Parallel Opportunities**:
- Setup: 4 tasks can run in parallel
- Foundational: 5 tasks can run in parallel
- User Story 1: 10 tasks can run in parallel across entity/view/stats work
- User Story 2: 9 tasks can run in parallel (Docker methods + messages)
- User Story 3: 5 tasks can run in parallel
- User Story 4: 2 tasks can run in parallel
- User Story 5: 5 tasks can run in parallel
- Polish: 8 tasks can run in parallel

**Estimated Time to MVP** (US1-4 complete):
- Single developer: 18-20 hours (per quickstart.md)
- Team of 3: 8-10 hours (parallel execution after Foundational)

**Independent Test Criteria**:
- **US1**: Dashboard displays 4 services with status icons and live CPU/memory bars
- **US2**: Pressing 's' on apache stops it (gray), pressing again starts (green)
- **US3**: Pressing 'l' opens logs, 'f' enables follow, new requests appear in real-time
- **US4**: Pressing 'D' shows warning, typing "yes" removes all volumes and containers
- **US5**: Pressing 'p' shows project list, selecting switches context, dashboard reloads

**Format Validation**: ✅ All 150 tasks follow strict checklist format:
- ✅ Checkbox: `- [ ]` prefix on every task
- ✅ Task ID: T001-T150 sequential
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
