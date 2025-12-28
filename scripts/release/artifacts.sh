#!/usr/bin/env bash
# artifacts.sh - Package release artifacts
# Usage: artifacts.sh <version>

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DIST_DIR="${REPO_ROOT}/dist"
VERSION="${1:-}"

# Validate arguments
if [[ -z "${VERSION}" ]]; then
    echo "Usage: $0 <version>" >&2
    echo "  version: Version number without 'v' prefix (e.g., 2.0.0)" >&2
    exit 1
fi

# Clean and create dist directory
echo "📦 Preparing distribution directory..."
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

# Create temporary directory for archive contents
TEMP_DIR=$(mktemp -d)
ARCHIVE_NAME="20i-stack-v${VERSION}"
ARCHIVE_DIR="${TEMP_DIR}/${ARCHIVE_NAME}"
mkdir -p "${ARCHIVE_DIR}"

echo "📋 Copying files to archive..."

# Copy main files
cp "${REPO_ROOT}/docker-compose.yml" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/README.md" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/LICENSE" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/CHANGELOG.md" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/CONTRIBUTING.md" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/.env.example" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/20i-gui" "${ARCHIVE_DIR}/"
cp "${REPO_ROOT}/zsh-example-script.zsh" "${ARCHIVE_DIR}/"

# Copy directories
cp -r "${REPO_ROOT}/docker" "${ARCHIVE_DIR}/"
cp -r "${REPO_ROOT}/config" "${ARCHIVE_DIR}/"
cp -r "${REPO_ROOT}/scripts" "${ARCHIVE_DIR}/"

# Exclude release scripts from archive (they're for development)
rm -rf "${ARCHIVE_DIR}/scripts/release"

# Copy source if it exists
if [[ -d "${REPO_ROOT}/src" ]]; then
    cp -r "${REPO_ROOT}/src" "${ARCHIVE_DIR}/"
fi

# Create VERSION file in archive
echo "${VERSION}" > "${ARCHIVE_DIR}/VERSION"

# Create archive
echo "🗜️  Creating archive..."
cd "${TEMP_DIR}"
tar -czf "${DIST_DIR}/${ARCHIVE_NAME}.tar.gz" "${ARCHIVE_NAME}"
cd "${REPO_ROOT}"

# Generate checksums
echo "🔐 Generating checksums..."
cd "${DIST_DIR}"
sha256sum "${ARCHIVE_NAME}.tar.gz" > checksums.sha256

# Create standalone install script
echo "📝 Creating install script..."
cat > "${DIST_DIR}/install.sh" << 'INSTALL_SCRIPT_EOF'
#!/usr/bin/env bash
# 20i Stack Quick Installer
# Downloads and extracts the latest release

set -euo pipefail

VERSION="${1:-latest}"
GITHUB_REPO="peternicholls/20i-stack"
INSTALL_DIR="${2:-${HOME}/20i-stack}"

echo "🚀 20i Stack Installer"
echo "━━━━━━━━━━━━━━━━━━━━━━"

# Determine download URL
if [[ "${VERSION}" == "latest" ]]; then
    echo "📥 Fetching latest release information..."
    RELEASE_URL="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    DOWNLOAD_URL=$(curl -sL "${RELEASE_URL}" | grep "browser_download_url.*tar.gz" | cut -d'"' -f4)
    VERSION=$(curl -sL "${RELEASE_URL}" | grep '"tag_name"' | cut -d'"' -f4 | sed 's/^v//')
else
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/20i-stack-v${VERSION}.tar.gz"
fi

echo "📦 Downloading 20i Stack v${VERSION}..."
TEMP_DIR=$(mktemp -d)
cd "${TEMP_DIR}"

curl -LO "${DOWNLOAD_URL}"
ARCHIVE_NAME=$(basename "${DOWNLOAD_URL}")

# Verify checksum if available
CHECKSUM_URL="${DOWNLOAD_URL/%.tar.gz/\/checksums.sha256}"
if curl -sfLO "${CHECKSUM_URL}" 2>/dev/null; then
    echo "🔐 Verifying checksum..."
    sha256sum -c checksums.sha256 --ignore-missing || {
        echo "❌ Checksum verification failed!" >&2
        exit 1
    }
    echo "✅ Checksum verified"
fi

# Extract archive
echo "📂 Extracting to ${INSTALL_DIR}..."
mkdir -p "${INSTALL_DIR}"
tar -xzf "${ARCHIVE_NAME}" -C "${INSTALL_DIR}" --strip-components=1

# Cleanup
cd "${HOME}"
rm -rf "${TEMP_DIR}"

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "  1. cd ${INSTALL_DIR}"
echo "  2. cp .env.example .env"
echo "  3. Edit .env with your configuration"
echo "  4. Run: ./20i-gui"
echo ""
echo "For shell integration, add to your ~/.zshrc:"
echo "  source ${INSTALL_DIR}/zsh-example-script.zsh"
INSTALL_SCRIPT_EOF

chmod +x "${DIST_DIR}/install.sh"

# Add install.sh to checksums
cd "${DIST_DIR}"
sha256sum install.sh >> checksums.sha256

# Cleanup
rm -rf "${TEMP_DIR}"

echo ""
echo "✅ Release artifacts created successfully!"
echo ""
echo "📦 Archive: ${ARCHIVE_NAME}.tar.gz"
echo "📄 Install script: install.sh"
echo "🔐 Checksums: checksums.sha256"
echo ""
echo "Files in ${DIST_DIR}:"
ls -lh "${DIST_DIR}"
