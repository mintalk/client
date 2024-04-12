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
	input         *elements.Input
	list          *elements.List
	serverCache   *cache.ServerCache
	Connector     *network.Connector
	ActiveChannel uint
	ChannelOpened bool
}

func NewChannelPanel(connector *network.Connector, serverCache *cache.ServerCache) (*ChannelPanel, error) {
	panel, err := elements.NewPanel(3, 1)
	if err != nil {
		return nil, err
	}
	channelPanel := &ChannelPanel{Panel: panel, Connector: connector, ChannelOpened: false}
	channelPanel.input = elements.NewInput(1, channelPanel.sendMessage)
	channelPanel.Add(channelPanel.input)
	channelPanel.list = elements.NewList(1, 1)
	channelPanel.Add(channelPanel.list)
	channelPanel.serverCache = serverCache
	channelPanel.serverCache.AddListener(channelPanel.updateListData)
	return channelPanel, nil
}

func (panel *ChannelPanel) Draw(window *gc.Window) error {
	panel.input.Active = panel.Active && panel.ChannelOpened
	err := panel.Panel.Draw(window)
	if err != nil {
		return err
	}
	channel, ok := panel.serverCache.Channels[panel.ActiveChannel]
	if ok {
		panel.Window().MovePrint(0, 2, channel.Name)
	}
	return nil
}

func (panel *ChannelPanel) Resize() {
	width, height := panel.RealWidth, panel.RealHeight
	panel.Panel.Resize(width, height)
	panel.input.Resize(width - 2)
	panel.input.Move(1, height-2)
	panel.list.Resize(width-2, height-3)
	panel.list.Move(1, 1)
	panel.Connector.LoadMessages(width-2, panel.ActiveChannel)
}

func (panel *ChannelPanel) MoveChannel(channel uint) {
	panel.ChannelOpened = true
	panel.serverCache.GetChannelCache(channel).Listeners = nil
	panel.ActiveChannel = channel
	panel.serverCache.GetChannelCache(channel).AddListener(panel.updateListData)
	panel.updateListData()
}

func (panel *ChannelPanel) sendMessage(message string) {
	panel.Connector.SendMessage(message, panel.ActiveChannel)
}

func (panel *ChannelPanel) updateListData() {
	panel.list.Data = make([]fmt.Stringer, 0)
	for _, message := range panel.serverCache.GetChannelCache(panel.ActiveChannel).GetMessages() {
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
