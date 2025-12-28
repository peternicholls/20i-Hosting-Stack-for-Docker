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
## Copilot Coding Agent Rules (High Priority)

These rules are written for GitHub Copilot Coding Agent. Follow them exactly.

### Scope and safety
- Treat the GitHub Issue description as the task contract.
- Only modify files explicitly listed in the Issue under “Files to create/modify”. If the Issue does not list files, stop and ask for clarification in the PR description.
- Do not implement out-of-scope features, “future phases”, refactors, or drive-by cleanups.
- Keep PRs small and reviewable: one work packet per PR.

### TUI development rules (001-stack-manager-tui)
- Work inside `tui/` for all TUI tasks unless the Issue explicitly includes non-TUI paths.
- Bubble Tea Elm Architecture:
  - Never block in `Update()`.
  - All I/O must be performed via `tea.Cmd`.
  - UI rendering must be a pure `View()` function.
- Styling:
  - Use Lipgloss for all styling.
  - Do not use raw ANSI codes.
- Paths and configuration:
  - Never hard-code paths to `docker-compose.yml`.
  - Always use environment-driven configuration (`STACK_FILE`, `STACK_HOME`) with executable-relative fallback detection.
  - Validate the stack file exists before any compose operation.
- Platform:
  - Use `runtime.GOARCH` for architecture decisions.
  - Use ARM64-native phpMyAdmin image on ARM, default image on x86, with env override.

### Quality gates
- Add or update tests for all non-trivial logic (table-driven tests preferred).
- Run these commands before marking a PR ready:
  - From `tui/`: `go test ./...`
  - If a Makefile target exists for tests, run it too.
- Ensure all exported Go symbols you add or modify have godoc comments.
- Prefer clear, user-friendly error messages (avoid surfacing raw Docker errors directly to the UI).

### PR completion checklist
Include a short checklist in the PR description:
- [ ] Tests: `go test ./...` (from `tui/`) passed
- [ ] Scope limited to the Issue
- [ ] Godoc added for exported symbols
- [ ] User-visible errors are friendly
<!-- MANUAL ADDITIONS END -->
