package elements

import gc "github.com/rthornton128/goncurses"

type Text struct {
	Length    int
	Markable  bool
	MarkStart int
	MarkEnd   int
}

func NewText(length int) *Text {
	return &Text{length, false, 0, 0}
}

func Update(key gc.Key) {
}

func Draw(window *gc.Window, x, y int) {
}
