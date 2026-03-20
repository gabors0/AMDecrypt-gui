package app

import (
	"archive/zip"
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed amd_pip_requirements.txt
var pipRequirements []byte

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context) {
	fmt.Println("Welcome to AMDecrypt-gui!")
	a.ctx = ctx
}

func DomReady(a *App, ctx context.Context) {
	a.EmitLog("Welcome to AMDecrypt-gui!")
}

func (a *App) EmitLog(msg string) {
	wailsRuntime.EventsEmit(a.ctx, "log", msg)
}

func (a *App) WhichCmd(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return "Error: not found"
	}
	return path
}

func (a *App) RunCmd(command string) string {
	if strings.TrimSpace(command) == "" {
		return "Error: empty command"
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	hideWindow(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}

func (a *App) GetAppDataDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "amdecrypt-gui")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func (a *App) SetupAmd() {
	if runtime.GOOS == "windows" {
		a.EmitLog("[ERROR] Windows is not yet supported...")
		return
	}

	// ------------------- download
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Error getting app data dir: " + err.Error())
		return
	}

	a.EmitLog("[INFO] Attempting download!")
	resp, err := http.Get("https://github.com/WorldObservationLog/AppleMusicDecrypt/archive/refs/heads/v2.zip")
	if err != nil {
		a.EmitLog("[ERROR] Error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.EmitLog(fmt.Sprintf("[ERROR] unexpected status: %s", resp.Status))
		return
	}

	zipPath := filepath.Join(appDataDir, "amd.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		a.EmitLog("[ERROR] Error creating zip file: " + err.Error())
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		a.EmitLog("[ERROR] Error saving zip: " + err.Error())
		return
	}

	a.EmitLog("[SUCCESS] Downloaded to: " + zipPath)

	// ------------------- extract
	a.EmitLog("[INFO] Extracting zip...")
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		a.EmitLog("[ERROR] Error opening zip: " + err.Error())
		return
	}
	defer r.Close()

	for _, file := range r.File {
		// strip the top-level directory from the zip (e.g. "AppleMusicDecrypt-v2/")
		parts := strings.SplitN(file.Name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			continue
		}
		destPath := filepath.Join(appDataDir, "amd", filepath.FromSlash(parts[1]))

		// guard against zip-slip
		if !strings.HasPrefix(destPath, filepath.Clean(appDataDir)+string(os.PathSeparator)) {
			a.EmitLog("[ERROR] Zip-slip attempt blocked: " + file.Name)
			return
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				a.EmitLog("[ERROR] mkdir: " + err.Error())
				return
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			a.EmitLog("[ERROR] mkdir: " + err.Error())
			return
		}

		out, err := os.Create(destPath)
		if err != nil {
			a.EmitLog("[ERROR] create file: " + err.Error())
			return
		}

		rc, err := file.Open()
		if err != nil {
			out.Close()
			a.EmitLog("[ERROR] open zip entry: " + err.Error())
			return
		}

		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			a.EmitLog("[ERROR] extracting file: " + err.Error())
			return
		}
	}

	a.EmitLog("[SUCCESS] Extracted to: " + filepath.Join(appDataDir, "amd"))

	if err := os.Remove(zipPath); err != nil {
		a.EmitLog("[ERROR] Failed to delete zip: " + err.Error())
		return
	}
	a.EmitLog("[INFO] Deleted zip archive")

	// ------------------- create config
	exConfigPath := filepath.Join(appDataDir, "amd", "config.example.toml")
	configPath := filepath.Join(appDataDir, "amd", "config.toml")
	src, err := os.ReadFile(exConfigPath)
	if err != nil {
		a.EmitLog("[ERROR] Failed to read config.example.toml: " + err.Error())
		return
	}
	if err := os.WriteFile(configPath, src, 0644); err != nil {
		a.EmitLog("[ERROR] Failed to write config.toml: " + err.Error())
		return
	}
	a.EmitLog("[SUCCESS] Created config.toml from example")

	// ------------------- write requirements file
	reqPath := filepath.Join(appDataDir, "amd", "requirements.txt")
	if err := os.WriteFile(reqPath, pipRequirements, 0644); err != nil {
		a.EmitLog("[ERROR] Failed to write requirements.txt: " + err.Error())
		return
	}
	a.EmitLog("[INFO] Written requirements.txt")

	// ------------------- create venv
	a.EmitLog("[INFO] Creating Python venv...")
	venvPath := filepath.Join(appDataDir, "amd", "venv")
	out, err := exec.Command("python3", "-m", "venv", venvPath).CombinedOutput()
	if err != nil {
		a.EmitLog("[ERROR] Failed to create venv: " + err.Error() + "\n" + string(out))
		return
	}
	a.EmitLog("[SUCCESS] Venv created at: " + venvPath)

	// ------------------- install pip requirements
	a.EmitLog("[INFO] Installing pip requirements...")
	pipBin := filepath.Join(venvPath, "bin", "pip")
	out, err = exec.Command(pipBin, "install", "-r", reqPath).CombinedOutput()
	if err != nil {
		a.EmitLog("[ERROR] pip install failed: " + err.Error() + "\n" + string(out))
		return
	}
	a.EmitLog("[SUCCESS] Pip requirements installed")
	a.EmitLog("[SUCCESS] AMDecrypt setup complete!")
}

func (a *App) RemoveAmd() {
	// ------------------- remove amd dir
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Error getting app data dir: " + err.Error())
		return
	}

	amdDir := filepath.Join(appDataDir, "amd")
	if err := os.RemoveAll(amdDir); err != nil {
		a.EmitLog("[ERROR] Failed to remove amd dir: " + err.Error())
		return
	}
	a.EmitLog("[SUCCESS] Removed amd dir at: " + amdDir)
}
