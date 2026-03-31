package main

import (
	"context"
	"embed"
	"amdecrypt-gui/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

var appobj *app.App

func main() {
	// Create an instance of the app structure
	appobj = app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "AMDecrypt-gui",
		Width:     1024,
		Height:    768,
		MinWidth:  450,
		MinHeight: 650,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        startup,
		OnDomReady:       domReady,
		OnShutdown:       shutdown,
		Linux: &linux.Options{
			ProgramName: "amdecrypt-gui",
		},
		Bind: []interface{}{
			appobj,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func startup(ctx context.Context) {
	app.Startup(appobj, ctx)
}

func domReady(ctx context.Context) {
	app.DomReady(appobj, ctx)
}

func shutdown(ctx context.Context) {
	appobj.StopAmd()
}
