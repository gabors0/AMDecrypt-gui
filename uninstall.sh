#!/usr/bin/env bash
set -e

APP_NAME="amdecrypt-gui"

echo "Uninstalling $APP_NAME..."

sudo rm -f "/usr/local/bin/$APP_NAME"
rm -f "$HOME/.local/share/icons/hicolor/256x256/apps/$APP_NAME.png"
rm -f "$HOME/.local/share/applications/$APP_NAME.desktop"

update-desktop-database "$HOME/.local/share/applications/" 2>/dev/null || true

echo "AMDecrypt-gui is now uninstalled!"
