package elements

import (
	"fmt"
	"sort"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type Tree struct {
	Nodes     []*TreeNode
	X         int
	Y         int
	Cursor    int
	Width     int
	Height    int
	Expansion int
	Active    bool
	click     bool
	length    int
}

type TreeNode struct {
	Item     TreeItem
	Children []*TreeNode
}

type TreeItem struct {
	Value   fmt.Stringer
	OnClick func()
}

func NewTreeNode(item TreeItem) *TreeNode {
	return &TreeNode{Item: item, Children: make([]*TreeNode, 0)}
}

func NewTree() *Tree {
	return &Tree{}
}

func (tree *Tree) Update(key gc.Key) {
	if !tree.Active {
		return
	}
	if key == gc.KEY_DOWN {
		if tree.Cursor < tree.length-1 {
			tree.Cursor++
		}
	} else if key == gc.KEY_UP {
		if tree.Cursor > 0 {
			tree.Cursor--
		}
	} else if key == gc.KEY_ENTER || key == gc.KEY_RETURN || key == gc.KEY_RIGHT {
		tree.click = true
	}
}

func (tree *Tree) Move(x, y int) {
	tree.X = x
	tree.Y = y
}

func (tree *Tree) Resize(width, height int) {
	tree.Width = width
	tree.Height = height
}

func (tree *Tree) Expand(value int) {
	tree.Expansion = value
}

func (tree *Tree) Draw(window *gc.Window) {
	nodes := make([]*TreeNode, len(tree.Nodes))
	copy(nodes, tree.Nodes)
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i] == nil || nodes[j] == nil {
			return true
		}
		return nodes[i].Item.Value.String() < nodes[j].Item.Value.String()
	})
	offset := 0
	for _, node := range nodes {
		if node == nil {
			continue
		}
		cursor := tree.Cursor - offset
		if !tree.Active {
			cursor = -1
		}
		length := node.draw(window, tree.X, tree.Y+offset, tree.X, tree.Width, tree.Height, false, false, tree.Expansion, cursor, tree.click)
		offset += length
	}
	tree.length = offset
	if tree.click {
		tree.click = false
	}
}

func (node *TreeNode) draw(window *gc.Window, x, y, originX, width, height int, child bool, last bool, expand int, cursor int, click bool) int {
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Item.Value.String() < node.Children[j].Item.Value.String()
	})
	if originX < x {
		window.MovePrint(y, originX, strings.Repeat(" ", x-originX))
	}
	window.Move(y, x)
	if child {
		if last {
			window.AddChar(gc.ACS_LLCORNER)
		} else {
			window.AddChar(gc.ACS_LTEE)
		}
		for i := 0; i < expand; i++ {
			window.MoveAddChar(y, x+i+1, gc.ACS_HLINE)
		}
		window.Move(y, x+expand+1)
	}
	lineY, lineX := window.CursorYX()
	line := node.Item.Value.String()
	printedLineLength := len(line)
	if printedLineLength > width {
		printedLineLength = width
	}
	if cursor == 0 {
		if click && node.Item.OnClick != nil {
			go node.Item.OnClick()
		}
		window.AttrOn(gc.A_REVERSE)
	}
	window.Print(line[:printedLineLength])
	if cursor == 0 {
		window.AttrOff(gc.A_REVERSE)
	}
	if len(line) < width {
		window.MovePrint(lineY, lineX+printedLineLength, strings.Repeat(" ", width-printedLineLength))
	}
	childrenOffset := 0
	if child {
		childrenOffset = expand + 1
	}
	length := 1
	for idx, childTree := range node.Children {
		if length >= height {
			break
		}
		window.Move(y+length, x)
		childLast := idx == len(node.Children)-1
		childLength := childTree.draw(window, x+childrenOffset, y+length, x, width-childrenOffset-expand-1, height-length, true, childLast, expand, cursor-length, click)
		if !childLast {
			for i := 1; i < childLength; i++ {
				window.MoveAddChar(y+length+i, x+childrenOffset, gc.ACS_VLINE)
			}
		}
		length += childLength
	}
	return length
}
