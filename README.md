<div align=center>
    <img src="/build/appicon.png" height="128">
    <h2>AMDecrypt-gui</h2>
    <h4>A cross-platform app made to easily install and use <a href="https://github.com/WorldObservationLog/AppleMusicDecrypt">AppleMusicDecrypt</a>, currently in development</h4>
</div>

> [!IMPORTANT]
> This project **uses (but does not modify, bundle, or embed)** [AppleMusicDecrypt](https://github.com/WorldObservationLog/AppleMusicDecrypt) and [wrapper-manager](https://github.com/WorldObservationLog/wrapper-manager) both made by @WorldObservationLog and released under the **AGPL-3.0 license**. AMDecrypt-gui is released under the MIT license.

### Features
- Made with Wails
- Checks dependencies needed to run AppleMusicDecrypt/wrapper-manager
- Automatically installs and sets up (or removes) AppleMusicDecrpyt (clone, venv, pip install, etc.)
- Automatically installs and sets up (or removes) wrapper-manager (clone, docker compose up)
- Opens AppleMusicDecrypt in a seperate window with one click
- Command builder

### Todo
- Windows support
- Automatically update both modules
- Upload pre-compiled binaries

### Installation

<details>
<summary><b>Linux</b></summary>

- Have Go, Node/npm and Wails installed *(on Linux, libgtk-3-dev and libwebkit2gtk-4.1-dev are probably required to build.)*
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
- Easy install: run the bootstrap installer
```bash
curl -fsSL https://raw.githubusercontent.com/gabors0/AMDecrypt-gui/main/install.sh | sh
```
This clones the repo into a temporary directory, builds it, installs the app, and removes the build files afterward.

- Manual install: build/install from a local clone
```bash
git clone https://github.com/gabors0/AMDecrypt-gui.git
cd AMDecrypt-gui
wails build #-tags webkit2_41 for some linux systems
sudo ./scripts/install-built.sh
```

</details>

<details>
<summary><b>Uninstall</b></summary>

- Clone the repo if not done already
- Run the uninstall script inside the project folder
```bash
sudo ./uninstall.sh
```

</details>

### Platform Compatibility

| Platform | Supported?
|----------|--------|
| Linux (tested on Fedora & Arch) | ✔ |
| Windows  | Planned |
| macOS    | - |

<details>
    <summary><b>Note on macOS</b></summary>
    Builds, but AMD can't be started using the gui. Installation is also probably broken. I don't own a mac so this is untested, however if you stumble upon this and need mac support, feel free to open a PR!
</details>

### Screenshots
Main screen             |  Command builder
:-------------------------:|:-------------------------:
<img src="/frontend/src/assets/images/screenshot1.png" alt="main screen" width="350"> | <img src="/frontend/src/assets/images/screenshot2.png" alt="command builder" width="350">
