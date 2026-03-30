#!/bin/sh
# Install script for copypasta
# Usage: curl -sSL https://raw.githubusercontent.com/guigui42/copypasta/main/install.sh | sh

set -e

REPO="guigui42/copypasta"

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  arm64)   ARCH="arm64" ;;
  aarch64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$OS" != "darwin" ]; then
  echo "Error: copypasta only works on macOS (uses system clipboard APIs)."
  exit 1
fi

# Get latest release tag
LATEST=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
  echo "Error: Could not determine latest release."
  exit 1
fi

VERSION="${LATEST#v}"
FILENAME="copypasta_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

INSTALL_DIR="/usr/local/bin"
TMP_DIR=$(mktemp -d)

echo "Downloading copypasta ${LATEST} for ${OS}/${ARCH}..."
curl -sSL "$URL" -o "${TMP_DIR}/${FILENAME}"

echo "Installing to ${INSTALL_DIR}..."
tar -xzf "${TMP_DIR}/${FILENAME}" -C "$TMP_DIR"
install -m 755 "${TMP_DIR}/copypasta" "${INSTALL_DIR}/copypasta"

rm -rf "$TMP_DIR"

echo "✓ copypasta ${LATEST} installed to ${INSTALL_DIR}/copypasta"
echo "  Run 'copypasta' after copying text to clean it up!"
