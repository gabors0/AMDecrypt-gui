<div align=center>
    <img src="/build/appicon.png" height="128">
    <h2>AMDecrypt-gui</h2>
    <h4>A cross-platform app made to easily install and use <a href="https://github.com/WorldObservationLog/AppleMusicDecrypt">AppleMusicDecrypt</a>, currently in development</h4>
</div>

> [!NOTE]
> This project **uses (but does not modify, bundle, or embed)** [AppleMusicDecrypt](https://github.com/WorldObservationLog/AppleMusicDecrypt) and [wrapper-manager](https://github.com/WorldObservationLog/wrapper-manager) both made by @WorldObservationLog and released under the **AGPL-3.0 license**. AMDecrypt-gui is released under the MIT license.

### Features
- Made with Wails
- Checks dependencies needed to run AppleMusicDecrypt/wrapper-manager
- Automatically installs and sets up (or removes) AppleMusicDecrpyt (clone, venv, pip install, etc.)
- Opens AppleMusicDecrypt in a seperate window with one click
- Command builder

### Todo
- Windows support
- Ability to manage wrapper manager too
- Automatically update both modules
- Upload pre-compiled binaries

### Installation
1. Have Go and Wails installed *(on Linux, libgtk-3-dev and libwebkit2gtk-4.1-dev are probably required to build.)*
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
2. Clone the repo
```bash
git clone https://github.com/gabors0/AMDecrypt-gui.git
cd AMDecrypt-gui
```
3. Build
```bash
wails build #-tags webkit2_41 for some linux systems
```
4. Install using install.sh (or run directly from build/bin/AMDecrypt-gui)
```bash
sudo ./install.sh
```

### Platform Compatibility

| Platform | Supported?
|----------|--------|
| Linux (tested on Fedora & Arch) | ✔ |
| Windows  | - |
| macOS    | - |

### Screenshots
Main screen             |  Command builder
:-------------------------:|:-------------------------:
<img src="/frontend/src/assets/images/screenshot1.png" alt="main screen" width="350"> | <img src="/frontend/src/assets/images/screenshot2.png" alt="command builder" width="350">
