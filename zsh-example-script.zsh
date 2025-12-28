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

# Function to sanitize project names for Docker Compose
sanitize_project_name() {
    local name="$1"
    # Lowercase
    name="$(echo "$name" | tr '[:upper:]' '[:lower:]')"
    # Replace any sequence of invalid chars with single hyphen
    name="$(echo "$name" | sed -E 's/[^a-z0-9]+/-/g')"
    # Collapse consecutive hyphens/underscores into single hyphen
    name="$(echo "$name" | sed -E 's/[-_]+/-/g')"
    # Trim leading/trailing hyphens using parameter expansion
    name="${name##-}"
    name="${name%%-}"
    # Ensure it starts with a letter or number
    if [[ ! "$name" =~ ^[a-z0-9] ]]; then
        name="p${name}"
    fi
    # Fallback
    if [[ -z "$name" ]]; then
        name="project"
    fi
    echo "$name"
}

# Function to start 20i stack
20i-up() {
    local PROJECT_DIR="$(pwd)"
    local PROJECT_NAME="$(basename "$PROJECT_DIR")"
    local SAFE_PROJECT_NAME="$(sanitize_project_name "$PROJECT_NAME")"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    # Check if stack directory exists
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    echo "🚀 Starting 20i stack for project: $PROJECT_NAME"
    if [[ "$SAFE_PROJECT_NAME" != "$PROJECT_NAME" ]]; then
        echo "📛 Normalized project name: $SAFE_PROJECT_NAME"
    fi
    echo "📁 Code directory: $CODE_DIR"
    
    # Set project name based on current directory
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$SAFE_PROJECT_NAME}"
    export CODE_DIR="$PROJECT_DIR"
    
    # Source optional per-project overrides
    [[ -f .20i-local ]] && source .20i-local
    
    docker compose -f "$STACK_FILE" up -d "$@"
}

# Function to stop 20i stack
20i-down() {
    local PROJECT_DIR="$(pwd)"
    local PROJECT_NAME="$(basename "$PROJECT_DIR")"
    local SAFE_PROJECT_NAME="$(sanitize_project_name "$PROJECT_NAME")"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$SAFE_PROJECT_NAME}"
    
    echo "🛑 Stopping 20i stack..."
    docker compose -f "$STACK_FILE" down "$@"
}

# Function to show 20i stack status
20i-status() {
    local PROJECT_DIR="$(pwd)"
    local PROJECT_NAME="$(basename "$PROJECT_DIR")"
    local SAFE_PROJECT_NAME="$(sanitize_project_name "$PROJECT_NAME")"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$SAFE_PROJECT_NAME}"
    
    echo "📊 20i stack status:"
    docker compose -f "$STACK_FILE" ps
}

# Function to view 20i stack logs
20i-logs() {
    local PROJECT_DIR="$(pwd)"
    local PROJECT_NAME="$(basename "$PROJECT_DIR")"
    local SAFE_PROJECT_NAME="$(sanitize_project_name "$PROJECT_NAME")"
    local STACK_FILE="${STACK_FILE:-$STACK_HOME/docker-compose.yml}"
    
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
    
    # Export CODE_DIR to satisfy docker-compose.yml requirements
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$SAFE_PROJECT_NAME}"
    
    docker compose -f "$STACK_FILE" logs -f "$@"
}

# Aliases for convenience
alias 20i='20i-status'
alias dcu='20i-up'
alias dcd='20i-down'

# GUI script shortcut - interactive menu for 20i stack management
function 20i-gui() {
    "$STACK_HOME/20i-gui" "$@"
}
