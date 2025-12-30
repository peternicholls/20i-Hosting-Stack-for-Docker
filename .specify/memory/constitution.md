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

### Default Branch Protection (NON-NEGOTIABLE)
The default branch (`main`/`master`) MUST be protected:
- No direct pushes (PRs only)
- At least 1 human approval required
- Required status checks MUST pass (lint, compose validation, smoke test)
- Branch must be up to date before merge

**Rationale**: This constitution is only enforceable if the repository prevents bypassing review and checks.

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

### Commit Hygiene
- Small, focused commits with clear messages
- Commits SHOULD be atomic and made often: each commit represents one logical change that can be understood, reviewed, and reverted independently
- Conventional commits MUST be used for merged PR titles (and SHOULD be used for individual commits): `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`
- Avoid mixing unrelated changes in a single commit
- Commit messages MUST describe intent and scope (avoid ambiguous messages like `update`)

**What “atomic” means in practice (guidance)**:
Examples of acceptable atomic commits include:
- Documentation-only changes (README, CHANGELOG, inline comments)
- Introducing or modifying a single environment variable and its references
- Adding or adjusting configuration for one service
- A single bug fix with no unrelated refactoring
- A focused refactor that does not change behaviour

Non-atomic examples (to avoid):
- Feature changes mixed with formatting or lint cleanups
- Multiple services changed for unrelated reasons
- Refactors combined with functional changes without clear separation

This guidance exists to support review quality, rollback safety, and effective use of CI and `git bisect`.

### Repository Hygiene (REQUIRED)
The repository MUST include and maintain:
- `.gitignore` covering Docker artifacts, logs, local overrides, and common OS/IDE files
- `.gitattributes` enforcing consistent text handling (recommended: `text=auto eol=lf`)
- `.editorconfig` for consistent whitespace and indentation
- `LICENSE`
- `CODEOWNERS` (even if it maps to a single maintainer)

**Rationale**: Consistency and portability reduce friction for contributors and prevent accidental commits of local state.

### Secrets & Local State (NON-NEGOTIABLE)
- Secrets MUST NOT be committed under any circumstances.
- `.env`, `.20i-local`, and any local secret files MUST be ignored by Git.
- Only template configuration MAY be committed (e.g., `.env.example`, documented YAML defaults).
- CI SHOULD include a basic secret scan and MUST fail if obvious credentials are detected.

**Rationale**: The stack is designed for portability and reuse. Leaking credentials or committing local state undermines security and reproducibility.

### Pull Request Standards (REQUIRED)
All PRs MUST:
- Have a single clear purpose (no drive-by refactors mixed with feature work)
- Include a clear description of changes and motivation
- Confirm manual testing steps were performed (see Testing Requirements)
- Update documentation for user-facing changes
- Include a CHANGELOG entry for notable changes

Automated reviews (Copilot, Gemini) may provide suggestions but human review has final authority.

### Releases, Tags, and Versioning
- Every constitution amendment and notable user-facing change SHOULD be tagged with `vX.Y.Z`.
- Release notes SHOULD be derived from `CHANGELOG` entries.
- Breaking changes MUST include migration notes.

**Rationale**: Tags provide stable reference points for users and make rollback and diagnosis practical.

### CI Enforcement of Constitutional Invariants
CI MUST, at minimum:
- Validate Docker Compose configuration (`docker compose config`)
- Run a minimal smoke test that brings the stack up and verifies core services start
- Fail on prohibited hard-coded paths or credentials patterns where feasible

**Rationale**: Manual review catches many issues, but CI is the consistent backstop.

### Submodules and Vendored Code
- Git submodules SHOULD be avoided unless there is strong justification.
- Vendored third-party code MUST retain license headers and attribution.

**Rationale**: Submodules and untracked licensing obligations are common sources of long-term maintenance and legal risk.

### Dependency Update Discipline
- Base images and key dependency versions SHOULD be reviewed on a regular cadence (recommended: monthly).
- Security updates may be expedited, but MUST be documented and included in CHANGELOG.

**Rationale**: This stack is infrastructure. Silent drift or stale images are reliability and security liabilities.

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
