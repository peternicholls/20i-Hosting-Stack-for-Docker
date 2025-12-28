# Implementation Plan: Release Versioning and Workflow

**Branch**: `004-release-workflow` | **Date**: 2025-12-28 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-release-workflow/spec.md`

## Summary

Implement a fully automated release workflow using GitHub Actions that creates semantic versioned releases with auto-generated changelogs, release artifacts, and comprehensive validation gates. The system will enforce conventional commits, prevent invalid releases, and provide maintainers with a streamlined release process via GitHub Actions UI.

## Technical Context

**Language/Version**: YAML (GitHub Actions), Bash 5.x  
**Primary Dependencies**: GitHub Actions, conventional-changelog-cli, gh CLI, git  
**Storage**: N/A (Git tags, GitHub Releases as storage)  
**Testing**: GitHub Actions workflow validation, ShellCheck for scripts  
**Target Platform**: GitHub-hosted runners (ubuntu-latest)  
**Project Type**: Single project - DevOps/Infrastructure automation  
**Performance Goals**: Release completion in <10 minutes (per SC-001)  
**Constraints**: GitHub Actions rate limits, artifact size limits (2GB per file), 6-hour workflow timeout  
**Scale/Scope**: Single repository, maintainer-triggered releases, ~4 artifacts per release

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Compliance Notes |
|-----------|--------|------------------|
| I. Environment-Driven Configuration | ✅ PASS | Version in `VERSION` file, workflow inputs for overrides, no hard-coded values |
| II. Multi-Platform First | ✅ PASS | Artifacts support AMD64/ARM64, workflow runs on ubuntu-latest (standard) |
| III. Path Independence | ✅ PASS | All scripts use `$GITHUB_WORKSPACE`, relative paths within repo |
| IV. Centralized Defaults | ✅ PASS | VERSION file as single source, workflow inputs for release-time overrides |
| V. User Experience & Feedback | ✅ PASS | Preflight summary, emoji status indicators, clear error messages |
| VI. Documentation as First-Class | ✅ PASS | Auto-generated CHANGELOG, README badges, release notes |
| VII. Version Consistency | ✅ PASS | VERSION file synced to git tags, CHANGELOG, docker-compose labels |
| Commit Hygiene (Dev Workflow) | ✅ PASS | Conventional commits enforced via PR validation workflow |

**Constitution Gate**: ✅ PASSED - All principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/004-release-workflow/
├── plan.md              # This file
├── research.md          # Phase 0: Tooling decisions and best practices
├── data-model.md        # Phase 1: VERSION format, config schema
├── quickstart.md        # Phase 1: Maintainer release guide
├── contracts/           # Phase 1: Workflow API specifications
│   ├── release-workflow.yml    # Main release workflow interface
│   ├── validate-pr.yml         # PR validation interface
│   └── events.md               # Workflow dispatch events
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
.github/
├── workflows/
│   ├── release.yml              # Main release workflow (manual trigger)
│   ├── validate-pr.yml          # PR validation (commit format, tests)
│   ├── changelog-preview.yml    # Generate changelog preview on PR
│   └── auto-release.yml         # Optional: auto-release on main push
├── PULL_REQUEST_TEMPLATE.md     # PR template with conventional commit guide
├── ISSUE_TEMPLATE/
│   └── release-request.yml      # Release request issue template
└── release.yml                  # Release drafter config (optional)

config/
├── stack-vars.yml               # Existing - unchanged
└── release/
    ├── config.yml               # Release categories, labels, templates
    └── changelog-template.hbs   # Handlebars template for CHANGELOG

scripts/
├── setup-local.sh               # Existing - unchanged
└── release/
    ├── validate.sh              # Pre-release validation checks
    ├── changelog.sh             # Generate/update CHANGELOG
    ├── version.sh               # Version bumping (major/minor/patch)
    ├── artifacts.sh             # Package release artifacts
    └── publish.sh               # Publish to GitHub Releases

VERSION                          # Single source of truth (e.g., "1.0.0")
CHANGELOG.md                     # Existing - auto-updated by workflow
README.md                        # Existing - add version badge
```

**Structure Decision**: Single project with `.github/workflows/` for CI/CD automation and `scripts/release/` for reusable release utilities. This follows the existing `scripts/` convention and keeps release logic testable outside GitHub Actions.

## Complexity Tracking

> No constitution violations - section not required.
