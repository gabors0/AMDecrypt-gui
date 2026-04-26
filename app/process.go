package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var errNoTerminal = errors.New("no terminal emulator found (tried gnome-terminal, konsole, xfce4-terminal, xterm)")

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
	runCmd := fmt.Sprintf("cd %q && %q main.py; echo; echo 'Exited, press Enter to close.'; read", amdDir, pythonBin)

	// Build terminal command: <terminal> <termArgs...> bash -c "<shellCmd>"
	args := append(termArgs, "bash", "-c", runCmd)
	cmd := exec.Command(termBin, args...)
	cmd.Dir = amdDir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	a.EmitLog(fmt.Sprintf("[DEBUG] Command: %s %v", termBin, args))

	if err := cmd.Start(); err != nil {
		a.EmitLog("[ERROR] Failed to start terminal: " + err.Error())
		return fmt.Errorf("failed to start terminal: %w", err)
	}

	a.amdCmd = cmd
	a.amdRunning = true
	a.amdDone = make(chan struct{})

	wailsRuntime.EventsEmit(a.ctx, "amd:started")
	a.EmitLog("[SUCCESS] AMD started in external terminal (" + termBin + ")")

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
	a.mu.Lock()
	if !a.amdRunning || a.amdCmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.amdDone
	pgid := a.amdCmd.Process.Pid
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
	if !a.amdRunning || a.amdCmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.amdDone
	pgid := a.amdCmd.Process.Pid
	a.mu.Unlock()

	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	<-done
	return nil
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
	a.mu.Lock()
	if !a.wmRunning || a.wmCmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.wmDone
	pgid := a.wmCmd.Process.Pid
	a.mu.Unlock()

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

func (a *App) KillWm() error {
	a.mu.Lock()
	if !a.wmRunning || a.wmCmd == nil {
		a.mu.Unlock()
		return nil
	}
	done := a.wmDone
	pgid := a.wmCmd.Process.Pid
	a.mu.Unlock()

	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	<-done
	return nil

}
