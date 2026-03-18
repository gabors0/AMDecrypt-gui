package app

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context) {
	fmt.Println("App Startup")
	a.ctx = ctx
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
