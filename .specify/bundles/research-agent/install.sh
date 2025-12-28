#!/usr/bin/env bash
# SpecKit Research Agent Installer
# Version: 1.0.0
# Installs the research agent into a SpecKit project

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory (where this bundle is)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUNDLE_DIR="$SCRIPT_DIR"

# Detect project root (look for .git or .specify)
find_project_root() {
    local current="$PWD"
    while [[ "$current" != "/" ]]; do
        if [[ -d "$current/.git" ]] || [[ -d "$current/.specify" ]]; then
            echo "$current"
            return 0
        fi
        current="$(dirname "$current")"
    done
    echo "$PWD"  # Default to current directory
}

PROJECT_ROOT="${PROJECT_ROOT:-$(find_project_root)}"

# Configuration
INSTALL_MODE="${INSTALL_MODE:-interactive}"  # interactive | auto
SKIP_CONFIG="${SKIP_CONFIG:-false}"
PROJECT_TYPE="${PROJECT_TYPE:-}"

# Print functions
print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}  ${GREEN}SpecKit Research Agent Installer v1.0.0${NC}               ${BLUE}║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_step() {
    echo -e "${BLUE}▸${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo -e "  $1"
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."
    
    local missing=()
    
    # Check for required commands
    command -v git >/dev/null 2>&1 || missing+=("git")
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        print_error "Missing required tools: ${missing[*]}"
        print_info "Please install missing tools and try again"
        exit 1
    fi
    
    # Warn about optional tools
    if ! command -v yq >/dev/null 2>&1; then
        print_warning "yq not found (optional but recommended)"
        print_info "Install with: brew install yq (macOS) or apt install yq (Linux)"
    fi
    
    print_success "Prerequisites OK"
    echo ""
}

# Create directory structure
create_directories() {
    print_step "Creating directory structure..."
    
    cd "$PROJECT_ROOT"
    
    mkdir -p .github/agents
    mkdir -p .specify/config
    mkdir -p .specify/context
    mkdir -p .specify/docs
    mkdir -p .specify/examples/research-configs
    mkdir -p .specify/scripts/bash
    mkdir -p .specify/templates
    
    print_success "Directories created"
    echo ""
}

