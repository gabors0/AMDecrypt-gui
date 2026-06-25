//go:build windows

package app

func (a *App) DetectTerminal() string {
	return ""
}
