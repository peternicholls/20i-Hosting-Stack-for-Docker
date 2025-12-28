# Feature Specification: Self-Contained CLI Binary

**Feature Branch**: `005-cli-binary`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🔴 Critical  
**Input**: User description: "Single-file executable for global or per-project installation with commands for init, start, stop, status, logs"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Initialize Stack in Project (Priority: P1)

As a developer, I want to run a single command to initialize the 20i stack in my project directory so that I can start developing immediately without manual setup.

**Why this priority**: Initialization is the entry point to using the stack; without it, no other commands work.

**Independent Test**: Run `20i init` in an empty directory, verify `.20i-config.yml` is created and necessary compose files are downloaded.

**Acceptance Scenarios**:

1. **Given** an empty project directory, **When** the user runs `20i init`, **Then** `.20i-config.yml` is created with default settings
2. **Given** `20i init` is run, **When** initialization completes, **Then** all required Docker Compose files are present
3. **Given** a directory already initialized, **When** the user runs `20i init`, **Then** the system prompts for confirmation before overwriting

---

### User Story 2 - Start and Stop Stack (Priority: P1)

As a developer, I want to start and stop the stack with simple commands so that I can manage my development environment efficiently.

**Why this priority**: Start/stop are the core daily operations; essential for basic functionality.

**Independent Test**: Run `20i start` in initialized project, verify all containers are running; run `20i stop`, verify all containers are stopped.

**Acceptance Scenarios**:

1. **Given** an initialized project, **When** the user runs `20i start`, **Then** all stack services start within 60 seconds
2. **Given** a running stack, **When** the user runs `20i stop`, **Then** all containers stop gracefully
3. **Given** a custom port specified, **When** the user runs `20i start --port 8080`, **Then** the web server is accessible on port 8080

---

### User Story 3 - View Stack Status (Priority: P2)

As a developer, I want to see the status of all stack services so that I can diagnose issues quickly.

**Why this priority**: Status visibility is essential for troubleshooting but not required for basic operation.

**Independent Test**: Run `20i status` with stack running, verify output shows container states, ports, and uptime.

**Acceptance Scenarios**:

1. **Given** a running stack, **When** the user runs `20i status`, **Then** all service names, states, and ports are displayed
2. **Given** no stack running, **When** the user runs `20i status`, **Then** the system shows "No stack running in this directory"
3. **Given** multiple stacks running in different directories, **When** the user runs `20i status --all`, **Then** all running stacks are listed

---

### User Story 4 - View Service Logs (Priority: P2)

As a developer, I want to view logs for specific services so that I can debug application issues.

**Why this priority**: Log access is important for debugging but not required for basic stack operation.

**Independent Test**: Run `20i logs apache` with stack running, verify Apache access/error logs are streamed to terminal.

**Acceptance Scenarios**:

1. **Given** a running stack, **When** the user runs `20i logs apache`, **Then** Apache logs are streamed to the terminal
2. **Given** a running stack, **When** the user runs `20i logs` (no service specified), **Then** logs from all services are interleaved
3. **Given** the `--follow` flag, **When** the user runs `20i logs --follow`, **Then** new log entries are streamed in real-time

---

### User Story 5 - Destroy Stack and Volumes (Priority: P3)

As a developer, I want to completely remove a stack including volumes so that I can start fresh or clean up after a project.

**Why this priority**: Cleanup is important for maintenance but less frequently used than start/stop.

**Independent Test**: Run `20i destroy` in project with running stack, verify all containers are removed and volumes are deleted.

**Acceptance Scenarios**:

1. **Given** a running stack, **When** the user runs `20i destroy`, **Then** the system prompts for confirmation
2. **Given** confirmation provided, **When** destroy completes, **Then** all containers and volumes are removed
3. **Given** the `--force` flag, **When** the user runs `20i destroy --force`, **Then** no confirmation prompt is shown

---

### User Story 6 - Install CLI Globally (Priority: P3)

As a developer, I want to install the CLI globally using a single command so that I can use it across all my projects.

**Why this priority**: Easy installation improves adoption but users can also use per-project installation.

**Independent Test**: Run the curl installer command, verify `20i` is available in PATH and responds to `20i version`.

**Acceptance Scenarios**:

1. **Given** a macOS system, **When** the user runs the curl installer, **Then** `20i` binary is installed to `/usr/local/bin`
2. **Given** Homebrew is installed, **When** the user runs `brew install peternicholls/tap/20i-stack`, **Then** the CLI is installed and available
3. **Given** the CLI is installed, **When** the user runs `20i version`, **Then** CLI version and stack version are displayed

---

### Edge Cases

- What happens when Docker is not installed or not running?
- How does the CLI handle permission issues (e.g., port 80 requires sudo)?
- What happens when `.20i-config.yml` is corrupted or invalid?
- How does the CLI behave when run outside a project directory?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST be distributed as a single self-contained executable
- **FR-002**: CLI MUST support `init`, `start`, `stop`, `status`, `logs`, `destroy`, `update`, and `version` commands
- **FR-003**: CLI MUST work from any project directory with an initialized stack
- **FR-004**: CLI MUST embed all required stack configuration files
- **FR-005**: CLI MUST support macOS (Intel and ARM) and Linux (x64 and ARM64)
- **FR-006**: CLI MUST be installable via curl one-liner, Homebrew, and per-project
- **FR-007**: CLI MUST validate Docker availability on first run
- **FR-008**: CLI MUST provide clear error messages with actionable suggestions
- **FR-009**: CLI MUST support `--help` flag for all commands with usage examples

### Key Entities

- **CLI Binary**: Self-contained executable with embedded configuration files
- **Project Configuration**: `.20i-config.yml` file storing project-specific settings
- **Stack Instance**: Running set of containers managed by the CLI for a specific project

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: New users can install CLI and start first stack within 5 minutes
- **SC-002**: `20i init` completes in under 10 seconds on typical connection
- **SC-003**: `20i start` brings up full stack in under 60 seconds (first run may be longer due to image pulls)
- **SC-004**: CLI binary size is under 10MB for each platform
- **SC-005**: No Docker Compose knowledge required for basic usage
- **SC-006**: CLI works on 95%+ of developer machines (macOS 10.15+, Ubuntu 20.04+, Debian 10+)

## Assumptions

- Docker Desktop or Docker Engine is installed and running
- Users have basic terminal familiarity
- Internet connectivity available for initial image downloads
- File system permissions allow writing to project directory
