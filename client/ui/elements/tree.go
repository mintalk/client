package elements

import (
	"fmt"
	"strings"

	gc "github.com/rthornton128/goncurses"
)

type Tree struct {
	Nodes     []*TreeNode
	X         int
	Y         int
	Width     int
	Height    int
	Expansion int
}

type TreeNode struct {
	Item     fmt.Stringer
	Children []*TreeNode
}

func NewTreeNode(item fmt.Stringer) *TreeNode {
	return &TreeNode{Item: item, Children: make([]*TreeNode, 0)}
}

func NewTree() *Tree {
	return &Tree{}
}

func (tree *Tree) Update(key gc.Key) {}

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
	offset := 0
	for _, node := range tree.Nodes {
		length := node.draw(window, tree.X, tree.Y+offset, tree.X, tree.Width, tree.Height, false, false, tree.Expansion)
		offset += length
	}
}

func (node *TreeNode) draw(window *gc.Window, x, y, originX, width, height int, child bool, last bool, expand int) int {
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
	line := node.Item.String()
	if len(line) < width {
		line += strings.Repeat(" ", width-len(line))
	}
	window.Print(line[:width])
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
		childLength := childTree.draw(window, x+childrenOffset, y+length, x, width-childrenOffset-expand-1, height-length, true, childLast, expand)
		if !childLast {
			for i := 1; i < childLength; i++ {
				window.MoveAddChar(y+length+i, x+childrenOffset, gc.ACS_VLINE)
			}
		}
		length += childLength
	}
	return length
}
