package main

import (
	"log/slog"
	"mintalk/server/app"
	"mintalk/server/config"
)

func main() {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		return
	}

	app := app.NewApp(config)
	err = app.ReadArgs()
	if err != nil {
		slog.Error("failed to read args", "err", err)
		return
	}
	app.Init()
	app.Run()
}
