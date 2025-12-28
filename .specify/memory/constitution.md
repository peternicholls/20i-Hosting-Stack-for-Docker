# 20i Stack Constitution

## Core Principles

### I. Environment-Driven Configuration (NON-NEGOTIABLE)
All configuration MUST be externalized via environment variables, `.env` files, or central YAML config (`config/stack-vars.yml`). Hard-coded credentials, paths, or environment-specific values are prohibited. This ensures portability across machines, contributors, and deployment environments.

**Rationale**: The project exists to be portable and reusable. Hard-coding defeats this core purpose and creates security risks.

### II. Multi-Platform First
The stack MUST support both Intel/AMD64 and ARM64 (Apple Silicon) architectures. Image selections and build processes MUST offer platform-specific optimization while maintaining cross-platform compatibility as the default. Users MUST be able to choose optimal variants through environment variables or GUI prompts.

**Rationale**: Developer machines span multiple architectures. The stack should perform well on all platforms without requiring arcane workarounds.

### III. Path Independence
All scripts and tooling MUST use `STACK_FILE` and `STACK_HOME` detection patterns. No absolute paths may be hard-coded. Project name sanitization MUST ensure Docker Compose compatibility (lowercase, alphanumeric with hyphens, leading letter/number).

**Rationale**: Users clone projects into arbitrary locations. Path-independent design allows the stack to run from any directory without modification.

### IV. Centralized Defaults with Override Hierarchy
Default values MUST be defined in `config/stack-vars.yml` when applicable. The override hierarchy is (highest to lowest priority):
1. Explicit environment variables
2. Per-project `.20i-local` file
3. Workspace `.env` file
4. Central YAML (`config/stack-vars.yml`)
5. Compose/Dockerfile fallback defaults

**Rationale**: Provides sensible defaults while allowing per-project and per-developer customization without conflict.

### V. User Experience & Feedback
CLI and GUI tools MUST provide clear, actionable feedback:
- Echo effective configuration before actions (preflight summaries)
- Display project names (both original and sanitized) during operations
- Use emoji/icons for visual clarity (🚀 starting, ✅ success, ❌ error, ⚙️ config)
- Confirm destructive operations with explicit prompts

**Rationale**: Transparent operations build user confidence and reduce debugging time. Clear feedback prevents mistakes and aids troubleshooting.

### VI. Documentation as First-Class Artifact
Every feature MUST be documented in at least two places:
1. Inline code comments explaining "why" (not just "what")
2. User-facing docs (README, GUI-HELP, or CHANGELOG)

Breaking changes MUST update CHANGELOG following [Keep a Changelog](https://keepachangelog.com/) format.

**Rationale**: The project is a developer tool. Good documentation reduces friction for new contributors and users.

### VII. Version Consistency
All version variables (PHP_VERSION, MYSQL_VERSION, etc.) MUST be kept in sync across:
- `config/stack-vars.yml`
- `.env.example`
- `docker-compose.yml` defaults
- Dockerfile ARG defaults
- README documentation

**Rationale**: Version drift causes confusion and inconsistent environments. Single source of truth prevents this.

## Architecture Constraints

### Stack Composition
The stack consists of four services in a fixed dependency order:
1. **MariaDB**: Database server (port 3306)
2. **Apache/PHP-FPM**: PHP processing engine (internal port 9000)
3. **Nginx**: Reverse proxy and web server (port 80/custom)
4. **phpMyAdmin**: Database management UI (port 8081/custom)

This composition is based on 20i's shared hosting environment and MUST NOT be fundamentally altered without strong justification and community consensus.

### Service Isolation
Each project creates isolated container instances using `COMPOSE_PROJECT_NAME`. Volume mounts use `CODE_DIR` to inject project code. The stack itself remains centralized; only the project code and configuration vary per-instance.

### Image Selection Standards
- Base images MUST use official Docker Hub images or well-maintained alternatives
- Alpine Linux variants MUST be preferred for smaller footprint
- Multi-arch images MUST be documented; ARM-native alternatives MUST be offered where performance matters (e.g., phpMyAdmin)

## Development Workflow

### Branch Strategy
- `master` / `main`: Stable, production-ready code
- Feature branches: Named descriptively (`add-redis-support`, `fix-port-detection`)
- PRs merge to default branch after review

### Testing Requirements
Before opening a PR, contributors MUST:
1. Run `scripts/setup-local.sh` in a fresh clone
2. Execute `docker compose up -d` and verify all services start
3. Test shell integration (`20i-up`, `20i-down`) from a sample project directory
4. If GUI changes: Test both CLI (`20i-gui`) and AppleScript workflows

### Code Review Standards
All PRs require:
- Clear description of changes and motivation
- Manual testing confirmation
- No hard-coded paths or credentials
- Updated documentation for user-facing changes
- CHANGELOG entry for notable changes

Automated reviews (Copilot, Gemini) provide suggestions but human review has final authority.

### Commit Hygiene
- Small, focused commits with clear messages
- Conventional commit format encouraged: `feat:`, `fix:`, `docs:`, `refactor:`
- Avoid mixing unrelated changes in a single commit

## Governance

### Constitution Authority
This constitution supersedes all other practices, guides, and documentation. In case of conflict, constitution principles prevail.

### Amendment Process
1. Propose amendment via GitHub issue with rationale
2. Community discussion period (minimum 7 days)
3. Maintainer approval required
4. Update version following semantic versioning:
   - **MAJOR**: Breaking changes to principles or workflow
   - **MINOR**: New principles or substantial additions
   - **PATCH**: Clarifications, typo fixes, minor refinements
5. Update `LAST_AMENDED_DATE` to amendment date

### Compliance
All pull requests MUST be verified for constitutional compliance during review. Maintainers may request changes to align with principles before merge.

### Exceptions
Exceptions to principles require:
- Documented justification in PR description
- Explicit approval from repository owner
- Technical or business rationale that outweighs principle adherence

**Version**: 1.0.0 | **Ratified**: 2025-12-28 | **Last Amended**: 2025-12-28
