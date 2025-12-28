#!/usr/bin/env bash
# changelog-preview.sh - Generate changelog preview from commits
# Usage: changelog-preview.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Get commits since last release tag
get_commits_since_last_release() {
    local last_tag
    last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    
    if [[ -z "${last_tag}" ]]; then
        # No tags yet, get all commits
        git log --pretty=format:"%s" HEAD
    else
        # Get commits since last tag
        git log --pretty=format:"%s" "${last_tag}..HEAD"
    fi
}

# Parse conventional commits into sections
parse_commits() {
    local features=()
    local fixes=()
    local docs=()
    local perf=()
    local other=()
    
    while IFS= read -r commit; do
        case "${commit}" in
            feat:*|feat\(*\):*)
                features+=("${commit#feat*: }")
                ;;
            fix:*|fix\(*\):*)
                fixes+=("${commit#fix*: }")
                ;;
            docs:*|docs\(*\):*)
                docs+=("${commit#docs*: }")
                ;;
            perf:*|perf\(*\):*)
                perf+=("${commit#perf*: }")
                ;;
            *)
                # Include other conventional commit types
                other+=("${commit}")
                ;;
        esac
    done < <(get_commits_since_last_release)
    
    # Output formatted preview
    echo "### Changes in this PR"
    echo ""
    
    if [[ ${#features[@]} -gt 0 ]]; then
        echo "#### Features"
        printf '- %s\n' "${features[@]}"
        echo ""
    fi
    
    if [[ ${#fixes[@]} -gt 0 ]]; then
        echo "#### Bug Fixes"
        printf '- %s\n' "${fixes[@]}"
        echo ""
    fi
    
    if [[ ${#docs[@]} -gt 0 ]]; then
        echo "#### Documentation"
        printf '- %s\n' "${docs[@]}"
        echo ""
    fi
    
    if [[ ${#perf[@]} -gt 0 ]]; then
        echo "#### Performance"
        printf '- %s\n' "${perf[@]}"
        echo ""
    fi
    
    if [[ ${#other[@]} -gt 0 ]]; then
        echo "#### Other Changes"
        printf '- %s\n' "${other[@]}"
        echo ""
    fi
    
    if [[ ${#features[@]} -eq 0 && ${#fixes[@]} -eq 0 && ${#docs[@]} -eq 0 && ${#perf[@]} -eq 0 && ${#other[@]} -eq 0 ]]; then
        echo "_No conventional commits found since last release._"
        echo ""
    fi
}

# Main execution
cd "${REPO_ROOT}"
parse_commits
