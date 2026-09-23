#!/bin/sh
set -e

REPO="vervstack/verv"
INSTALL_DIR="${HOME}/.local/bin"
BIN_NAME="verv"

uname_os=$(uname -s)
case "$uname_os" in
    Linux)
        OS="linux"
        ;;
    Darwin)
        OS="darwin"
        ;;
    *)
        echo "Error: unsupported OS '$uname_os'" >&2
        exit 1
        ;;
esac

uname_arch=$(uname -m)
case "$uname_arch" in
    x86_64 | amd64)
        ARCH="amd64"
        ;;
    arm64 | aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: unsupported architecture '$uname_arch'" >&2
        exit 1
        ;;
esac

ASSET="verv_${OS}_${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
TMP_FILE=$(mktemp)

echo "Downloading ${URL}"

if command -v wget >/dev/null 2>&1; then
    if ! wget -q -O "$TMP_FILE" "$URL"; then
        echo "Error: failed to download ${URL}" >&2
        rm -f "$TMP_FILE"
        exit 1
    fi
elif command -v curl >/dev/null 2>&1; then
    if ! curl -fsSL -o "$TMP_FILE" "$URL"; then
        echo "Error: failed to download ${URL}" >&2
        rm -f "$TMP_FILE"
        exit 1
    fi
else
    echo "Error: neither wget nor curl is available" >&2
    rm -f "$TMP_FILE"
    exit 1
fi

chmod +x "$TMP_FILE"

mkdir -p "$INSTALL_DIR"

DEST="${INSTALL_DIR}/${BIN_NAME}"

if ! mv "$TMP_FILE" "$DEST"; then
    echo "Error: failed to install to ${DEST}" >&2
    rm -f "$TMP_FILE"
    exit 1
fi

echo "Installed ${BIN_NAME} to ${DEST}. Run '${BIN_NAME} --version' to confirm."

case ":${PATH}:" in
    *":${INSTALL_DIR}:"*)
        ;;
    *)
        echo "Note: ${INSTALL_DIR} is not on your PATH. Add it, e.g.:" >&2
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\"" >&2
        ;;
esac
