This pull request makes the project more friendly for people who clone it:

- Use `STACK_FILE` / `STACK_HOME` detection so helpers don't depend on absolute paths
- Read DB/passwords from `.env` or environment variables (no hard-coded secrets)
- Make macOS Automator/workflow scripts optional and overrideable via `STACK_FILE`
- Add `LICENSE` (MIT), `CONTRIBUTING.md`, and `scripts/setup-local.sh`
