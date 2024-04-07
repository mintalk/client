package panels

import gc "github.com/rthornton128/goncurses"

type Panel struct {
	*gc.Panel

	Name string

	Width  int
	Height int
}

func NewPanel(width, height int, name string) (*Panel, error) {
	window, err := gc.NewWindow(0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	window.Box(0, 0)
	panel := gc.NewPanel(window)
	return &Panel{panel, name, width, height}, nil
}

func (panel *Panel) Draw(active bool) error {
	if active {
		err := panel.Window().ColorOn(1)
		if err != nil {
			return err
		}
		err = panel.Window().Box(0, 0)
		if err != nil {
			return err
		}
		err = panel.Window().ColorOff(1)
		if err != nil {
			return err
		}
	} else {
		if err := panel.Window().Box(0, 0); err != nil {
			return err
		}
	}

	return nil
}

func (panel *Panel) ShowList(list []string) {
	for i := 0; i < len(list); i++ {
		panel.Window().Move(i+1, 1)
		panel.Window().Print(list[i])
	}
}

func (panel *Panel) ShowName() {
	panel.Window().Move(1, 1)
	panel.Window().Print(panel.Name)
}
