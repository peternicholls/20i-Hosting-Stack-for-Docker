#!/usr/bin/env bash
# version.sh - Read version from release-please manifest
# Usage: version.sh [get]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
MANIFEST_FILE="${REPO_ROOT}/.release-please-manifest.json"

# Function to get current version from manifest
get_version() {
    if [[ ! -f "${MANIFEST_FILE}" ]]; then
        echo "Error: Manifest file not found at ${MANIFEST_FILE}" >&2
        exit 1
    fi
    
    # Extract version from manifest (assumes single package at ".")
    version=$(grep -o '"\.": "[^"]*"' "${MANIFEST_FILE}" | cut -d'"' -f4)
    
    if [[ -z "${version}" ]]; then
        echo "Error: Could not parse version from manifest" >&2
        exit 1
    fi
    
    echo "${version}"
}

# Main execution
case "${1:-get}" in
    get)
        get_version
        ;;
    *)
        echo "Usage: $0 [get]" >&2
        echo "  get - Print current version from manifest (default)" >&2
        exit 1
        ;;
esac
