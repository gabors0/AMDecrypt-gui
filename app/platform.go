package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func venvPythonPath(venvPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venvPath, "Scripts", "python.exe")
	}
	return filepath.Join(venvPath, "bin", "python")
}

func venvPipPath(venvPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venvPath, "Scripts", "pip.exe")
	}
	return filepath.Join(venvPath, "bin", "pip")
}

func pythonVenvCommand(venvPath string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("py"); err == nil {
			return exec.Command("py", "-3", "-m", "venv", venvPath)
		}
		return exec.Command("python", "-m", "venv", venvPath)
	}

	if _, err := exec.LookPath("python3"); err == nil {
		return exec.Command("python3", "-m", "venv", venvPath)
	}
	return exec.Command("python", "-m", "venv", venvPath)
}

func (a *App) IsAmdInstalled() bool {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return false
	}
	pythonBin := venvPythonPath(filepath.Join(appDataDir, "amd", "venv"))
	if _, err := os.Stat(pythonBin); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(appDataDir, "amd", "main.py")); err != nil {
		return false
	}
	return true
}

func (a *App) IsWmInstalled() bool {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(appDataDir, "wrapper-manager", "docker-compose.yml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(appDataDir, "wrapper-manager", "compose.yml")); err == nil {
		return true
	}
	return false
}
