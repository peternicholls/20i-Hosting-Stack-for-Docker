# ─────────────────────────────────────────────────────────────────────────────
# Docker shortcuts for 20i-style local stack
# ─────────────────────────────────────────────────────────────────────────────

# 20i stack configuration
# Prefer explicit STACK_FILE (full path to docker-compose.yml). If not set, derive from STACK_HOME or default.
if [[ -n "${STACK_FILE:-}" ]]; then
    # STACK_FILE is explicitly set, derive STACK_HOME from it
    STACK_HOME="${STACK_HOME:-$(cd "$(dirname "$STACK_FILE")" 2>/dev/null && pwd)}"
else
    # Default: use STACK_HOME or fall back to $HOME/docker/20i-stack
    STACK_HOME="${STACK_HOME:-$HOME/docker/20i-stack}"
    STACK_FILE="$STACK_HOME/docker-compose.yml"
fi

# Function to start 20i stack
20i-up() {
    local PROJECT_DIR="$(pwd)"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    # Check if stack directory exists
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Set project name based on current directory
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$(basename "$PROJECT_DIR")}"
    export CODE_DIR="$PROJECT_DIR"
    
    # Source optional per-project overrides
    [[ -f .20i-local ]] && source .20i-local
    
    echo "🚀 Starting 20i stack for project: $COMPOSE_PROJECT_NAME"
    echo "📁 Code directory: $CODE_DIR"
    
    docker compose -f "$STACK_FILE" up -d "$@"
}

# Function to stop 20i stack
20i-down() {
    local PROJECT_DIR="$(pwd)"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$(basename "$PROJECT_DIR")}"
    
    echo "🛑 Stopping 20i stack..."
    docker compose -f "$STACK_FILE" down "$@"
}

# Function to show 20i stack status
20i-status() {
    local PROJECT_DIR="$(pwd)"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$(basename "$PROJECT_DIR")}"
    
    echo "📊 20i stack status:"
    docker compose -f "$STACK_FILE" ps
}

# Function to view 20i stack logs
20i-logs() {
    local PROJECT_DIR="$(pwd)"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$(basename "$PROJECT_DIR")}"
    
    docker compose -f "$STACK_FILE" logs -f "$@"
}

# Aliases for convenience
alias 20i='20i-status'
alias dcu='20i-up'
alias dcd='20i-down'
# GUI script shortcut - interactive menu for 20i stack management
20i-gui() {
    "$STACK_HOME/20i-gui" "$@"
}
