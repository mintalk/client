package app

import (
	"log/slog"
	"mintalk/server/secure"
	"net"
)

type App struct {
	config *Config
}

func NewApp(config *Config) *App {
	return &App{config: config}
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

	message, err := secure.Recieve3Pass(conn)
	if err != nil {
		slog.Warn("failed to receive message", err)
		return
	}

	slog.Info(string(message))
}
