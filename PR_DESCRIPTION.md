This PR updates the project to be friendlier for people cloning it. Key changes:

- Use `STACK_FILE` / `STACK_HOME` detection so helper scripts no longer rely on hard-coded absolute paths
- Move secrets to environment variables (read from `.env`) and remove hard-coded DB passwords
- Make macOS Automator/workflow scripts optional and overrideable via `STACK_FILE` (still backward compatible)
- Add `LICENSE` (MIT), `CONTRIBUTING.md`, `scripts/setup-local.sh`, and a `CHANGELOG`

How to test:
1. Clone the repo and run `scripts/setup-local.sh`
2. `docker compose up -d` should start the stack with env values from `.env`
3. Source `zsh-example-script.zsh` and call `20i-up` from a project folder

If you're happy with this, please merge into `main`.