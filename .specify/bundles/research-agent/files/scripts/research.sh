#!/usr/bin/env bash

# Research Infrastructure Management Script
# Project-agnostic version - works with any project using speckit.research
#
# This script provides utilities for the research agent and research workflow.
#
# Usage: ./research.sh [COMMAND] [OPTIONS]
#
# COMMANDS:
#   status              Show status of all research areas
#   select [criteria]   Select next research area based on criteria
#   validate [id]       Validate research area prerequisites
#   update [id]         Update research area status
#   stats               Show research statistics
#   init                Initialize research configuration
#   help                Show this help message
#
# OPTIONS:
#   --json              Output in JSON format
#   --id ID             Specify research area ID
#   --status STATUS     New status (not_started|in_progress|complete|blocked)
#   --hours HOURS       Actual hours spent
#   --version VERSION   Filter by target version
#
# EXAMPLES:
#   ./research.sh init                     # Create initial configuration
#   ./research.sh status --json
#   ./research.sh select --json --version 1.0.0
#   ./research.sh validate --id R01
#   ./research.sh update --id R01 --status in_progress

set -e

# Source common functions
SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/common.sh"

# Get repository root
REPO_ROOT=$(get_repo_root)

# Configuration file (project-specific)
CONFIG_FILE="$REPO_ROOT/.specify/config/research.yaml"

# Load configuration
load_config() {
    if [[ ! -f "$CONFIG_FILE" ]]; then
        echo "ERROR: No research configuration found at $CONFIG_FILE" >&2
        echo "" >&2
        echo "Run './research.sh init' to create configuration." >&2
        exit 1
    fi
    
    # Check if yq is available
    if ! command -v yq &> /dev/null; then
        echo "ERROR: yq is required but not installed." >&2
        echo "Install with: brew install yq (macOS) or see https://github.com/mikefarah/yq" >&2
        exit 1
    fi
    
    # Extract paths from config
    RESEARCH_DIR="$REPO_ROOT/$(yq '.paths.research_root' "$CONFIG_FILE")"
    RESEARCH_PLAN="$REPO_ROOT/$(yq '.paths.research_plan' "$CONFIG_FILE")"
    DECISION_LOG="$REPO_ROOT/$(yq '.paths.decision_log' "$CONFIG_FILE")"
    CONSTITUTION_PATH=$(yq '.paths.constitution' "$CONFIG_FILE")
    
    # Only set CONSTITUTION if it's defined
    if [[ "$CONSTITUTION_PATH" != "null" ]]; then
        CONSTITUTION="$REPO_ROOT/$CONSTITUTION_PATH"
    else
        CONSTITUTION=""
    fi
}

# Parse command line arguments
COMMAND=""
JSON_MODE=false
RESEARCH_ID=""
NEW_STATUS=""
ACTUAL_HOURS=""
TARGET_VERSION=""

parse_args() {
    if [[ $# -eq 0 ]]; then
        COMMAND="help"
        return
    fi
    
    COMMAND="$1"
    shift
    
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --json)
                JSON_MODE=true
                shift
                ;;
            --id)
                RESEARCH_ID="$2"
                shift 2
                ;;
            --status)
                # shellcheck disable=SC2034
                NEW_STATUS="$2"
                shift 2
                ;;
            --hours)
                # shellcheck disable=SC2034
                ACTUAL_HOURS="$2"
                shift 2
                ;;
            --version)
                TARGET_VERSION="$2"
                shift 2
                ;;
            --help|-h)
                COMMAND="help"
                shift
                ;;
            *)
                echo "ERROR: Unknown option '$1'. Use --help for usage information." >&2
                exit 1
                ;;
        esac
    done
}

