# Feature Specification: Interactive Setup Wizard

**Feature Branch**: `010-setup-wizard`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Guided setup for first-time users with project type detection, preferences, and automatic configuration"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Guided First-Time Setup (Priority: P1)

As a first-time user, I want to be guided through setup with interactive prompts so that I can configure the stack without reading documentation.

**Why this priority**: First-time experience is critical for adoption; a smooth wizard reduces friction for new users.

**Independent Test**: Run `20i init` in an empty directory without any flags, verify interactive prompts guide the user through configuration.

**Acceptance Scenarios**:

1. **Given** a first-time user runs `20i init`, **When** the command starts, **Then** an interactive wizard begins with welcome message
2. **Given** the wizard is running, **When** user completes all prompts, **Then** configuration files are generated based on choices
3. **Given** the wizard completes, **When** user runs `20i start`, **Then** stack starts with the configured settings

---

### User Story 2 - Automatic Project Type Detection (Priority: P2)

As a developer with an existing project, I want the wizard to detect my project type so that it suggests appropriate configuration.

**Why this priority**: Auto-detection reduces manual choices and shows intelligence that builds trust.

**Independent Test**: Run `20i init` in a directory with `composer.json` containing Laravel, verify wizard suggests Laravel template.

**Acceptance Scenarios**:

1. **Given** a directory with Laravel `composer.json`, **When** wizard runs, **Then** it suggests Laravel template
2. **Given** a directory with WordPress files, **When** wizard runs, **Then** it suggests WordPress template
3. **Given** an empty directory, **When** wizard runs, **Then** it asks user to select project type manually

---

### User Story 3 - Port Conflict Detection (Priority: P3)

As a developer with multiple projects, I want the wizard to detect port conflicts so that I don't have to troubleshoot startup failures.

**Why this priority**: Port conflicts are a common frustration; proactive detection prevents issues.

**Independent Test**: Run `20i init` while another service uses port 80, verify wizard warns about conflict and suggests alternative.

**Acceptance Scenarios**:

1. **Given** port 80 is in use, **When** wizard runs, **Then** it warns about the conflict
2. **Given** port conflict detected, **When** user is prompted, **Then** alternative port is suggested
3. **Given** no port conflicts, **When** wizard runs, **Then** default ports are suggested without warnings

---

### User Story 4 - Optional Services Selection (Priority: P4)

As a developer, I want to select optional services during setup so that my stack includes everything I need from the start.

**Why this priority**: Service selection during init is convenient but not essential for basic setup.

**Independent Test**: Run wizard, select Redis and Mailhog, verify both services are enabled in generated configuration.

**Acceptance Scenarios**:

1. **Given** wizard reaches services step, **When** user selects Redis, **Then** Redis is marked enabled in configuration
2. **Given** wizard reaches services step, **When** user selects multiple services, **Then** all selected services are enabled
3. **Given** user skips optional services, **When** wizard completes, **Then** only core services are configured

---

### User Story 5 - Skip Wizard with Defaults (Priority: P5)

As an experienced user, I want to skip the wizard and use defaults so that I can set up quickly without prompts.

**Why this priority**: Experienced users shouldn't be slowed down by interactive prompts.

**Independent Test**: Run `20i init --skip-wizard`, verify stack is configured with defaults without any prompts.

**Acceptance Scenarios**:

1. **Given** the `--skip-wizard` flag, **When** the user runs `20i init --skip-wizard`, **Then** no prompts are shown
2. **Given** `--skip-wizard` is used, **When** initialization completes, **Then** default configuration is applied
3. **Given** both `--skip-wizard` and `--template laravel`, **When** init runs, **Then** Laravel template is used without prompts

---

### Edge Cases

- What happens when the terminal doesn't support interactive input (CI/CD pipelines)?
- How does the wizard handle invalid input (non-numeric port, invalid PHP version)?
- What happens when user cancels wizard mid-way (Ctrl+C)?
- How does the wizard behave in non-TTY environments?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Wizard MUST guide users through PHP version selection
- **FR-002**: Wizard MUST guide users through database configuration (name, credentials)
- **FR-003**: Wizard MUST guide users through port selection with conflict detection
- **FR-004**: Wizard MUST guide users through optional services selection
- **FR-005**: Wizard MUST detect project type from existing files when possible
- **FR-006**: Wizard MUST validate all user inputs before proceeding
- **FR-007**: Wizard MUST support `--skip-wizard` flag to use defaults
- **FR-008**: Wizard MUST provide clear feedback and progress indication
- **FR-009**: Wizard MUST work in non-interactive mode for CI/CD environments

### Key Entities

- **Setup Wizard**: Interactive command-line interface that guides users through configuration choices
- **Project Detection**: Logic that analyzes directory contents to identify framework/project type
- **Port Conflict Checker**: Utility that scans system for ports in use and reports conflicts

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Non-technical users can complete setup without reading documentation
- **SC-002**: Wizard completes in under 2 minutes for typical setup
- **SC-003**: Project type detection is accurate for 90%+ of Laravel, WordPress, and Symfony projects
- **SC-004**: Port conflicts are detected and reported before they cause startup failures
- **SC-005**: 95%+ of first-time users successfully complete wizard and start stack
- **SC-006**: Experienced users can skip wizard and initialize in under 5 seconds

## Assumptions

- Users have basic terminal familiarity (can type responses and press enter)
- Terminal supports standard input/output for interactive prompts
- Common frameworks have identifiable project files (composer.json, wp-config.php, etc.)
- Port scanning is available and permitted on the host system
