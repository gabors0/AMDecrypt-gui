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
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed amd_pip_requirements.txt
var pipRequirements []byte

const Version = "0.1.4"

type App struct {
	ctx        context.Context
	mu         sync.Mutex
	amdCmd     *exec.Cmd
	amdRunning bool
	amdDone    chan struct{}
	wmCmd      *exec.Cmd
	wmRunning  bool
	wmDone     chan struct{}
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context) {
	fmt.Println("Welcome to AMDecrypt-gui!")
	a.ctx = ctx
}

func DomReady(a *App, ctx context.Context) {
	a.EmitLog("[INFO] Welcome to AMDecrypt-gui!")
}

func (a *App) OpenAppDataDir() error {
	dir, err := a.GetAppDataDir()
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	return cmd.Start()
}

func (a *App) OpenDownloadsDir() error {
	dir, err := a.GetAppDataDir()
	if err != nil {
		return err
	}
	downloadsDir := filepath.Join(dir, "amd", "downloads")
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", downloadsDir)
	case "darwin":
		cmd = exec.Command("open", downloadsDir)
	default:
		cmd = exec.Command("xdg-open", downloadsDir)
	}
	return cmd.Start()
}

func (a *App) EmitLog(msg string) {
	wailsRuntime.EventsEmit(a.ctx, "log", msg)
}

func (a *App) GetVersion() string {
	return Version
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

func (a *App) SetupBento4() {
	if runtime.GOOS == "windows" {
		a.EmitLog("[ERROR] Installing bento4 automatically isn't available on Windows!")
		return
	}

	a.EmitLog("[INFO] Installing Bento4 toolkits...")
	tmpDir := os.TempDir()
	repoDir := filepath.Join(tmpDir, "Bento4")
	buildDir := filepath.Join(repoDir, "cmakebuild")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to resolve home directory: " + err.Error())
		return
	}

	installPrefix := strings.TrimSpace(os.Getenv("PREFIX"))
	if installPrefix == "" {
		installPrefix = filepath.Join(homeDir, ".local")
	}

	emitOutput := func(out []byte) {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			a.EmitLog(line)
		}
	}

	runStep := func(label string, cmd *exec.Cmd) bool {
		a.EmitLog("[INFO] " + label)
		hideWindow(cmd)
		out, err := cmd.CombinedOutput()
		emitOutput(out)
		if err != nil {
			a.EmitLog("[ERROR] " + label + " failed: " + err.Error())
			return false
		}
		return true
	}

	if err := os.RemoveAll(repoDir); err != nil {
		a.EmitLog("[ERROR] Failed to clean old Bento4 dir: " + err.Error())
		return
	}

	if !runStep("Cloning Bento4 repository...", exec.Command("git", "clone", "--depth=1", "https://github.com/axiomatic-systems/Bento4.git", repoDir)) {
		return
	}

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		a.EmitLog("[ERROR] Failed to create cmakebuild dir: " + err.Error())
		return
	}

	cmakeCmd := exec.Command(
		"cmake",
		"-DCMAKE_BUILD_TYPE=Release",
		"-DCMAKE_INSTALL_PREFIX="+installPrefix,
		"..",
	)
	cmakeCmd.Dir = buildDir
	if !runStep("Configuring Bento4 with CMake", cmakeCmd) {
		return
	}

	jobs := runtime.NumCPU()
	if jobs < 1 {
		jobs = 1
	}
	makeCmd := exec.Command("make", fmt.Sprintf("-j%d", jobs))
	makeCmd.Dir = buildDir
	if !runStep(fmt.Sprintf("Building Bento4 with %d parallel jobs", jobs), makeCmd) {
		return
	}

	a.EmitLog("[INFO] Installing Bento4 to user prefix: " + installPrefix)
	installCmd := exec.Command("make", "install")
	installCmd.Dir = buildDir
	if !runStep("Installing Bento4", installCmd) {
		return
	}

	bentoPath := filepath.Join(installPrefix, "bin", "mp4decrypt")
	if _, err := os.Stat(bentoPath); err != nil {
		pathLookup, lookErr := exec.LookPath("mp4decrypt")
		if lookErr != nil {
			a.EmitLog("[ERROR] Bento4 installed but mp4decrypt was not found at " + bentoPath + " or on PATH")
			return
		}
		bentoPath = pathLookup
	}

	bentoDir := filepath.Dir(bentoPath)
	if err := a.updateBento4Settings(bentoPath, bentoDir); err != nil {
		a.EmitLog("[WARN] Bento4 installed, but failed to update settings.jsonc: " + err.Error())
	} else {
		a.EmitLog("[INFO] Updated Bento4 uninstall settings in settings.jsonc")
	}

	if err := os.RemoveAll(repoDir); err != nil {
		a.EmitLog("[WARN] Bento4 installed, but cleanup failed: " + err.Error())
		return
	}

	a.EmitLog("[SUCCESS] Bento4 installed successfully!")
}

