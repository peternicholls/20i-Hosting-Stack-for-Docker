# Feature Specification: 20i Stack Manager TUI

**Feature Branch**: `001-stack-manager-tui`  
**Created**: 2025-12-28  
**Updated**: 2025-12-28  
**Status**: Draft  
**Priority**: 🔴 Critical  
**Input**: User description: "Full-featured terminal UI for managing 20i stack: containers, config, logs, monitoring, projects"

---

## Overview

A professional terminal user interface (TUI) built with Bubble Tea framework to replace and enhance the existing 20i-gui bash script. This MVP replicates all 20i-gui functionality (start/stop/restart/status/logs/destroy) with a modern, keyboard-driven interface following best practices from lazydocker, lazygit, k9s, and gh-dash.

**Phase 1 Scope** (MVP - this spec):
- Dashboard view with 3-panel layout (service list | detail | help footer)
- Container lifecycle operations (start, stop, restart, remove)
- Real-time status and resource monitoring
- Log viewer with follow mode
- Project selection and switching

**Phase 2+** (future specs):
- Configuration editor (replaces manual .20i-local editing)
- Image management (pull, remove, list)
- Advanced monitoring with graphs
- Custom commands and plugins

**Design Principles** (from research):
- **Panel-based layout** (not tabs) - see multiple contexts simultaneously
- **Component composition** - each view is a standalone Bubble Tea model
- **Keyboard-first** - vim bindings + arrow keys, max 3 keystrokes to any action
- **Real-time updates** - background goroutines, never block UI
- **Progressive disclosure** - show essentials first, details on demand
- **Context-aware help** - footer shows current view's shortcuts

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Dashboard Overview (Priority: P0 - Core)

**Replaces**: 20i-gui "View Status" command

As a developer, I want to see an at-a-glance dashboard of all 20i stack services with status, CPU, and memory so that I can quickly assess stack health without running `docker ps`.

**Why this priority**: The dashboard is the entry point and primary view - users spend 80% of time here.

**Independent Test**: Start TUI from project directory with running stack, verify 4 services display with green status, CPU bars, memory usage.

**Acceptance Scenarios**:

1. **Given** a running 20i stack, **When** the TUI opens, **Then** the dashboard shows all 4 services (apache, mariadb, nginx, phpmyadmin) with color-coded status
2. **Given** containers are running, **When** viewing the dashboard, **Then** CPU and memory metrics update every 2s without blocking user input
3. **Given** a stopped container, **When** viewing the dashboard, **Then** it shows gray/dim with "Stopped" label
4. **Given** the dashboard, **When** I select a service and press Enter, **Then** the detail panel shows ports, image, uptime, container ID

---

### User Story 2 - Container Lifecycle (Priority: P0 - Core)

**Replaces**: 20i-gui "Start Stack", "Stop Stack", "Restart Stack" commands

As a developer, I want to start, stop, restart individual services or the entire stack so that I can control my development environment without leaving the TUI.

**Why this priority**: Lifecycle management is the primary use case - matches 20i-gui menu options 1-3.

**Independent Test**: Select apache service, press `s` to stop, verify status changes to gray "Stopped"; press `s` again, verify starts and turns green.

**Acceptance Scenarios**:

1. **Given** the dashboard with apache running, **When** I press `s`, **Then** apache container stops and status updates to "Stopped" with gray color
2. **Given** apache stopped, **When** I press `s`, **Then** it starts and status shows "Running" with green color
3. **Given** any service selected, **When** I press `r`, **Then** the container restarts (stop + start) with "Restarting..." feedback
4. **Given** the dashboard, **When** I press `S` (shift-s), **Then** all stack containers stop with confirmation prompt
5. **Given** the dashboard, **When** I press `R` (shift-r), **Then** entire stack restarts (docker compose restart)

---

### User Story 3 - Log Viewer (Priority: P0 - Core)

**Replaces**: 20i-gui "View Logs" command

