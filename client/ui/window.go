package ui

import (
	"log/slog"
	"mintalk/client/ui/panels"
	"os"
	"os/signal"
	"syscall"

	gc "github.com/rthornton128/goncurses"
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
	window := &Window{ncursesWindow, nil, nil, nil, NewUIState(), false}
	window.State.ActiveTab = TabChannel
	return window, nil
}

func (window *Window) Create() error {
	window.Keypad(true)
	gc.Echo(false)
	gc.CBreak(true)
	window.Timeout(0)
	gc.Cursor(0)

	InitColors()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go window.CloseListener(sigc)

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

	window.Resize(window.MaxYX())
	return nil
}

func (window *Window) CloseListener(sigc <-chan os.Signal) {
	for {
		select {
		case <-sigc:
			window.running = false
			window.Close()
			os.Exit(0)
		}
	}
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
	char := window.GetChar()
	if char == 0 {
		return
	}
	window.channel.Update(char)
	window.channels.Update(char)
	switch char {
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

	return nil
}

func (window *Window) Resize(height, width int) {
	window.layout.Update(width, height, 0, 0)
	window.channel.Resize()
}