func (a *App) RemoveBento4() {
	if runtime.GOOS == "windows" {
		a.EmitLog("[ERROR] Removing bento4 automatically isn't available on Windows!")
		return
	}

	settings, err := a.GetSettings()
	if err != nil {
		a.EmitLog("[ERROR] Failed to read settings.jsonc: " + err.Error())
		return
	}

	managedPath := strings.TrimSpace(settings.Bento4.Mp4decryptPath)
	managedDir := strings.TrimSpace(settings.Bento4.BinDir)
	if managedPath == "" || managedDir == "" {
		a.EmitLog("[INFO] Bento4 uninstall skipped: no app-managed install path in settings.jsonc")
		return
	}

	currentPath, err := exec.LookPath("mp4decrypt")
	if err != nil {
		a.EmitLog("[INFO] mp4decrypt is not currently in PATH, nothing to remove")
		return
	}

	if filepath.Clean(currentPath) != filepath.Clean(managedPath) {
		a.EmitLog("[INFO] Bento4 uninstall skipped: only app-installed Bento4 can be removed")
		a.EmitLog("[INFO] managed=" + managedPath)
		a.EmitLog("[INFO] current=" + currentPath)
		return
	}

	binaries := []string{
		"mp4compact", "mp4dash", "mp4dashclone", "mp4dcfpackager", "mp4decrypt",
		"mp4dump", "mp4edit", "mp4encrypt", "mp4extract", "mp4fragment",
		"mp4hls", "mp4iframeindex", "mp4info", "mp4rtphint",
	}

	removed := 0
	for _, bin := range binaries {
		binPath := filepath.Join(managedDir, bin)
		if err := os.Remove(binPath); err != nil {
			if !os.IsNotExist(err) {
				a.EmitLog("[WARN] Failed to remove " + binPath + ": " + err.Error())
			}
			continue
		}
		removed++
		a.EmitLog("[INFO] Removed " + binPath)
	}

	if removed == 0 {
		a.EmitLog("[INFO] No managed Bento4 binaries were removed")
		return
	}

	if err := a.updateBento4Settings("", ""); err != nil {
		a.EmitLog("[WARN] Removed Bento4 binaries, but failed to update settings.jsonc: " + err.Error())
	} else {
		a.EmitLog("[INFO] Cleared Bento4 uninstall paths in settings.jsonc")
	}

	a.EmitLog("[SUCCESS] Removed managed Bento4 binaries!")
}

func (a *App) updateBento4Settings(mp4decryptPath string, binDir string) error {
	settings, err := a.GetSettings()
	if err != nil {
		return err
	}

	settings.Bento4.Mp4decryptPath = mp4decryptPath
	settings.Bento4.BinDir = binDir

	return a.SaveSettings(settings)
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
	a.EmitLog("[SUCCESS] Pip requirements installed!")
	a.EmitLog("[SUCCESS] AMDecrypt setup complete!")
}

type InstanceConfig struct {
	URL    string `json:"url"`
	Secure bool   `json:"secure"`
}

func (a *App) GetInstanceConfig() (*InstanceConfig, error) {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(appDataDir, "amd", "config.toml"))
	if err != nil {
		a.EmitLog("[ERROR] Failed to read config.toml: " + err.Error())
		return nil, err
	}
	cfg := &InstanceConfig{}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "url") && !strings.HasPrefix(trimmed, "#") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				cfg.URL = strings.Trim(strings.TrimSpace(parts[1]), "\"")
				break
			}
		}
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "secure") && !strings.HasPrefix(trimmed, "#") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				cfg.Secure = strings.TrimSpace(parts[1]) == "true"
				break
			}
		}
	}
	a.EmitLog(fmt.Sprintf("[INFO] Read instance config: url=%s, secure=%v", cfg.URL, cfg.Secure))
	return cfg, nil
}

func (a *App) SetInstanceConfig(url string, secure bool) error {
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return err
	}
	configPath := filepath.Join(appDataDir, "amd", "config.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		a.EmitLog("[ERROR] Failed to read config.toml: " + err.Error())
		return err
	}
	lines := strings.Split(string(data), "\n")
	secureStr := "false"
	if secure {
		secureStr = "true"
	}
	urlDone, secureDone := false, false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !urlDone && strings.HasPrefix(trimmed, "url") && !strings.HasPrefix(trimmed, "#") {
			lines[i] = fmt.Sprintf("url = \"%s\"", url)
			urlDone = true
		}
		if !secureDone && strings.HasPrefix(trimmed, "secure") && !strings.HasPrefix(trimmed, "#") {
			lines[i] = fmt.Sprintf("secure = %s", secureStr)
			secureDone = true
		}
		if urlDone && secureDone {
			break
		}
	}
	if err := os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		a.EmitLog("[ERROR] Failed to write config.toml: " + err.Error())
		return err
	}
	a.EmitLog(fmt.Sprintf("[INFO] Updated instance config: url=%s, secure=%v", url, secure))
	return nil
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

func (a *App) SetupWm() {
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
	resp, err := http.Get("https://github.com/WorldObservationLog/wrapper-manager/archive/refs/heads/main.zip")
	if err != nil {
		a.EmitLog("[ERROR] Error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.EmitLog(fmt.Sprintf("[ERROR] unexpected status: %s", resp.Status))
		return
	}

	zipPath := filepath.Join(appDataDir, "wrapper-manager.zip")
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
		// strip the top-level directory from the zip
		parts := strings.SplitN(file.Name, "/", 2)
		if len(parts) < 2 || parts[1] == "" {
			continue
		}
		destPath := filepath.Join(appDataDir, "wrapper-manager", filepath.FromSlash(parts[1]))

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

	a.EmitLog("[SUCCESS] Extracted to: " + filepath.Join(appDataDir, "wrapper-manager"))

	if err := os.Remove(zipPath); err != nil {
		a.EmitLog("[ERROR] Failed to delete zip: " + err.Error())
		return
	}
	a.EmitLog("[INFO] Deleted zip archive")
}

func (a *App) RemoveWm() {
	// ------------------- remove wrapper-manager dir
	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Error getting app data dir: " + err.Error())
		return
	}

	wmDir := filepath.Join(appDataDir, "wrapper-manager")
	if err := os.RemoveAll(wmDir); err != nil {
		a.EmitLog("[ERROR] Failed to remove wrapper-manager dir: " + err.Error())
		return
	}
	a.EmitLog("[SUCCESS] Removed wrapper-manager dir at: " + wmDir)
}
