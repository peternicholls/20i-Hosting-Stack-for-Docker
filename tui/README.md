# 20i Stack Manager - Terminal UI

A modern, keyboard-driven Terminal User Interface (TUI) for managing the 20i Docker stack, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Quick Start

### The 20i-Stack Workflow

The TUI replicates the `20i-gui` workflow in a modern terminal interface:

**MVP (Phase 3a) - Single Project Mode**:
1. **Navigate to your web project**: `cd ~/my-website/`
2. **Launch the TUI**: `20i-stack-manager`
3. **Left Panel** shows:
   - Current project directory (`~/my-website/`)
   - Project name (`my-website`)
   - Stack status (Not Running / Running / Starting)
4. **Right Panel** shows (dynamic):
   - **Pre-flight**: Checks for `public_html/` folder
   - **Starting**: Live `docker compose up` output
   - **Running**: Stack status table (containers, ports, CPU%, memory)
5. **Bottom Panel** shows:
   - Available commands and keyboard shortcuts
   - Status messages and errors

**Future (Phase 4+) - Multi-Project Browser**:
- Left panel will list ALL available projects
- Navigate between projects with `↑/↓`
- Right panel shows selected project's stack status
- Manage multiple running stacks simultaneously

### Build and Run

```bash
# Build from source:
cd /path/to/20i-stack/tui/
make build          # Builds binary to bin/20i-stack-manager

# Install globally:
make install        # Installs to $GOPATH/bin as '20i-stack-manager'

# Run from your web project:
cd ~/my-website/
20i-stack-manager   # Manages stack for 'my-website'
```

### Development

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage  # Generates coverage.html

# Clean build artifacts
make clean

# Format and vet code
gofmt -w internal/ && go vet ./...
```

## Features

### Target MVP (Phase 3a - Single Project Stack Manager)

**Left Panel - Current Project Info**:
- 📂 **Directory Display** - Shows current working directory (`$PWD`)
- 📛 **Project Name** - Displays normalized project name
- 🚦 **Stack Status** - Shows if stack is running/stopped/starting

**Right Panel - Dynamic Context View**:
- ✅ **Pre-flight Check** - Validates `public_html/` folder exists
- 📦 **Template Installation** - Offers to create `public_html/` from demo template
- 📊 **Startup Output** - Shows live `docker compose up` output during stack start
- 📈 **Stack Status Table** - Displays running containers (like Docker Desktop):
  - Container name, status, image, ports, CPU%, memory
  - Live resource metrics and health indicators

**Bottom Panel - Commands/Instructions**:
- 💬 **Status Messages** - Operation feedback and error messages
- ⌨️ **Keyboard Shortcuts** - Available commands for current context
- 📝 **Instructions** - Contextual help and next steps

**Stack Operations**:
- 🚀 **Start Stack** - Runs `docker compose up -d` for current project
- 🛑 **Stop Stack** - Stops all containers for current project
- 🔁 **Restart Stack** - Restarts entire project stack
- 🗑️ **Destroy Stack** - Stops and removes volumes (with confirmation)

**Future (Phase 4+) - Multi-Project Browser**:
- 📁 Left panel lists ALL available web projects
- 🔀 Navigate between projects with keyboard
- 🎯 Right panel shows selected project's stack
- 🌐 Manage multiple running stacks simultaneously

### Current Implementation Status

✅ **UI Framework** - Bubble Tea model working with responsive layout  
✅ **Container Operations** - Start/stop/restart individual containers working  
✅ **Keyboard Navigation** - `↑/↓` or `j/k` navigation implemented  
⚠️ **Project Browser** - **NOT YET IMPLEMENTED** (needs directory listing)  
⚠️ **Pre-flight Checks** - **NOT YET IMPLEMENTED** (public_html/ validation)  
⚠️ **Template Installation** - **NOT YET IMPLEMENTED**  
⚠️ **Stack Lifecycle** - **NOT YET IMPLEMENTED** (compose up/down for projects)  
⚠️ **Three-Panel Layout** - **NOT YET IMPLEMENTED** (currently two-panel)  
⚠️ **Resource Metrics** - **NOT YET IMPLEMENTED** (CPU/memory stats)  

### Keyboard Controls (Planned)

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up (service list) |
| `↓` / `j` | Move selection down (service list) |
| `S` | **Start stack** for current project |
| `T` | **Stop stack** for current project |
| `R` | **Restart stack** for current project |
| `s` | Start individual selected service |
| `t` | Stop individual selected service |
| `r` | Restart individual selected service |
| `q` / `ctrl+c` | Quit application |

**Note**: Current implementation only supports individual container actions. Stack-level operations are planned.

### Planned (Future Phases)

⏳ **Stack Operations** - Stop/restart entire stack (`S`/`R` on no selection)  
⏳ **Detail Panel** - Container logs, environment, ports, volumes  
⏳ **Tab Navigation** - Switch between detail views (`tab`)  
⏳ **Multi-Stack** - Manage multiple projects/compose files  
⏳ **Project Settings** - Edit .env variables from TUI  

## Architecture

### Technology Stack

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.3.10+ - The Elm Architecture for TUI
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal layouts and colors
- **Docker**: [Docker SDK for Go](https://github.com/docker/docker) - Container operations
- **Testing**: Go standard testing with mocks

### Design Patterns

- **Elm Architecture**: Model-Update-View pattern for state management
- **Generic Command Pattern**: Reusable async command executor (ADR-004)
- **Three-Panel Layout**: Project info (left) + Service list (right) + Commands (bottom)
- **Centralized Errors**: Consistent Docker error formatting (ADR-006)
- **Project-Aware Operations**: Context-driven stack management based on `$PWD`

### Implementation Roadmap

**Phase 3a - Single Project Stack Manager** (MVP - Next Priority):
- [ ] Detect current working directory (`$PWD`) as project root
- [ ] Display project info in left panel (name, directory, status)
- [ ] Three-panel layout implementation
- [ ] `public_html/` folder detection
- [ ] Template installation from `demo-site-folder/`
- [ ] `docker compose up -d` for current project
- [ ] Live compose output streaming to right panel
- [ ] Stack status table (containers, ports, basic metrics)
- [ ] Environment variable handling (`CODE_DIR`, `COMPOSE_PROJECT_NAME`)

**Phase 3b - Enhanced Stack Management**:
- [ ] Resource metrics (CPU%, memory, network)
- [ ] Stack destruction with volume removal
- [ ] Port conflict detection and resolution
- [ ] phpMyAdmin architecture selection (ARM vs x86)
- [ ] `.20i-local` configuration editor

**Phase 4 - Multi-Project Browser**:
- [ ] Directory listing in left panel (multiple projects)
- [ ] Project navigation (`↑/↓` between projects)
- [ ] Show all running stacks across projects
- [ ] Switch active project context
- [ ] Manage multiple stacks simultaneously

**What Already Works** ✅:
- Bubble Tea UI framework and responsive layout rendering
- Container listing from Docker daemon
- Individual container start/stop/restart actions (foundation for stack ops)
- Keyboard navigation and selection
- Async command execution with status updates

### Project Structure

```
tui/
├── main.go                    # Entry point
├── Makefile                   # Build automation
├── go.mod                     # Dependencies
├── bin/                       # Build output
│   └── 20i-stack-manager
├── internal/
│   ├── app/                   # Application core
│   │   ├── messages.go        # Bubble Tea messages
│   │   └── root.go           # Root model/controller
│   ├── docker/                # Docker client wrapper
│   │   ├── client.go
│   │   └── types.go
│   ├── ui/                    # Shared UI components
│   │   ├── colors.go
│   │   ├── styles.go
│   │   └── status.go
│   └── views/                 # View components
│       └── dashboard/
│           ├── dashboard.go   # Dashboard model
│           ├── service_list.go
│           └── dashboard_test.go
└── tests/
    ├── integration/           # Integration tests
    ├── mocks/                 # Test mocks
    └── unit/                  # Unit tests
