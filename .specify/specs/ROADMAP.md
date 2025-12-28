# 20i Stack Roadmap

**Status**: Draft | **Version**: 1.0.0 | **Last Updated**: 2025-12-28

This roadmap outlines planned improvements, features, and architectural enhancements for the 20i Stack project. Items are organized by priority and complexity.

---

## Legend

- 🔴 **Critical** - Core functionality or security
- 🟡 **High** - Significant value add or user-requested
- 🟢 **Medium** - Nice to have, quality of life
- ⚪ **Low** - Future exploration

**Status Icons**: 📋 Planned | 🚧 In Progress | ✅ Complete | ⏸️ Deferred

---

## 1. Configuration & Defaults Management

### 1.1 Smart Version Override System 🟡 📋
**Goal**: Allow per-project PHP version overrides with automatic validation

**Features**:
- Detect available PHP versions from Docker Hub API
- Validate `PHP_VERSION` in `.20i-local` against available tags
- Warn if using EOL PHP version with upgrade suggestions
- GUI dropdown for version selection (8.1, 8.2, 8.3, 8.4, 8.5)

**Files Affected**: `20i-gui`, `zsh-example-script.zsh`, `20i-stack-manager.scpt`

**Acceptance Criteria**:
- User sets `PHP_VERSION=8.3` in `.20i-local` → Apache container builds with PHP 8.3
- Invalid version (e.g., `PHP_VERSION=9.0`) → Clear error message with available options
- GUI shows only valid/available PHP versions

---

### 1.2 Extended Central YAML Configuration 🟢 📋
**Goal**: Add more stack variables to `config/stack-vars.yml`

**New Variables**:
```yaml
# Network & Ports
HOST_PORT: "80"
MYSQL_PORT: "3306"
PMA_PORT: "8081"

# Database Credentials (overridable for security)
MYSQL_ROOT_PASSWORD: "root"  # Default only
MYSQL_DATABASE: "devdb"
MYSQL_USER: "devuser"
MYSQL_PASSWORD: "devpass"

# Optional Services Toggle
ENABLE_REDIS: false
ENABLE_MAILHOG: false
```

**Rationale**: Centralized defaults reduce `.env.example` verbosity and make global updates easier.

**Acceptance Criteria**:
- All scripts load extended YAML variables
- README documents new variables
- Backward compatible with existing `.env` files

---

## 2. Automation & CI/CD

### 2.1 Automated Version Update Pipeline 🟡 📋
**Goal**: Keep PHP, MariaDB, Nginx versions current without manual intervention

**Implementation**:
- GitHub Actions workflow (scheduled weekly)
- Query Docker Hub API for latest stable versions
- Update `config/stack-vars.yml` if newer versions available
- Run test suite (automated stack start/stop)
- Create PR with version bump if tests pass

**Files**:
- `.github/workflows/version-check.yml`
- `.github/scripts/check-versions.sh`

**Acceptance Criteria**:
- Workflow runs weekly on Sunday 00:00 UTC
- PR created only if new version available and tests pass
- PR includes CHANGELOG entry and rationale

**Bonus**: Notify via GitHub Discussions or issue when major version available

---

### 2.2 Release Versioning & Workflow 🟡 📋
**Goal**: Formal release process with semantic versioning

**Features**:
- Automated CHANGELOG generation from conventional commits
- GitHub Release creation with notes
- Tag format: `vMAJOR.MINOR.PATCH` (e.g., `v2.0.0`)
- Release artifacts:
  - Standalone CLI binary (see 3.1)
  - Docker images for custom builds
  - Installation scripts

**Workflow**:
1. Maintainer triggers release workflow via GitHub Actions
2. Workflow validates CHANGELOG is updated
3. Creates git tag and GitHub Release
4. Builds and attaches artifacts
5. Updates `master` and `main` (if diverged)

**Files**:
- `.github/workflows/release.yml`
- `.github/scripts/generate-changelog.sh`

**Acceptance Criteria**:
- `git tag v2.0.0` → Triggers release workflow
- GitHub Release page shows changelog, artifacts
- All tests pass before release finalization

---

## 3. CLI & Distribution

### 3.1 Self-Contained CLI Binary 🔴 📋
**Goal**: Single-file executable for global or per-project installation

**Features**:
- Standalone binary with embedded stack files
- Minimal dependencies (Docker, Bash/Zsh)
- Commands:
  ```bash
  20i init [path]          # Initialize stack in directory
  20i start [--port 8080]  # Start stack
  20i stop                 # Stop stack
  20i status               # Show running stacks
  20i logs [service]       # Tail logs
  20i destroy              # Remove stack + volumes
  20i update               # Pull latest stack config
  20i version              # Show CLI + stack versions
  ```

