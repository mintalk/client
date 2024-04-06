package ui

import (
	gc "github.com/rthornton128/goncurses"
)

type Window struct {
	*gc.Window
	layout   *Layout
	channel  *ChannelPanel
	channels *ChannelsPanel
	State    *UIState
}

func NewWindow() (*Window, error) {
	ncursesWindow, err := gc.Init()
	if err != nil {
		return nil, err
	}
	window := &Window{ncursesWindow, nil, nil, nil, NewUIState()}
	return window, nil
}

func (window *Window) Create() error {
	window.Keypad(true)

	var err error
	window.channel, err = NewChannelPanel()
	if err != nil {
		return err
	}
	window.channels, err = NewChannelsPanel()
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
	gc.End()
}

func (window *Window) Run() {
	for {
		window.Draw()
		gc.Echo(window.State.Mode == ModeInsert)
		switch window.GetChar() {
		case 'q':
			if window.State.Mode == ModeNormal {
				return
			}
		case '\n':
			window.State.Mode = ModeInsert
		case gc.KEY_ESC:
			window.State.Mode = ModeNormal
		}
	}
}

func (window *Window) Resize(height, width int) {
	window.layout.Update(width, height, 0, 0)
	window.Draw()
}

func (window *Window) Draw() {
	window.channel.Window().Move(1, 1)
	window.channel.Window().Print("Channel")
	window.channels.Window().Move(1, 1)
	window.channels.Window().Print("Channels")
	gc.UpdatePanels()
	gc.Update()
}
