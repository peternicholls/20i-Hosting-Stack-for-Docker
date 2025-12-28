# R01: Real-World TUI Excellence Study - Research Findings

**Research ID**: R01  
**Status**: Complete  
**Last Updated**: 2025-12-28  
**Researcher**: Agent  
**Time Spent**: 3.5 hours  

---

## Executive Summary

After analyzing 4 excellent TUI applications (lazydocker, lazygit, k9s, gh-dash) and Bubble Tea framework patterns, clear design principles emerge that separate great TUIs from messy, unusable ones. The research reveals that successful TUIs share common architectural patterns: **simple, consistent navigation** (panels + keyboard shortcuts), **progressive information disclosure** (master-detail views), **real-time updates without disruption** (separate goroutines with throttling), and **discoverable help systems** (context-aware hints in footer).

**Key Finding**: The biggest cause of "messy TUIs that don't work" is **violating the single-responsibility principle per view** and **poor state management**. Great TUIs use Bubble Tea's Elm Architecture properly: each view is a self-contained component with its own Model/Update/View, composed together by a root coordinator. Failed TUIs try to cram everything into one massive update function with tangled state.

**Critical Recommendations for 20i Stack Manager**:
1. **Use panel-based layout** (not tabs) - lazydocker/k9s pattern with 3 fixed panels: navigation list (left), main content (center), detail/help (bottom)
2. **Keyboard-first navigation** with vim-style bindings + arrow keys as fallback
3. **Real-time updates in background** with 1-2s refresh, never block user input
4. **Context-aware help in footer** showing current view's shortcuts
5. **Start simple, expand gradually** - Dashboard only in v1, add features incrementally

---

## Application Analysis

### Additional Docker TUIs Considered

