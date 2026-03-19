package app

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

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
