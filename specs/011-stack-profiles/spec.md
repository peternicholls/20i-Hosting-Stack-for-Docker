# Feature Specification: Stack Profiles

**Feature Branch**: `011-stack-profiles`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Save and switch between multiple configurations for different development scenarios"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Save Current Configuration as Profile (Priority: P1)

As a developer, I want to save my current stack configuration as a named profile so that I can return to it later.

**Why this priority**: Saving profiles is the foundation - without it, profiles cannot be created or switched.

**Independent Test**: Configure stack with PHP 8.3 and Redis enabled, run `20i profile save dev-php83`, verify profile file is created with current settings.

**Acceptance Scenarios**:

1. **Given** a configured stack, **When** the user runs `20i profile save dev-php83`, **Then** current configuration is saved to a named profile
2. **Given** a profile name already exists, **When** saving with same name, **Then** system prompts for confirmation to overwrite
3. **Given** the `--force` flag, **When** saving an existing profile name, **Then** profile is overwritten without prompt

---

### User Story 2 - Switch Between Profiles (Priority: P1)

As a developer testing compatibility, I want to switch between profiles so that I can quickly change PHP versions or service configurations.

**Why this priority**: Profile switching is the primary use case that delivers value.

**Independent Test**: Create profiles for PHP 8.3 and PHP 8.4, switch between them, verify stack restarts with correct PHP version each time.

**Acceptance Scenarios**:

1. **Given** profiles `dev-php83` and `dev-php84` exist, **When** the user runs `20i profile use dev-php84`, **Then** configuration is switched to PHP 8.4 settings
2. **Given** a profile is activated, **When** stack is restarted, **Then** it uses the activated profile's configuration
3. **Given** profile includes enabled services (Redis), **When** profile is activated, **Then** those services are enabled/disabled accordingly

---

### User Story 3 - List Available Profiles (Priority: P2)

As a developer, I want to see all my saved profiles so that I can choose which one to use.

**Why this priority**: Listing enables discovery of available profiles for selection.

**Independent Test**: Create multiple profiles, run `20i profile list`, verify all profiles are shown with basic info.

**Acceptance Scenarios**:

1. **Given** multiple profiles exist, **When** the user runs `20i profile list`, **Then** all profile names are displayed
2. **Given** profile listing, **When** viewing output, **Then** each profile shows key settings (PHP version, enabled services)
3. **Given** a profile is currently active, **When** listing profiles, **Then** the active profile is marked/highlighted

---

### User Story 4 - Delete Profiles (Priority: P3)

As a developer, I want to delete profiles I no longer need so that I can keep my profile list clean.

**Why this priority**: Deletion is maintenance functionality, less critical than core profile operations.

**Independent Test**: Run `20i profile delete dev-old`, verify profile is removed from storage.

**Acceptance Scenarios**:

1. **Given** a profile exists, **When** the user runs `20i profile delete dev-old`, **Then** the profile is removed
2. **Given** the currently active profile, **When** user tries to delete it, **Then** system warns and requires confirmation
3. **Given** a non-existent profile name, **When** trying to delete, **Then** error message indicates profile not found

---

### User Story 5 - Profile Metadata and History (Priority: P4)

As a developer, I want profiles to track when they were created and last used so that I can identify stale profiles.

**Why this priority**: Metadata is helpful for management but not essential for core functionality.

**Independent Test**: Create a profile, use it, run `20i profile list --details`, verify creation and last-used timestamps are shown.

**Acceptance Scenarios**:

1. **Given** a profile is created, **When** viewing profile details, **Then** creation timestamp is shown
2. **Given** a profile is activated, **When** viewing profile details later, **Then** last-used timestamp is updated
3. **Given** `--details` flag on list command, **When** listing profiles, **Then** extended metadata is displayed

---

### Edge Cases

- What happens when switching profiles while stack is running?
- How does the system handle corrupted profile files?
- What happens when a profile references services that are no longer available?
- How does the system handle profile migration when configuration format changes?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST support `20i profile save <name>` to save current configuration as named profile
- **FR-002**: CLI MUST support `20i profile use <name>` to switch to a saved profile
- **FR-003**: CLI MUST support `20i profile list` to display all saved profiles
- **FR-004**: CLI MUST support `20i profile delete <name>` to remove a profile
- **FR-005**: Profiles MUST capture PHP version, enabled services, port mappings, and database settings
- **FR-006**: Profiles MUST store creation timestamp and last-used timestamp
- **FR-007**: Profile storage MUST be in a dedicated `.20i-profiles/` directory
- **FR-008**: Active profile MUST be tracked in `.20i-config.yml`
- **FR-009**: Profile switching MUST not require stack restart unless configuration differs

### Key Entities

- **Profile**: Named configuration snapshot containing PHP version, services, ports, and other settings
- **Profile Storage**: Directory (`.20i-profiles/`) containing profile files as YAML
- **Active Profile**: Currently selected profile tracked in project configuration

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Profile switching completes in under 5 seconds (excluding stack restart time)
- **SC-002**: Profiles correctly restore all saved configuration settings 100% of the time
- **SC-003**: Developers can switch between PHP versions without editing files
- **SC-004**: Profile list command provides enough information to choose the right profile
- **SC-005**: Profile storage uses less than 1KB per profile

## Assumptions

- Users have limited number of profiles per project (typically under 10)
- Profile names follow simple naming conventions (alphanumeric, hyphens)
- Configuration format is stable between minor versions
- File system storage is reliable for profile persistence
