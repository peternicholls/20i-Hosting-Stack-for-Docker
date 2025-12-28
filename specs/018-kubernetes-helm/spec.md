# Feature Specification: Kubernetes Helm Charts

**Feature Branch**: `018-kubernetes-helm`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: ⚪ Low (Future Exploration)  
**Input**: User description: "Production deployment option with Helm chart mirroring local stack structure"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Deploy Stack to Kubernetes Cluster (Priority: P1)

As a developer transitioning to production, I want to deploy my 20i stack to Kubernetes so that I can use familiar configuration in a production environment.

**Why this priority**: Deployment is the primary value proposition - enabling production use of the stack.

**Independent Test**: Run `helm install` with 20i chart, verify all services deploy to Kubernetes cluster and are accessible.

**Acceptance Scenarios**:

1. **Given** a Kubernetes cluster, **When** the user runs `helm install 20i-stack ./chart`, **Then** PHP, MariaDB, and Nginx pods are created
2. **Given** successful deployment, **When** viewing pods, **Then** all pods reach "Running" state
3. **Given** deployment completes, **When** accessing the service URL, **Then** the PHP application responds

---

### User Story 2 - Configure via Helm Values (Priority: P2)

As a DevOps engineer, I want to customize the deployment via Helm values so that I can adjust resources and settings for my environment.

**Why this priority**: Configurability is essential for production use across different environments.

**Independent Test**: Create custom `values.yaml` with increased replicas, deploy with values, verify pod count matches configuration.

**Acceptance Scenarios**:

1. **Given** custom `values.yaml` with `replicas: 3`, **When** deploying chart, **Then** 3 PHP pods are created
2. **Given** custom PHP version in values, **When** deploying, **Then** pods use specified PHP version
3. **Given** custom resource limits in values, **When** viewing pods, **Then** resource requests/limits match configuration

---

### User Story 3 - Mirror Local Stack Structure (Priority: P3)

As a developer, I want the Helm chart to mirror my local stack so that production behavior matches development.

**Why this priority**: Parity between local and production reduces deployment surprises.

**Independent Test**: Compare local Docker Compose services to Helm chart services, verify same services with same relationships.

**Acceptance Scenarios**:

1. **Given** local stack with Apache, MariaDB, phpMyAdmin, **When** viewing Helm chart, **Then** same services are defined
2. **Given** enabled optional service (Redis) locally, **When** deploying with Redis enabled, **Then** Redis pod is created
3. **Given** local environment variables, **When** viewing Helm templates, **Then** equivalent ConfigMaps/Secrets are created

---

### User Story 4 - Include Production-Ready Defaults (Priority: P4)

As a DevOps engineer, I want production-ready defaults in the Helm chart so that deployments are secure and scalable out of the box.

**Why this priority**: Good defaults reduce misconfiguration risks in production.

**Independent Test**: Deploy chart with defaults, verify security contexts, resource limits, and health checks are configured.

**Acceptance Scenarios**:

1. **Given** default deployment, **When** viewing pod specs, **Then** security contexts prevent privilege escalation
2. **Given** default deployment, **When** viewing pod specs, **Then** resource requests and limits are defined
3. **Given** default deployment, **When** viewing service specs, **Then** health checks are configured for all services

---

### Edge Cases

- What happens when cluster doesn't support required Kubernetes version?
- How does the chart handle persistent storage in different cloud providers?
- What happens when ingress controller is not installed?
- How does the chart handle secrets management in production?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Helm chart MUST deploy PHP-FPM, MariaDB, and Nginx services
- **FR-002**: Helm chart MUST support optional services (Redis, Mailhog) via values
- **FR-003**: Helm chart MUST support configurable PHP version
- **FR-004**: Helm chart MUST include production-ready security contexts
- **FR-005**: Helm chart MUST include configurable resource requests/limits
- **FR-006**: Helm chart MUST include health checks for all services
- **FR-007**: Helm chart MUST support Ingress configuration for external access
- **FR-008**: Helm chart MUST support persistent volume claims for database storage
- **FR-009**: Documentation MUST explain migration from local to Kubernetes deployment

### Key Entities

- **Helm Chart**: Packaged Kubernetes application containing templates, values, and dependencies
- **Values File**: Configuration file for customizing Helm chart deployment
- **Kubernetes Deployment**: Running instance of the stack in a Kubernetes cluster

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Helm chart deploys successfully to major cloud providers (AWS EKS, GCP GKE, Azure AKS)
- **SC-002**: Deployment completes in under 5 minutes on typical clusters
- **SC-003**: Production defaults pass common security scanning tools
- **SC-004**: Application behavior in Kubernetes matches local Docker behavior
- **SC-005**: Documentation enables deployment without Kubernetes expertise

## Assumptions

- Users have access to a Kubernetes cluster with Helm installed
- Persistent storage is available in the target cluster
- Ingress controller is installed for external access
- Users understand basic Kubernetes concepts (pods, services, deployments)
