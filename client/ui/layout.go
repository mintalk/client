package ui

import "mintalk/client/ui/panels"

type Direction uint

const (
	Horizontal Direction = iota
	Vertical
)

type Layout struct {
	Panel     *panels.Panel
	Child     *Layout
	Direction Direction
}

func (layout *Layout) GetWidth() (width int) {
	width += layout.Panel.Width
	if layout.Child != nil && layout.Direction == Horizontal {
		width += layout.Child.GetWidth()
	}
	return
}

func (layout *Layout) GetHeight() (height int) {
	height += layout.Panel.Width
	if layout.Child != nil && layout.Direction == Vertical {
		height += layout.Child.GetHeight()
	}
	return
}

func (layout *Layout) Update(maxWidth, maxHeight, offsetX, offsetY int) {
	widthSum := layout.Panel.Width
	heightSum := layout.Panel.Height
	if layout.Child != nil {
		if layout.Direction == Horizontal {
			widthSum += layout.Child.GetWidth()
		} else if layout.Direction == Vertical {
			heightSum += layout.Child.GetHeight()
		}
	}

	availableWidth := float64(maxWidth - offsetX)
	availableHeight := float64(maxHeight - offsetY)

	widthFraction := float64(layout.Panel.Width) / float64(widthSum)
	heightFraction := float64(layout.Panel.Height) / float64(heightSum)

	panelWidth := widthFraction * availableWidth
	panelHeight := heightFraction * availableHeight

	layout.Panel.Window().Resize(int(panelHeight), int(panelWidth))
	layout.Panel.Move(offsetY, offsetX)
	layout.Panel.Window().Box(0, 0)

	if layout.Child != nil {
		if layout.Direction == Horizontal {
			offsetX += int(panelWidth)
		} else if layout.Direction == Vertical {
			offsetY += int(panelHeight)
		}
		layout.Child.Update(maxWidth, maxHeight, offsetX, offsetY)
	}
}
