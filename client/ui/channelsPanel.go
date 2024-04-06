package ui

type ChannelsPanel struct {
	*Panel
}

func NewChannelsPanel() (*ChannelsPanel, error) {
	panel, err := NewPanel(1, 1)
	return &ChannelsPanel{panel}, err
}
