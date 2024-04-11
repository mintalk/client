package network

import (
	"log/slog"
	"mintalk/client/cache"
	"mintalk/client/secure"
	"net"
	"time"
)

type Connector struct {
	Host     string
	conn     net.Conn
	session  string
	sender   chan map[string]interface{}
	receiver chan map[string]interface{}
}

func NewConnector(host string) (*Connector, error) {
	connector := &Connector{Host: host}
	var err error
	connector.conn, err = net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return connector, nil
}

func (connector *Connector) Start(username, password string) error {
	if err := connector.Auth(username, password); err != nil {
		return err
	}
	connector.sender = make(chan map[string]interface{})
	connector.receiver = make(chan map[string]interface{})
	go connector.Send(connector.sender)
	go connector.Receive(connector.receiver)
	return nil
}

func (connector *Connector) Run(channelCache *cache.ChannelCache, serverCache *cache.ServerCache) {
	for {
		data := <-connector.receiver
		switch data["action"].(string) {
		case "message":
			var messageTime time.Time
			err := messageTime.UnmarshalText([]byte(data["time"].([]byte)))
			if err != nil {
				slog.Error("failed to parse time", "err", err)
				continue
			}
			message := cache.Message{
				Sender:   data["uid"].(uint),
				Contents: data["text"].(string),
				Time:     messageTime,
			}
			channelCache.AddMessage(data["mid"].(uint), message)
		case "user":
			serverCache.AddUser(data["uid"].(uint), data["name"].(string))
		case "fetchmsg":
			for _, messageData := range data["messages"].([]string) {
				message, err := Decode([]byte(messageData))
				if err != nil {
					slog.Error("failed to decode message", "err", err)
					continue
				}
				var messageTime time.Time
				err = messageTime.UnmarshalText([]byte(message["time"].([]byte)))
				if err != nil {
					slog.Error("failed to parse time", "err", err)
					continue
				}
				messageItem := cache.Message{
					Sender:   message["uid"].(uint),
					Contents: message["text"].(string),
					Time:     messageTime,
				}
				channelCache.AddMessage(message["mid"].(uint), messageItem)
			}
		case "fetchgroup":
			for _, groupData := range data["groups"].([]string) {
				group, err := Decode([]byte(groupData))
				if err != nil {
					slog.Error("failed to decode group", "err", err)
					continue
				}
				serverCache.AddGroup(group["gid"].(uint), cache.ServerGroup{
					Name: group["name"].(string), Parent: group["parent"].(uint), HasParent: group["hasParent"].(bool),
				})
			}
		case "fetchchannel":
			for _, channelData := range data["channels"].([]string) {
				channel, err := Decode([]byte(channelData))
				if err != nil {
					slog.Error("failed to decode channel", "err", err)
					continue
				}
				serverCache.AddChannel(channel["cid"].(uint), cache.ServerChannel{
					Name: channel["name"].(string), Group: channel["group"].(uint),
				})
			}
		}
	}
}

func (connector *Connector) LoadUser(uid uint) {
	connector.sender <- map[string]interface{}{
		"action": "user",
		"uid":    uid,
	}
}

func (connector *Connector) LoadMessages(limit int) {
	connector.sender <- map[string]interface{}{
		"action": "fetchmsg",
		"limit":  limit,
	}
}

func (connector *Connector) LoadGroups() {
	connector.sender <- map[string]interface{}{
		"action": "fetchgroup",
	}
}

func (connector *Connector) LoadChannels() {
	connector.sender <- map[string]interface{}{
		"action": "fetchchannel",
	}
}

func (connector *Connector) SendMessage(text string) {
	connector.sender <- map[string]interface{}{
		"action": "message",
		"text":   text,
	}
}

func (connector *Connector) Receive(received chan<- map[string]interface{}) {
	for {
		rawData, err := secure.ReceiveAES(connector.conn, connector.session)
		if err != nil {
			slog.Error("failed to receive data", "err", err)
			continue
		}
		data, err := Decode(rawData)
		if err != nil {
			slog.Error("failed to decode received data", "err", err)
			continue
		}
		received <- data
	}
}

func (connector *Connector) Send(data <-chan map[string]interface{}) {
	for {
		sendData := <-data
		rawData, err := Encode(sendData)
		if err != nil {
			slog.Error("failed to encode data", "err", err)
			continue
		}
		if err := secure.SendAES(connector.conn, rawData, connector.session); err != nil {
			slog.Error("failed to send data", "err", err)
			continue
		}
	}
}

func (connector *Connector) Close() {
	connector.conn.Close()
}
