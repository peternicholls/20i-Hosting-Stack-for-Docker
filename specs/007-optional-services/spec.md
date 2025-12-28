# Feature Specification: Optional Service Modules

**Feature Branch**: `007-optional-services`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Easy addition of common development services like Redis, Mailhog, Elasticsearch via CLI enable/disable commands"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enable Optional Service via CLI (Priority: P1)

As a developer, I want to enable a service like Redis with a single command so that I can add caching to my project without manually editing configuration files.

**Why this priority**: Simple enablement is the core value proposition - reducing complexity of adding services.

**Independent Test**: Run `20i enable redis`, restart the stack, verify Redis container is running and accessible on port 6379.

**Acceptance Scenarios**:

1. **Given** a running stack without Redis, **When** the user runs `20i enable redis`, **Then** Redis service is added to configuration
2. **Given** Redis is enabled, **When** the stack restarts, **Then** Redis container starts and is accessible on port 6379
3. **Given** `20i enable mailhog`, **When** the stack restarts, **Then** Mailhog SMTP is available on port 1025 and web UI on port 8025

---

### User Story 2 - Disable Optional Service via CLI (Priority: P2)

As a developer, I want to disable a service I no longer need so that I can free up resources and reduce stack complexity.

**Why this priority**: Disabling services completes the lifecycle management of optional services.

**Independent Test**: Run `20i disable redis` on stack with Redis enabled, restart stack, verify Redis container is not running.

**Acceptance Scenarios**:

1. **Given** Redis is enabled in the stack, **When** the user runs `20i disable redis`, **Then** Redis is removed from configuration
2. **Given** Redis is disabled, **When** the stack restarts, **Then** no Redis container is started
3. **Given** the `--remove-data` flag, **When** disabling a service, **Then** associated volumes are also removed

---

### User Story 3 - List Available and Enabled Services (Priority: P3)

As a developer, I want to see which optional services are available and which are currently enabled so that I can make informed decisions.

**Why this priority**: Visibility into service state helps users manage their stack effectively.

**Independent Test**: Run `20i services`, verify output shows all available services with their enabled/disabled status.

**Acceptance Scenarios**:

1. **Given** a stack with Redis enabled, **When** the user runs `20i services`, **Then** Redis shows as "enabled" and other services show as "available"
2. **Given** no optional services enabled, **When** the user runs `20i services`, **Then** all services show as "available" with descriptions
3. **Given** a service is running, **When** the user runs `20i services --status`, **Then** running services show health and port information

---

### User Story 4 - Persist Service Configuration (Priority: P4)

As a developer, I want my enabled services to be remembered across stack restarts so that I don't need to re-enable them each time.

**Why this priority**: Persistence is expected behavior but is dependent on basic enable/disable working first.

**Independent Test**: Enable Redis, stop stack, start stack again, verify Redis is still running without re-enabling.

**Acceptance Scenarios**:

1. **Given** Redis is enabled, **When** the user runs `20i stop` then `20i start`, **Then** Redis starts automatically
2. **Given** service configuration in `.20i-config.yml`, **When** viewing the file, **Then** enabled services are listed
3. **Given** project is cloned fresh, **When** `.20i-config.yml` exists, **Then** enabled services start with the stack

---

### Edge Cases

- What happens when enabling a service that conflicts with an existing port?
- How does the system handle enabling a service that requires another service (dependencies)?
- What happens when disabling a service that other services depend on?
- How does the system behave when service image pull fails?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST support `20i enable <service>` command to enable optional services
- **FR-002**: CLI MUST support `20i disable <service>` command to disable optional services
- **FR-003**: CLI MUST support `20i services` command to list available and enabled services
- **FR-004**: System MUST persist service enablement in `.20i-config.yml`
- **FR-005**: System MUST support Redis service (caching, port 6379)
- **FR-006**: System MUST support Mailhog service (email testing, SMTP 1025, UI 8025)
- **FR-007**: System MUST support Elasticsearch service (search, port 9200)
- **FR-008**: System MUST support RabbitMQ service (messaging, port 5672, UI 15672)
- **FR-009**: System MUST support MinIO service (S3 storage, port 9000)
- **FR-010**: Services MUST start automatically with stack when enabled

### Key Entities

- **Optional Service**: A container service that can be enabled/disabled per-project, with attributes: name, description, ports, dependencies
- **Service Configuration**: Persistent record of enabled services in `.20i-config.yml`
- **Service Module**: Modular compose definition for each optional service

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Enabling a service takes under 30 seconds (excluding image pull time)
- **SC-002**: Disabling a service takes under 10 seconds
- **SC-003**: Service configuration persists across 100% of stack restarts
- **SC-004**: Zero manual file editing required to enable/disable services
- **SC-005**: All services start with stack in under 60 seconds after enable
- **SC-006**: Service documentation accessible via `20i help services`

## Assumptions

- Service images are available from public Docker registries
- Default ports don't conflict with host services (can be customized)
- Services have reasonable default configurations for development
- Users understand service purposes from brief descriptions
