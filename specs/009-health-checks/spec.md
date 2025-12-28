# Feature Specification: Health Checks and Auto-Restart

**Feature Branch**: `009-health-checks`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟡 High  
**Input**: User description: "Ensure services are truly ready before marking stack as started with container health checks and optional auto-restart"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Wait for Services to Be Ready (Priority: P1)

As a developer, I want the `20i start` command to wait until all services are truly ready so that I can immediately start using the stack without connection errors.

**Why this priority**: This is the core value - preventing premature "ready" state that leads to failed connections.

**Independent Test**: Run `20i start`, measure time from command to "ready" message, then immediately connect to web server and database - both should succeed.

**Acceptance Scenarios**:

1. **Given** a fresh stack start, **When** the user runs `20i start`, **Then** the command waits until all health checks pass before reporting "ready"
2. **Given** MariaDB is still initializing, **When** PHP tries to connect, **Then** the connection succeeds because start waits for MariaDB health check
3. **Given** the `--no-wait` flag, **When** the user runs `20i start --no-wait`, **Then** the command returns immediately without waiting for health checks

---

### User Story 2 - Display Clear Health Status (Priority: P2)

As a developer, I want to see the health status of all services so that I can quickly identify which service is having issues.

**Why this priority**: Visibility into health status enables faster troubleshooting.

**Independent Test**: Run `20i status` while stack is starting, verify each service shows its current health state (starting, healthy, unhealthy).

**Acceptance Scenarios**:

1. **Given** a running stack, **When** the user runs `20i status`, **Then** health status (healthy/unhealthy) is shown for each service
2. **Given** a service is unhealthy, **When** viewing status, **Then** a clear error message indicates the problem
3. **Given** health check is in progress, **When** viewing status, **Then** "starting" or "waiting" state is displayed

---

### User Story 3 - Auto-Restart Unhealthy Containers (Priority: P3)

As a developer, I want unhealthy containers to automatically restart so that temporary issues resolve without manual intervention.

**Why this priority**: Auto-restart improves resilience but is less critical than basic health checking.

**Independent Test**: Configure auto-restart, simulate container failure, verify container restarts automatically and becomes healthy.

**Acceptance Scenarios**:

1. **Given** auto-restart is enabled, **When** a container becomes unhealthy, **Then** the container is automatically restarted
2. **Given** auto-restart limit of 3 attempts, **When** a container fails 3 times, **Then** it stays stopped and an alert is shown
3. **Given** auto-restart is disabled, **When** a container becomes unhealthy, **Then** it remains in unhealthy state until manual intervention

---

### User Story 4 - Configurable Health Check Parameters (Priority: P4)

As an advanced user, I want to configure health check intervals and thresholds so that I can tune the behavior for my environment.

**Why this priority**: Configuration is for advanced users; defaults should work for most cases.

**Independent Test**: Set custom health check interval in config, restart stack, verify health checks run at the configured interval.

**Acceptance Scenarios**:

1. **Given** default configuration, **When** health checks run, **Then** they use sensible defaults (10s interval, 3 retries)
2. **Given** custom interval of 5s in config, **When** health checks run, **Then** they execute every 5 seconds
3. **Given** custom retry count of 5, **When** a service fails, **Then** 5 attempts are made before marking unhealthy

---

### Edge Cases

- What happens when health check command itself fails or times out?
- How does the system handle services without defined health checks?
- What happens when all containers become unhealthy simultaneously?
- How does auto-restart interact with `20i stop` command?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: All service containers MUST have defined health check configurations
- **FR-002**: `20i start` MUST wait for all health checks to pass before reporting "ready"
- **FR-003**: `20i status` MUST display health status for each service
- **FR-004**: System MUST support optional auto-restart of unhealthy containers
- **FR-005**: System MUST support `--no-wait` flag to skip health check waiting
- **FR-006**: System MUST display clear error messages when services fail health checks
- **FR-007**: Health check parameters MUST be configurable (interval, timeout, retries)
- **FR-008**: Auto-restart MUST have configurable attempt limits to prevent infinite loops
- **FR-009**: System MUST log health check failures for troubleshooting

### Key Entities

- **Health Check**: Configuration defining how to verify service readiness, with attributes: check command, interval, timeout, retries
- **Service Health State**: Current status of a service (starting, healthy, unhealthy)
- **Auto-Restart Policy**: Rules for automatic container restart, with attributes: enabled flag, max attempts, cooldown period

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Stack reports "ready" only when all services are accessible (zero false positives)
- **SC-002**: Health check adds less than 10 seconds to startup for healthy services
- **SC-003**: Unhealthy services are identified within 30 seconds of failure
- **SC-004**: Auto-restart resolves 80%+ of transient failures without user intervention
- **SC-005**: Health status is visible within 1 second of running `20i status`
- **SC-006**: Zero "connection refused" errors when using stack immediately after "ready" message

## Assumptions

- Services have predictable startup times within configured health check windows
- Health check commands are lightweight and don't impact service performance
- Container orchestration supports health check and restart policies
- Network connectivity between containers is reliable
