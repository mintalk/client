package ui

import (
	"log/slog"
	"mintalk/client/ui/panels"

	gc "github.com/rthornton128/goncurses"
)

type Window struct {
	*gc.Window
	layout   *Layout
	channel  *panels.ChannelPanel
	channels *panels.ChannelsPanel
	State    *UIState
	Input    *Input
	running  bool
}

func NewWindow() (*Window, error) {
	ncursesWindow, err := gc.Init()
	if err != nil {
		return nil, err
	}
	err = gc.StartColor()
	if err != nil {
		return nil, err
	}
	err = gc.InitPair(1, gc.C_CYAN, 0)
	if err != nil {
		return nil, err
	}
	window := &Window{ncursesWindow, nil, nil, nil, NewUIState(), nil, false}
	return window, nil
}

func (window *Window) Create() error {
	window.Keypad(true)
	gc.Echo(false)
	gc.CBreak(true)
	window.Timeout(0)
	gc.Cursor(0)

	var err error
	window.channel, err = panels.NewChannelPanel()
	if err != nil {
		return err
	}
	window.channels, err = panels.NewChannelsPanel()
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

	window.Input = NewInput(20)

	window.Resize(window.MaxYX())
	return nil
}

func (window *Window) Close() {
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
	window.Input.Active = true
	char := window.GetChar()
	if char == 0 {
		return
	}
	window.Input.Update(char)
	switch char {
	case 'q':
		if window.State.Mode == ModeNormal {
			window.running = false
			return
		}
	case '\n':
		window.State.Mode = ModeInsert
	case gc.KEY_TAB:
		if window.State.ActiveTab == TabChannels {
			window.State.ActiveTab = TabChannel
		} else {
			window.State.ActiveTab = TabChannels
		}
	case gc.KEY_ESC:
		window.State.Mode = ModeNormal
	}
}

func (window *Window) Draw() error {
	if err := window.channel.Draw(window.State.ActiveTab == TabChannel); err != nil {
		return err
	}
	if err := window.channels.Draw(window.State.ActiveTab == TabChannels); err != nil {
		return err
	}

	gc.UpdatePanels()
	err := gc.Update()
	if err != nil {
		return err
	}

	window.Input.Draw(window.channel.Window(), 1, 1)

	return nil
}

func (window *Window) Resize(height, width int) {
	window.layout.Update(width, height, 0, 0)
}
