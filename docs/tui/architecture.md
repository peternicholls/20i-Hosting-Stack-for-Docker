# 20i Stack Manager TUI - Architecture

## Overview

The 20i Stack Manager TUI is built using the **Elm Architecture** pattern via the Bubble Tea framework, ensuring predictable state management and a clean separation of concerns.

## Technology Stack

### Core Dependencies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** v1.3.10+ - The Elm Architecture for TUI
  - Provides Model-Update-View pattern
  - Event-driven message passing
  - Built-in terminal I/O handling
  
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling and layouts
  - Flexbox-style layouts
  - Color and style management
  - Responsive design primitives

- **[Docker SDK for Go](https://github.com/docker/docker)** - Docker API client
  - Container lifecycle operations
  - Stats streaming
  - Compose integration

### Design Patterns

#### The Elm Architecture

All state management follows the Elm Architecture (TEA):

```
┌─────────────────────────────────────────┐
│  User Input (keyboard, mouse, timer)   │
└──────────────┬──────────────────────────┘
               │
               ▼
         ┌──────────┐
         │   Msg    │  (tea.Msg - immutable message)
         └─────┬────┘
               │
               ▼
         ┌──────────┐
         │  Update  │  (pure function: Model, Msg -> Model, Cmd)
         └─────┬────┘
               │
               ▼
         ┌──────────┐
         │  Model   │  (application state)
         └─────┬────┘
               │
               ▼
         ┌──────────┐
         │   View   │  (pure function: Model -> string)
         └─────┬────┘
               │
               ▼
         ┌──────────┐
         │ Terminal │  (rendered output)
         └──────────┘
```

**Key Principles**:
1. **Immutability** - Models are never mutated, always copied
2. **Pure Functions** - Update and View are deterministic
3. **Message Passing** - All events are messages
4. **Command Pattern** - Side effects via `tea.Cmd`

## Project Structure

```
tui/
├── main.go                      # Entry point, program initialization
├── internal/
│   ├── app/                     # Application core
│   │   ├── root.go             # RootModel - top-level coordinator
│   │   ├── root_test.go        # RootModel tests
│   │   ├── messages.go         # Custom message types
│   │   └── messages_test.go    # Message tests
│   ├── docker/                  # Docker client abstraction
│   │   ├── client.go           # Docker API wrapper
│   │   ├── client_test.go      # Client tests
│   │   ├── compose.go          # Docker Compose operations
│   │   └── compose_test.go     # Compose tests
│   ├── project/                 # Project detection and management
│   │   ├── detector.go         # Project/template detection
│   │   ├── detector_test.go    # Detector tests
│   │   ├── template.go         # Template installation
│   │   ├── template_test.go    # Template tests
│   │   ├── sanitize.go         # Project name sanitization
│   │   └── types.go            # Project data types
│   ├── ui/                      # Shared UI components and styles
│   │   ├── styles.go           # Lipgloss style definitions
│   │   ├── styles_test.go      # Style tests
│   │   ├── components.go       # Reusable UI components
│   │   ├── components_test.go  # Component tests
│   │   └── layout.go           # Layout calculations
│   └── views/                   # View-specific models
│       └── dashboard/
│           ├── dashboard.go       # Dashboard model and update logic
│           ├── dashboard_test.go  # Dashboard tests
│           └── service_list.go    # Service list rendering
└── tests/
    ├── mocks/                   # Mock implementations
    │   └── docker_mock.go       # Mock Docker client
    └── integration/             # Integration tests
        └── phase3a_test.go      # Phase 3a integration tests
```

## Core Components

### RootModel (app/root.go)

**Purpose**: Top-level application coordinator

**Responsibilities**:
- Initialize Docker client
- Route messages to appropriate view models
- Handle global shortcuts (`?`, `q`, `Esc`)
- Manage view switching (dashboard, help, projects)
- Coordinate window size updates

**State**:
```go
type RootModel struct {
    dockerClient *docker.Client  // Shared Docker client
    activeView   string           // Current view: "dashboard", "help", "projects"
    dashboard    dashboard.Model  // Dashboard view model
    width, height int             // Terminal dimensions
    lastError    error            // Global error state
}
```

**Key Methods**:
- `Init() tea.Cmd` - Initialize and load initial data
- `Update(tea.Msg) (tea.Model, tea.Cmd)` - Handle messages
- `View() string` - Render current view

### DashboardModel (views/dashboard/dashboard.go)

**Purpose**: Container management view

**Responsibilities**:
- Display container list
- Handle container selection
- Execute container operations (start/stop/restart)
- Refresh container list
- Display status messages

**State**:
```go
type Model struct {
    containers    []docker.Container  // Current container list
    selectedIndex int                 // Selected container index
    projectName   string               // Current project
    dockerClient  *docker.Client       // Docker API client
    width, height int                  // Panel dimensions
    lastError     error                // Operation errors
    lastStatusMsg string               // Status messages
}
```

### Message Types (app/messages.go)

Custom message types for event-driven communication:

- **ContainerActionMsg** - User requests action (start/stop/restart)
- **ContainerActionResultMsg** - Action completed (success/failure)
- **ComposeActionMsg** - Stack-level action request
- **ComposeActionResultMsg** - Stack action result
- **ErrorMsg** - Error occurred
- **SuccessMsg** - Operation succeeded
- **StackStatusMsg** - Stack operation status
- **StackOutputMsg** - Streaming compose output

## Data Flow

### Example: Starting a Container

1. **User Input**: User presses `s` on a container
2. **Key Message**: Bubble Tea creates `tea.KeyMsg{String: "s"}`
3. **Update**: Dashboard.Update() processes the key:
   ```go
   case "s":
       return m, containerActionCmd(client, containerID, "start", serviceName)
   ```
4. **Command Execution**: `containerActionCmd` runs asynchronously:
   ```go
   func containerActionCmd(...) tea.Cmd {
       return func() tea.Msg {
           err := client.StartContainer(containerID)
           return containerActionResultMsg{success: err == nil, ...}
       }
   }
   ```
5. **Result Message**: Command returns `containerActionResultMsg`
6. **Update Again**: Dashboard.Update() handles result:
   ```go
   case containerActionResultMsg:
       m.lastStatusMsg = msg.message
       if msg.success {
           return m, loadContainersCmd(...)  // Refresh list
       }
   ```
7. **View Update**: View() re-renders with new state

## Design Decisions

### ADR-002: Two-Panel Layout

**Decision**: Use left panel for service list (30%), right panel for details (70%)

**Rationale**:
- Service list is narrower, details need more space
- Mirrors common TUI patterns (file managers, email clients)
- Scales well to terminal sizes

**Implementation**: `lipgloss.JoinHorizontal()` for side-by-side panels

### ADR-004: Generic Command Pattern

**Decision**: Single `containerActionCmd()` for all container operations

**Rationale**:
- Reduces code duplication
- Consistent error handling
- Easier to test and maintain

**Implementation**:
```go
func containerActionCmd(client, id, action, service) tea.Cmd {
    return func() tea.Msg {
        var err error
        switch action {
        case "start":
            err = client.StartContainer(id)
        case "stop":
            err = client.StopContainer(id, timeout)
        case "restart":
            err = client.RestartContainer(id, timeout)
        }
        return containerActionResultMsg{...}
    }
}
```

### ADR-006: Centralized Error Formatting

**Decision**: Single `formatDockerError()` function for all Docker errors

**Rationale**:
- Consistent user experience
- User-friendly messages (not raw Docker errors)
- Centralized error pattern detection

**Implementation**: Pattern matching on error strings:
- "port is already allocated" → Port conflict message
- "timeout" → Timeout message with retry hint
- "No such container" → Not found message with refresh hint
- "permission denied" → Permission guidance

### ADR-007: Graceful Docker Unavailability

**Decision**: TUI starts even without Docker, shows helpful error

**Rationale**:
- Better user experience than crash
- Allows checking help, understanding requirements
- Supports CI/testing without Docker

**Implementation**: 
- `NewRootModel()` handles `docker.NewClient()` error
- Sets `dockerClient = nil`
- Views check for nil and show appropriate message

## Testing Strategy

### Unit Tests

**Location**: Next to implementation files (`*_test.go`)

**Coverage**:
- Model initialization
- Update logic for all message types
- View rendering (basic checks)
- Helper functions

**Example**:
```go
func TestDashboardUpdate_StartContainer(t *testing.T) {
    model := NewModel(mockClient, "test-project")
    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
    
    updated, cmd := model.Update(msg)
    
    if cmd == nil {
        t.Error("expected command to start container")
    }
}
```

### Integration Tests

**Location**: `tests/integration/`

**Coverage**:
- Full workflows (detect → start → status → stop)
- View switching
- Error handling

**CI-Safe Design**:
- Build tag: `//go:build !integration_docker`
- Tests skip when Docker unavailable
- Use mocks by default
- Live Docker tests gated behind `integration_docker` tag

**Running**:
```bash
# Default (no Docker required)
go test ./tests/integration/...

# With live Docker (future)
go test -tags=integration_docker ./tests/integration/...
```

### Mock Strategy

**Mock Docker Client** (`tests/mocks/docker_mock.go`):
- Implements same interface as real client
- Returns predictable test data
- No side effects
- Fast execution

## Phase Roadmap

### Phase 3a (Current) - MVP Container Management

✅ **Completed**:
- RootModel and DashboardModel wiring
- Container list and navigation
- Individual container operations (start/stop/restart)
- Help modal
- Mouse support
- Error handling and formatting
- CI-safe integration tests
- Documentation (user guide, troubleshooting, architecture)

### Phase 3b - Enhanced Stack Management

🚧 **Planned**:
- Stack-level operations (start/stop entire stack)
- Resource metrics (CPU%, memory)
- Stack destruction with volume removal
- Live compose output streaming
- `.20i-local` config editing

### Phase 4 - Multi-Project Browser

🔮 **Future**:
- Project listing and discovery
- Project switching
- Multiple stack management
- Project creation workflow

## Performance Considerations

### Async Operations

All I/O operations are asynchronous via `tea.Cmd`:
- Docker API calls don't block the UI
- User can continue interacting during operations
- Updates happen via message passing

### Refresh Strategy

**Current**: Refresh on demand after operations
**Future**: Periodic refresh with configurable interval

### Resource Usage

- Single Docker client shared across views
- Container list cached in model
- Stats streaming only for selected containers (Phase 3b+)

## Extension Points

### Adding a New View

1. Create model in `internal/views/<viewname>/`
2. Implement `Init()`, `Update()`, `View()`
3. Add to `RootModel.activeView` routing
4. Add keyboard shortcut in `RootModel.Update()`
5. Add view rendering in `RootModel.View()`

### Adding a New Container Operation

1. Add message type in `messages.go`
2. Add case in `dashboard.Update()`
3. Create command function (follow `containerActionCmd` pattern)
4. Add keyboard shortcut
5. Update footer and help modal

### Adding Docker Compose Operations

1. Implement in `internal/docker/compose.go`
2. Create message types in `messages.go`
3. Add command functions in dashboard
4. Wire up keyboard shortcuts
5. Handle streaming output (Phase 3b)

## References

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [The Elm Architecture](https://guide.elm-lang.org/architecture/)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss/tree/master/examples)
- [Docker SDK for Go](https://pkg.go.dev/github.com/docker/docker)