**Installation Methods**:
```bash
# Global (Homebrew)
brew install peternicholls/tap/20i-stack

# Global (curl)
curl -fsSL https://raw.githubusercontent.com/peternicholls/20i-stack/master/install.sh | bash

# Per-project
curl -fsSL https://get.20i.dev/cli | bash -s -- --local
```

**Implementation Options**:
1. **Go binary** - Cross-platform, single file, embedded assets
2. **Bash bundle** - Shell script with base64-encoded compose files
3. **Python + PyInstaller** - Easier YAML manipulation, larger binary

**Recommended**: Go binary for portability and performance

**Files**:
- `cli/` directory with Go source
- `install.sh` - Installation script
- `.github/workflows/build-cli.yml` - Build automation
- `docs/CLI.md` - CLI usage guide

**Acceptance Criteria**:
- `20i init` creates `.20i-config.yml` and downloads compose files
- `20i start` works from any project directory
- Installer works on macOS (Intel + ARM) and Linux (x64 + ARM64)
- No Docker Compose knowledge required for basic usage

---

### 3.2 Docker Distribution (Pre-built Images) 🟢 📋
**Goal**: Publish pre-built PHP-FPM images to Docker Hub/GHCR

**Benefits**:
- Faster stack startup (no build step)
- Consistent environments across users
- Reduced local disk usage

**Images**:
```
peternicholls/20i-php-fpm:8.3
peternicholls/20i-php-fpm:8.4
peternicholls/20i-php-fpm:8.5
peternicholls/20i-php-fpm:latest → 8.5
```

**Update Compose**:
```yaml
apache:
  image: ${PHP_IMAGE:-peternicholls/20i-php-fpm:8.5}
  # Fallback to build if not available
  build:
    context: ./docker/apache
    args:
      PHP_VERSION: ${PHP_VERSION:-8.5}
```

**Acceptance Criteria**:
- Images published on every release
- Multi-arch support (AMD64 + ARM64)
- README documents pre-built vs custom build options

---

## 4. Features & Extensions

### 4.1 Optional Service Modules 🟢 📋
**Goal**: Easy addition of common development services

**Services**:
- **Redis** - Caching layer
  ```bash
  20i enable redis
  # Adds redis service to compose, exposes port 6379
  ```
- **Mailhog** - Email testing
  ```bash
  20i enable mailhog
  # SMTP on 1025, UI on 8025
  ```
- **Elasticsearch** - Search engine
- **RabbitMQ** - Message queue
- **MinIO** - S3-compatible storage

**Implementation**:
- Modular compose files: `docker/modules/redis.yml`
- CLI merges modules into main compose
- GUI shows checkboxes for optional services

**Acceptance Criteria**:
- `20i enable redis` → Redis container starts with stack
- `20i disable redis` → Redis removed cleanly
- Service config persists in `.20i-config.yml`

---

### 4.2 Project Templates 🟢 📋
**Goal**: Quick-start templates for common frameworks

**Templates**:
- Laravel (with queue worker, scheduler)
- WordPress (with WP-CLI)
- Symfony
- Plain PHP (existing default)

**Usage**:
```bash
20i init --template laravel
# Scaffolds composer.json, .env.example, artisan shortcuts
```

**Files**:
- `templates/laravel/`
- `templates/wordpress/`
- Template metadata in `templates.yml`

---

### 4.3 Health Checks & Auto-Restart 🟡 📋
**Goal**: Ensure services are truly ready before marking stack as "started"

**Features**:
- Dockerfile `HEALTHCHECK` directives
- Compose `healthcheck` definitions
- CLI waits for healthy status before exit
- Auto-restart unhealthy containers (configurable)

**Example**:
```yaml
apache:
  healthcheck:
    test: ["CMD", "php-fpm-healthcheck"]
    interval: 10s
    timeout: 3s
    retries: 3
```

**Acceptance Criteria**:
- `20i start` waits until all services healthy
- Unhealthy service shows clear error message
- Optional flag `--no-wait` skips health check

---

## 5. Developer Experience

### 5.1 Interactive Setup Wizard 🟢 📋
**Goal**: Guided setup for first-time users

**Flow**:
```bash
20i init
```
1. Detect project type (Laravel, WordPress, generic PHP)
2. Ask for PHP version preference
3. Ask for database name/credentials
4. Ask for web port (with conflict detection)
5. Ask for optional services (Redis, Mailhog)
6. Generate `.20i-config.yml` and `.env`
7. Start stack automatically

**Acceptance Criteria**:
- Non-technical users can complete setup without reading docs
- Wizard validates inputs (port conflicts, valid PHP version)
- Can skip wizard with `--skip-wizard` flag

