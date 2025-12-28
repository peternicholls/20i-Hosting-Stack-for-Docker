# Feature Specification: Smart Version Override System

**Feature Branch**: `001-smart-version-override`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟡 High  
**Input**: User description: "Allow per-project PHP version overrides with automatic validation"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Set Custom PHP Version Per Project (Priority: P1)

As a developer working on multiple PHP projects, I want to specify a PHP version in my project's `.20i-local` file so that each project uses the correct PHP version regardless of the global default.

**Why this priority**: This is the core functionality - without per-project version specification, the feature has no value.

**Independent Test**: Create a `.20i-local` file with `PHP_VERSION=8.3`, start the stack, and verify the Apache container runs PHP 8.3 by checking `php -v` output.

**Acceptance Scenarios**:

1. **Given** a project with `.20i-local` containing `PHP_VERSION=8.3`, **When** the stack starts, **Then** the Apache container builds and runs with PHP 8.3
2. **Given** a project with no `.20i-local` file, **When** the stack starts, **Then** the Apache container uses the global default PHP version from `config/stack-vars.yml`
3. **Given** a project with `.20i-local` containing `PHP_VERSION=8.4`, **When** the stack starts, **Then** the Apache container builds and runs with PHP 8.4

---

### User Story 2 - Validate PHP Version Against Available Versions (Priority: P2)

As a developer, I want the system to validate my chosen PHP version against available Docker Hub tags so that I receive immediate feedback if I specify an invalid version.

**Why this priority**: Validation prevents frustrating build failures and provides a better developer experience.

**Independent Test**: Set `PHP_VERSION=9.0` in `.20i-local`, attempt to start the stack, and verify a clear error message is displayed with available options.

**Acceptance Scenarios**:

1. **Given** a `.20i-local` with `PHP_VERSION=9.0` (non-existent), **When** the stack starts, **Then** the system displays an error message listing available PHP versions
2. **Given** a `.20i-local` with `PHP_VERSION=8.3` (valid), **When** the stack starts, **Then** validation passes silently and the stack starts normally
3. **Given** network unavailable for Docker Hub API, **When** the stack starts with a custom PHP version, **Then** validation is skipped with a warning and the build proceeds

---

### User Story 3 - Warn About EOL PHP Versions (Priority: P3)

As a developer, I want to be warned if I'm using an end-of-life PHP version so that I can plan upgrades proactively.

**Why this priority**: This is a quality-of-life improvement that helps developers stay secure without blocking their work.

**Independent Test**: Set `PHP_VERSION=7.4` (EOL), start the stack, and verify a warning message appears suggesting an upgrade path.

**Acceptance Scenarios**:

1. **Given** a `.20i-local` with `PHP_VERSION=7.4` (EOL), **When** the stack starts, **Then** a warning is displayed suggesting upgrade to a supported version
2. **Given** a `.20i-local` with `PHP_VERSION=8.3` (supported), **When** the stack starts, **Then** no EOL warning is displayed
3. **Given** the `--quiet` flag is passed, **When** the stack starts with an EOL version, **Then** warnings are suppressed

---

### User Story 4 - Select PHP Version from GUI (Priority: P4)

As a developer using the macOS GUI app, I want to select a PHP version from a dropdown menu so that I don't need to manually edit configuration files.

**Why this priority**: GUI integration is an enhancement that improves accessibility but isn't essential for core functionality.

**Independent Test**: Open the 20i Stack Manager app, select PHP 8.4 from the dropdown, start the stack, and verify PHP 8.4 is running.

**Acceptance Scenarios**:

1. **Given** the GUI app is open, **When** the user clicks the PHP version dropdown, **Then** a list of available PHP versions (8.1, 8.2, 8.3, 8.4, 8.5) is displayed
2. **Given** the user selects PHP 8.4 from the dropdown, **When** the user clicks Start, **Then** the stack starts with PHP 8.4
3. **Given** the user has a `.20i-local` with a PHP version, **When** the GUI loads, **Then** the dropdown shows the currently configured version as selected

---

### Edge Cases

- What happens when `.20i-local` contains an empty `PHP_VERSION=` value?
- How does the system handle malformed version strings like `PHP_VERSION=8.x` or `PHP_VERSION=latest`?
- What happens if Docker Hub API rate limits are exceeded during validation?
- How does the system behave when `.20i-local` has conflicting entries (e.g., `PHP_VERSION` defined twice)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST read `PHP_VERSION` from project-local `.20i-local` file if present
- **FR-002**: System MUST fall back to `config/stack-vars.yml` default when no local override exists
- **FR-003**: System MUST validate specified PHP version against Docker Hub available tags
- **FR-004**: System MUST display clear error message listing available versions when invalid version specified
- **FR-005**: System MUST display warning when EOL PHP version is detected (7.4, 8.0)
- **FR-006**: System MUST support PHP versions 8.1, 8.2, 8.3, 8.4, and 8.5
- **FR-007**: GUI MUST display dropdown with available PHP versions
- **FR-008**: GUI MUST persist selected PHP version to `.20i-local` file
- **FR-009**: System MUST skip validation gracefully when offline, with warning message

### Key Entities

- **PHP Version**: A version string (e.g., "8.3") representing the PHP runtime version, with attributes: version number, EOL status, Docker Hub availability
- **Local Configuration**: Project-specific overrides stored in `.20i-local` file, with attributes: PHP version, other future overrides
- **Global Configuration**: Default values stored in `config/stack-vars.yml`, applies when no local override exists

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can set a custom PHP version and see it reflected in running container within 60 seconds of stack start
- **SC-002**: Invalid PHP version errors display within 5 seconds of stack start attempt
- **SC-003**: 100% of supported PHP versions (8.1-8.5) successfully build and run
- **SC-004**: EOL warnings display for all PHP versions with security support ended
- **SC-005**: GUI version selection takes fewer than 3 clicks from app launch to version change
- **SC-006**: Zero user confusion about which PHP version is running (clear feedback in logs/GUI)

## Assumptions

- Docker Hub API remains freely accessible for version tag queries
- PHP Docker images follow semantic versioning pattern (8.x format)
- Users have internet connectivity for initial version validation (offline fallback exists)
- EOL dates follow official PHP release schedule

## Files Affected

- `20i-gui` - Add PHP version dropdown and selection logic
- `zsh-example-script.zsh` - Add version reading and validation
- `20i-stack-manager.scpt` - Add AppleScript version selection UI
- `config/stack-vars.yml` - Document default PHP version
