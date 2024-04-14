package elements

import (
	gc "github.com/mintalk/goncurses"
)

type Panel struct {
	*gc.Panel

	Width      int
	Height     int
	RealWidth  int
	RealHeight int
	Active     bool

	Components []Element
}

func NewPanel(width, height int) (*Panel, error) {
	window, err := gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	window.Box(0, 0)
	panel := gc.NewPanel(window)
	components := make([]Element, 0)
	return &Panel{panel, width, height, 0, 0, false, components}, nil
}

func (panel *Panel) Draw(window *gc.Window) error {
	var color int16 = 2
	if panel.Active {
		color = 1
	}
	panel.Window().ColorOn(color)
	if err := panel.Window().Box(0, 0); err != nil {
		return err
	}

	for _, component := range panel.Components {
		component.Draw(panel.Window())
	}

	return nil
}

func (panel *Panel) Update(key gc.Key) {
	for _, component := range panel.Components {
		component.Update(key)
	}
}

func (panel *Panel) Add(component Element) {
	panel.Components = append(panel.Components, component)
}

func (panel *Panel) Resize(width, height int) {
	panel.RealWidth = width
	panel.RealHeight = height
	panel.Window().Resize(height, width)
}
