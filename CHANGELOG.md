# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Automated release workflow with release-please
- Release validation scripts (version.sh, validate.sh)
- Release artifact packaging and distribution

## [1.0.0] - 2025-12-28

Initial release with generic project setup.

### Added
- `LICENSE` (MIT), `CONTRIBUTING.md` and `scripts/setup-local.sh`
- Project name sanitization for Docker Compose compliance
- Destroy stack option with confirmation prompt (`20i-destroy`)
- Multi-platform support (Intel/AMD and Apple Silicon)
- Interactive GUI menu system (`20i-gui`)
- Shell integration with convenient aliases

### Changed
- Use `STACK_FILE` / `STACK_HOME` detection in helper scripts instead of hard-coded paths
- Read secrets from `.env` / environment variables (no hard-coded passwords)
- Make macOS Automator / workflow scripts honor `STACK_FILE` and be overrideable
- Container detection regex updated to support project names with hyphens

[Unreleased]: https://github.com/peternicholls/20i-stack/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/peternicholls/20i-stack/releases/tag/v1.0.0
