package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) StartAmd() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return fmt.Errorf("AMD is already running")
	}

	appDataDir, err := a.GetAppDataDir()
	if err != nil {
		return fmt.Errorf("failed to get app data dir: %w", err)
	}

	var pythonPath string
	if runtime.GOOS == "windows" {
		pythonPath = filepath.Join(appDataDir, "amd", "venv", "Scripts", "python.exe")
	} else {
		pythonPath = filepath.Join(appDataDir, "amd", "venv", "bin", "python")
	}

	if _, err := os.Stat(pythonPath); err != nil {
		return fmt.Errorf("python not found at %s: %w", pythonPath, err)
	}

	amdDir := filepath.Join(appDataDir, "amd")
	mainPy := filepath.Join(amdDir, "main.py")

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, pythonPath, mainPy)
	cmd.Dir = amdDir
	cmd.Env = append(os.Environ(),
		"TERM=dumb",
		"PYTHONUNBUFFERED=1",
		"LOGURU_COLORIZE=false",
	)
	hideWindow(cmd)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start AMD: %w", err)
	}

	a.cmd = cmd
	a.stdin = stdin
	a.cancel = cancel
	a.running = true
	a.done = make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			wailsRuntime.EventsEmit(a.ctx, "amd:stdout", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			wailsRuntime.EventsEmit(a.ctx, "amd:stderr", scanner.Text())
		}
	}()

	go func() {
		_ = cmd.Wait()
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		wailsRuntime.EventsEmit(a.ctx, "amd:stopped")
		close(a.done)
	}()

	wailsRuntime.EventsEmit(a.ctx, "amd:started")
	return nil
}

func (a *App) StopAmd() error {
	a.mu.Lock()

	if !a.running {
		a.mu.Unlock()
		return nil
	}

	if a.stdin != nil {
		_, _ = a.stdin.Write([]byte("exit\n"))
	}

	cancel := a.cancel
	done := a.done
	a.mu.Unlock()

	// Wait for the single wait-goroutine in StartAmd to finish cleanup
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		if cancel != nil {
			cancel()
		}
		// Wait again briefly for the goroutine to finish after cancel
		<-done
	}

	return nil
}

func (a *App) SendInput(text string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running || a.stdin == nil {
		return fmt.Errorf("AMD is not running")
	}

	_, err := a.stdin.Write([]byte(text + "\n"))
	return err
}

func (a *App) KillAmd() error {
	a.mu.Lock()

	if !a.running {
		a.mu.Unlock()
		return nil
	}

	cancel := a.cancel
	done := a.done
	a.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	// Wait for the wait-goroutine to finish cleanup so the frontend gets amd:stopped
	<-done

	return nil
}

func (a *App) IsAmdRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.running
}