```

## Development Guidelines

### ADR Compliance

All code follows the [Architectural Decision Records](../specs/001-stack-manager-tui/) in the spec directory:

- **ADR-002**: Two-panel layout (service list left, details right)
- **ADR-004**: Generic command function for Docker operations
- **ADR-006**: Centralized error message formatting
- **ADR-007**: Startup resilience (graceful Docker daemon handling)

### Testing

- **Unit Tests**: Test models, update logic, and view rendering
- **Integration Tests**: Test with mock Docker client
- **Coverage Target**: >80% for core business logic

```bash
# Run specific test files
go test ./internal/views/dashboard/... -v

# Run with race detection
go test -race ./...

# Generate coverage report
make test-coverage
open coverage.html
```

### Code Style

- Run `gofmt` before committing
- Use `go vet` to catch common issues
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Keep functions focused and testable

## Troubleshooting

### Build Issues

```bash
# Verify Go version (requires 1.21+)
go version

# Download dependencies
go mod download

# Clean and rebuild
make clean
make build
```

### Runtime Issues

**"Cannot connect to Docker daemon"**
- Ensure Docker Desktop is running
- Check Docker socket permissions: `ls -la /var/run/docker.sock`

**"No containers found"**
- Start the 20i stack first: `cd .. && docker compose up -d`
- Verify containers exist: `docker ps -a`

**UI rendering issues**
- Ensure terminal supports 256 colors
- Try a different terminal emulator (iTerm2, Alacritty recommended)
- Check terminal size: `tput cols` and `tput lines` (minimum 80x24)

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

### Current Development Phase

**Phase 3 - Container Lifecycle MVP** (In Progress)
- Focus: Individual container start/stop/restart
- Status: Core features implemented, additional tests pending
- Next: Stack-wide operations (stop all, restart all)

See [tasks.md](../specs/001-stack-manager-tui/tasks.md) for detailed task breakdown.

## License

MIT License - See [LICENSE](../LICENSE) for details.

## References

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss/tree/master/examples)
- [Docker SDK Documentation](https://pkg.go.dev/github.com/docker/docker)
- [Phase 3 Roadmap](../specs/001-stack-manager-tui/PHASE3-ROADMAP.md)