---

### 5.2 Stack Profiles 🟢 📋
**Goal**: Save and switch between multiple configurations

**Usage**:
```bash
# Save current config as profile
20i profile save dev-php83

# Switch profiles
20i profile use dev-php85

# List profiles
20i profile list
```

**Storage**: `.20i-profiles/` directory with named YAML files

**Acceptance Criteria**:
- Switch between PHP 8.3 and 8.5 without editing files
- Profiles include service enablement (Redis on/off)
- Profile metadata shows last used date

---

### 5.3 Performance Metrics & Insights 🟢 📋
**Goal**: Show resource usage and performance tips

**Features**:
- `20i stats` shows CPU, memory, disk usage per service
- Startup time tracking
- Recommendations:
  - "MariaDB using 80% memory → Consider increasing limit"
  - "Stack took 45s to start → Enable pre-built images"

**Acceptance Criteria**:
- `20i stats` works without additional tools
- Insights are actionable and accurate

---

## 6. Security & Compliance

### 6.1 Secrets Management Integration 🟡 📋
**Goal**: Avoid plain-text credentials in `.env`

**Options**:
1. **Docker Secrets** (Swarm mode)
2. **External secret providers** (1Password CLI, AWS Secrets Manager)
3. **Encrypted .env** with `age` or `sops`

**Example**:
```yaml
# .20i-config.yml
secrets:
  provider: 1password
  vault: Development
  items:
    - MYSQL_ROOT_PASSWORD: op://Development/MySQL/password
```

**Acceptance Criteria**:
- No plain-text passwords in git-tracked files
- Seamless integration with existing workflows
- Fallback to `.env` if secrets unavailable

---

### 6.2 Dependency Scanning 🟡 📋
**Goal**: Detect vulnerable base images and dependencies

**Implementation**:
- GitHub Actions with Trivy or Grype
- Scan PHP-FPM Dockerfile and base images
- Fail PR if critical vulnerabilities found

**Acceptance Criteria**:
- Weekly scans on `master`
- PR comments show vulnerability summary
- Documented upgrade path for CVEs

---

## 7. Documentation & Community

### 7.1 Video Tutorials 🟢 📋
**Goal**: Visual guides for common tasks

**Topics**:
- Getting started (5 min)
- CLI installation and usage (10 min)
- Adding optional services (8 min)
- Troubleshooting guide (12 min)

**Platform**: YouTube, embedded in README

---

### 7.2 Community Templates Repository 🟢 📋
**Goal**: User-contributed project templates

**Structure**:
```
templates/
  community/
    drupal/
    magento/
    shopware/
```

**Contribution flow**: PR to `templates/community/` with template YAML and README

---

## 8. Future Explorations ⚪

### 8.1 GUI Desktop App
**Goal**: Native macOS/Windows app for stack management

**Tech**: Electron or Tauri for cross-platform

**Features**:
- Visual service toggles
- Live log viewer
- Port conflict resolution
- Built-in terminal

---

### 8.2 Kubernetes/Helm Charts
**Goal**: Production deployment option

**Use Case**: Transition from local dev to production K8s cluster

**Deliverable**: Helm chart mirroring local stack structure

---

### 8.3 Telemetry & Analytics (Opt-in)
**Goal**: Understand usage patterns to improve UX

**Data Collected**:
- OS and architecture
- PHP version popularity
- Command usage frequency
- Anonymous error reports

**Privacy**: Fully opt-in, open-source collection, no PII

---

## Priority Matrix

| Priority | Items |
|----------|-------|
| **Phase 1 (Q1 2026)** | 3.1 CLI Binary, 2.2 Release Workflow, 4.3 Health Checks |
| **Phase 2 (Q2 2026)** | 2.1 Auto Version Updates, 1.1 Smart Overrides, 4.1 Optional Services |
| **Phase 3 (Q3 2026)** | 3.2 Pre-built Images, 5.1 Setup Wizard, 6.2 Dependency Scanning |
| **Phase 4 (Q4 2026)** | 4.2 Project Templates, 5.2 Stack Profiles, 7.1 Video Tutorials |
| **Future** | 8.1 Desktop GUI, 8.2 K8s Charts, 8.3 Telemetry |

---

## Contributing to Roadmap

Have ideas? Open a GitHub Discussion or Issue with:
- **Problem**: What pain point does this solve?
- **Proposal**: How would it work?
- **Priority**: Why should this be prioritized?

Maintainers review quarterly and update this roadmap.

---

**Next Steps**: Review with community, create GitHub Projects for Phase 1 items, assign owners.
