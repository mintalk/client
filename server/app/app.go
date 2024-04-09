package app

import (
	"log/slog"
	"mintalk/server/config"
	"mintalk/server/db"
	"mintalk/server/network"
)

type App struct {
	config   *config.Config
	database *db.Connection
	server   *network.Server
}

func NewApp(config *config.Config) *App {
	return &App{config: config}
}

func (app *App) Init() error {
	var err error
	app.database, err = db.NewConnection(app.config)
	if err != nil {
		return err
	}
	err = app.database.Setup()
	if err != nil {
		return err
	}
	app.server = network.NewServer(app.database, app.config)
	return nil
}

func (app *App) Run() {
	err := app.server.Run()
	if err != nil {
		slog.Error("server failed", "err", err)
		return
	}
}
