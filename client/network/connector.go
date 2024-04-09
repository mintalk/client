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
	sender   chan NetworkData
	receiver chan NetworkData
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
	connector.sender = make(chan NetworkData)
	connector.receiver = make(chan NetworkData)
	go connector.Send(connector.sender)
	go connector.Receive(connector.receiver)
	return nil
}

func (connector *Connector) Run(channelCache *cache.ChannelCache, serverCache *cache.ServerCache) {
	for {
		data := <-connector.receiver
		switch data["action"].(string) {
		case "message":
			rawTime := data["time"].([]byte)
			var messageTime time.Time
			err := messageTime.GobDecode(rawTime)
			if err != nil {
				slog.Error("failed to decode time", "err", err)
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
		case "fetch":
			for _, message := range data["messages"].([]NetworkData) {
				messageItem := cache.Message{
					Sender:   message["uid"].(uint),
					Contents: message["text"].(string),
					Time:     message["time"].(time.Time),
				}
				channelCache.AddMessage(message["mid"].(uint), messageItem)
			}
		}
	}
}

func (connector *Connector) LoadUser(uid uint) {
	connector.sender <- NetworkData{
		"action": "user",
		"uid":    uid,
	}
}

func (connector *Connector) LoadMessages(limit int) {
	connector.sender <- NetworkData{
		"action": "fetch",
		"limit":  limit,
	}
}

func (connector *Connector) SendMessage(text string) {
	connector.sender <- NetworkData{
		"action": "message",
		"text":   text,
	}
}

func (connector *Connector) Receive(received chan<- NetworkData) {
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

func (connector *Connector) Send(data <-chan NetworkData) {
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
