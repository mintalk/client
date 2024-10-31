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
}

func NewWindow() *Window {
	window := &Window{Application: tview.NewApplication()}
	return window
}

func (window *Window) Create(connector *network.Connector, serverCache *cache.ServerCache) error {
	window.channel = panels.NewChannel(window.Application, connector, serverCache)
	window.channels = panels.NewChannels(window.Application, connector, window.channel, serverCache)

	flex := tview.NewFlex()
	flex.AddItem(window.channels, 0, 1, false)
	flex.AddItem(window.channel, 0, 4, true)

	return window.SetRoot(flex, true).EnableMouse(true).Run()
}
