package panels

type ChannelsPanel struct {
	*Panel
}

func NewChannelsPanel() (*ChannelsPanel, error) {
	panel, err := NewPanel(1, 1, "ChannelPanels")
	return &ChannelsPanel{panel}, err
}