package network

import (
	"log/slog"
	"mintalk/client/secure"
	"net"
)

type Connector struct {
	Host    string
	conn    net.Conn
	session string
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

func (connector *Connector) Receive(received chan<- NetworkData) {
	for {
		rawData, err := secure.ReceiveAES(connector.conn, connector.session)
		if err != nil {
			slog.Error("failed to receive data", err)
			continue
		}
		data, err := Decode(rawData)
		if err != nil {
			slog.Error("failed to decode received data", err)
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
			slog.Error("failed to encode data", err)
			continue
		}
		if err := secure.SendAES(connector.conn, rawData, connector.session); err != nil {
			slog.Error("failed to send data", err)
			continue
		}
	}
}

func (connector *Connector) Close() {
	connector.conn.Close()
}
