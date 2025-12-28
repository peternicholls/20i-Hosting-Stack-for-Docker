#!/usr/bin/env bash
# validate.sh - Pre-release validation checks
# Usage: validate.sh [--version <version>] [--changelog] [--tags] [--all]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
VERSION_SCRIPT="${SCRIPT_DIR}/version.sh"
CHANGELOG_FILE="${REPO_ROOT}/CHANGELOG.md"
MANIFEST_FILE="${REPO_ROOT}/.release-please-manifest.json"

# Color/emoji output
SUCCESS="✅"
FAILURE="❌"

# Exit code tracker
EXIT_CODE=0

# Validate version format (semver)
validate_version_format() {
    local version="${1}"
    
    # Regex for semantic versioning (with optional pre-release)
    local semver_regex='^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+(\.[0-9]+)?)?$'
    
    if [[ "${version}" =~ ${semver_regex} ]]; then
        echo "${SUCCESS} Version format valid: ${version}"
        return 0
    else
        echo "${FAILURE} Invalid version format: ${version}" >&2
        echo "  Expected: MAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH-PRERELEASE" >&2
        return 1
    fi
}

# Validate CHANGELOG has entry for version
validate_changelog() {
    local version="${1}"
    
    if [[ ! -f "${CHANGELOG_FILE}" ]]; then
        echo "${FAILURE} CHANGELOG.md not found" >&2
        return 2
    fi
    
    # Check for version entry in CHANGELOG
    if grep -q "## \[${version}\]" "${CHANGELOG_FILE}"; then
        echo "${SUCCESS} CHANGELOG has entry for ${version}"
        return 0
    else
        echo "${FAILURE} CHANGELOG missing entry for version ${version}" >&2
        echo "  Add a section: ## [${version}] - $(date +%Y-%m-%d)" >&2
        return 2
    fi
}

# Validate no duplicate git tags exist
validate_no_duplicate_tags() {
    local version="${1}"
    local tag="v${version}"
    
    if git rev-parse "${tag}" >/dev/null 2>&1; then
        echo "${FAILURE} Git tag already exists: ${tag}" >&2
        echo "  Delete with: git tag -d ${tag} && git push origin :refs/tags/${tag}" >&2
        return 1
    else
        echo "${SUCCESS} No duplicate tag for ${version}"
        return 0
    fi
}

# Validate required files exist
validate_required_files() {
    local files=(
        "CHANGELOG.md"
        "docker-compose.yml"
        "release-please-config.json"
        ".release-please-manifest.json"
    )
    
    local missing=0
    for file in "${files[@]}"; do
        if [[ ! -f "${REPO_ROOT}/${file}" ]]; then
            echo "${FAILURE} Required file missing: ${file}" >&2
            missing=1
        fi
    done
    
    if [[ ${missing} -eq 0 ]]; then
        echo "${SUCCESS} All required files present"
        return 0
    else
        return 3
    fi
}

# Main execution
main() {
    local run_version=false
    local run_changelog=false
    local run_tags=false
    local run_all=false
    local version=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --version)
                run_version=true
                if [[ -n "${2:-}" && "${2}" != --* ]]; then
                    version="$2"
                    shift
                fi
                shift
                ;;
            --changelog)
                run_changelog=true
                shift
                ;;
            --tags)
                run_tags=true
                shift
                ;;
            --all)
                run_all=true
                shift
                ;;
            *)
                echo "Usage: $0 [--version <version>] [--changelog] [--tags] [--all]" >&2
                exit 1
                ;;
        esac
    done
    
    # If no version provided, get from manifest
    if [[ -z "${version}" ]] && { ${run_version} || ${run_changelog} || ${run_tags} || ${run_all}; }; then
        version=$("${VERSION_SCRIPT}" get)
    fi
    
    # Run all validations if --all specified
    if ${run_all}; then
        run_version=true
        run_changelog=true
        run_tags=true
    fi
    
    # Execute validations
    if ${run_version}; then
        validate_version_format "${version}" || EXIT_CODE=$?
    fi
    
    if ${run_changelog}; then
        validate_changelog "${version}" || EXIT_CODE=$?
    fi
    
    if ${run_tags}; then
        validate_no_duplicate_tags "${version}" || EXIT_CODE=$?
    fi
    
    # Always validate required files if any validation runs
    if ${run_version} || ${run_changelog} || ${run_tags} || ${run_all}; then
        validate_required_files || EXIT_CODE=$?
    fi
    
    # If nothing specified, show usage
    if ! ${run_version} && ! ${run_changelog} && ! ${run_tags} && ! ${run_all}; then
        echo "Usage: $0 [--version <version>] [--changelog] [--tags] [--all]" >&2
        echo "  --version <ver>  Validate version format" >&2
        echo "  --changelog      Validate CHANGELOG has entry" >&2
        echo "  --tags           Validate no duplicate git tags" >&2
        echo "  --all            Run all validations" >&2
        exit 1
    fi
    
    exit ${EXIT_CODE}
}

main "$@"
