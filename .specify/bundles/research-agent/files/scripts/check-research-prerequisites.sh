#!/usr/bin/env bash
# Research Prerequisites Checker
# Validates environment before research agent execution
# Similar to check-prerequisites.sh used by other speckit agents

set -euo pipefail

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Default values
OUTPUT_FORMAT="text"
VERBOSE=false

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --json)
            OUTPUT_FORMAT="json"
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Validate research environment prerequisites"
            echo ""
            echo "Options:"
            echo "  --json          Output in JSON format"
            echo "  --verbose, -v   Verbose output"
            echo "  --help, -h      Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Change to project root
cd "$PROJECT_ROOT"

# Initialize check results
ERRORS=()
WARNINGS=()
INFO=()

# Check 1: Configuration file exists
CONFIG_FILE=".specify/config/research.yaml"
if [[ ! -f "$CONFIG_FILE" ]]; then
    ERRORS+=("Configuration file not found: $CONFIG_FILE")
    ERRORS+=("Run: .specify/scripts/bash/research.sh init")
else
    INFO+=("Configuration file: $CONFIG_FILE")
fi

# Check 2: Required tools
command -v yq >/dev/null 2>&1 || {
    ERRORS+=("yq not found - install with: brew install yq (macOS) or apt install yq (Linux)")
}

command -v git >/dev/null 2>&1 || {
    WARNINGS+=("git not found - some features may not work")
}

# Check 3: Templates exist
RESEARCH_TEMPLATE=".specify/templates/research-template.yaml"
FEEDBACK_TEMPLATE=".specify/templates/research-feedback-template.md"

if [[ ! -f "$RESEARCH_TEMPLATE" ]]; then
    ERRORS+=("Research template not found: $RESEARCH_TEMPLATE")
else
    INFO+=("Research template: $RESEARCH_TEMPLATE")
fi

if [[ ! -f "$FEEDBACK_TEMPLATE" ]]; then
    WARNINGS+=("Feedback template not found: $FEEDBACK_TEMPLATE")
else
    INFO+=("Feedback template: $FEEDBACK_TEMPLATE")
fi

# Check 4: Load and validate config (if it exists)
if [[ -f "$CONFIG_FILE" ]]; then
    # Extract paths from config
    RESEARCH_ROOT=$(yq '.paths.research_root' "$CONFIG_FILE" 2>/dev/null || echo "null")
    
    if [[ "$RESEARCH_ROOT" == "null" ]]; then
        ERRORS+=("Invalid config: paths.research_root not defined")
    else
        INFO+=("Research root: $RESEARCH_ROOT")
        
        # Check if research root exists
        if [[ ! -d "$RESEARCH_ROOT" ]]; then
            WARNINGS+=("Research directory does not exist: $RESEARCH_ROOT")
            WARNINGS+=("Will be created when first research area is initialized")
        fi
    fi
    
    # Check constitution (optional)
    CONSTITUTION=$(yq '.paths.constitution' "$CONFIG_FILE" 2>/dev/null || echo "null")
    if [[ "$CONSTITUTION" != "null" ]] && [[ -n "$CONSTITUTION" ]]; then
        if [[ ! -f "$CONSTITUTION" ]]; then
            WARNINGS+=("Constitution file specified but not found: $CONSTITUTION")
        else
            INFO+=("Constitution: $CONSTITUTION")
        fi
    fi
    
    # Check decision criteria weights
    WEIGHTS_CHECK=$(yq '.decision_criteria.weights' "$CONFIG_FILE" 2>/dev/null || echo "null")
    if [[ "$WEIGHTS_CHECK" == "null" ]]; then
        ERRORS+=("Invalid config: decision_criteria.weights not defined")
    fi
    
    # Check schema version
    SCHEMA_VERSION=$(yq '.metadata.schema_version' "$CONFIG_FILE" 2>/dev/null || echo "null")
    if [[ "$SCHEMA_VERSION" == "null" ]]; then
        WARNINGS+=("No schema version in config - expected 1.0.0")
    elif [[ "$SCHEMA_VERSION" != "1.0.0" ]]; then
        WARNINGS+=("Schema version $SCHEMA_VERSION may be incompatible (expected 1.0.0)")
    fi
fi

# Check 5: Git repository (for commit tracking)
if git rev-parse --git-dir > /dev/null 2>&1; then
    INFO+=("Git repository: $(git rev-parse --show-toplevel)")
else
    WARNINGS+=("Not in a git repository - commit tracking will not work")
fi

# Output results
if [[ "$OUTPUT_FORMAT" == "json" ]]; then
    # JSON output
    cat <<EOF
{
  "status": "$([ ${#ERRORS[@]} -eq 0 ] && echo "ok" || echo "error")",
  "errors": [$(printf '"%s",' "${ERRORS[@]}" | sed 's/,$//')],
  "warnings": [$(printf '"%s",' "${WARNINGS[@]}" | sed 's/,$//')],
  "info": [$(printf '"%s",' "${INFO[@]}" | sed 's/,$//')],
  "project_root": "$PROJECT_ROOT",
  "config_file": "$CONFIG_FILE",
  "research_root": "${RESEARCH_ROOT:-null}",
  "constitution": "${CONSTITUTION:-null}"
}
EOF
else
    # Text output
    echo "Research Prerequisites Check"
    echo "============================"
    echo ""
    
    if [[ ${#ERRORS[@]} -gt 0 ]]; then
        echo -e "${RED}ERRORS:${NC}"
        for err in "${ERRORS[@]}"; do
            echo -e "  ${RED}✗${NC} $err"
        done
        echo ""
    fi
    
    if [[ ${#WARNINGS[@]} -gt 0 ]]; then
        echo -e "${YELLOW}WARNINGS:${NC}"
        for warn in "${WARNINGS[@]}"; do
            echo -e "  ${YELLOW}⚠${NC} $warn"
        done
        echo ""
    fi
    
    if [[ "$VERBOSE" == "true" ]] || [[ ${#ERRORS[@]} -eq 0 && ${#WARNINGS[@]} -eq 0 ]]; then
        echo -e "${GREEN}INFO:${NC}"
        for info in "${INFO[@]}"; do
            echo -e "  ${GREEN}✓${NC} $info"
        done
        echo ""
    fi
    
    if [[ ${#ERRORS[@]} -eq 0 ]]; then
        echo -e "${GREEN}✓ All prerequisites met${NC}"
        echo ""
        echo "Ready to research! Usage:"
        echo "  @speckit.research \"Your research question\""
        echo "  @speckit.research R01  # Research specific area"
    else
        echo -e "${RED}✗ Prerequisites check failed${NC}"
        echo ""
        echo "Fix the errors above before running research agent"
    fi
fi

# Exit with appropriate code
if [[ ${#ERRORS[@]} -gt 0 ]]; then
    exit 1
else
    exit 0
fi
