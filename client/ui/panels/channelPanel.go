package panels

import (
	"fmt"
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui/elements"

	gc "github.com/rthornton128/goncurses"
)

type ChannelPanel struct {
	*elements.Panel
	input        *elements.Input
	list         *elements.List
	channelCache *cache.ChannelCache
	serverCache  *cache.ServerCache
	Connector    *network.Connector
}

func NewChannelPanel(connector *network.Connector, channelCache *cache.ChannelCache, serverCache *cache.ServerCache) (*ChannelPanel, error) {
	panel, err := elements.NewPanel(3, 1)
	if err != nil {
		return nil, err
	}
	channelPanel := &ChannelPanel{Panel: panel, Connector: connector}
	channelPanel.input = elements.NewInput(1, channelPanel.sendMessage)
	channelPanel.Add(channelPanel.input)
	channelPanel.list = elements.NewList(1, 1)
	channelPanel.Add(channelPanel.list)
	channelPanel.channelCache = channelCache
	channelPanel.channelCache.AddListener(channelPanel.updateListData)
	channelPanel.serverCache = serverCache
	channelPanel.serverCache.AddListener(channelPanel.updateListData)
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
	panel.Connector.LoadMessages(width - 2)
}

func (panel *ChannelPanel) sendMessage(message string) {
	panel.Connector.SendMessage(message)
}

func (panel *ChannelPanel) updateListData() {
	panel.list.Data = make([]fmt.Stringer, 0)
	for _, message := range panel.channelCache.GetMessages() {
		username, ok := panel.serverCache.Users[message.Sender]
		if ok {
			message.Username = username
		} else {
			panel.Connector.LoadUser(message.Sender)
		}
		panel.list.Add(message)
	}
	panel.list.ProcessData()
}
