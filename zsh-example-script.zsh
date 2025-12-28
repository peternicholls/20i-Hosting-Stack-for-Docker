# ─────────────────────────────────────────────────────────────────────────────
# Docker shortcuts for 20i-style local stack
# ─────────────────────────────────────────────────────────────────────────────

# 20i stack configuration
# Prefer explicit STACK_FILE (full path to docker-compose.yml). If not set, derive from STACK_HOME or default.
STACK_FILE="${STACK_FILE:-${STACK_HOME:-$HOME/docker/20i-stack}/docker-compose.yml}"
STACK_HOME="${STACK_HOME:-$(cd "$(dirname "$STACK_FILE")" >/dev/null 2>&1 && pwd)}"

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
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    echo "🛑 Stopping 20i stack..."
    docker compose -f "$STACK_FILE" down "$@"
}

# Function to show 20i stack status
20i-status() {
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    echo "📊 20i stack status:"
    docker compose -f "$STACK_FILE" ps
}

# Function to view 20i stack logs
20i-logs() {
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
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
