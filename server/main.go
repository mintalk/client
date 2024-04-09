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
	app.Init()
	app.Run()
}
