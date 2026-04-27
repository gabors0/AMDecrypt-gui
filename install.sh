#!/bin/sh
set -eu

REPO_URL="${AMDECRYPT_GUI_REPO:-https://github.com/gabors0/AMDecrypt-gui.git}"
REPO_REF="${AMDECRYPT_GUI_REF:-}"
WAILS_BUILD_TAGS="${WAILS_BUILD_TAGS:-webkit2_41}"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/amdecrypt-gui-install.XXXXXX")"
APP_NAME="amdecrypt-gui"
SOURCE_DIR=""

if [ -n "${SUDO_USER:-}" ] && [ "$SUDO_USER" != "root" ]; then
    INVOKING_USER="$SUDO_USER"
else
    INVOKING_USER="$(id -un)"
fi

INVOKING_HOME="$(getent passwd "$INVOKING_USER" | cut -d: -f6)"

cleanup() {
    if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}

trap cleanup EXIT

require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "Missing required command: $1" >&2
        exit 1
    fi
}

resolve_wails() {
    if command -v wails >/dev/null 2>&1; then
        command -v wails
        return 0
    fi

    if [ -n "${GOBIN:-}" ] && [ -x "${GOBIN}/wails" ]; then
        printf '%s\n' "${GOBIN}/wails"
        return 0
    fi

    if [ -n "${INVOKING_HOME:-}" ] && [ -x "${INVOKING_HOME}/go/bin/wails" ]; then
        printf '%s\n' "${INVOKING_HOME}/go/bin/wails"
        return 0
    fi

    return 1
}

run_as_root() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    else
        sudo "$@"
    fi
}

install_for_user() {
    mode="$1"
    source_path="$2"
    target_path="$3"
    target_dir="$(dirname "$target_path")"

    if [ "$(id -u)" -eq 0 ]; then
        run_as_root install -d -o "$INVOKING_USER" -g "$INVOKING_USER" "$target_dir"
        run_as_root install -m "$mode" -o "$INVOKING_USER" -g "$INVOKING_USER" "$source_path" "$target_path"
    else
        mkdir -p "$target_dir"
        install -m "$mode" "$source_path" "$target_path"
    fi
}

install_built_app() {
    repo_dir="$1"
    binary="$repo_dir/build/bin/$APP_NAME"
    icon="$repo_dir/build/appicon.png"
    install_bin="/usr/local/bin/$APP_NAME"

    if [ -z "$INVOKING_HOME" ]; then
        echo "Failed to resolve home directory for user: $INVOKING_USER" >&2
        exit 1
    fi

    install_icon="$INVOKING_HOME/.local/share/icons/hicolor/256x256/apps/$APP_NAME.png"
    install_desktop="$INVOKING_HOME/.local/share/applications/$APP_NAME.desktop"
    desktop_tmp="$TMP_DIR/$APP_NAME.desktop"

    if [ ! -f "$binary" ]; then
        echo "Binary not found: $binary" >&2
        echo "Build failed or did not produce the expected output." >&2
        exit 1
    fi

    echo "Installing $APP_NAME..."
    echo "Target user: $INVOKING_USER"
    echo "Binary path: $install_bin"
    echo "Desktop file: $install_desktop"

    run_as_root install -Dm755 "$binary" "$install_bin"
    install_for_user 644 "$icon" "$install_icon"

    cat > "$desktop_tmp" <<EOF
[Desktop Entry]
Name=AMDecrypt-GUI
Comment=GUI for AppleMusicDecrypt
Exec=$install_bin
Icon=$install_icon
Type=Application
Terminal=false
Categories=Utility;AudioVideo;
StartupWMClass=amdecrypt-gui
StartupNotify=true
EOF
    install_for_user 644 "$desktop_tmp" "$install_desktop"

    gtk-update-icon-cache "$INVOKING_HOME/.local/share/icons/hicolor/" 2>/dev/null || true
    update-desktop-database "$INVOKING_HOME/.local/share/applications/" 2>/dev/null || true
    xdg-desktop-menu forceupdate 2>/dev/null || true
}

require_cmd git
require_cmd go
require_cmd npm

if ! WAILS_BIN="$(resolve_wails)"; then
    echo "Missing required command: wails" >&2
    echo "Install it with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" >&2
    echo "If you run this script with sudo, make sure wails is installed for $INVOKING_USER." >&2
    exit 1
fi

if [ -f "./wails.json" ]; then
    SOURCE_DIR="$(pwd)"
    echo "Using local checkout at $SOURCE_DIR..."
else
    echo "Cloning $REPO_URL..."
    git clone --depth 1 "$REPO_URL" "$TMP_DIR/repo"
    SOURCE_DIR="$TMP_DIR/repo"
fi

cd "$SOURCE_DIR"

if [ -n "$REPO_REF" ] && [ "$SOURCE_DIR" = "$TMP_DIR/repo" ]; then
    echo "Checking out $REPO_REF..."
    git fetch --depth 1 origin "$REPO_REF"
    git checkout --detach FETCH_HEAD
fi

echo "Building AMDecrypt-gui..."
if [ -n "$WAILS_BUILD_TAGS" ]; then
    "$WAILS_BIN" build -tags "$WAILS_BUILD_TAGS"
else
    "$WAILS_BIN" build
fi

echo "Installing built application..."
install_built_app "$SOURCE_DIR"

echo "AMDecrypt-gui is now installed (or updated)."
