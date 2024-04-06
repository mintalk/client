package app

import (
	"log/slog"
	"mintalk/client/secure"
	"net"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		slog.Error("could not connect to server", err)
	}
	defer conn.Close()
	prime, err := secure.RandomPrime(1024)
	if err != nil {
		slog.Error("could not generate prime", err)
	}
	err = secure.Send3Pass(conn, []byte("hello"), prime)
	if err != nil {
		slog.Error("could not send message", err)
	}
}