As a developer, I want to view live container logs with follow mode so that I can debug issues in real-time without running `docker compose logs -f`.

**Why this priority**: Log viewing is critical for debugging - one of the 6 core 20i-gui commands.

**Independent Test**: Select apache, press `l`, verify logs panel opens showing last 100 lines; press `f` to enable follow, make web request, verify new log line appears.

**Acceptance Scenarios**:

1. **Given** the dashboard with apache selected, **When** I press `l`, **Then** the bottom panel shows last 100 lines of apache logs
2. **Given** the log viewer open, **When** I press `f`, **Then** follow mode enables and new log entries auto-scroll
3. **Given** follow mode enabled, **When** new log lines appear, **Then** viewport auto-scrolls to bottom
4. **Given** the log viewer, **When** I press `Esc` or `q`, **Then** log panel closes and returns to dashboard
5. **Given** the log viewer, **When** I press `/`, **Then** a search input appears and filters logs as I type

---

### User Story 4 - Destroy Stack (Priority: P0 - Core)

**Replaces**: 20i-gui "Destroy Stack" command (menu option 6)

As a developer, I want to destroy a stack (stop containers + remove volumes) so that I can clean up projects or reset database state.

**Why this priority**: One of the 6 core 20i-gui commands - needed for cleanup and fresh starts.

**Independent Test**: Press `D` (shift-d), verify confirmation modal shows "⚠️ This will REMOVE ALL VOLUMES and data", type "yes", verify stack destroyed.

**Acceptance Scenarios**:

1. **Given** the dashboard with running stack, **When** I press `D` (shift-d), **Then** a confirmation modal appears with warning about data loss
2. **Given** the destroy confirmation, **When** I type "yes", **Then** `docker compose down -v` runs and all containers + volumes are removed
3. **Given** the destroy confirmation, **When** I press `Esc` or type anything else, **Then** operation cancels and modal closes
4. **Given** stack destroyed, **When** operation completes, **Then** success message shows and dashboard updates to show no running containers

---

### User Story 5 - Project Switcher (Priority: P1 - Important)

**Replaces**: 20i-gui multi-project selection (shown when stopping/viewing logs)

As a developer, I want to see all detected 20i projects and switch between them so that I can manage multiple projects from one TUI instance.

**Why this priority**: Power users have multiple projects - 20i-gui already supports this via project selection prompts.

**Independent Test**: Press `:projects` or `p`, verify list shows current directory and other projects with docker-compose.yml, select different project, verify dashboard switches.

**Acceptance Scenarios**:

1. **Given** the dashboard, **When** I press `p`, **Then** a project list modal shows all directories with docker-compose.yml in parent/sibling folders
2. **Given** the project list, **When** I select a project and press Enter, **Then** TUI switches context and dashboard reloads with that project's containers
3. **Given** the project list, **When** viewing, **Then** current project is highlighted and marked with `[current]`
4. **Given** a project has running containers, **When** shown in list, **Then** it displays running container count (e.g., "myproject (4 running)")
5. **Given** the project list, **When** I press `Esc`, **Then** modal closes without switching projects

---

### Deferred Stories (Phase 2+)

These features are intentionally excluded from MVP to keep scope focused:

- **Configuration Editor**: Edit .20i-local, .env, stack-vars.yml files in TUI (currently done manually)
- **PHP Version Validation**: Validate PHP_VERSION against Docker Hub tags with EOL warnings
- **Extended YAML Config**: Support HOST_PORT, MYSQL_PORT, ENABLE_REDIS, etc. in stack-vars.yml
- **Resource Monitoring**: Graphs/charts for CPU, memory, network over time
- **Image Management**: Pull, remove, list Docker images
- **Custom Commands**: User-defined shortcuts and scripts

**Rationale**: The existing 20i-gui doesn't have these features, so they're enhancements rather than baseline functionality. MVP focuses on matching 20i-gui plus modern TUI patterns.

---

### Edge Cases

