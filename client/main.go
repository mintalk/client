package main

import "mintalk/client/app"

func main() {
	app := app.NewApp()
	app.ReadArgs()
	app.Run()
}
