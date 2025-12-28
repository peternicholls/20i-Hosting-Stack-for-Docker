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

# Helper to ensure stack file exists
_20i_check_stack_file() {
    if [[ ! -f "$STACK_FILE" ]]; then
        echo "❌ Error: Docker compose file not found at $STACK_FILE"
        return 1
    fi
}

# Helper to set up project environment variables
_20i_setup_project_env() {
    local PROJECT_DIR="$(pwd)"
    local PROJECT_NAME="$(basename "$PROJECT_DIR")"
    local SAFE_PROJECT_NAME="$(sanitize_project_name "$PROJECT_NAME")"
    
    # Export environment variables
    export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$SAFE_PROJECT_NAME}"
    export CODE_DIR="${CODE_DIR:-$PROJECT_DIR}"
    
    # Store values for functions to access
    _20I_PROJECT_NAME="$PROJECT_NAME"
    _20I_SAFE_PROJECT_NAME="$SAFE_PROJECT_NAME"
}

# Function to start 20i stack
20i-up() {
    _20i_check_stack_file || return 1
    _20i_setup_project_env
    
    # Source optional per-project overrides
    [[ -f .20i-local ]] && source .20i-local
    
    echo "🚀 Starting 20i stack for project: $_20I_PROJECT_NAME"
    if [[ "$_20I_SAFE_PROJECT_NAME" != "$_20I_PROJECT_NAME" ]]; then
        echo "📛 Normalized project name: $_20I_SAFE_PROJECT_NAME"
    fi
    echo "📁 Code directory: $CODE_DIR"
    
    docker compose -f "$STACK_FILE" up -d "$@"
}

# Function to stop 20i stack
20i-down() {
    _20i_check_stack_file || return 1
    _20i_setup_project_env
    
    echo "🛑 Stopping 20i stack..."
    docker compose -f "$STACK_FILE" down "$@"
}

# Function to show 20i stack status
20i-status() {
    _20i_check_stack_file || return 1
    _20i_setup_project_env
    
    echo "📊 20i stack status:"
    docker compose -f "$STACK_FILE" ps
}

# Function to view 20i stack logs
20i-logs() {
    _20i_check_stack_file || return 1
    _20i_setup_project_env
    
    docker compose -f "$STACK_FILE" logs -f "$@"
}

# Function to destroy 20i stack (stop and remove volumes)
20i-destroy() {
    _20i_check_stack_file || return 1
    _20i_setup_project_env
    
    echo "⚠️  WARNING: This will destroy the stack for project: $_20I_SAFE_PROJECT_NAME"
    echo "    - Stop all containers"
    echo "    - Remove all volumes (database data will be lost!)"
    echo "    - Remove networks"
    echo ""
    read -p "Are you sure? (type 'yes' to confirm): " confirmation
    
    if [[ "$confirmation" == "yes" ]]; then
        echo "💥 Destroying 20i stack..."
        docker compose -f "$STACK_FILE" down -v
        echo "✅ Stack destroyed: $_20I_SAFE_PROJECT_NAME"
    else
        echo "❌ Destroy cancelled"
        return 1
    fi
}

# Aliases for convenience
alias 20i='20i-status'
alias dcu='20i-up'
alias dcd='20i-down'

# GUI script shortcut - interactive menu for 20i stack management
function 20i-gui() {
    "$STACK_HOME/20i-gui" "$@"
}
