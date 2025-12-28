# Feature Specification: Docker Distribution Pre-built Images

**Feature Branch**: `006-prebuilt-images`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟢 Medium  
**Input**: User description: "Publish pre-built PHP-FPM images to Docker Hub/GHCR for faster stack startup"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Start Stack with Pre-built Images (Priority: P1)

As a developer, I want the stack to use pre-built images by default so that startup is fast without waiting for image builds.

**Why this priority**: Faster startup directly improves developer experience and is the core benefit of pre-built images.

**Independent Test**: Run `20i start` with default configuration, verify images are pulled from registry (not built locally) and stack is ready in under 30 seconds.

**Acceptance Scenarios**:

1. **Given** default stack configuration, **When** the user runs `20i start`, **Then** pre-built PHP-FPM images are pulled from registry
2. **Given** images already cached locally, **When** the user runs `20i start`, **Then** stack starts in under 15 seconds
3. **Given** a new PHP version is specified, **When** the user runs `20i start`, **Then** the corresponding pre-built image is pulled

---

### User Story 2 - Multi-Architecture Support (Priority: P2)

As a developer using Apple Silicon or ARM-based Linux, I want pre-built images to support my architecture so that I get native performance.

**Why this priority**: ARM support is essential for modern Mac users and cloud deployments.

**Independent Test**: Pull the image on ARM64 machine, verify it runs natively without emulation warnings.

**Acceptance Scenarios**:

1. **Given** an ARM64 machine (M1/M2 Mac), **When** the stack starts, **Then** ARM64 native image is pulled and runs without emulation
2. **Given** an x64 machine, **When** the stack starts, **Then** AMD64 image is pulled and runs natively
3. **Given** manifest inspection, **When** viewing image details, **Then** both AMD64 and ARM64 architectures are listed

---

### User Story 3 - Fall Back to Local Build (Priority: P3)

As a developer, I want the option to build images locally so that I can customize the image or work offline.

**Why this priority**: Fallback ensures users aren't blocked if registry is unavailable or customization is needed.

**Independent Test**: Set configuration to use local build, run `20i start`, verify Dockerfile is built locally.

**Acceptance Scenarios**:

1. **Given** `USE_PREBUILT: false` in configuration, **When** the user runs `20i start`, **Then** images are built from local Dockerfile
2. **Given** registry is unavailable, **When** pulling fails, **Then** system prompts to use local build as fallback
3. **Given** custom Dockerfile modifications, **When** user builds locally, **Then** modifications are included in the running container

---

### User Story 4 - Image Versioning Aligned with Stack Releases (Priority: P4)

As a maintainer, I want images published for each stack release so that users can pin to specific versions.

**Why this priority**: Version pinning enables reproducible environments and safe upgrades.

**Independent Test**: Pull image tagged with specific version (e.g., `v2.0.0`), verify it matches the stack configuration from that release.

**Acceptance Scenarios**:

1. **Given** release v2.0.0 is published, **When** images are built, **Then** images are tagged with `v2.0.0` and `latest`
2. **Given** `PHP_VERSION: 8.4`, **When** viewing available images, **Then** `peternicholls/20i-php-fpm:8.4` is available
3. **Given** `latest` tag, **When** pulling, **Then** it points to the most recent stable PHP version (e.g., 8.5)

---

### Edge Cases

- What happens when pre-built image doesn't exist for requested PHP version?
- How does the system handle registry authentication if images are private?
- What happens when local and pre-built images have different capabilities?
- How are images invalidated when security patches are needed?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST publish pre-built PHP-FPM images to Docker Hub and/or GitHub Container Registry
- **FR-002**: System MUST support multiple PHP versions (8.1, 8.2, 8.3, 8.4, 8.5)
- **FR-003**: System MUST build multi-architecture images (AMD64 and ARM64)
- **FR-004**: System MUST tag images with version numbers and `latest`
- **FR-005**: System MUST fall back to local build when pre-built images are unavailable
- **FR-006**: System MUST update images on every stack release
- **FR-007**: Documentation MUST explain pre-built vs custom build options
- **FR-008**: System MUST provide image checksums/digests for verification

### Key Entities

- **Pre-built Image**: Docker image published to registry, ready to pull without local build
- **Image Tag**: Version identifier (e.g., `8.4`, `v2.0.0`, `latest`) for pulling specific image versions
- **Multi-arch Manifest**: Docker manifest listing available architectures for a single image tag

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: First stack startup time reduced by 80% compared to local build (from ~5 minutes to ~1 minute)
- **SC-002**: Subsequent startups complete in under 15 seconds with cached images
- **SC-003**: 100% of supported PHP versions have pre-built images available
- **SC-004**: Images available for both AMD64 and ARM64 architectures
- **SC-005**: Image size is under 500MB for each PHP version
- **SC-006**: Images are updated within 24 hours of stack releases

## Assumptions

- Docker Hub or GHCR provides sufficient bandwidth for image pulls
- Users have internet connectivity for initial image pull
- Multi-arch build process is reliable and reproducible
- Image tag naming follows Docker conventions
