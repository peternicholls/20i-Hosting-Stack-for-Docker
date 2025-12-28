# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

Initial changes to make the project more generic and safe for cloning by others.

### Added
- `LICENSE` (MIT), `CONTRIBUTING.md` and `scripts/setup-local.sh`
- Project name sanitization for Docker Compose compliance
- Destroy stack option with confirmation prompt (`20i-destroy`)

### Changed
- Use `STACK_FILE` / `STACK_HOME` detection in helper scripts instead of hard-coded paths
- Read secrets from `.env` / environment variables (no hard-coded passwords)
- Make macOS Automator / workflow scripts honor `STACK_FILE` and be overrideable
- Container detection regex updated to support project names with hyphens
