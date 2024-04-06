package ui

import (
	gc "github.com/rthornton128/goncurses"
)

type State uint

const (
	STATE_NORMAL State = iota
	STATE_INSERT
)

type Window struct {
	*gc.Window
	panels []*Tab
	State  State
}

func NewWindow() (*Window, error) {
	ncursesWindow, err := gc.Init()
	if err != nil {
		return nil, err
	}
	window := &Window{ncursesWindow, []*Tab{}, STATE_NORMAL}
	return window, nil
}

func (window *Window) Create() {
	window.Keypad(true)

	window.panels = []*Tab{
		MakeTab(0.3, 1),
		MakeTab(1, 1),
	}

	window.Resize(window.MaxYX())
}

func (window *Window) Close() {
	gc.End()
}

func (window *Window) Run() {
	for {
		window.panels[0].WriteText(map[State]string{STATE_NORMAL: "normal", STATE_INSERT: "insert"}[window.State])
		window.Draw()
		gc.Echo(window.State == STATE_INSERT)
		switch window.GetChar() {
		case 'q':
			if window.State == STATE_NORMAL {
				return
			}
		case '\n':
			window.State = STATE_INSERT
		case gc.KEY_ESC:
			window.State = STATE_NORMAL
		}
	}
}

func (window *Window) Resize(th int, tw int) {
	lastx := 0
	for i := 0; i < len(window.panels); i++ {
		tab := window.panels[i]
		w := int(float64(tw) * tab.widthPercent)
		h := int(float64(th) * tab.heightPercent)

		tab.Window().Resize(h, w)
		tab.Window().MoveWindow(0, lastx)
		tab.WriteText("Test")
		window.panels[i].Window().Box(0, 0)

		tw -= w
		lastx = w
	}
	window.Draw()
}

func (window *Window) Draw() {
	gc.UpdatePanels()
	gc.Update()
}
