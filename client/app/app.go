package app

import (
	"log/slog"
	"mintalk/client/network"
	"mintalk/client/ui"
	"os"
)

type App struct {
	Host      string
	Username  string
	Password  string
	connector *network.Connector
}

func NewApp() *App {
	return &App{}
}

func (app *App) ReadArgs() {
	app.Host = os.Args[1]
	app.Username = os.Args[2]
	app.Password = os.Args[3]
}

func (app *App) Run() {
	var err error
	app.connector, err = network.NewConnector(app.Host)
	if err != nil {
		slog.Error("could not connect to host", err)
		return
	}
	defer app.connector.Close()
	err = app.connector.Auth(app.Username, app.Password)
	if err != nil {
		slog.Error("failed to authenticate", err)
	}

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
