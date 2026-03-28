#!/usr/bin/env bash
set -e

BINARY="build/bin/amdecrypt-gui"
ICON="build/appicon.png"
APP_NAME="amdecrypt-gui"
INSTALL_BIN="/usr/local/bin/$APP_NAME"
INSTALL_ICON="$HOME/.local/share/icons/hicolor/256x256/apps/$APP_NAME.png"
INSTALL_DESKTOP="$HOME/.local/share/applications/$APP_NAME.desktop"

if [ ! -f "$BINARY" ]; then
    echo "Binary not found: $BINARY"
    echo "Run 'wails build' first."
    exit 1
fi

echo "Installing $APP_NAME..."

sudo install -Dm755 "$BINARY" "$INSTALL_BIN"

mkdir -p "$(dirname "$INSTALL_ICON")"
cp "$ICON" "$INSTALL_ICON"

mkdir -p "$(dirname "$INSTALL_DESKTOP")"
cat > "$INSTALL_DESKTOP" <<EOF
[Desktop Entry]
Name=AMDecrypt GUI
Comment=GUI for AppleMusicDecrypt
Exec=$INSTALL_BIN
Icon=$APP_NAME
Type=Application
Categories=Utility;AudioVideo;
EOF

update-desktop-database "$HOME/.local/share/applications/" 2>/dev/null || true

echo "AMDecrypt-gui is now installed (or updated)!"
