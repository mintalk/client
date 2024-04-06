package cli

import gc "github.com/rthornton128/goncurses"

type Tab struct {
	*gc.Panel

	widthPercent  float32
	heightPercent float32

	active bool
}

func MakeTab(widthPercent float32, heightPercent float32) *Tab {
	window, _ := gc.NewWindow(0, 0, 0, 0)
	window.Box(0, 0)
	panel := gc.NewPanel(window)
	return &Tab{
		panel,
		widthPercent,
		heightPercent,
		false,
	}
}

func (tab *Tab) Resize(th int, tw int) {

}

func (tab *Tab) WriteText(text string) {
	tab.Window().Move(1, 1)
	tab.Window().Print(text)
}