# Copy files
copy_files() {
    print_step "Installing files..."
    
    cd "$PROJECT_ROOT"
    
    # Agent definition
    if [[ -f "$BUNDLE_DIR/files/agents/speckit.research.agent.md" ]]; then
        cp "$BUNDLE_DIR/files/agents/speckit.research.agent.md" .github/agents/
        print_success "Agent definition → .github/agents/"
    elif [[ -f "$BUNDLE_DIR/../../agents/speckit.research.agent.md" ]]; then
        # Handle case where we're running from within SafeDownload
        cp "$BUNDLE_DIR/../../agents/speckit.research.agent.md" .github/agents/
        print_success "Agent definition → .github/agents/"
    else
        # Try to find it in parent directories
        if [[ -f "../../../.github/agents/speckit.research.agent.md" ]]; then
            cp "../../../.github/agents/speckit.research.agent.md" .github/agents/
            print_success "Agent definition → .github/agents/"
        else
            print_error "Agent definition not found"
            return 1
        fi
    fi
    
    # Templates
    cp "$BUNDLE_DIR"/files/templates/*.{yaml,md} .specify/templates/ 2>/dev/null || true
    print_success "Templates → .specify/templates/"
    
    # Scripts
    cp "$BUNDLE_DIR"/files/scripts/*.sh .specify/scripts/bash/ 2>/dev/null || true
    chmod +x .specify/scripts/bash/*.sh
    print_success "Scripts → .specify/scripts/bash/"
    
    # Documentation
    cp "$BUNDLE_DIR"/files/docs/*.md .specify/docs/ 2>/dev/null || true
    print_success "Documentation → .specify/docs/"
    
    # Example configurations
    cp "$BUNDLE_DIR"/files/examples/*.yaml .specify/examples/research-configs/ 2>/dev/null || true
    print_success "Example configs → .specify/examples/research-configs/"
    
    # Context file
    cp "$BUNDLE_DIR"/files/context/*.yaml .specify/context/ 2>/dev/null || true
    print_success "Context → .specify/context/"
    
    echo ""
}

# Interactive configuration setup
setup_configuration() {
    if [[ "$SKIP_CONFIG" == "true" ]]; then
        print_warning "Skipping configuration (you'll need to set up .specify/config/research.yaml manually)"
        echo ""
        return 0
    fi
    
    print_step "Setting up configuration..."
    echo ""
    
    if [[ "$INSTALL_MODE" == "interactive" ]] && [[ -z "$PROJECT_TYPE" ]]; then
        echo "What type of project is this?"
        echo "  1) Generic/Default (balanced)"
        echo "  2) Web Application (security & scalability focused)"
        echo "  3) Library/SDK (API design focused)"
        echo "  4) Infrastructure/DevOps (reliability focused)"
        echo ""
        read -p "Enter choice [1-4]: " choice
        
        case $choice in
            1) PROJECT_TYPE="default" ;;
            2) PROJECT_TYPE="web-app" ;;
            3) PROJECT_TYPE="library" ;;
            4) PROJECT_TYPE="infrastructure" ;;
            *) PROJECT_TYPE="default" ;;
        esac
    fi
    
    PROJECT_TYPE="${PROJECT_TYPE:-default}"
    
    # Copy selected configuration template
    cp ".specify/examples/research-configs/${PROJECT_TYPE}.yaml" .specify/config/research.yaml
    print_success "Configuration created from ${PROJECT_TYPE}.yaml template"
    
    echo ""
    print_info "Edit .specify/config/research.yaml to customize for your project:"
    print_info "  - Project name and type"
    print_info "  - File paths"
    print_info "  - Decision criteria weights"
    print_info "  - Performance gates"
    echo ""
}

# Verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    cd "$PROJECT_ROOT"
    
    local errors=()
    
    # Check critical files
    [[ -f ".github/agents/speckit.research.agent.md" ]] || errors+=("Agent definition missing")
    [[ -f ".specify/templates/research-template.yaml" ]] || errors+=("Research template missing")
    [[ -f ".specify/scripts/bash/research.sh" ]] || errors+=("Helper script missing")
    [[ -f ".specify/config/research.yaml" ]] || errors+=("Configuration missing")
    
    if [[ ${#errors[@]} -gt 0 ]]; then
        print_error "Installation verification failed:"
        for err in "${errors[@]}"; do
            print_info "  - $err"
        done
        echo ""
        return 1
    fi
    
    print_success "Installation verified"
    echo ""
    
    # Run prerequisites check if available
    if [[ -f ".specify/scripts/bash/check-research-prerequisites.sh" ]]; then
        print_step "Running prerequisites check..."
        echo ""
        if .specify/scripts/bash/check-research-prerequisites.sh; then
            echo ""
            print_success "Prerequisites check passed"
        else
            echo ""
            print_warning "Prerequisites check found issues (see above)"
            print_info "These can be fixed after installation"
        fi
        echo ""
    fi
}

# Print next steps
print_next_steps() {
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${NC}  ${GREEN}Installation Complete!${NC}                                   ${GREEN}║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Next Steps:${NC}"
    echo ""
    echo "1. Customize configuration:"
    echo -e "   ${YELLOW}vim .specify/config/research.yaml${NC}"
    echo ""
    echo "2. Verify prerequisites:"
    echo -e "   ${YELLOW}.specify/scripts/bash/check-research-prerequisites.sh${NC}"
    echo ""
    echo "3. Read the setup guide:"
    echo -e "   ${YELLOW}cat .specify/docs/research-agent-setup.md${NC}"
    echo ""
    echo "4. Start researching:"
    echo -e "   ${YELLOW}@speckit.research \"Your research question\"${NC}"
    echo ""
    echo -e "${BLUE}Documentation:${NC}"
    echo "  - Setup Guide: .specify/docs/research-agent-setup.md"
    echo "  - Framework: .specify/docs/research-framework.md"
    echo "  - Integration: .specify/docs/agent-integration.md"
    echo ""
    echo -e "${BLUE}Configuration Templates:${NC}"
    echo "  - Default: .specify/examples/research-configs/default.yaml"
    echo "  - Web App: .specify/examples/research-configs/web-app.yaml"
    echo "  - Library: .specify/examples/research-configs/library.yaml"
    echo "  - Infrastructure: .specify/examples/research-configs/infrastructure.yaml"
    echo ""
}

# Main installation flow
main() {
    print_header
    
    echo "Installing to: $PROJECT_ROOT"
    echo ""
    
    check_prerequisites
    create_directories
    copy_files
    setup_configuration
    verify_installation
    print_next_steps
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --auto)
            INSTALL_MODE="auto"
            shift
            ;;
        --project-type)
            PROJECT_TYPE="$2"
            shift 2
            ;;
        --skip-config)
            SKIP_CONFIG="true"
            shift
            ;;
        --project-root)
            PROJECT_ROOT="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Install SpecKit Research Agent into a project"
            echo ""
            echo "Options:"
            echo "  --auto                  Non-interactive mode"
            echo "  --project-type TYPE     Project type (default|web-app|library|infrastructure)"
            echo "  --skip-config           Skip configuration setup"
            echo "  --project-root PATH     Install to specific directory"
            echo "  --help, -h              Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run installation
main
