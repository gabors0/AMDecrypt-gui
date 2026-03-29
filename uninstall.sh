#!/usr/bin/env bash
set -e

APP_NAME="amdecrypt-gui"

if [ -n "${SUDO_USER:-}" ] && [ "$SUDO_USER" != "root" ]; then
    TARGET_USER="$SUDO_USER"
else
    TARGET_USER="$(id -un)"
fi

TARGET_HOME="$(getent passwd "$TARGET_USER" | cut -d: -f6)"
if [ -z "$TARGET_HOME" ]; then
    echo "Failed to resolve home directory for user: $TARGET_USER"
    exit 1
fi

echo "Uninstalling $APP_NAME..."
echo "Target user: $TARGET_USER"

if [ "$EUID" -eq 0 ]; then
    rm -f "/usr/local/bin/$APP_NAME"
else
    sudo rm -f "/usr/local/bin/$APP_NAME"
fi

rm -f "$TARGET_HOME/.local/share/icons/hicolor/256x256/apps/$APP_NAME.png"
rm -f "$TARGET_HOME/.local/share/applications/$APP_NAME.desktop"

gtk-update-icon-cache "$TARGET_HOME/.local/share/icons/hicolor/" 2>/dev/null || true
update-desktop-database "$TARGET_HOME/.local/share/applications/" 2>/dev/null || true
xdg-desktop-menu forceupdate 2>/dev/null || true

echo "AMDecrypt-gui is now uninstalled!"
