#!/usr/bin/env bash
# Bundle Verification Script
# Verifies integrity of bundled files

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHECKSUM_FILE="$SCRIPT_DIR/checksums.txt"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "Verifying SpecKit Research Agent Bundle..."
echo ""

if [[ ! -f "$CHECKSUM_FILE" ]]; then
    echo -e "${RED}✗${NC} Checksum file not found: $CHECKSUM_FILE"
    exit 1
fi

# Change to bundle directory
cd "$SCRIPT_DIR"

# Verify checksums
echo "Checking file integrity..."
if shasum -a 256 -c checksums.txt 2>&1 | grep -q "FAILED"; then
    echo ""
    echo -e "${RED}✗${NC} Checksum verification FAILED"
    echo "Some files may have been modified or corrupted"
    shasum -a 256 -c checksums.txt
    exit 1
else
    verified_count=$(wc -l < checksums.txt | tr -d ' ')
    echo -e "${GREEN}✓${NC} All $verified_count files verified"
fi

echo ""

# Check for required files
echo "Checking required files..."
required_files=(
    "files/agents/speckit.research.agent.md"
    "files/templates/research-template.yaml"
    "files/templates/research-feedback-template.md"
    "files/scripts/research.sh"
    "files/scripts/check-research-prerequisites.sh"
    "files/docs/research-framework.md"
    "files/docs/research-agent-setup.md"
    "files/docs/agent-integration.md"
    "files/examples/default.yaml"
    "files/examples/web-app.yaml"
    "files/examples/library.yaml"
    "files/examples/infrastructure.yaml"
    "files/context/research-agent-context.yaml"
)

missing=()
for file in "${required_files[@]}"; do
    if [[ ! -f "$file" ]]; then
        missing+=("$file")
    fi
done

if [[ ${#missing[@]} -gt 0 ]]; then
    echo -e "${RED}✗${NC} Missing required files:"
    for file in "${missing[@]}"; do
        echo "  - $file"
    done
    exit 1
else
    echo -e "${GREEN}✓${NC} All required files present"
fi

echo ""

# Check executable permissions
echo "Checking executable permissions..."
executable_files=(
    "install.sh"
    "files/scripts/research.sh"
    "files/scripts/check-research-prerequisites.sh"
)

permission_issues=()
for file in "${executable_files[@]}"; do
    if [[ ! -x "$file" ]]; then
        permission_issues+=("$file")
    fi
done

if [[ ${#permission_issues[@]} -gt 0 ]]; then
    echo -e "${YELLOW}⚠${NC} Files missing executable permission:"
    for file in "${permission_issues[@]}"; do
        echo "  - $file"
    done
    echo ""
    echo "Fix with: chmod +x ${permission_issues[*]}"
    echo ""
else
    echo -e "${GREEN}✓${NC} All executables have correct permissions"
fi

echo ""
echo -e "${GREEN}Bundle verification complete!${NC}"
echo ""
echo "To install, run: ./install.sh"
