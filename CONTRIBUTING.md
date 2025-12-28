Thanks for considering contributing! Please follow these guidelines:

- Fork the repository and create a feature branch.
- Make small, focused commits with clear messages.
- Run `docker compose up -d` and ensure the stack starts successfully before opening a PR.
- Open a pull request against the `main` branch and describe your change.

Development tips:
- Copy `.env.example` to `.env` and edit values locally (`cp .env.example .env`).
- To test with a different compose file, set `STACK_FILE` to the full path of the `docker-compose.yml` you want to use.

Thanks — maintainers will review PRs and request changes as needed.