- What happens when Docker daemon is not running? (show error screen with retry button)
- How does the TUI handle terminal resize events? (SIGWINCH - recalculate layout, minimum 80x24)
- What happens when a container action fails? (show inline error for 3s, maintain UI state)
- What happens if user starts TUI from non-project directory? (show error: "No docker-compose.yml found")
- What happens when project has no running containers? (dashboard shows empty list with hint "Press S to start stack")
- What happens when log buffer exceeds 10,000 lines? (truncate oldest lines, show "[truncated]" indicator)
- How does project switcher detect 20i stacks vs generic Docker Compose? (check for .20i-local or apache+mariadb+nginx services)

---

## Requirements *(mandatory)*

### Functional Requirements

**Core TUI Framework**
- **FR-001**: TUI MUST use Bubble Tea v1.3.10+ framework with Elm Architecture (Model-Update-View pattern)
- **FR-002**: TUI MUST use Bubbles components (list.Model, viewport.Model) for service list and log viewer
- **FR-003**: TUI MUST use Lipgloss for ALL styling (no raw ANSI codes)
- **FR-004**: TUI MUST run in alternate screen mode with clean restoration on exit (restore cursor, clear alt screen)
- **FR-005**: TUI MUST use 3-panel layout: service list (left 25%) | detail/logs (right 75%) | help footer (2 lines)

**Dashboard & Service List**
- **FR-010**: Dashboard MUST show all project containers with name, status (icon+color), CPU%, memory% in left panel
- **FR-011**: Service list MUST use Bubbles list.Model component with vim-style navigation (j/k, arrow keys)
- **FR-012**: Dashboard MUST auto-refresh container stats every 2s using background goroutine + tea.Tick
- **FR-013**: Selected service MUST show detail panel with: image, ports, uptime, container ID, volumes (read-only view)
- **FR-014**: Status indicators MUST use: 🟢 Green="Running", ⚪ Gray="Stopped", 🟡 Yellow="Restarting", 🔴 Red="Error"

**Container Lifecycle**
- **FR-020**: System MUST support single-service operations: start (`s`), stop (`s` toggle), restart (`r`)
- **FR-021**: System MUST support whole-stack operations: stop all (`S`), restart all (`R`), start all (if all stopped)
- **FR-022**: Destroy stack (`D`) MUST show confirmation modal requiring "yes" to proceed, run `docker compose down -v`
- **FR-023**: All operations MUST provide inline feedback: "Starting apache..." → "✅ apache started" (3s duration)
- **FR-024**: Failed operations MUST show error inline: "❌ Failed to start apache: port 80 in use" (persist until dismissed)

**Log Viewer**
- **FR-030**: Log viewer MUST open in bottom panel (replacing detail) when `l` pressed, use Bubbles viewport.Model
- **FR-031**: Log viewer MUST load last 100 lines on open, support follow mode toggle (`f`), auto-scroll when following
- **FR-032**: Log viewer MUST support search filter (`/`) that highlights matches as user types
- **FR-033**: Log buffer MUST cap at 10,000 lines per container, auto-truncate oldest lines
- **FR-034**: Log viewer MUST close with `Esc` or `q`, return to detail panel

**Project Switching**
- **FR-040**: Project switcher (`p`) MUST show modal list of directories with docker-compose.yml (search up 2 levels from current)
- **FR-041**: Project list MUST mark current project with `[current]` indicator
- **FR-042**: Project list MUST show running container count per project (e.g., "myproject (4 running)")
- **FR-043**: Switching projects MUST clear all state (logs, stats cache) and reload from new project directory
- **FR-044**: 20i stack detection: docker-compose.yml exists AND (services include apache+mariadb+nginx OR .20i-local exists)

**Navigation & Input**
- **FR-050**: Global shortcuts MUST work from any view: `?`=help, `q`=quit, `p`=projects, `:cmd`=command mode
- **FR-051**: Service list navigation MUST support: `↑/k`=up, `↓/j`=down, `Enter`=show detail, `Tab`=cycle panels
- **FR-052**: Help modal (`?`) MUST show context-aware shortcuts for current active panel
- **FR-053**: Footer MUST always show top 8 shortcuts for current context (e.g., "?:help  s:start/stop  r:restart  l:logs  q:quit")

