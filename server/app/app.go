package app

import (
	"flag"
	"fmt"
	"log/slog"
	"mintalk/server/config"
	"mintalk/server/db"
	"mintalk/server/input"
	"mintalk/server/network"
	"os"
)

type App struct {
	config   *config.Config
	database *db.Connection
	server   *network.Server
	console  *input.Console
	debug    bool
}

func NewApp(config *config.Config) *App {
	return &App{config: config}
}

func (app *App) ReadArgs() error {
	flag.Usage = func() {
		fmt.Println("Usage: [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.BoolVar(&app.debug, "d", false, "run in debug mode")

	flag.Parse()
	return nil
}

func (app *App) Init() error {
	level := slog.LevelInfo
	if app.debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

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
	app.console = input.NewConsole(app.database, app.server)
	return nil
}

func (app *App) Run() {
	go func() {
		err := app.server.Run()
		if err != nil {
			slog.Error("server failed", "err", err)
			return
		}
	}()
	app.console.InputLoop()
}
