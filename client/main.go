package main

import (
	"mintalk/client/app"
)

func main() {
	app := app.NewApp()
	err := app.ReadArgs()
	if err != nil {
		return
	}
	app.Run()
}
