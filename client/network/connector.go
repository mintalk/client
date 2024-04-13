package network

import (
	"io"
	"log/slog"
	"mintalk/client/cache"
	"mintalk/client/secure"
	"net"
)

type Connector struct {
	Host        string
	conn        net.Conn
	session     string
	sender      chan NetworkData
	receiver    chan NetworkData
	serverCache *cache.ServerCache
	closer      func()
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

func (connector *Connector) Run(serverCache *cache.ServerCache) {
	connector.serverCache = serverCache
	for {
		data := <-connector.receiver
		connector.HandleResponse(data)
	}
}

func (connector *Connector) Receive(received chan<- NetworkData) {
	for {
		rawData, err := secure.ReceiveAES(connector.conn, connector.session)
		if err != nil {
			if err == io.EOF {
				slog.Info("connection closed")
				connector.Close()
				break
			}
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

func (connector *Connector) CloseListener(closer func()) {
	connector.closer = closer
}

func (connector *Connector) Close() {
	connector.conn.Close()
	if connector.closer != nil {
		connector.closer()
	}
}
