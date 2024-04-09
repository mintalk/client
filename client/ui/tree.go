package ui

import (
	"mintalk/client/ui/elements"

	gc "github.com/rthornton128/goncurses"
)

type Tree struct {
	Item     string
	Children []Tree
}

func (tree *Tree) Draw(panel *elements.Panel, x, y, expand int) {
	tree.draw(panel, x, y, false, false, expand)
}

func (tree *Tree) draw(panel *elements.Panel, x, y int, child bool, last bool, expand int) int {
	panel.Window().Move(y, x)
	if child {
		if last {
			panel.Window().AddChar(gc.ACS_LLCORNER)
		} else {
			panel.Window().AddChar(gc.ACS_LTEE)
		}
		for i := 0; i < expand; i++ {
			panel.Window().MoveAddChar(y, x+i+1, gc.ACS_HLINE)
		}
		panel.Window().Move(y, x+expand+1)
	}
	panel.Window().Print(tree.Item)
	childrenX := x
	if child {
		childrenX += expand + 1
	}
	length := 1
	for idx, childTree := range tree.Children {
		panel.Window().Move(y+length, x)
		childLast := idx == len(tree.Children)-1
		childLength := childTree.draw(panel, childrenX, y+length, true, childLast, expand)
		if !childLast {
			for i := 1; i < childLength; i++ {
				panel.Window().MoveAddChar(y+length+i, childrenX, gc.ACS_VLINE)
			}
		}
		length += childLength
	}
	return length
}
