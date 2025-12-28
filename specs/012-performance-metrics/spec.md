# Feature Specification: Performance Metrics and Insights

**Feature Branch**: `012-performance-metrics`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Show resource usage, startup time, and actionable recommendations for optimization"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Resource Usage per Service (Priority: P1)

As a developer, I want to see CPU and memory usage for each service so that I can identify resource-hungry containers.

**Why this priority**: Resource visibility is the core value - understanding where resources are consumed.

**Independent Test**: Run `20i stats` with stack running, verify CPU and memory usage is displayed for each service.

**Acceptance Scenarios**:

1. **Given** a running stack, **When** the user runs `20i stats`, **Then** CPU and memory usage for each container is displayed
2. **Given** stats output, **When** viewing each service, **Then** current usage and limits (if set) are shown
3. **Given** a service using high resources, **When** stats are displayed, **Then** the high-usage service is highlighted or flagged

---

### User Story 2 - Track Startup Time (Priority: P2)

As a developer, I want to see how long stack startup takes so that I can identify slow components and optimize.

**Why this priority**: Startup time affects daily developer experience; tracking enables optimization.

**Independent Test**: Run `20i start`, verify startup time is displayed upon completion and recorded for future reference.

**Acceptance Scenarios**:

1. **Given** a stack starting, **When** startup completes, **Then** total startup time is displayed
2. **Given** startup time tracking, **When** user runs `20i stats`, **Then** last startup time and average are shown
3. **Given** individual service startup, **When** viewing detailed stats, **Then** per-service startup time is available

---

### User Story 3 - Receive Actionable Recommendations (Priority: P3)

As a developer, I want to receive recommendations for improving performance so that I can optimize my stack without research.

**Why this priority**: Recommendations add value by translating data into action, but require metrics first.

**Independent Test**: Configure MariaDB with low memory limit, run `20i stats`, verify recommendation is shown to increase limit.

**Acceptance Scenarios**:

1. **Given** MariaDB using 80%+ of memory limit, **When** stats are shown, **Then** recommendation to increase limit is displayed
2. **Given** startup takes over 60 seconds, **When** stats are shown, **Then** recommendation to use pre-built images is displayed
3. **Given** disk usage is high, **When** stats are shown, **Then** recommendation to prune unused images/volumes is displayed

---

### User Story 4 - View Historical Metrics (Priority: P4)

As a developer, I want to see historical resource usage so that I can identify trends and intermittent issues.

**Why this priority**: History is valuable for troubleshooting but less critical than real-time stats.

**Independent Test**: Run stack for extended period, run `20i stats --history`, verify historical data is displayed.

**Acceptance Scenarios**:

1. **Given** stack has been running for 1 hour, **When** the user runs `20i stats --history`, **Then** historical metrics are shown
2. **Given** historical data available, **When** viewing history, **Then** data is shown in time-series format (table or simple graph)
3. **Given** no historical data (first run), **When** history is requested, **Then** message indicates not enough data yet

---

### Edge Cases

- What happens when container runtime doesn't support stats API?
- How does the system handle services that restart frequently (skewing averages)?
- What happens when historical data storage becomes too large?
- How does the system handle stats collection when stack is stopped?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST support `20i stats` command to display real-time resource usage
- **FR-002**: System MUST show CPU usage (percentage) per service
- **FR-003**: System MUST show memory usage (current/limit) per service
- **FR-004**: System MUST show disk usage for volumes
- **FR-005**: System MUST track and display startup time for stack and individual services
- **FR-006**: System MUST provide actionable recommendations based on metrics
- **FR-007**: Recommendations MUST be contextual and specific to observed issues
- **FR-008**: System SHOULD track historical metrics for trend analysis
- **FR-009**: Stats command MUST work without additional tools/dependencies

### Key Entities

- **Service Metrics**: Real-time resource usage data for a container (CPU, memory, network, disk)
- **Startup Metrics**: Time measurements for stack and service startup
- **Recommendation**: Actionable suggestion based on observed metrics, with context and action

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `20i stats` displays resource usage within 2 seconds
- **SC-002**: Recommendations are accurate and actionable 80%+ of the time
- **SC-003**: Developers can identify resource bottlenecks within 30 seconds
- **SC-004**: Startup time improvements are trackable after implementing recommendations
- **SC-005**: Stats collection has negligible impact on container performance (<1% overhead)

## Assumptions

- Container runtime provides stats API (standard in Docker)
- Resource usage data is accurate and available in real-time
- Historical storage requirements are reasonable (hourly aggregates)
- Recommendations are based on common developer scenarios
