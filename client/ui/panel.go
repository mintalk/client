package ui

import gc "github.com/rthornton128/goncurses"

type Panel struct {
	*gc.Panel

	name string

	width  int
	height int
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

type Tree struct {
	Item     string
	Children []Tree
}

func (panel *Panel) ShowTree(tree Tree, y, xoffset int) {
	panel.showTreeRec(tree, y, xoffset)
}

func (panel *Panel) showTreeRec(tree Tree, y, xoffset int) int {
	panel.Window().Move(y, 1)
	panel.Window().Print("|")
	panel.Window().Move(y, xoffset)
	panel.Window().Print("|-" + tree.Item)
	c := 0
	for i := 0; i < len(tree.Children); i++ {
		c += panel.showTreeRec(tree.Children[i], y+1+c, xoffset+2)
	}
	return c
}

func (panel *Panel) ShowList(list []string) {
	for i := 0; i < len(list); i++ {
		panel.Window().Move(i+1, 1)
		panel.Window().Print(list[i])
	}
}

func (panel *Panel) ShowName() {
	panel.Window().Move(1, 1)
	panel.Window().Print(panel.name)
}
