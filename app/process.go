package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var errNoTerminal = errors.New("no terminal emulator found (tried gnome-terminal, konsole, xfce4-terminal, xterm)")

func (a *App) StartAmd() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		a.EmitLog("[WARNING] AMD is already running")
		return errors.New("AMD is already running")
	}

	a.EmitLog("[INFO] Starting AMD...")

	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		a.EmitLog("[ERROR] Failed to get app data dir: " + err.Error())
		return fmt.Errorf("failed to get app data dir: %w", err)
	}

	amdDir := filepath.Join(appDataDir, "amd")
	pythonBin := filepath.Join(amdDir, "venv", "bin", "python")

	if _, err := os.Stat(pythonBin); err != nil {
		a.EmitLog("[ERROR] Python venv not found at " + pythonBin)
		return fmt.Errorf("python not found at %s: %w", pythonBin, err)
	}

	settings, _ := a.GetSettings()
	termBin, termArgs, err := findTerminal(settings.Terminal)
	if err != nil {
		a.EmitLog("[ERROR] " + err.Error())
		return err
	}
	a.EmitLog("[INFO] Using terminal: " + termBin)

	// Build the shell command that runs inside the terminal.
	shellCmd := fmt.Sprintf("cd %q && %q main.py; echo; echo 'Exited, press Enter to close.'; read", amdDir, pythonBin)

	// Build terminal command: <terminal> <termArgs...> bash -c "<shellCmd>"
	args := append(termArgs, "bash", "-c", shellCmd)
	cmd := exec.Command(termBin, args...)
	cmd.Dir = amdDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	a.EmitLog(fmt.Sprintf("[DEBUG] Command: %s %v", termBin, args))

	if err := cmd.Start(); err != nil {
		a.EmitLog("[ERROR] Failed to start terminal: " + err.Error())
		return fmt.Errorf("failed to start terminal: %w", err)
	}

	a.cmd = cmd
	a.running = true
	a.done = make(chan struct{})

	wailsRuntime.EventsEmit(a.ctx, "amd:started")
	a.EmitLog("[SUCCESS] AMD started in external terminal (" + termBin + ")")

	go func() {
		waitErr := cmd.Wait()
		a.mu.Lock()
		a.running = false
		a.cmd = nil
		close(a.done)
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
	a.mu.Lock()
	if !a.running || a.cmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.done
	pgid := a.cmd.Process.Pid
	a.mu.Unlock()

	// Send SIGINT to the process group
	_ = syscall.Kill(-pgid, syscall.SIGINT)

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		// Force kill the process group
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
		<-done
		return nil
	}
}

func (a *App) KillAmd() error {
	a.mu.Lock()
	if !a.running || a.cmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.done
	pgid := a.cmd.Process.Pid
	a.mu.Unlock()

	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	<-done
	return nil
}

func (a *App) IsAmdRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.running
}
