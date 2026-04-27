#!/bin/sh
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY="$REPO_DIR/build/bin/amdecrypt-gui"
ICON="$REPO_DIR/build/appicon.png"
APP_NAME="amdecrypt-gui"
INSTALL_BIN="/usr/local/bin/$APP_NAME"

if [ -n "${SUDO_USER:-}" ] && [ "$SUDO_USER" != "root" ]; then
    TARGET_USER="$SUDO_USER"
else
    TARGET_USER="$(id -un)"
fi

TARGET_HOME="$(getent passwd "$TARGET_USER" | cut -d: -f6)"
if [ -z "$TARGET_HOME" ]; then
    echo "Failed to resolve home directory for user: $TARGET_USER" >&2
    exit 1
fi

INSTALL_ICON="$TARGET_HOME/.local/share/icons/hicolor/256x256/apps/$APP_NAME.png"
INSTALL_DESKTOP="$TARGET_HOME/.local/share/applications/$APP_NAME.desktop"

if [ ! -f "$BINARY" ]; then
    echo "Binary not found: $BINARY" >&2
    echo "Run 'wails build' first." >&2
    exit 1
fi

echo "Installing $APP_NAME..."
echo "Target user: $TARGET_USER"
echo "Binary path: $INSTALL_BIN"
echo "Desktop file: $INSTALL_DESKTOP"

install -Dm755 "$BINARY" "$INSTALL_BIN"
mkdir -p "$(dirname "$INSTALL_ICON")"
install -Dm644 "$ICON" "$INSTALL_ICON"

mkdir -p "$(dirname "$INSTALL_DESKTOP")"
cat > "$INSTALL_DESKTOP" <<EOF
[Desktop Entry]
Name=AMDecrypt-GUI
Comment=GUI for AppleMusicDecrypt
Exec=$INSTALL_BIN
Icon=$INSTALL_ICON
Type=Application
Terminal=false
Categories=Utility;AudioVideo;
StartupWMClass=amdecrypt-gui
StartupNotify=true
EOF

if [ "$(id -u)" -eq 0 ] && [ "$TARGET_USER" != "root" ]; then
    chown "$TARGET_USER":"$TARGET_USER" "$INSTALL_ICON" "$INSTALL_DESKTOP"
fi

gtk-update-icon-cache "$TARGET_HOME/.local/share/icons/hicolor/" 2>/dev/null || true
update-desktop-database "$TARGET_HOME/.local/share/applications/" 2>/dev/null || true
xdg-desktop-menu forceupdate 2>/dev/null || true

echo "AMDecrypt-gui is now installed (or updated)."
