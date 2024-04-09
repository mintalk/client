package panels

import (
	"fmt"
	"mintalk/client/cache"
	"mintalk/client/ui/elements"
	"time"

	gc "github.com/rthornton128/goncurses"
)

type ChannelPanel struct {
	*elements.Panel
	input *elements.Input
	list  *elements.List
	cache *cache.ChannelCache
}

func NewChannelPanel() (*ChannelPanel, error) {
	panel, err := elements.NewPanel(3, 1)
	if err != nil {
		return nil, err
	}
	channelPanel := &ChannelPanel{Panel: panel}
	channelPanel.input = elements.NewInput(1, channelPanel.sendMessage)
	channelPanel.Add(channelPanel.input)
	channelPanel.list = elements.NewList(1, 1)
	channelPanel.Add(channelPanel.list)
	channelPanel.cache = cache.NewChannelCache()
	return channelPanel, nil
}

func (panel *ChannelPanel) Draw(window *gc.Window) error {
	panel.input.Active = panel.Active
	return panel.Panel.Draw(window)
}

func (panel *ChannelPanel) Resize() {
	width, height := panel.RealWidth, panel.RealHeight
	panel.Panel.Resize(width, height)
	panel.input.Resize(width - 2)
	panel.input.Move(1, height-2)
	panel.list.Resize(width-2, height-3)
	panel.list.Move(1, 1)
}

func (panel *ChannelPanel) sendMessage(message string) {
	panel.cache.AddMessage(cache.Message{Sender: "me", Contents: message, Time: time.Now()})
	panel.updateListData()
}

func (panel *ChannelPanel) updateListData() {
	panel.list.Data = make([]fmt.Stringer, 0)
	for _, message := range panel.cache.Messages {
		panel.list.Add(message)
	}
	panel.list.ProcessData()
}
