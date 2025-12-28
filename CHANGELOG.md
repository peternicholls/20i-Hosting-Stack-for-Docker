# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2025-12-28)


### Features

* add automated release workflow with versioning and artifacts ([876a3c5](https://github.com/peternicholls/20i-Hosting-Stack-for-Docker/commit/876a3c5c80786ea42597ef5f8d52e55ec00341bc))


### Bug Fixes

* add id-token permission to release workflow ([23f1325](https://github.com/peternicholls/20i-Hosting-Stack-for-Docker/commit/23f1325046a0afa28f1cb6e8c77a6469b6482a3a))
* update workflows to run on master branch instead of main ([dd83c16](https://github.com/peternicholls/20i-Hosting-Stack-for-Docker/commit/dd83c1671fa9f9f76175091c80e833c04ddba6a4))

## [Unreleased]

### Added - Phase 3a TUI Features
- **Terminal UI (TUI)** - Modern keyboard-driven interface for stack management
- **Container Management** - Start/stop/restart individual containers via keyboard
- **Dashboard View** - Two-panel layout with service list and status information
- **Help Modal** - Press `?` to see all available commands
- **Mouse Support** - Full mouse interaction support (click, scroll)
- **Keyboard Navigation** - Vim-style (`j/k`) and arrow key navigation
- **Error Handling** - User-friendly Docker error messages with recovery hints
- **Documentation** - Complete user guide, troubleshooting, and architecture docs
- **Integration Tests** - CI-safe test suite with Docker mocking
- Automated release workflow with release-please
- Release validation scripts (version.sh, validate.sh)
- Release artifact packaging and distribution

### Documentation
- `/docs/tui/user-guide.md` - Complete TUI usage guide
- `/docs/tui/troubleshooting.md` - Common issues and solutions
- `/docs/tui/architecture.md` - Technical design and patterns
- Updated `tui/README.md` with Phase 3a features
- Updated main `README.md` with TUI quick start

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
