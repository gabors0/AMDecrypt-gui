//go:build !windows && !darwin

package app

import (
	"os"
	"os/exec"
)

type terminalInfo struct {
	bin  string
	args []string // args before the shell command
}

// knownTerminals lists terminals with their correct exec-arg syntax.
var knownTerminals = []terminalInfo{
	{"gnome-terminal", []string{"--"}},
	{"kgx", []string{"--"}},           // GNOME Console (Fedora default)
	{"ptyxis", []string{"--"}},        // Ptyxis (newer Fedora/GNOME)
	{"konsole", []string{"-e"}},
	{"xfce4-terminal", []string{"-e"}},
	{"xterm", []string{"-e"}},
	{"alacritty", []string{"-e"}},
	{"kitty", []string{"--"}},
	{"wezterm", []string{"start", "--"}},
}

// argsForTerminal returns the exec-args for a terminal binary.
// For unknown terminals, falls back to "-e".
func argsForTerminal(bin string) []string {
	for _, t := range knownTerminals {
		if t.bin == bin {
			return t.args
		}
	}
	return []string{"-e"}
}

// DetectTerminal auto-detects the best available terminal (no saved override).
// Returns the binary name, or "" if none found.
func DetectTerminal() string {
	// 1. $TERMINAL env var
	if t := os.Getenv("TERMINAL"); t != "" {
		if _, err := exec.LookPath(t); err == nil {
			return t
		}
	}
	// 2. Debian/Ubuntu alternatives
	if _, err := exec.LookPath("x-terminal-emulator"); err == nil {
		return "x-terminal-emulator"
	}
	// 3. XDG standard (newer systems)
	if _, err := exec.LookPath("xdg-terminal-exec"); err == nil {
		return "xdg-terminal-exec"
	}
	// 4. Hardcoded list
	for _, t := range knownTerminals {
		if _, err := exec.LookPath(t.bin); err == nil {
			return t.bin
		}
	}
	return ""
}

// findTerminal returns the terminal binary and its exec-args.
// If override is non-empty and on PATH, it is used first.
func findTerminal(override string) (bin string, args []string, err error) {
	if override != "" {
		if _, err := exec.LookPath(override); err == nil {
			return override, argsForTerminal(override), nil
		}
	}
	bin = DetectTerminal()
	if bin == "" {
		return "", nil, errNoTerminal
	}
	return bin, argsForTerminal(bin), nil
}

// DetectTerminalRPC is the Wails-exposed version of DetectTerminal.
func (a *App) DetectTerminal() string {
	return DetectTerminal()
}
