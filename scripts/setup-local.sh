#!/usr/bin/env bash
set -euo pipefail

# Small helper to bootstrap a local development copy
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

if [[ ! -f .env ]]; then
  cp .env.example .env
  echo "Created .env from .env.example — edit .env if you need to change secrets or ports."
else
  echo ".env already exists — leaving it in place."
fi

echo "Next steps:"
echo "  1) Edit .env if you want to override defaults (e.g., MYSQL_PASSWORD)."
echo "  2) (Optional) set COMPOSE_PROJECT_NAME by running: export COMPOSE_PROJECT_NAME=\$(basename \$(pwd))"
echo "  3) Start the stack: docker compose up -d"
echo ""
echo "If you intend to use the macOS Automator workflows, they will honor the STACK_FILE env var, or fall back to:"
echo "  \$HOME/docker/20i-stack/docker-compose.yml"