**Error Handling**
- **FR-060**: When Docker daemon unreachable, MUST show error screen with "Docker not running" + retry button (`r`)
- **FR-061**: System MUST auto-retry Docker connection every 5s in background, show "Retrying..." indicator
- **FR-062**: Terminal < 80x24 MUST show error: "Terminal too small (need 80x24, got {w}x{h})" with resize hint
- **FR-063**: Missing docker-compose.yml MUST show error: "Not a Docker Compose project. Run from project directory."
- **FR-064**: All Docker API errors MUST be user-friendly (e.g., "port 80 in use" not "bind: address already in use")

---

### Component Architecture

**Following Bubble Tea + Research Best Practices**

```go
// Root model composes all views (Elm Architecture)
type RootModel struct {
    activeView   string            // "dashboard" | "logs" | "help" | "projects"
    dashboard    DashboardModel    // Main view (default)
    help         HelpModel         // Modal overlay
    projects     ProjectListModel  // Modal overlay
    dockerClient *docker.Client    // Shared Docker API wrapper
    err          error             // Global error state
}

func (m RootModel) Init() tea.Cmd {
    // Start background stats updater
    return tea.Batch(
        m.dashboard.Init(),
        tickEvery(2 * time.Second),  // Stats refresh
    )
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Global shortcuts
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "?":
            m.activeView = "help"
            return m, nil
        case "p":
            m.activeView = "projects"
            return m, m.projects.Load()
        }
        
        // Delegate to active view
        if m.activeView == "help" {
            updated, cmd := m.help.Update(msg)
            m.help = updated.(HelpModel)
            if m.help.closed {
                m.activeView = "dashboard"
            }
            return m, cmd
        }
        // ... delegate to other views
        
    case statsMsg:  // From background goroutine
        m.dashboard.UpdateStats(msg.stats)
        return m, tickEvery(2 * time.Second)  // Schedule next update
    }
    
    // Delegate to dashboard by default
    updated, cmd := m.dashboard.Update(msg)
    m.dashboard = updated.(DashboardModel)
    return m, cmd
}

func (m RootModel) View() string {
    // Modals overlay dashboard
    base := m.dashboard.View()
    if m.activeView == "help" {
        return overlayModal(base, m.help.View())
    }
    if m.activeView == "projects" {
        return overlayModal(base, m.projects.View())
    }
    return base
}

// Dashboard model (3-panel layout)
type DashboardModel struct {
    serviceList    list.Model       // Bubbles list (left panel)
    detailPanel    DetailPanel      // Custom component (right top)
    logPanel       *viewport.Model  // Bubbles viewport (right bottom, optional)
    logVisible     bool
    stats          map[string]Stats // Container stats cache
    width, height  int
}

func (m DashboardModel) View() string {
    leftPanel := m.renderServiceList()
    
    var rightPanel string
    if m.logVisible {
        rightPanel = lipgloss.JoinVertical(
            lipgloss.Top,
            m.renderDetail(),     // 30% height
            m.renderLogs(),       // 70% height
        )
    } else {
        rightPanel = m.renderDetail()  // 100% height
    }
    
    main := lipgloss.JoinHorizontal(
        lipgloss.Left,
        leftPanel,   // 25% width
        rightPanel,  // 75% width
    )
    
    footer := m.renderFooter()
    
    return lipgloss.JoinVertical(
        lipgloss.Top,
        main,
        footer,
    )
}
```

