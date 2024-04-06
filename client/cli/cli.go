package cli

import (
	gc "github.com/rthornton128/goncurses"
)

type CLI struct {
	stdsrc *gc.Window
	panels [2]*Tab
}

func InitCli() *CLI {
	stdscr, _ := gc.Init()
	gc.Echo(false)
	stdscr.Keypad(true)

	panels := [2]*Tab{
		MakeTab(0.3, 1),
		MakeTab(1, 1),
	}
	ret := &CLI{
		stdscr,
		panels,
	}

	ret.Resize(stdscr.MaxYX())

	return ret
}

func (cli *CLI) Close() {
	gc.End()
}

func (cli *CLI) Run() {
	cli.Draw()
	for {
		if cli.stdsrc.GetChar() == 'q' {
			return
		}
	}
}

func (cli *CLI) Resize(th int, tw int) {
	lastx := 0
	for i := 0; i < len(cli.panels); i++ {
		tab := cli.panels[i]
		w := int(float32(tw) * tab.widthPercent)
		h := int(float32(th) * tab.heightPercent)

		tab.Window().Resize(h, w)
		tab.Window().MoveWindow(0, lastx)
		tab.WriteText("Test")
		cli.panels[i].Window().Box(0, 0)

		tw -= w
		lastx = w
	}
	cli.Draw()
}

func (cli *CLI) Draw() {
	gc.UpdatePanels()
	gc.Update()
}