From [awesome-tuis](https://github.com/rothgar/awesome-tuis#dockerlxck8s), several other Docker/container TUIs exist:

**Evaluated but not deeply analyzed** (reasons noted):
- **ctop** (~15k stars) - Container metrics dashboard, similar scope to lazydocker but less feature-rich
- **dive** (~46k stars) - Image layer explorer (different use case - image analysis, not runtime management)
- **dockly** (~3.8k stars) - Node.js-based, less active maintenance
- **dry** (~3k stars) - Older project, gocui-based like lazydocker
- **oxker** (~700 stars) - Rust TUI, simpler feature set
- **sen** (~900 stars) - Python urwid-based, archived project
- **dtop** - Multi-host monitoring (different scope)
- **dprs** - Similar to lazydocker but newer/less mature
- **ducker** - "k9s-inspired" but for Docker (worth noting the k9s influence)
- **Pocker**, **Podman-tui** - Podman-specific (not Docker)

**Why lazydocker remains primary reference**:
1. Most active development (48k+ stars, regular updates)
2. Best UX feedback from community ("the lazier way")
3. Same author as lazygit (consistent design patterns)
4. Most comprehensive feature set for Docker Compose
5. Widely recommended in Docker community

**Pattern confirmed**: The Docker TUI ecosystem converges on similar patterns (3-panel layouts, vim bindings, real-time stats, log streaming) - validating our research findings.

---

### lazydocker

**Overview**: A terminal UI for Docker/Docker Compose management with ~48k GitHub stars. Built with gocui (not Bubble Tea), focuses on "the lazier way" to manage Docker without memorizing commands. Created by Jesse Duffield (same author as lazygit), emphasizing consistent UX patterns.

**Navigation Patterns**:
- **3-panel layout**: Services list (left), main content (center), stats/info (bottom)
- **Panel switching**: Tab key cycles panels, numbers 1-8 for quick actions
- **View context switching**: Type `:` to enter command mode, then resource name (`:containers`, `:images`, `:volumes`)
- **Breadcrumbs**: Shows navigation path (Project > Containers > nginx)
- **No nested menus**: Everything is max 2 keystrokes away

**Information Hierarchy**:
- **At-a-glance dashboard**: Service name, status (color-coded), CPU %, Memory %, Image
- **Progressive disclosure**: Select service → see logs/stats/config in main panel
- **Status indicators**: Color + icon (🟢 running, 🔴 stopped, 🟡 restarting)
- **Grouped by compose project**: Shows only relevant containers for current directory

**Visual Design**:
- **Minimal borders**: Only where needed to separate panels
- **Color for meaning**: Green=running, Red=error, Yellow=warning, Gray=stopped
- **Monospace alignment**: Columns align perfectly (names, CPU%, Memory%)
- **Compact density**: Information-dense but not overwhelming

**Strengths**:
- **Instant comprehension**: See all services + status in <1 second
- **Fast keyboard navigation**: Muscle memory develops quickly (Tab, numbers)
- **Log following works well**: Real-time logs without UI flicker
- **Contextual actions**: Menu shows only relevant actions for selected item

**Weaknesses**:
- **Can't resize panels**: Fixed panel sizes sometimes waste space
- **Limited search/filter**: Hard to find specific container in long lists
- **Config editing is basic**: No syntax highlighting or validation

**Key Takeaways**:
- ✅ **Panel layout is superior to tabs** - see multiple contexts simultaneously
- ✅ **Color-coded status** is faster to scan than text labels
- ✅ **Command mode** (`:` prefix) is discoverable and powerful
- ⚠️ **Keep action menus short** (<10 items) or they become overwhelming

---

### lazygit

**Overview**: Terminal UI for git with ~70k stars. Widely considered the gold standard for TUI design - complex workflows made simple, extensive keyboard shortcuts, excellent help system. Also by Jesse Duffield, shares design DNA with lazydocker.

**Navigation Patterns**:
- **5-panel layout**: Status/files/branches (left stack), main diff/log (center), commit message (bottom when needed)
- **Panel focus**: Tab/Shift-Tab cycles, numbers 1-5 jump directly
- **Modal workflows**: Press `s` to stage → see immediate feedback → continue working
- **History navigation**: `[` and `]` go back/forward through command history
- **Filter mode**: `/` to filter current view, works across all contexts

**Information Hierarchy**:
- **Most important at top**: Current branch, ahead/behind commits, modified files count
- **Grouped by state**: Unstaged files, staged files, stashed changes (separate sections)
- **Commit graph**: Visual branch/merge representation with colors per author
- **Minimal text, maximum info**: Single-character status codes (M=modified, A=added, D=deleted)

**Visual Design**:
- **Extreme color discipline**: Each status (new/modified/deleted/conflict) has one consistent color
- **Box drawing characters**: Subtle panel borders using `│`, `─`, `┌`, `┐`
- **Highlighted current item**: Full-width selection bar, high contrast
- **Commit graph beauty**: Tree diagram uses Unicode box characters for merge lines

**Strengths**:
- **Discoverability**: Press `?` anywhere → see ALL shortcuts for current view
- **No context switching**: Stage lines, rebase, resolve conflicts - all in TUI
- **Undo/redo**: Press `z` to undo last action (using reflog) - HUGE confidence booster
- **Progressive learning curve**: Arrow keys + Enter work, vim bindings optional
- **Consistent patterns**: Same keys do similar things across views (d=delete, e=edit)

**Weaknesses**:
- **Initial complexity**: 5 panels overwhelming for newcomers (but help system mitigates)
- **Modal interactions** can confuse: Sometimes Enter=select, sometimes Enter=expand

**Key Takeaways**:
- ✅ **Help system is critical** - `?` to show shortcuts makes TUI self-documenting
- ✅ **Undo capability** gives users confidence to experiment
- ✅ **Consistent key bindings** across views (d=delete everywhere)
- ✅ **Status bar footer** shows current context and hints
- ⚠️ **Too many panels** can overwhelm - start with 3, add more only if needed

---

### k9s

**Overview**: Kubernetes TUI with ~32k stars. Real-time cluster monitoring, handles massive scale (1000s of pods), excellent performance. Uses its own TUI framework (not Bubble Tea).

**Navigation Patterns**:
- **Resource-first**: Type `:pod` to view pods, `:svc` for services, `:ns` for namespaces
- **Context switching**: Type `:ctx` to switch clusters, `:ns` to switch namespaces
- **Drill-down**: Select pod → press Enter → see containers → select container → see logs
- **Quick filter**: `/` to filter by name, `/-l` to filter by label
- **Breadcrumbs at top**: Shows cluster > namespace > resource type > item

**Information Hierarchy**:
- **Table-first**: Everything is a table (Name, Ready, Status, Restarts, Age)
- **Color-coded status**: Green=running, Red=error, Yellow=pending
- **Live metrics**: CPU/Memory bars update in real-time (1s intervals)
- **Hierarchical drilling**: Pod → Container → Logs (always one level deep)

**Visual Design**:
- **Header bar**: Bright color, shows cluster/namespace/view
- **Table with alternating rows**: Subtle gray/white for readability
- **Progress bars**: Horizontal bars for CPU/Memory (▓▓▓░░░)
- **Skins system**: Customizable color themes (30+ built-in themes)

**Strengths**:
- **Handles scale**: 1000+ pods without lag (virtualized scrolling)
- **Real-time updates**: Stats refresh every 1s, doesn't block input
- **Powerful filtering**: Regex, label selectors, fuzzy find all work
- **Log following**: Tail multiple pods simultaneously (split view)

**Weaknesses**:
- **Steep learning curve**: Kubernetes concepts required (not TUI's fault)
- **Too many keybindings**: Over 50 shortcuts (but `?` help mitigates)
- **Resource-heavy**: High CPU when monitoring many pods

**Key Takeaways**:
- ✅ **Virtualized scrolling** essential for long lists (render only visible rows)
- ✅ **Background updates** with goroutines + channels pattern works perfectly
- ✅ **Filter-first design**: Let users hide noise, don't paginate
- ✅ **Header bar context**: Always show where you are (cluster/namespace/view)
- ⚠️ **Performance matters**: Real-time updates must not block UI

---

### gh-dash (Bubble Tea Reference)

**Overview**: GitHub dashboard TUI with ~10k stars. **Built with Bubble Tea** - serves as excellent reference for modern framework patterns. Clean, focused, does one thing well.

**Navigation Patterns**:
- **Tab-based sections**: PRs, Issues, Workflows (Tab/Shift-Tab to navigate)
- **List-detail pattern**: PR list (top), PR detail (bottom)
- **Keyboard shortcuts**: Single-key actions (c=checkout, o=open in browser, r=refresh)
- **Search mode**: `/` to search, Esc to clear

**Information Hierarchy**:
- **Scannable lists**: PR #, Title, Author, Status (color-coded)
- **Smart grouping**: "My PRs", "Needs Review", "All PRs" (user-configurable)
- **Detail on demand**: Select PR → see description, comments, checks below
- **Status indicators**: ✓=approved, ⏸=pending, ✗=failed

**Visual Design**:
- **Lipgloss styling**: Professional borders, colors, spacing
- **Glamour markdown**: Renders markdown beautifully in terminal
- **Minimalist chrome**: Thin borders, focus on content
- **Dark mode friendly**: Uses ANSI colors that adapt to terminal theme

**Strengths**:
- **Bubble Tea best practices**: Clean Model/Update/View separation per component
- **Component composition**: Each tab is a Bubble Tea model, composed by root model
- **Responsive**: Handles terminal resize gracefully (SIGWINCH)
- **Fast**: GraphQL queries cached, updates async

**Weaknesses**:
- **Simpler use case**: GitHub only (easier than Docker/Git multi-tool)
- **Limited real-time**: Polls every 30s, not continuous streaming

**Key Takeaways**:
- ✅ **Bubble Tea Component Pattern**: Root model contains tab models, delegates Update/View
- ✅ **Bubbles components work**: Use list, viewport, textinput from Bubbles library
- ✅ **Lipgloss for all styling**: Never use raw ANSI codes
- ✅ **One goroutine per async task**: Use `tea.Cmd` for HTTP requests, file I/O
- ✅ **Model is immutable**: Update returns new model, never mutate in place

---

## Pattern Analysis

### Navigation Patterns - What Works

**Command Mode is King** (`:resource-name`):
- **Why it works**: Type `:pod` is faster than Menu→Resources→Pods
- **Discoverability**: Help text can list available commands
- **Flexibility**: `:pod nginx` with args, `:pod /filter` with search
- **Used by**: lazygit, lazydocker, k9s (all top TUIs)

**Panel vs Tab Layouts**:
- **Panels win for context**: See list + detail + stats simultaneously (lazydocker pattern)
- **Tabs work for workflows**: Different tasks need full screen (gh-dash pattern)
- **Hybrid**: k9s uses panels within command-mode views
- **For 20i**: **Use panels** - Docker management needs context (containers + logs + stats)

**Keyboard Navigation Hierarchy**:
1. **Arrow keys + Enter** (universal, beginners)
2. **Vim bindings** (h/j/k/l, optional but beloved by power users)
3. **Single-key actions** (s=start, d=delete, l=logs)
4. **Tab key** (cycle panels/tabs)
5. **Numbers** (1-9 for quick jumps)
6. **Command mode** (`:` prefix for advanced)

**Filter/Search Patterns**:
- **Inline filter**: `/text` filters current view (k9s, lazygit)
- **Live search**: Updates as you type (don't require Enter)
- **Clear filter**: Esc or Ctrl-C returns to full view
- **Smart matching**: Fuzzy find (substring, acronym) better than exact match

### Layout Strategies

**3-Panel Golden Ratio**:
```
┌────────────┬─────────────────────────────┐
│ Navigation │ Main Content                │
│ List       │ (Logs/Details/Config)       │
│ 20% width  │ 80% width                   │
│            │                             │
│            ├─────────────────────────────┤
│            │ Footer/Help                 │
│            │ (Context + Shortcuts)       │
└────────────┴─────────────────────────────┘
```
- **Left panel**: 20-25% width, scrollable list
- **Main panel**: 75-80% width, detail/logs/stats
- **Footer**: 2-3 lines, context + key hints
- **Header**: Optional, cluster/project/version info

**Table Design Best Practices**:
- **Column headers**: Bold or different color, aligned with data
- **Alternating rows**: Subtle (every other row slightly gray)
- **Selected row**: High contrast background (not just color)
- **Status column**: Color + icon (🟢 + "Running")
- **Sortable**: Click header or key (s) to cycle sort
- **Resizable**: Auto-fit to terminal width, truncate with `...`

**Real-Time Update Patterns**:
```go
// GOOD: Background goroutine sends messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case statsMsg:  // Sent from background ticker
        m.stats = msg.stats
        return m, waitForStats()  // Schedule next update
    }
}

// BAD: Blocking call in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.stats = docker.GetStats()  // BLOCKS UI!
    return m, nil
}
```
- **Use `tea.Tick`** for periodic updates (1-2s intervals)
- **Use channels** to receive data from Docker API goroutines
- **Throttle updates**: Max 30 FPS rendering (33ms min between frames)
- **Cancel on exit**: Close channels when view switches

### Keyboard Interaction Patterns

**Global Shortcuts** (work everywhere):
- `?` - Help/shortcuts for current view
- `q` or `Ctrl-C` - Quit (with confirmation if dirty)
- `Tab`/`Shift-Tab` - Cycle panels/tabs
- `/` - Search/filter current view
- `Esc` - Cancel/back/clear filter
- `:` - Command mode (like vim)

**Context Shortcuts** (depend on view):
- List views: `j/k` or arrows, `Enter` to select, `d` to delete
- Log views: `f` to follow/unfollow, `g/G` for top/bottom, `/<search>`
- Config views: `e` to edit, `s` to save, `Esc` to cancel

**Action Confirmation**:
- **Destructive**: `d` to delete → show "Delete container 'nginx'? (y/N)" → `y` to confirm
- **Non-destructive**: `s` to start → immediate action, show "Starting..." feedback
- **Bulk actions**: `Space` to mark items, `Shift-D` to delete all marked

### Help Systems

**Context-Aware Footer** (always visible):
```
?:help  Tab:next  s:start  r:restart  l:logs  d:delete  q:quit
```
- **Current view shortcuts** only (don't overwhelm)
- **Most important first** (?,Tab,q always leftmost)
- **Visual separators**: Use `│` or colors to group related actions

**Help Modal** (`?` key):
- **Overlay**: Semi-transparent, doesn't lose context
- **Grouped**: Navigation | Container Actions | Logs | Global
- **Searchable**: Can type to filter shortcuts
- **Exit**: Same key `?` toggles, or `Esc`

**Progressive Disclosure**:
- **Hints on empty states**: "No containers. Press `r` to refresh or `?` for help"
- **Tooltips on hover** (if mouse enabled)
- **Status messages**: "Container nginx started successfully"

---

## Bubble Tea Framework Patterns

### Architecture - Elm's Model-Update-View

**The Golden Rule**: Each component is a standalone Bubble Tea program.

```go
// GOOD: Composable components
type RootModel struct {
    activeView  string
    dashboard   dashboardModel  // Standalone Bubble Tea model
    containers  containersModel // Standalone Bubble Tea model
    logs        logsModel       // Standalone Bubble Tea model
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "tab" {
            // Switch active view
            return m, nil
        }
    }
    
    // Delegate to active view
    switch m.activeView {
    case "dashboard":
        updated, cmd := m.dashboard.Update(msg)
        m.dashboard = updated.(dashboardModel)
        return m, cmd
    // ...
    }
}

// BAD: Monolithic god-object
type Model struct {
    containers []Container
    logs []string
    selectedContainer int
    selectedLog int
    statsUpdateTimer time.Time
    // ... 50 more fields
}
```

### State Management

**Immutable Updates**:
```go
// GOOD: Return new model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.cursor++  // Modify copy
    return m, nil  // Return modified copy
}

// BAD: Mutation without return
func (m *model) Update(msg tea.Msg) {
    m.cursor++  // Bubble Tea won't know state changed
}
```

**Shared State via Messages**:
```go
// Component A sends message
return m, func() tea.Msg {
    return containerStartedMsg{name: "nginx"}
}

// Component B receives it
case containerStartedMsg:
    m.statusText = msg.name + " started"
```

### Component Organization

**File Structure**:
```
tui/
  main.go                 # Entry point
  internal/
    app/
      app.go              # Root model
      router.go           # View switching logic
    views/
      dashboard/
        dashboard.go      # Model + Update + View
        update.go         # Update method (if large)
        commands.go       # tea.Cmd factories
      containers/
        containers.go
        list.go           # Bubbles list component
      logs/
        logs.go
        viewport.go       # Bubbles viewport component
    docker/
      client.go           # Docker API wrapper
      stats.go            # Real-time stats goroutine
    ui/
      styles.go           # Lipgloss styles
      components.go       # Reusable UI components
```

**Bubbles Components** (use these, don't reinvent):
- `list.Model` - Scrollable lists with selection
- `viewport.Model` - Scrollable text (logs, diffs)
- `textinput.Model` - Text input fields
- `table.Model` - Tables with sorting, filtering
- `progress.Model` - Progress bars, spinners

### Performance Considerations

**Render Optimization**:
```go
// GOOD: Only render when needed
func (m model) View() string {
    if m.hidden {
        return ""  // Don't render hidden views
    }
    if !m.dirty {
        return m.cachedView  // Cache unchanged views
    }
    m.cachedView = m.render()
    m.dirty = false
    return m.cachedView
}

// BAD: Always render everything
func (m model) View() string {
    return lipgloss.JoinVertical(
        m.header(),      // Expensive
        m.allContainers(),  // Expensive
        m.allLogs(),     // VERY expensive
    )
}
```

**Virtual Scrolling** for long lists:
```go
// Only render visible items
visibleStart := m.scroll
visibleEnd := min(m.scroll + m.height, len(m.items))
for i := visibleStart; i < visibleEnd; i++ {
    s += renderItem(m.items[i])
}
```

**Throttle Updates**:
```go
// Deduplicate rapid messages
func throttle(interval time.Duration) tea.Cmd {
    return tea.Tick(interval, func(t time.Time) tea.Msg {
        return tick MSG{}
    })
}
```

---

## Recommendations for 20i Stack Manager TUI

### Phase 1: MVP (v1.0.0) - Keep It Simple

**What to Build**:
1. **Dashboard View Only**:
   - 3-panel layout: Service list | Service detail | Footer help
   - Show 4 services: apache, mariadb, nginx, phpmyadmin
   - Basic actions: start, stop, restart (s, x, r keys)
   - Real-time status updates (2s interval)
   - Logs viewer (select service, press `l`)

**What to Skip** (add later):
- ❌ Tab navigation (use panels only)
- ❌ Config editing (do in v1.1)
- ❌ Project switching (single project only)
- ❌ Image management (not critical)
- ❌ Custom commands/plugins
- ❌ Mouse support (keyboard first)

**Why**: Get core loop working perfectly before adding features. Messy TUIs happen when too many half-baked features are crammed in.

### Architecture Recommendations

**Component Structure**:
```
RootModel
  ├─ DashboardModel (active by default)
  │    ├─ ServiceListModel (left panel, Bubbles list component)
  │    ├─ ServiceDetailModel (main panel, custom)
  │    └─ LogViewerModel (toggleable, Bubbles viewport)
  └─ HelpModel (modal, shown with `?`)
```

**State Management**:
- **Docker state**: Background goroutine polls every 2s, sends `statsMsg`
- **User input**: Handled in Update(), returns commands
- **View state**: Each component has its own model (cursor, scroll, selection)

**Key Bindings** (start with these):
```
Global:
  ?      - Help
  q      - Quit
  Ctrl-C - Quit (no confirmation)
  Tab    - Cycle panels (list → detail → logs → list)

Service List Panel:
  ↑/k    - Move up
  ↓/j    - Move down
  Enter  - Show detail
  l      - View logs
  s      - Start/stop service
  r      - Restart service
  d      - Remove container (with confirmation)

Log Viewer (when active):
  f      - Toggle follow mode
  /      - Search logs
  Esc    - Close logs, return to detail
```

### Visual Design

**Color Palette** (use Lipgloss constants):
```go
var (
    ColorRunning   = lipgloss.Color("10")  // Green
    ColorStopped   = lipgloss.Color("8")   // Gray
    ColorError     = lipgloss.Color("9")   // Red
    ColorWarning   = lipgloss.Color("11")  // Yellow
    ColorAccent    = lipgloss.Color("12")  // Blue
    ColorBorder    = lipgloss.Color("8")   // Gray
    ColorSelected  = lipgloss.Color("13")  // Magenta
)
```

**Layout** (80x24 minimum, 120x40 recommended):
```
┌─ 20i Stack Manager ────────────────────────────────────┐
│                                                         │
│ ┌─ Services ────┐ ┌─ Service: apache ───────────────┐ │
│ │ █ apache      │ │ Status: Running                  │ │
│ │   mariadb     │ │ Image:  php:8.2-apache           │ │
│ │   nginx       │ │ CPU:    45%  ▓▓▓▓▓░░░░░          │ │
│ │   phpmyadmin  │ │ Memory: 128MB / 512MB            │ │
│ └───────────────┘ │                                  │ │
│                   │ Ports: 80:8080                   │ │
│                   │ Up: 2h 34m                       │ │
│                   └──────────────────────────────────┘ │
│                                                         │
│ Tab:panels  s:start/stop  r:restart  l:logs  ?:help    │
└─────────────────────────────────────────────────────────┘
```

### Error Handling

**Graceful Failures**:
- **Docker daemon down**: Show error screen with "Docker not running. Start Docker and press `r` to retry."
- **Container action fails**: Show inline error for 3s: "❌ Failed to start apache: port 80 in use"
- **Terminal too small**: Show "Terminal too small (need 80x24, got 60x20). Resize and press any key."

**Recovery Patterns**:
- **Auto-retry Docker connection**: Try every 5s when daemon down
- **Stale data indicator**: Show "Last updated 30s ago" if Docker API slow
- **Partial data**: Show what's available, mark missing data as "Loading..."

---

## Anti-Patterns to Avoid

### What Makes TUIs Messy and Unusable

**1. God Object Syndrome**:
```go
// BAD: Everything in one massive struct
type AppModel struct {
    containers, images, volumes, networks, logs, configs, projects, stats...
    containerCursor, imageCursor, volumeCursor, logsCursor...
    containerView, imageView, configView, logsView...
    // 100 more fields
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 500 line switch statement trying to handle everything
}
```
**Why it fails**: Impossible to reason about state, bugs cascade, adding features breaks everything.

**2. Blocking I/O in Update**:
```go
// BAD: Freezes UI for 2 seconds
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.containers = docker.ListContainers()  // HTTP call blocks
    return m, nil
}
```
**Why it fails**: User presses key, UI freezes, user thinks app crashed, presses more keys, chaos.

**3. Too Many Navigation Levels**:
```
Main Menu → Docker → Containers → Apache → Actions → Start → Confirm → Success → Back → Back → Back → Back → Back
```
**Why it fails**: 7 steps to start a container. Users give up, use CLI instead.

**4. Inconsistent Key Bindings**:
- Dashboard: `d` = delete
- Logs view: `d` = download
- Config view: `d` = discard changes
**Why it fails**: Muscle memory causes disasters (delete instead of discard).

**5. Information Overload**:
```
Container: nginx-proxy-8a3f2d1 | Image: nginx:1.21.6-alpine-perl | Status: Running | Health: healthy | 
Ports: 80:8080/tcp, 443:8443/tcp | CPU: 0.03% | Mem: 12.5MB/256MB | Net I/O: 1.2kB/3.4kB | 
Block I/O: 0B/8kB | PIDs: 3 | Created: 2024-12-28T10:23:45Z | Started: 2024-12-28T10:23:46Z | 
Labels: com.docker.compose.project=20i, com.docker.compose.service=nginx...
```
**Why it fails**: Can't find important info (Status) in wall of text.

**6. No Visual Hierarchy**:
- Everything same color
- No spacing/grouping
- No bold/highlights
- Walls of text

**Why it fails**: Takes 10+ seconds to find status of one service.

**7. Hidden Functionality**:
- No help system
- Shortcuts not documented
- Features accessible only via obscure key combos
- No hints or tooltips

**Why it fails**: Users don't know what's possible, think app is limited.

**8. Poor Error Messages**:
```
Error: 1
Error: operation failed
Error: [object Object]
```
**Why it fails**: User has no idea what went wrong or how to fix it.

---

## Open Questions

1. **Docker Compose integration**: Should TUI parse `docker-compose.yml` for service names/ports, or just use Docker labels?
   - **Research needed**: Compare lazydocker's approach vs. reading compose file directly
   
2. **Config file format**: YAML vs. TOML vs. JSON for TUI settings (keybindings, colors, refresh interval)?
   - **Recommendation**: YAML (matches Docker Compose, easier for users to edit)

3. **Multi-project support**: How to detect/list projects? (presence of `docker-compose.yml`, `.20i-local`, or git repos?)
   - **Phase 2 feature**: Not critical for v1.0

4. **Logs storage**: Keep in memory (limited to 10k lines) or write to temp files?
   - **Research needed**: Test memory usage with 4 containers × 10k lines

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2025-12-28 | Research initiated | Agent |
| 2025-12-28 | Analyzed lazydocker, lazygit, k9s, gh-dash | Agent |
| 2025-12-28 | Completed findings and recommendations | Agent |

---

### k9s

**Overview**:

**Navigation Patterns**:

**Information Hierarchy**:

**Visual Design**:

**Strengths**:

**Weaknesses**:

**Key Takeaways**:

---

### gh-dash (Bubble Tea)

**Overview**:

**Navigation Patterns**:

**Information Hierarchy**:

**Visual Design**:

**Strengths**:

**Weaknesses**:

**Key Takeaways**:

---

## Pattern Analysis

### Navigation Patterns

### Layout Strategies

### Real-Time Updates

### Keyboard Interactions

### Help Systems

---

## Bubble Tea Framework Patterns

### Application Architecture

### State Management

### Component Organization

### Performance Considerations

---

## Recommendations

### Core Design Principles

### Navigation Strategy

### Visual Hierarchy

### Error Handling

### Help & Discoverability

---

## Anti-Patterns to Avoid

---

## Open Questions

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2025-12-28 | Research initiated | Agent |
