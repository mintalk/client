package ui

import gc "github.com/rthornton128/goncurses"

type Panel struct {
	*gc.Panel
	width  int
	height int
}

func NewPanel(width, height int) (*Panel, error) {
	window, err := gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	window.Box(0, 0)
	panel := gc.NewPanel(window)
	return &Panel{panel, width, height}, nil
}
