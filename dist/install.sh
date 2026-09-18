#!/usr/bin/env bash
# ==============================================================================
# X-Parity Universal Installer
# Automatically detects OS and Architecture, installs x-parity to your PATH.
# ==============================================================================

set -e

REPO="AppeiYA/x-parity"
BINARY_NAME="x-parity"

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
    linux*)     OS="linux" ;;
    darwin*)    OS="darwin" ;;
    msys*|cygwin*|mingw*) OS="windows" ;;
    *)
        echo "Error: Unsupported operating system: ${OS}"
        exit 1
        ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64|amd64)   ARCH="amd64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    *)
        echo "Error: Unsupported CPU architecture: ${ARCH}"
        exit 1
        ;;
esac

TARGET_BINARY="${BINARY_NAME}-${OS}-${ARCH}"
if [ "${OS}" = "windows" ]; then
    TARGET_BINARY="${TARGET_BINARY}.exe"
fi

echo "==> Detected Platform: ${OS}/${ARCH}"

# 3. Determine Installation Directory
if [ -z "${INSTALL_DIR}" ]; then
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "${INSTALL_DIR}" ]; then
        INSTALL_DIR="${HOME}/.local/bin"
    fi
fi
mkdir -p "${INSTALL_DIR}"

DEST_FILE="${INSTALL_DIR}/${BINARY_NAME}"
if [ "${OS}" = "windows" ]; then
    DEST_FILE="${DEST_FILE}.exe"
fi

# 4. Source the Binary (Local repository check or GitHub download)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "${SCRIPT_DIR}/${TARGET_BINARY}" ]; then
    echo "==> Installing from local distribution: ${SCRIPT_DIR}/${TARGET_BINARY}"
    cp "${SCRIPT_DIR}/${TARGET_BINARY}" "${DEST_FILE}"
elif [ -f "${SCRIPT_DIR}/dist/${TARGET_BINARY}" ]; then
    echo "==> Installing from local dist folder: ${SCRIPT_DIR}/dist/${TARGET_BINARY}"
    cp "${SCRIPT_DIR}/dist/${TARGET_BINARY}" "${DEST_FILE}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${TARGET_BINARY}"
    echo "==> Downloading ${TARGET_BINARY} from GitHub Releases..."
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "${DOWNLOAD_URL}" -o "${DEST_FILE}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "${DEST_FILE}" "${DOWNLOAD_URL}"
    else
        echo "Error: Neither curl nor wget was found. Please download manually."
        exit 1
    fi
fi

chmod +x "${DEST_FILE}"

# 5. Verify Installation
echo "=================================================="
echo "✓ X-Parity installed successfully to:"
echo "  ${DEST_FILE}"
echo "=================================================="

# Check if INSTALL_DIR is in PATH
case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) ;;
    *)
        echo ""
        echo "Notice: ${INSTALL_DIR} is not in your PATH."
        echo "Add it to your shell profile (~/.bashrc or ~/.zshrc):"
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
        ;;
esac

echo "Quick Start:"
echo "  x-parity help"
echo "  x-parity capture -app demo -env local -out snap.json"
echo ""
