package app

import (
	"log/slog"
	"mintalk/server/config"
	"mintalk/server/db"
	"mintalk/server/network"
	"net"
	"time"
)

type App struct {
	config         *config.Config
	database       *db.Connection
	sessionManager *network.SessionManager
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
	app.sessionManager = network.NewSessionManager(time.Duration(app.config.SessionLifetime)*time.Minute, 32)
	return nil
}

func (app *App) Run() {
	listener, err := net.Listen("tcp", app.config.Host)
	if err != nil {
		slog.Error("failed to create listener", err)
		return
	}
	defer listener.Close()
	slog.Info("listening on " + app.config.Host)
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Warn("failed to accept connection", err)
			continue
		}
		go app.handleClient(conn)
	}
}

func (app *App) handleClient(conn net.Conn) {
	defer conn.Close()
	executor := network.ProtocolExecutor{Conn: conn, Database: app.database, SessionManager: app.sessionManager}
	err := executor.Run()
	if err != nil {
		slog.Warn("connection failed", err)
		return
	}
}
