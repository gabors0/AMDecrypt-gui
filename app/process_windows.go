//go:build windows

package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) IsAmdRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.amdRunning
}

func (a *App) IsWmRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.wmRunning
}

func streamToLog(a *App, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		a.EmitLog(s.Text())
	}
}

func (a *App) StartAmd() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.amdRunning {
		a.EmitLog("[ERROR] AMD is already running")
		return errors.New("AMD is already running")
	}

	a.EmitLog("[INFO] Starting AMD...")

	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return fmt.Errorf("failed to get app data dir: %w", err)
	}

	amdDir := filepath.Join(appDataDir, "amd")
	pythonBin := venvPythonPath(filepath.Join(amdDir, "venv"))
	if _, err := os.Stat(pythonBin); err != nil {
		a.EmitLog("[ERROR] Python venv not found at " + pythonBin)
		return fmt.Errorf("python not found at %s: %w", pythonBin, err)
	}

	cmd := powershellAmdCommand(amdDir, pythonBin, "main.py", "Exited, press Enter to close.")
	cmd.Dir = amdDir

	if err := cmd.Start(); err != nil {
		a.EmitLog("[ERROR] Failed to start AMD: " + err.Error())
		return fmt.Errorf("failed to start AMD: %w", err)
	}

	a.amdCmd = cmd
	a.amdRunning = true
	a.amdDone = make(chan struct{})

	wailsRuntime.EventsEmit(a.ctx, "amd:started")
	a.EmitLog("[SUCCESS] AMD started in PowerShell")

	go func() {
		waitErr := cmd.Wait()
		a.mu.Lock()
		a.amdRunning = false
		a.amdCmd = nil
		close(a.amdDone)
		a.mu.Unlock()
		wailsRuntime.EventsEmit(a.ctx, "amd:stopped")
		if waitErr != nil {
			a.EmitLog("[INFO] AMD process exited: " + waitErr.Error())
		} else {
			a.EmitLog("[INFO] AMD process exited cleanly")
		}
	}()

	return nil
}

func (a *App) StopAmd() error {
	return a.stopWindowsProcess(&a.amdCmd, &a.amdRunning, a.amdDone, "AMD")
}

func (a *App) KillAmd() error {
	return a.killWindowsProcess(&a.amdCmd, &a.amdRunning, a.amdDone, "AMD")
}

func (a *App) LoginAmd() error {
	return a.runAmdTool("login.py", "login")
}

func (a *App) LogoutAmd() error {
	return a.runAmdTool("logout.py", "logout")
}

func (a *App) runAmdTool(script, label string) error {
	a.EmitLog("[INFO] Launching AMD " + label + "...")

	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return fmt.Errorf("failed to get app data dir: %w", err)
	}

	amdDir := filepath.Join(appDataDir, "amd")
	pythonBin := venvPythonPath(filepath.Join(amdDir, "venv"))
	scriptRel := filepath.Join("tools", script)

	if _, err := os.Stat(pythonBin); err != nil {
		a.EmitLog("[ERROR] Python venv not found at " + pythonBin)
		return fmt.Errorf("python not found at %s: %w", pythonBin, err)
	}
	if _, err := os.Stat(filepath.Join(amdDir, scriptRel)); err != nil {
		a.EmitLog("[ERROR] Script not found: " + scriptRel)
		return fmt.Errorf("script not found: %w", err)
	}

	cmd := powershellAmdCommand(amdDir, pythonBin, scriptRel, "Done, press Enter to close.")
	cmd.Dir = amdDir

	if err := cmd.Start(); err != nil {
		a.EmitLog("[ERROR] Failed to start AMD " + label + ": " + err.Error())
		return fmt.Errorf("failed to start AMD %s: %w", label, err)
	}

	go func() {
		waitErr := cmd.Wait()
		if waitErr != nil {
			a.EmitLog("[INFO] AMD " + label + " exited: " + waitErr.Error())
		} else {
			a.EmitLog("[SUCCESS] AMD " + label + " completed")
		}
	}()

	return nil
}

func powershellAmdCommand(workingDir string, pythonBin string, script string, closeMessage string) *exec.Cmd {
	command := fmt.Sprintf(
		"Set-Location -LiteralPath %s; & %s %s; Write-Host ''; Read-Host %s",
		powershellQuote(workingDir),
		powershellQuote(pythonBin),
		powershellQuote(script),
		powershellQuote(closeMessage),
	)
	return exec.Command("powershell", "-NoExit", "-ExecutionPolicy", "Bypass", "-Command", command)
}

func powershellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func (a *App) StartWm(verbose bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.wmRunning {
		a.EmitLog("[ERROR] wrapper-manager is already running")
		return errors.New("wrapper-manager is already running")
	}

	a.EmitLog("[INFO] Starting wrapper-manager...")

	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return fmt.Errorf("failed to get app data dir: %w", err)
	}

	wmDir := filepath.Join(appDataDir, "wrapper-manager")
	cmd := exec.Command("docker", "compose", "up")
	cmd.Dir = wmDir
	hideWindow(cmd)

	if verbose {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		go streamToLog(a, stdout)
		go streamToLog(a, stderr)
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if err := cmd.Start(); err != nil {
		a.EmitLog("[ERROR] Failed to start wrapper-manager: " + err.Error())
		return fmt.Errorf("failed to start wrapper-manager: %w", err)
	}

	a.wmCmd = cmd
	a.wmRunning = true
	a.wmDone = make(chan struct{})

	wailsRuntime.EventsEmit(a.ctx, "wm:started")
	a.EmitLog("[SUCCESS] wrapper-manager started")

	go func() {
		waitErr := cmd.Wait()
		a.mu.Lock()
		a.wmRunning = false
		a.wmCmd = nil
		close(a.wmDone)
		a.mu.Unlock()
		wailsRuntime.EventsEmit(a.ctx, "wm:stopped")
		if waitErr != nil {
			a.EmitLog("[INFO] wrapper-manager process exited: " + waitErr.Error())
		} else {
			a.EmitLog("[INFO] wrapper-manager process exited cleanly")
		}
	}()

	return nil
}

func (a *App) StopWm() error {
	return a.stopWindowsProcess(&a.wmCmd, &a.wmRunning, a.wmDone, "wrapper-manager")
}

func (a *App) KillWm() error {
	return a.killWindowsProcess(&a.wmCmd, &a.wmRunning, a.wmDone, "wrapper-manager")
}

func (a *App) stopWindowsProcess(cmdRef **exec.Cmd, running *bool, done chan struct{}, label string) error {
	a.mu.Lock()
	if !*running || *cmdRef == nil {
		a.mu.Unlock()
		return nil
	}
	pid := (*cmdRef).Process.Pid
	a.mu.Unlock()

	_ = exec.Command("taskkill", "/PID", fmt.Sprintf("%d", pid), "/T").Run()

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		a.EmitLog("[WARN] " + label + " did not stop gracefully, forcing termination")
		return a.killWindowsPid(pid, done)
	}
}

func (a *App) killWindowsProcess(cmdRef **exec.Cmd, running *bool, done chan struct{}, _ string) error {
	a.mu.Lock()
	if !*running || *cmdRef == nil {
		a.mu.Unlock()
		return nil
	}
	pid := (*cmdRef).Process.Pid
	a.mu.Unlock()

	return a.killWindowsPid(pid, done)
}

func (a *App) killWindowsPid(pid int, done chan struct{}) error {
	_ = exec.Command("taskkill", "/PID", fmt.Sprintf("%d", pid), "/T", "/F").Run()
	<-done
	return nil
}
