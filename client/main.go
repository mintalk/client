package main

import (
	"log/slog"
	"mintalk/client/app"
)

func main() {
	app := app.NewApp()
	err := app.ReadArgs()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	app.Run()
}
