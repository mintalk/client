package app

import (
	"flag"
	"fmt"
	"log/slog"
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Host        string
	Username    string
	Password    string
	connector   *network.Connector
	serverCache *cache.ServerCache
	debug       bool
}

func NewApp() *App {
	return &App{serverCache: cache.NewServerCache()}
}

func (app *App) ReadArgs() error {
	flag.Usage = func() {
		fmt.Println("Usage: [flags] <host> <username> <password>")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.BoolVar(&app.debug, "d", false, "run in debug mode")

	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		return fmt.Errorf("failed to read arguments")
	}
	app.Host = args[0]
	app.Username = args[1]
	app.Password = args[2]
	return nil
}

func (app *App) Run() {
	level := slog.LevelInfo
	if app.debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

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

	window := ui.NewWindow()

	app.connector.CloseListener(window.Stop)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc
		window.Stop()
	}()

	err = window.Create(app.connector, app.serverCache)
	if err != nil {
		slog.Error("could not create window", "err", err)
		return
	}
}
