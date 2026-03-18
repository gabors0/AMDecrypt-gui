package app

import (
	"context"
	"fmt"
	"os/exec"
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
    parts := strings.Fields(command)
    if len(parts) == 0 {
        return "Error: empty command"
    }
    cmd := exec.Command(parts[0], parts[1:]...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "Error: " + err.Error()
    }
    return strings.TrimSpace(string(output))
}
