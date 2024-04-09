package panels

import "mintalk/client/ui/elements"

type ChannelsPanel struct {
	*elements.Panel
}

func NewChannelsPanel() (*ChannelsPanel, error) {
	panel, err := elements.NewPanel(1, 1)
	return &ChannelsPanel{panel}, err
}