**File Structure**:
```
tui/
  main.go                    # Entry point, creates RootModel
  go.mod, go.sum             # Dependencies
  internal/
    app/
      root.go                # RootModel + routing logic
      messages.go            # Custom tea.Msg types (statsMsg, errorMsg)
    views/
      dashboard/
        dashboard.go         # DashboardModel
        service_list.go      # Service list panel
        detail.go            # Detail panel
        logs.go              # Log panel
      help/
        help.go              # HelpModel (modal)
      projects/
        projects.go          # ProjectListModel (modal)
    docker/
      client.go              # Docker SDK wrapper
      stats.go               # Background stats collector
      filters.go             # Container filtering by project
    ui/
      styles.go              # Lipgloss styles (colors, borders, layouts)
      components.go          # Reusable components (StatusIcon, ProgressBar)
      layout.go              # Panel sizing calculations
```

---

### Visual Design Specification

**Layout Dimensions** (based on research - 3-panel golden ratio):

```
Minimum: 80x24 characters
Recommended: 120x40 characters

┌─ 20i Stack Manager ─────────────────────────────────────────────┐
│ Project: myproject                           Docker ✓  CPU: 12%  │ <- Header (1 line)
├──────────────┬──────────────────────────────────────────────────┤
│ Services     │ Service: apache                                  │
│              │ Status: 🟢 Running                                │
│ 🟢 apache    │ Image:  php:8.2-apache                           │
│   mariadb    │ Uptime: 2h 34m                                   │
│   nginx      │                                                  │
│   phpmyadmin │ CPU:    45% ▓▓▓▓▓░░░░░ (0.9 cores)              │
│              │ Memory: 128MB/512MB ▓▓░░░░░░░░ (25%)             │
│ (25% width)  │                                                  │
│              │ Ports:  80:8080, 443:8443                        │
│              │ ID:     a3f2d1b8c9e4                             │
│              │                                                  │
│              │ (75% width, dynamic height based on log panel)   │
├──────────────┴──────────────────────────────────────────────────┤
│ ?:help  Tab:panels  s:start/stop  r:restart  l:logs  q:quit     │ <- Footer (1 line)
└──────────────────────────────────────────────────────────────────┘

With logs open (l key pressed):
┌─ 20i Stack Manager ─────────────────────────────────────────────┐
│ Project: myproject                           Docker ✓  CPU: 12%  │
├──────────────┬──────────────────────────────────────────────────┤
│ Services     │ Service: apache                                  │
│              │ Status: 🟢 Running  Uptime: 2h 34m                │
│ 🟢 apache    │ CPU: 45% ▓▓▓▓▓░░░░░  Memory: 128MB ▓▓░░░░░░░░  │
│   mariadb    ├──────────────────────────────────────────────────┤
│   nginx      │ Logs (following) - /search to filter            │
│   phpmyadmin │ [2025-12-28 10:23:45] GET /index.php 200         │
│              │ [2025-12-28 10:23:46] GET /styles.css 200        │
│ (25%)        │ [2025-12-28 10:23:47] POST /api/data 201         │
│              │ ...                                              │
│              │ (70% of right panel = ~50% total height)         │
├──────────────┴──────────────────────────────────────────────────┤
│ f:follow  /:search  Esc:close logs  ↑↓:scroll  q:quit           │
└──────────────────────────────────────────────────────────────────┘
```

**Color Palette** (Lipgloss named colors):

```go
var (
    // Status colors
    ColorRunning   = lipgloss.Color("10")  // Bright Green
    ColorStopped   = lipgloss.Color("8")   // Gray
    ColorRestart   = lipgloss.Color("11")  // Yellow
    ColorError     = lipgloss.Color("9")   // Bright Red
    
    // UI colors
    ColorAccent    = lipgloss.Color("12")  // Bright Blue (selected items)
    ColorBorder    = lipgloss.Color("8")   // Gray (panel borders)
    ColorText      = lipgloss.Color("7")   // White/Default
    ColorDim       = lipgloss.Color("8")   // Gray (secondary text)
    ColorHighlight = lipgloss.Color("13")  // Magenta (search matches)
    
    // Semantic colors
    ColorSuccess   = lipgloss.Color("10")  // Green
    ColorWarning   = lipgloss.Color("11")  // Yellow
    ColorDanger    = lipgloss.Color("9")   // Red
    ColorInfo      = lipgloss.Color("12")  // Blue
)
```

