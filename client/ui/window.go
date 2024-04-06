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
	err = gc.StartColor()
	if err != nil {
		return nil, err
	}
	err = gc.InitPair(1, gc.C_CYAN, 0)
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
		char := window.GetChar()
		switch char {
		case 'q':
			if window.State.Mode == ModeNormal {
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
}

func (window *Window) Resize(height, width int) {
	window.layout.Update(width, height, 0, 0)
}

func (window *Window) Draw() error {
	if err := window.channel.Draw(window.State.ActiveTab == TabChannel); err != nil {
		return err
	}
	window.channel.ShowTree(Tree{
		Item: "Folder1",
		Children: []Tree{
			{
				Item: "Folder 2",
				Children: []Tree{
					{
						Item:     "File 1",
						Children: make([]Tree, 0),
					},
				},
			},
			{
				Item:     "File 2",
				Children: make([]Tree, 0),
			},
		},
	}, 1, 1)
	if err := window.channels.Draw(window.State.ActiveTab == TabChannels); err != nil {
		return err
	}
	window.channels.ShowList([]string{
		"Ich",
		"bin",
		"besser",
		"als",
		"du",
		"---",
		"Nervermind",
		"Tree ->",
		"funktioniert",
		"noch",
		"nicht",
	})
	gc.UpdatePanels()
	return gc.Update()
}
