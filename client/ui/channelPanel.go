package ui

type ChannelPanel struct {
	*Panel
}

func NewChannelPanel() (*ChannelPanel, error) {
	panel, err := NewPanel(3, 1)
	return &ChannelPanel{panel}, err
}
