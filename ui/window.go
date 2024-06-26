package ui

import (
	"log/slog"
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui/panels"

	gc "github.com/mintalk/goncurses"
)

type Window struct {
	*gc.Window
	layout   *Layout
	channel  *panels.ChannelPanel
	channels *panels.ChannelsPanel
	State    *UIState
	running  bool
}

func NewWindow() (*Window, error) {
	stdscr, err := gc.Init()
	if err != nil {
		return nil, err
	}
	window := &Window{stdscr, nil, nil, nil, NewUIState(), false}
	window.State.ActiveTab = TabChannel
	return window, nil
}

func (window *Window) Create(connector *network.Connector, serverCache *cache.ServerCache) error {
	window.Keypad(true)
	gc.Echo(false)
	gc.CBreak(true)
	window.Timeout(0)
	gc.Cursor(0)

	InitColors()

	var err error
	window.channel, err = panels.NewChannelPanel(connector, serverCache)
	if err != nil {
		return err
	}
	window.channels, err = panels.NewChannelsPanel(connector, window.channel, serverCache)
	if err != nil {
		return err
	}

	window.layout = &Layout{
		Panel: window.channels.Panel,
		Child: &Layout{
			Panel: window.channel.Panel,
		},
		Direction: Horizontal,
	}

	window.Resize(window.MaxYX())
	return nil
}

func (window *Window) Close() {
	window.running = false
	gc.Cursor(1)
	gc.End()
}

func (window *Window) Run() {
	window.running = true
	for window.running {
		window.Update()
		err := window.Draw()
		if err != nil {
			slog.Error("error drawing window", err)
		}
	}
}

func (window *Window) Update() {
	window.channel.Active = window.State.ActiveTab == TabChannel
	window.channels.Active = window.State.ActiveTab == TabChannels

	char := window.GetChar()
	if char == gc.KEY_ENTER || char == gc.KEY_RETURN {
		window.State.ActiveTab = TabChannel
	} else if char == gc.KEY_TAB {
		if window.State.ActiveTab == TabChannels {
			window.State.ActiveTab = TabChannel
		} else {
			window.State.ActiveTab = TabChannels
		}
		return
	} else if char == gc.KEY_RESIZE {
		window.Resize(window.MaxYX())
	}
	window.channel.Update(char)
	window.channels.Update(char)
}

func (window *Window) Draw() error {
	window.channel.Active = window.State.ActiveTab == TabChannel
	if err := window.channel.Draw(window.Window); err != nil {
		return err
	}
	window.channels.Active = window.State.ActiveTab == TabChannels
	if err := window.channels.Draw(window.Window); err != nil {
		return err
	}

	gc.UpdatePanels()
	err := gc.Update()
	if err != nil {
		return err
	}

	window.Refresh()
	return nil
}

func (window *Window) Resize(height, width int) {
	window.layout.Update(width, height, 0, 0)
	window.channel.Resize()
	window.channels.Resize()
}
