package app

import (
	"flag"
	"fmt"
	"log/slog"
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui"
)

type App struct {
	Host        string
	Username    string
	Password    string
	connector   *network.Connector
	serverCache *cache.ServerCache
}

func NewApp() *App {
	return &App{serverCache: cache.NewServerCache()}
}

func (app *App) ReadArgs() error {
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		return fmt.Errorf("usage: <host> <username> <password>")
	}
	app.Host = args[0]
	app.Username = args[1]
	app.Password = args[2]
	return nil
}

func (app *App) Run() {
	var err error
	app.connector, err = network.NewConnector(app.Host)
	if err != nil {
		slog.Error("could not connect to host", "err", err)
		return
	}
	defer app.connector.Close()
	err = app.connector.Start(app.Username, app.Password)
	if err != nil {
		slog.Error("failed to connect", "err", err)
		return
	}

	go app.connector.Run(app.serverCache)

	app.connector.LoadChannels()

	window, err := ui.NewWindow()
	if err != nil {
		slog.Error("could not create window", "err", err)
		return
	}
	err = window.Create(app.connector, app.serverCache)
	if err != nil {
		slog.Error("could not create window", "err", err)
		return
	}
	defer window.Close()
	window.Run()
}
