package elements

import (
	"fmt"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type List struct {
	Data   []fmt.Stringer
	X      int
	Y      int
	Width  int
	Height int
	lines  []string
}

func NewList(x, y int) *List {
	return &List{make([]fmt.Stringer, 0), x, y, 1, 1, nil}
}

func (list *List) Add(data fmt.Stringer) {
	list.Data = append(list.Data, data)
}

func (list *List) Resize(width, height int) {
	list.Width = width
	list.Height = height
}

func (list *List) Move(x, y int) {
	list.X = x
	list.Y = y
}

func (list *List) Draw(window *gc.Window) {
	begin := 0
	if len(list.lines) > list.Height {
		begin = len(list.lines) - list.Height
	}
	for i, line := range list.lines[begin:] {
		if len(line) < list.Width {
			line = line + strings.Repeat(" ", list.Width-len(line))
		}
		window.MovePrint(list.Y+i, list.X, line)
	}
}

func (list *List) Update(key gc.Key) {
	// Do nothing
}

func (list *List) ProcessData() {
	list.lines = make([]string, 0)
	for _, data := range list.Data {
		lines := data.String()
		for _, line := range strings.Split(lines, "\n") {
			for begin := 0; begin < len(line); begin += list.Width {
				end := begin + list.Width
				if end > len(line) {
					end = len(line)
				}
				list.lines = append(list.lines, line[begin:end])
			}
		}
	}
}
