package app

import (
	"log/slog"
	"mintalk/client/ui"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	window, err := ui.NewWindow()
	if err != nil {
		slog.Error("could not create window", err)
		return
	}
	err = window.Create()
	if err != nil {
		slog.Error("could not create window", err)
		return
	}
	defer window.Close()
	window.Run()
}
