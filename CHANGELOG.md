# Changelog

This branch makes the project more generic and safe for cloning by others.

- Use `STACK_FILE` / `STACK_HOME` detection in helper scripts
- Read secrets from `.env` / environment variables (no hard-coded passwords)
- Make macOS Automator / workflow scripts honor `STACK_FILE` and be overrideable
- Add `LICENSE` (MIT), `CONTRIBUTING.md` and `scripts/setup-local.sh`
