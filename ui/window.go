package ui

import (
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui/panels"

	"github.com/rivo/tview"
)

type Window struct {
	*tview.Application
	channel  *panels.Channel
	channels *panels.Channels
	users *panels.Users
}

func NewWindow() *Window {
	window := &Window{Application: tview.NewApplication()}
	return window
}

func (window *Window) Create(connector *network.Connector, serverCache *cache.ServerCache) error {
	window.channel = panels.NewChannel(window.Application, connector, serverCache)
	window.channels = panels.NewChannels(window.Application, connector, window.channel, serverCache)
	window.users = panels.NewUsers(window.Application, serverCache)

	flex := tview.NewFlex()
	flex.AddItem(window.channels, 0, 1, false)
	flex.AddItem(window.channel, 0, 5, true)
	flex.AddItem(window.users, 0, 1, false)

	return window.SetRoot(flex, true).EnableMouse(true).Run()
}
