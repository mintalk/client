package app

import (
	"flag"
	"log/slog"
	"mintalk/client/network"
	"mintalk/client/ui"
)

type App struct {
	Host      string
	Username  string
	Password  string
	connector *network.Connector

	onlyUI bool // For testing purposes
}

func NewApp() *App {
	return &App{}
}

func (app *App) ReadArgs() {
	flag.BoolVar(&app.onlyUI, "u", false, "Run only the UI")
	flag.Parse()

	if app.onlyUI {
		return
	}
	args := flag.Args()
	app.Host = args[0]
	app.Username = args[1]
	app.Password = args[2]
}

func (app *App) Run() {
	if !app.onlyUI {
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