**Typography & Spacing**:
- **Panel borders**: Use lipgloss.Border (thin lines `│─┌┐└┘`)
- **Service list item**: `{icon} {name}` (e.g., `🟢 apache`)
- **Status icon + text**: Icon first, one space, then text
- **Progress bars**: Use Unicode blocks `▓` (filled) and `░` (empty), 10 chars wide
- **Padding**: 1 space inside panels, 0 between panels (border provides separation)
- **Line height**: Single-spaced (no blank lines within panels)

**Status Icons**:
- 🟢 Running (green circle)
- ⚪ Stopped (white/gray circle)  
- 🟡 Restarting (yellow circle)
- 🔴 Error/Unhealthy (red circle)
- ✅ Success (checkmark for actions)
- ❌ Error (X for failures)
- ⚠️  Warning (triangle for confirmations)

---

### Anti-Patterns to Avoid

**Based on research findings** - these cause "messy TUIs that don't work":

1. **❌ God Object Model**
   - DON'T: Put all state in one massive struct with 50+ fields
   - DO: Each view (Dashboard, Help, Projects) is a separate Bubble Tea model
   
2. **❌ Blocking I/O in Update**
   - DON'T: Call `docker.ListContainers()` directly in Update() method
   - DO: Use background goroutines + channels, send results as messages
   
3. **❌ Tab Navigation for Context Views**
   - DON'T: Use tabs to switch between Services/Logs/Stats (requires clicking through)
   - DO: Use panels so user sees service list + detail simultaneously
   
4. **❌ Inconsistent Key Bindings**
   - DON'T: `d` = delete in one view, `d` = download in another
   - DO: Establish global conventions (`s`=start/stop everywhere, `d`=delete everywhere)
   
5. **❌ No Visual Hierarchy**
   - DON'T: Everything same color/weight, no spacing, walls of text
   - DO: Use color for meaning (green=running), bold for emphasis, spacing for grouping
   
6. **❌ Hidden Features**
   - DON'T: Shortcuts only in help modal, no hints in UI
   - DO: Footer always shows top shortcuts for current view
   
7. **❌ Vague Errors**
   - DON'T: "Error: 1" or "operation failed"
   - DO: "❌ Failed to start apache: port 80 already in use"
   
8. **❌ Too Many Features in v1**
   - DON'T: Try to build config editor + monitoring + image management all at once
   - DO: MVP = dashboard + lifecycle + logs (match 20i-gui baseline)

---

### Non-Functional Requirements

- **NFR-001**: TUI MUST start and display dashboard in <2s on modern hardware (M1+, i5+)
- **NFR-002**: Panel focus switching MUST complete in <50ms (imperceptible lag)
- **NFR-003**: Container actions MUST show feedback within <300ms ("Starting..." message)
- **NFR-004**: Memory usage MUST stay <30MB with 4 services + 10k log lines per container
- **NFR-005**: TUI MUST handle terminal sizes from 80x24 (minimum) to 300x100 (practical maximum)
- **NFR-006**: Dashboard stats MUST refresh every 2s (configurable 1-5s range)
- **NFR-007**: Log buffer MUST NOT exceed 40MB total (4 containers × 10k lines × 1KB avg)
- **NFR-008**: Background goroutines MUST NOT block main UI thread (all I/O async)

### Keyboard Shortcuts

**Global (work from any view)**:
- `?` - Show help modal with all shortcuts
- `q` or `Ctrl-C` - Quit TUI (with confirmation if operations in progress)
- `p` - Open project switcher modal
- `Tab` - Cycle focus between panels (list → detail → logs)
- `Shift-Tab` - Cycle focus backwards