# Check prerequisites
check_prerequisites() {
    local errors=()
    
    # Check if config exists (unless we're running init)
    if [[ "$COMMAND" != "init" ]] && [[ ! -f "$CONFIG_FILE" ]]; then
        errors+=("Configuration not found: $CONFIG_FILE. Run './research.sh init' first.")
        
        if [[ ${#errors[@]} -gt 0 ]]; then
            for error in "${errors[@]}"; do
                echo "ERROR: $error" >&2
            done
            exit 1
        fi
        return
    fi
    
    # Skip other checks for init command
    if [[ "$COMMAND" == "init" ]]; then
        return
    fi
    
    # Load config
    load_config
    
    # Check research directory
    [[ ! -d "$RESEARCH_DIR" ]] && errors+=("Research directory not found: $RESEARCH_DIR")
    [[ ! -f "$RESEARCH_PLAN" ]] && errors+=("Research plan not found: $RESEARCH_PLAN")
    [[ ! -f "$DECISION_LOG" ]] && errors+=("Decision log not found: $DECISION_LOG")
    [[ ! -f "$CONSTITUTION" ]] && errors+=("Constitution not found: $CONSTITUTION")
    
    if [[ ${#errors[@]} -gt 0 ]]; then
        if $JSON_MODE; then
            printf '{"error":"prerequisites_missing","details":["%s"]}\n' "${errors[*]}"
        else
            echo "ERROR: Prerequisites not met:" >&2
            for err in "${errors[@]}"; do
                echo "  - $err" >&2
            done
            echo "" >&2
            echo "Run research infrastructure setup first." >&2
        fi
        exit 1
    fi
}

# Get research areas from YAML (simplified parsing)
get_research_areas() {
    # Extract research areas section and parse key fields
    # This is a simplified parser - for complex YAML, consider using yq
    awk '
    /^research_areas:/ { in_areas=1; next }
    /^[a-z_]+:/ && !/^  / { in_areas=0 }
    in_areas && /^  - id:/ { 
        gsub(/.*: *"?|"?$/, "", $0)
        id=$0
    }
    in_areas && /^    name:/ {
        gsub(/.*: *"?|"?$/, "", $0)
        name=$0
    }
    in_areas && /^    folder:/ {
        gsub(/.*: *"?|"?$/, "", $0)
        folder=$0
    }
    in_areas && /^    target_version:/ {
        gsub(/.*: *"?|"?$/, "", $0)
        version=$0
    }
    in_areas && /^    priority:/ {
        gsub(/.*: *"?|"?$/, "", $0)
        priority=$0
    }
    in_areas && /^    status:/ {
        gsub(/.*: *"?|"?$/, "", $0)
        status=$0
    }
    in_areas && /^    estimated_hours:/ {
        gsub(/.*: */, "", $0)
        est_hours=$0
    }
    in_areas && /^    actual_hours:/ {
        gsub(/.*: */, "", $0)
        act_hours=$0
        # Print complete record
        printf "%s|%s|%s|%s|%s|%s|%s|%s\n", id, name, folder, version, priority, status, est_hours, act_hours
    }
    ' "$RESEARCH_PLAN"
}

# Show status of all research areas
cmd_status() {
    check_prerequisites
    
    if $JSON_MODE; then
        echo '{"research_areas":['
        local first=true
        while IFS='|' read -r id name folder version priority status est_hours act_hours; do
            [[ -z "$id" ]] && continue
            $first || echo ","
            first=false
            printf '{"id":"%s","name":"%s","folder":"%s","version":"%s","priority":"%s","status":"%s","estimated_hours":%s,"actual_hours":%s}' \
                "$id" "$name" "$folder" "$version" "$priority" "$status" "$est_hours" "$act_hours"
        done < <(get_research_areas)
        echo ']}'
    else
        echo "Research Status"
        echo "==============="
        echo ""
        printf "%-5s %-30s %-8s %-8s %-12s %s\n" "ID" "Name" "Version" "Priority" "Status" "Hours"
        printf "%-5s %-30s %-8s %-8s %-12s %s\n" "----" "----" "-------" "--------" "------" "-----"
        while IFS='|' read -r id name folder version priority status est_hours act_hours; do
            [[ -z "$id" ]] && continue
            printf "%-5s %-30s %-8s %-8s %-12s %s/%s\n" \
                "$id" "${name:0:30}" "$version" "$priority" "$status" "$act_hours" "$est_hours"
        done < <(get_research_areas)
    fi
}

# Select next research area
cmd_select() {
    check_prerequisites
    
    local best_id=""
    local best_priority=99
    local best_version="9.9.9"
    
    while IFS='|' read -r id name folder version priority status est_hours act_hours; do
        [[ -z "$id" ]] && continue
        [[ "$status" != "not_started" ]] && continue
        
        # Filter by version if specified
        if [[ -n "$TARGET_VERSION" ]] && [[ "$version" != "$TARGET_VERSION" ]]; then
            continue
        fi
        
        # Convert priority to number (P0=0, P1=1, P2=2, P3=3)
        local pri_num=${priority#P}
        
        # Select if higher priority (lower number) or same priority but earlier version
        if [[ $pri_num -lt $best_priority ]] || \
           ([[ $pri_num -eq $best_priority ]] && [[ "$version" < "$best_version" ]]); then
            best_id="$id"
            best_priority=$pri_num
            best_version="$version"
        fi
    done < <(get_research_areas)
    
    if [[ -z "$best_id" ]]; then
        if $JSON_MODE; then
            echo '{"selected":null,"message":"No research areas available"}'
        else
            echo "No research areas available matching criteria"
        fi
        exit 0
    fi
    
    # Get full details for selected area
    while IFS='|' read -r id name folder version priority status est_hours act_hours; do
        if [[ "$id" == "$best_id" ]]; then
            if $JSON_MODE; then
                printf '{"selected":{"id":"%s","name":"%s","folder":"%s","version":"%s","priority":"%s","path":"%s"}}\n' \
                    "$id" "$name" "$folder" "$version" "$priority" "$RESEARCH_DIR/$folder"
            else
                echo "Selected: $id - $name"
                echo "  Version: $version"
                echo "  Priority: $priority"
                echo "  Path: $RESEARCH_DIR/$folder"
            fi
            break
        fi
    done < <(get_research_areas)
}

# Validate research area
cmd_validate() {
    check_prerequisites
    
    if [[ -z "$RESEARCH_ID" ]]; then
        echo "ERROR: --id required for validate command" >&2
        exit 1
    fi
    
    local found=false
    local folder=""
    
    while IFS='|' read -r id name f version priority status est_hours act_hours; do
        if [[ "$id" == "$RESEARCH_ID" ]]; then
            found=true
            folder="$f"
            break
        fi
    done < <(get_research_areas)
    
    if ! $found; then
        if $JSON_MODE; then
            printf '{"valid":false,"error":"Research area %s not found"}\n' "$RESEARCH_ID"
        else
            echo "ERROR: Research area $RESEARCH_ID not found" >&2
        fi
        exit 1
    fi
    
    local area_dir="$RESEARCH_DIR/$folder"
    local research_yaml="$area_dir/research.yaml"
    local findings_md="$area_dir/findings.md"
    
    local errors=()
    [[ ! -d "$area_dir" ]] && errors+=("Directory not found: $area_dir")
    [[ ! -f "$research_yaml" ]] && errors+=("research.yaml not found")
    [[ ! -f "$findings_md" ]] && errors+=("findings.md not found")
    
    if [[ ${#errors[@]} -gt 0 ]]; then
        if $JSON_MODE; then
            printf '{"valid":false,"id":"%s","errors":["%s"]}\n' "$RESEARCH_ID" "${errors[*]}"
        else
            echo "Validation FAILED for $RESEARCH_ID:"
            for err in "${errors[@]}"; do
                echo "  ✗ $err"
            done
        fi
        exit 1
    fi
    
    if $JSON_MODE; then
        printf '{"valid":true,"id":"%s","path":"%s","files":["research.yaml","findings.md"]}\n' \
            "$RESEARCH_ID" "$area_dir"
    else
        echo "Validation PASSED for $RESEARCH_ID:"
        echo "  ✓ $area_dir"
        echo "  ✓ research.yaml"
        echo "  ✓ findings.md"
    fi
}

# Show statistics
cmd_stats() {
    check_prerequisites
    
    local total=0 not_started=0 in_progress=0 complete=0 blocked=0
    local total_est=0 total_act=0
    
    while IFS='|' read -r id name folder version priority status est_hours act_hours; do
        [[ -z "$id" ]] && continue
        ((total++))
        total_est=$((total_est + est_hours))
        total_act=$((total_act + act_hours))
        
        case "$status" in
            not_started) ((not_started++)) ;;
            in_progress) ((in_progress++)) ;;
            complete) ((complete++)) ;;
            blocked) ((blocked++)) ;;
        esac
    done < <(get_research_areas)
    
    local pct=0
    [[ $total -gt 0 ]] && pct=$((complete * 100 / total))
    
    if $JSON_MODE; then
        printf '{"total":%d,"not_started":%d,"in_progress":%d,"complete":%d,"blocked":%d,"estimated_hours":%d,"actual_hours":%d,"completion_percentage":%d}\n' \
            "$total" "$not_started" "$in_progress" "$complete" "$blocked" "$total_est" "$total_act" "$pct"
    else
        echo "Research Statistics"
        echo "==================="
        echo ""
        echo "Areas:     $total total"
        echo "  Not Started: $not_started"
        echo "  In Progress: $in_progress"
        echo "  Complete:    $complete"
        echo "  Blocked:     $blocked"
        echo ""
        echo "Hours:     $total_act / $total_est estimated"
        echo "Progress:  $pct%"
    fi
}

# Show help
show_help() {
    cat << 'EOF'
Research Infrastructure Management Script

Usage: ./research.sh [COMMAND] [OPTIONS]

COMMANDS:
  init                Initialize research configuration
  status              Show status of all research areas
  select              Select next research area (highest priority, not started)
  validate            Validate research area has required files
  stats               Show research statistics
  help                Show this help message

OPTIONS:
  --json              Output in JSON format
  --id ID             Specify research area ID
  --version X.Y.Z     Filter by target version

EXAMPLES:
  # Initialize configuration
  ./research.sh init

  # Show all research areas
  ./research.sh status

  # Get next research area as JSON
  ./research.sh select --json

  # Get next research area for v1.0.0
  ./research.sh select --version 1.0.0

  # Validate R01 has required files
  ./research.sh validate --id R01

  # Get statistics
  ./research.sh stats --json

EOF
}

# Backward compatibility aliases
cmd_status() { research_status; }
cmd_select() { research_select; }
cmd_validate() { research_validate; }
cmd_stats() { research_stats; }
cmd_help() { show_help; }

# Main execution is at bottom of file
        cmd_select
        ;;
    validate)
        cmd_validate
        ;;
    stats)
        cmd_stats
        ;;
    help|--help|-h)
        cmd_help
        ;;
    *)
        echo "ERROR: Unknown command '$COMMAND'" >&2
        echo "Use './research.sh help' for usage information." >&2
        exit 1
        ;;
esac
