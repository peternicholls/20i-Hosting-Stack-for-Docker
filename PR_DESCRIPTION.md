This PR makes the stack generic and environment-driven with multi-platform support and central configuration.

Key changes
- Multi-platform: Intel/AMD and Apple Silicon (ARM64) supported.
- phpMyAdmin choice: GUI prompts for ARM-native (`arm64v8/phpmyadmin:latest`) or cross-platform (`phpmyadmin/phpmyadmin:latest`).
- Central variables: New `config/stack-vars.yml` with `PHP_VERSION`, `MYSQL_VERSION`, and `PMA_IMAGE`. Scripts read these defaults unless overridden by `.env` / `.20i-local` / environment.
- PHP 8.5: Standardized across Dockerfile, Compose, `.env.example`, and README.
- CLI preflight: `20i-gui` and `zsh-example-script.zsh` echo effective configuration before starting (helps debugging and review).
- Path handling: `STACK_FILE`/`STACK_HOME` detection in helpers; Automator workflows honor `STACK_FILE`.
- Docs: README opening updated for multi-platform, Topics added; GUI help documents image selection.

Files touched
- `docker/apache/Dockerfile`: default `ARG PHP_VERSION=8.5`.
- `docker-compose.yml`: `PHP_VERSION: ${PHP_VERSION:-8.5}`; phpMyAdmin image uses `PMA_IMAGE`.
- `.env.example`: updated to PHP 8.5 and documented phpMyAdmin image override.
- `config/stack-vars.yml`: centralized defaults (`PHP_VERSION`, `MYSQL_VERSION`, `PMA_IMAGE`).
- `20i-gui`, `zsh-example-script.zsh`, `20i-stack-manager.scpt`: load central YAML; show effective configuration; GUI image selection.
- `README.md`, `GUI-HELP.md`: multi-platform guidance, Topics, and GUI options.

Testing
1. `scripts/setup-local.sh` to generate `.env`.
2. `docker compose build apache && docker compose up -d` → stack starts using defaults / overrides.
3. `20i-gui` or `20i-up` → confirm preflight shows effective `PHP_VERSION`, `MYSQL_VERSION`, `PMA_IMAGE`, and ports.

Notes
- Precedence: `.20i-local` / `.env` / env vars override central YAML.
- Defaults remain sane if YAML is missing (Compose/Dockerfile provide fallbacks).

Topics
docker, docker-compose, php, php-fpm, nginx, mariadb, phpmyadmin, apple-silicon, arm64, cross-platform, development-environment, macos