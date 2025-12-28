# 20i Stack Manager TUI - User Guide

## Overview

The 20i Stack Manager TUI is a modern, keyboard-driven Terminal User Interface for managing Docker containers and stacks for the 20i development environment. Built with Go and Bubble Tea, it provides an efficient alternative to Docker Desktop for managing your local development stack.

## Getting Started

### Installation

#### From Source
```bash
cd /path/to/20i-stack/tui/
make build          # Builds binary to bin/20i-stack-manager
make install        # Installs to $GOPATH/bin
```

#### Verify Installation
```bash
20i-stack-manager --version   # (if implemented)
which 20i-stack-manager        # Shows installation path
```

### First Run

1. **Navigate to your web project directory**:
   ```bash
   cd ~/my-website/
   ```

2. **Launch the TUI**:
   ```bash
   20i-stack-manager
   ```

3. **What you'll see**:
   - **Main Panel**: List of containers/services for your project
   - **Status Panel**: Current selection details and helpful hints
   - **Footer**: Available keyboard commands

## Phase 3a Features

### Dashboard View

The dashboard is the main interface showing your Docker stack status.

**Left Panel - Service List (30%)**:
- Lists all containers for the current project
- Shows container status (running, stopped, etc.)
- Highlights selected container
- Navigate with `↑/↓` or `j/k`

**Right Panel - Status & Details (70%)**:
- Shows selected container information
- Displays status messages and errors
- Provides contextual help

**Footer**:
- Shows available keyboard shortcuts
- Context-aware command hints

### Container Operations

#### Navigate Containers
- `↑` or `k` - Move selection up
- `↓` or `j` - Move selection down

#### Control Individual Containers
- `s` - Start/Stop selected container (toggle)
- `r` - Restart selected container

### Global Commands

- `?` - Show help modal with all commands
- `q` or `Ctrl+C` - Quit the application
- `p` - Switch to projects view (Phase 4+)
- `Esc` - Return to dashboard from other views

### Help Modal

Press `?` at any time to see the help modal with all available commands:

```
Help - Keyboard Shortcuts

Global:
  q, ctrl+c    Quit
  ?            Show this help
  p            Projects
  esc          Back to dashboard

Dashboard:
  s            Start/stop container
  r            Restart container
  ↑/↓, k/j     Navigate containers
```

Press `Esc` to close the help modal.

## Workflows

### Viewing Container Status

1. Launch TUI in your project directory
2. Use `↑/↓` to browse containers
3. View status information in the right panel

### Starting/Stopping Containers

1. Navigate to the container using `↑/↓`
2. Press `s` to toggle start/stop
3. Wait for the operation to complete
4. Container list refreshes automatically

### Restarting a Container

1. Navigate to the container
2. Press `r` to restart
3. Status message confirms the operation

## Mouse Support

Phase 3a includes mouse support for enhanced usability:

- **Click** on containers to select them
- **Scroll** to navigate long lists
- Mouse gestures complement keyboard controls

Note: Keyboard shortcuts are the primary interface for efficiency.

## Tips and Best Practices

### Keyboard Efficiency
- Learn the basic shortcuts: `j/k` for navigation, `s` for start/stop, `r` for restart
- Use `?` whenever you forget a command
- `Esc` always returns you to the dashboard

### Terminal Setup
- Use a terminal with 256-color support for best visuals
- Minimum recommended size: 80x24
- Recommended terminals: iTerm2 (macOS), Alacritty, or Windows Terminal

### Workflow Integration
- Run the TUI in a dedicated terminal tab/pane
- Keep it running while you work on your project
- Quick access to container status without leaving the terminal

## Limitations (Phase 3a)

Current limitations that will be addressed in future phases:

- **Single Project Only**: TUI manages containers in current directory only
- **No Stack Operations**: Cannot start/stop entire stack yet (Phase 3b)
- **Limited Metrics**: CPU/memory stats not yet displayed (Phase 3b)
- **No Multi-Project**: Cannot browse or switch between projects (Phase 4)

## Next Steps

After getting familiar with Phase 3a:

1. Read [Troubleshooting Guide](troubleshooting.md) for common issues
2. Review [Architecture Documentation](architecture.md) to understand design
3. Explore keyboard shortcuts in the help modal (`?`)
4. Watch for Phase 3b updates with stack-level operations

## Feedback

Encounter issues or have suggestions? See [CONTRIBUTING.md](../../CONTRIBUTING.md) for how to report issues or contribute.
