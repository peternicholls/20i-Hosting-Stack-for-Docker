# Feature Specification: Project Templates

**Feature Branch**: `008-project-templates`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Quick-start templates for common frameworks like Laravel, WordPress, Symfony, and plain PHP"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Initialize with Framework Template (Priority: P1)

As a developer starting a new project, I want to initialize the stack with a framework-specific template so that I have optimal configuration and boilerplate for that framework.

**Why this priority**: Template initialization is the entry point - without it, templates provide no value.

**Independent Test**: Run `20i init --template laravel`, verify Laravel-specific compose overrides, environment variables, and convenience scripts are created.

**Acceptance Scenarios**:

1. **Given** an empty directory, **When** the user runs `20i init --template laravel`, **Then** Laravel-specific configuration is created
2. **Given** Laravel template is used, **When** viewing the configuration, **Then** queue worker and scheduler services are included
3. **Given** `20i init --template wordpress`, **When** initialization completes, **Then** WordPress-specific configuration with WP-CLI is created

---

### User Story 2 - List Available Templates (Priority: P2)

As a developer, I want to see what templates are available so that I can choose the right one for my project.

**Why this priority**: Template discovery helps users make informed choices before initialization.

**Independent Test**: Run `20i templates`, verify list shows all available templates with descriptions.

**Acceptance Scenarios**:

1. **Given** the CLI is installed, **When** the user runs `20i templates`, **Then** a list of available templates is displayed
2. **Given** template listing, **When** viewing each template, **Then** a brief description and included features are shown
3. **Given** a specific template, **When** the user runs `20i templates laravel --details`, **Then** detailed information about the template is shown

---

### User Story 3 - Use Plain PHP Template (Priority: P3)

As a developer working on a custom PHP project, I want a plain PHP template so that I get a clean setup without framework-specific overhead.

**Why this priority**: Plain PHP is the default/fallback template for users not using a framework.

**Independent Test**: Run `20i init` without template flag, verify a minimal PHP configuration is created.

**Acceptance Scenarios**:

1. **Given** no template specified, **When** the user runs `20i init`, **Then** plain PHP template is used by default
2. **Given** plain PHP template, **When** viewing configuration, **Then** only core services (Apache, MariaDB, phpMyAdmin) are configured
3. **Given** `20i init --template php`, **When** initialization completes, **Then** same result as no template flag

---

### User Story 4 - Template Includes Framework Conveniences (Priority: P4)

As a Laravel developer, I want artisan command shortcuts so that I can run common commands without entering the container.

**Why this priority**: Conveniences enhance developer experience but aren't required for basic functionality.

**Independent Test**: Initialize with Laravel template, run `20i artisan migrate`, verify command executes inside container and returns output.

**Acceptance Scenarios**:

1. **Given** Laravel template initialized, **When** the user runs `20i artisan migrate`, **Then** artisan command runs inside the PHP container
2. **Given** WordPress template initialized, **When** the user runs `20i wp plugin list`, **Then** WP-CLI command runs inside the PHP container
3. **Given** Symfony template initialized, **When** the user runs `20i console cache:clear`, **Then** Symfony console command executes

---

### Edge Cases

- What happens when using a template in a non-empty directory?
- How does the system handle template initialization when framework files already exist?
- What happens when a template requires services not yet available?
- How does the system handle corrupted or incomplete template downloads?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST support `20i init --template <name>` to initialize with specific template
- **FR-002**: CLI MUST support `20i templates` to list available templates
- **FR-003**: System MUST provide Laravel template with queue worker and scheduler
- **FR-004**: System MUST provide WordPress template with WP-CLI integration
- **FR-005**: System MUST provide Symfony template with console integration
- **FR-006**: System MUST provide plain PHP template as default
- **FR-007**: Templates MUST include framework-specific environment variables
- **FR-008**: Templates MUST include convenience commands for common framework operations
- **FR-009**: Template selection MUST be persisted in `.20i-config.yml`

### Key Entities

- **Project Template**: Predefined configuration set for a specific framework, includes compose overrides, environment variables, and convenience scripts
- **Template Metadata**: Information about a template including name, description, included services, and framework version compatibility
- **Convenience Command**: CLI shortcut that executes framework-specific commands inside containers

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Template initialization adds less than 10 seconds to base `init` time
- **SC-002**: Developers can start coding in their framework within 5 minutes of init
- **SC-003**: Framework-specific services start correctly on 100% of template initializations
- **SC-004**: Convenience commands reduce common operations to single CLI commands
- **SC-005**: Template documentation is accessible via `20i templates <name> --help`

## Assumptions

- Users know which framework they want to use before initialization
- Templates target current stable versions of each framework
- Framework-specific services use standard configurations
- Templates are updated when major framework versions are released
