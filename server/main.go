package main

import (
	"log/slog"
	"mintalk/server/app"
	"mintalk/server/db"
)

func main() {
	config, err := app.LoadConfig("config.yaml")
	if err != nil {
		slog.Error("failed to load config", err)
		return
	}
	database, err := db.NewConnetion(config)
	if err != nil {
		slog.Error("failed to connect database", err)
		return
	}
	database.Setup()

	app := app.NewApp(config)
	app.Run()
}
