package panels

import (
	"strings"

	"mintalk/client/cache"
	"mintalk/client/network"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Channel struct {
	*tview.Flex
	serverCache   *cache.ServerCache
	Connector     *network.Connector
	ActiveChannel uint
	ChannelOpened bool
	textView *tview.TextView
	inputField *tview.InputField
	app *tview.Application
}

func NewChannel(app *tview.Application, connector *network.Connector, serverCache *cache.ServerCache) *Channel {
	channel := &Channel{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
		Connector: connector,
		ChannelOpened: false,
		textView: tview.NewTextView().SetWordWrap(true).SetScrollable(true),
		inputField: tview.NewInputField().SetFieldWidth(0),
		app: app,
	}
	channel.AddItem(channel.textView, 0, 1, false)
	channel.AddItem(channel.inputField, 1, 1, true)
	channel.SetBorder(true)
	channel.serverCache = serverCache
	channel.serverCache.AddListener(channel.updateListData)
	channel.inputField.SetDoneFunc(func (key tcell.Key) {
		channel.sendMessage(channel.inputField.GetText())
		channel.inputField.SetText("")
	})
	channel.textView.SetFocusFunc(func () {
		channel.app.SetFocus(channel.inputField)
	})
	channel.textView.SetChangedFunc(func () {
		channel.textView.ScrollToEnd()
	})
	return channel
}

func (channel *Channel) MoveChannel(channelId uint) {
	channel.ChannelOpened = true
	channel.serverCache.GetChannelCache(channel.ActiveChannel).Listeners = nil
	channel.ActiveChannel = channelId
	channel.serverCache.GetChannelCache(channelId).AddListener(channel.updateListData)
	channel.Connector.LoadMessages(0, channelId)
	channelName := channel.serverCache.Channels[channelId].Name
	channel.SetTitle(channelName)
}

func (channel *Channel) sendMessage(message string) {
	channel.Connector.SendMessage(message, channel.ActiveChannel)
}

func (channel *Channel) updateListData() {
	list := make([]string, 0)
	for _, message := range channel.serverCache.GetChannelCache(channel.ActiveChannel).GetMessages() {
		username, ok := channel.serverCache.Users[message.Sender]
		if ok {
			message.Username = username
		} else {
			channel.Connector.LoadUser(message.Sender)
		}
		list = append(list, message.String())
	}
	text := strings.Join(list, "\n")
	channel.app.QueueUpdateDraw(func () {
		channel.textView.SetText(text)
	})
}
