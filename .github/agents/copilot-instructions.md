# 20i-stack Development Guidelines

Auto-generated from feature plans. Last updated: 2025-12-28

## Active Technologies
- Go 1.21+ + Bubble Tea v1.3.10+, Bubbles v1.0.0+, Lipgloss v1.0.0+, Docker SDK v27.0.0+ (001-stack-manager-tui)
- N/A (reads docker-compose.yml and .20i-local; no persistent state) (001-stack-manager-tui)

- **GitHub Actions** (004-release-workflow) - CI/CD automation
- **release-please** (004-release-workflow) - Automated versioning and changelog
- **Bash 5.x** (004-release-workflow) - Release scripts
- **Docker Compose** - Stack orchestration
- **PHP 8.5 / FPM** - Runtime environment
- **Nginx / MariaDB / phpMyAdmin** - Stack services

## Project Structure

```text
.github/
├── workflows/           # GitHub Actions workflows
│   ├── release.yml      # Main release workflow (release-please)
│   ├── validate-pr.yml  # PR validation (conventional commits)
│   └── changelog-preview.yml
├── PULL_REQUEST_TEMPLATE.md
└── agents/

config/
├── stack-vars.yml       # Central configuration
└── release/
    └── config.yml       # Release configuration

scripts/
├── setup-local.sh       # Local setup
└── release/
    ├── validate.sh      # Pre-release validation
    ├── changelog.sh     # Changelog generation
    ├── version.sh       # Version management
    └── artifacts.sh     # Build release artifacts

docker/                  # Docker configurations
specs/                   # Feature specifications
src/                     # Application source

VERSION                  # Single source of version truth
CHANGELOG.md             # Release history (Keep a Changelog format)
```

## Commands

```bash
# Version management
scripts/release/version.sh get         # Get current version
scripts/release/version.sh bump minor  # Bump minor version

# Release validation
scripts/release/validate.sh --all      # Run all validations

# Local development
./20i-gui                              # Interactive CLI menu
```

## Code Style

- **Shell Scripts**: ShellCheck compliant, POSIX-compatible where possible
- **Commits**: Conventional Commits format (feat:, fix:, docs:, etc.)
- **YAML**: 2-space indentation, quoted strings for versions

## Recent Changes
- 001-stack-manager-tui: Added Go 1.21+ + Bubble Tea v1.3.10+, Bubbles v1.0.0+, Lipgloss v1.0.0+, Docker SDK v27.0.0+

- 004-release-workflow: Added automated release versioning with release-please
  - GitHub Actions workflows for CI/CD
  - VERSION file as single source of truth
  - Auto-generated CHANGELOG from conventional commits
  - Release artifacts with SHA256 checksums

<!-- MANUAL ADDITIONS START -->
<!-- Add project-specific instructions below this line -->
<!-- MANUAL ADDITIONS END -->