**Service List Panel** (when focused):
- `↑` or `k` - Move selection up
- `↓` or `j` - Move selection down
- `Enter` - Show detail for selected service (focus detail panel)
- `s` - Start/stop selected service (toggle)
- `r` - Restart selected service
- `l` - Open logs for selected service (focus log panel)
- `S` (shift-s) - Stop all stack containers (with confirmation)
- `R` (shift-r) - Restart entire stack (docker compose restart)
- `D` (shift-d) - Destroy stack (down -v, requires "yes" confirmation)

**Log Panel** (when open and focused):
- `f` - Toggle follow mode (auto-scroll on new lines)
- `/` - Search/filter logs (type to filter, Esc to clear)
- `↑` or `k` - Scroll up (when not following)
- `↓` or `j` - Scroll down
- `g` - Jump to top of logs
- `G` (shift-g) - Jump to bottom of logs
- `Esc` or `q` - Close log panel, return to detail view

**Project Switcher Modal** (when open):
- `↑` or `k` - Move up in project list
- `↓` or `j` - Move down
- `Enter` - Switch to selected project
- `Esc` - Close modal without switching
- `/` - Filter projects by name

**Help Modal** (when open):
- `?` or `Esc` - Close help
- `↑/↓` or `j/k` - Scroll help content

**Design Rationale** (from research):
- **Vim bindings** (`j/k`) + arrow keys for accessibility
- **Single-key actions** (no Ctrl/Alt combos for common tasks)
- **Shift for "all"** operations (S=stop all, R=restart all, D=destroy all)
- **Consistent verbs**: `s`=start/stop everywhere, `r`=restart, `l`=logs
- **`/` for search** (universal pattern from vim, lazygit, k9s)
- **`Esc` always cancels/closes** (modal dialogs, log panel, search)

---

### Key Entities

- **Service**: A logical component of the 20i stack (apache, mariadb, nginx, phpmyadmin); maps 1:1 to a Docker Compose service definition
- **Container**: Runtime Docker container instance; attributes: ID, name, image, status ("running"|"stopped"|"restarting"|"error"), state, ports, stats (CPU%, memory%, network I/O)
- **Project**: Directory containing docker-compose.yml; identified by directory name; may have .20i-local file for project-specific config
- **Log Stream**: Container stdout/stderr output; tail mode (last N lines) or follow mode (real-time); max 10,000 lines buffered
- **Stats**: Resource metrics for a container: CPU% (0-400% on 4-core), memory (bytes used/limit), network RX/TX bytes; refreshed every 2s

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Performance** (quantitative):
- **SC-001**: TUI starts and shows dashboard in <2s from launch on modern hardware
- **SC-002**: Container start/stop operations complete in <5s (Docker API time, not TUI latency)
- **SC-003**: Log follow mode displays new entries within <500ms of generation
- **SC-004**: Stats refresh cycle completes in <200ms (does not block UI)
- **SC-005**: Panel switching (Tab key) feels instant (<50ms perceived lag)

**Usability** (qualitative but measurable via user testing):
- **SC-006**: New users can navigate to logs view in <10 seconds without reading docs (follow visual hints in footer)
- **SC-007**: Common workflow (check status → restart service → view logs) achievable in ≤3 keystrokes per step
- **SC-008**: Help system (`?`) shows all available shortcuts grouped by context (no hidden features)
- **SC-009**: Error messages are actionable ("port 80 in use" not "bind error")

**Reliability** (functional correctness):
- **SC-010**: 100% parity with 20i-gui features (start, stop, restart, status, logs, destroy)
- **SC-011**: Zero data loss from failed operations (atomic Docker API calls, clear error states)
- **SC-012**: TUI handles Docker daemon restart gracefully (shows error, auto-retries, reconnects)
- **SC-013**: Terminal resize never crashes or corrupts display (re-layout on SIGWINCH)

**Compatibility**:
- **SC-014**: Works with existing .20i-local files (no migration needed)
- **SC-015**: Runs on macOS (primary), Linux (secondary) - matches 20i-gui platform support
- **SC-016**: Requires only Go 1.21+ and Docker (no additional dependencies beyond 20i-gui)

