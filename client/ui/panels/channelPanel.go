package panels

type ChannelPanel struct {
	*Panel
}

func NewChannelPanel() (*ChannelPanel, error) {
	panel, err := NewPanel(3, 1, "ChannelPanel")
	return &ChannelPanel{panel}, err
}