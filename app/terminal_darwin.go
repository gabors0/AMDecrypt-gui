//go:build darwin

package app

import "errors"

var errDarwinTerminal = errors.New("launching AMD in an external terminal is not yet implemented on macOS")

func findTerminal(_ string) (string, []string, error) {
	return "", nil, errDarwinTerminal
}

func (a *App) DetectTerminal() string {
	return ""
}