**Comparison to 20i-gui** (upgrade value):
- **SC-017**: Faster workflow: TUI dashboard shows all services at once (vs 20i-gui multi-step menus)
- **SC-018**: Live updates: Stats refresh every 2s (vs 20i-gui static status snapshots)
- **SC-019**: Richer logs: Follow mode + search + scroll (vs 20i-gui `docker compose logs -f` passthrough)
- **SC-020**: Keyboard-first: All actions <3 keystrokes (vs 20i-gui requires mouse for dialog selection)

---

## Assumptions

- Go 1.21+ is available (or install instructions provided in tui/README.md)
- Docker daemon is running and accessible via default socket (unix:///var/run/docker.sock or similar)
- User has permissions to manage Docker containers (same as 20i-gui requirements)
- Terminal supports ANSI colors (8-color minimum, 256-color recommended)
- Terminal supports alternate screen mode (standard since 1980s - xterm, iTerm2, Terminal.app, etc.)
- Minimum terminal size: 80x24 characters (enforced with error message if smaller)
- Docker Compose v2+ installed (docker compose, not docker-compose)
- Projects follow standard structure: docker-compose.yml at root, services named apache/mariadb/nginx/phpmyadmin
- User runs TUI from project directory (cd into project, then run `20i-tui`)

---

## Dependencies

**Go Modules** (specific versions for reproducible builds):
- **Bubble Tea** v1.3.10+ - TUI framework ([github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea))
- **Bubbles** v1.0.0+ - TUI components: list, viewport, textinput ([github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles))
- **Lipgloss** v1.0.0+ - Terminal styling ([github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss))
- **Docker SDK for Go** v27.0.0+ - Docker API client ([github.com/docker/docker/client](https://github.com/docker/docker))

**Optional** (not in MVP):
- Cobra - CLI flags/subcommands (if we add `20i-tui --project /path` later)
- YAML v3 - Config parsing (Phase 2 when we add config editor)

**Build Requirements**:
- Go 1.21 or later
- No C dependencies (pure Go, cross-platform)

**Runtime Requirements**:
- Docker daemon accessible (same as 20i-gui)
- Docker Compose v2+

---

## Files Affected

**New Files** (all under `tui/` directory):
```
tui/
  main.go                           # Entry point
  go.mod                            # Go module definition
  go.sum                            # Dependency checksums
  README.md                         # Build and usage instructions
  Makefile                          # Build targets (build, install, clean)
  internal/
    app/
      root.go                       # RootModel (top-level app state)
      messages.go                   # Custom tea.Msg types
    views/
      dashboard/
        dashboard.go                # DashboardModel
        service_list.go             # Service list panel (Bubbles list)
        detail.go                   # Detail panel
        logs.go                     # Log panel (Bubbles viewport)
      help/
        help.go                     # Help modal
      projects/
        projects.go                 # Project switcher modal
    docker/
      client.go                     # Docker SDK wrapper
      stats.go                      # Background stats collector
      filters.go                    # Project/container filtering
    ui/
      styles.go                     # Lipgloss styles (colors, borders)
      components.go                 # Reusable components (StatusIcon, ProgressBar)
      layout.go                     # Panel sizing functions
```

**Modified Files**:
- `README.md` - Add "TUI Interface" section with install/usage
- `CHANGELOG.md` - Document new 20i-tui feature in v2.0.0
- `.gitignore` - Add `tui/20i-tui` (compiled binary)

**Unchanged Files** (important for compatibility):
- `20i-gui` - Existing GUI still works, can coexist
- `docker-compose.yml` - No changes to stack definition
- `config/stack-vars.yml` - No changes (Phase 2 will extend)
- `.20i-local` - No changes to format

**Installation**:
- Binary installed to `/usr/local/bin/20i-tui` (or user's $GOBIN)
- Symlink `tui` → `20i-tui` for shorter command
- No config files needed for MVP (all defaults hardcoded)
