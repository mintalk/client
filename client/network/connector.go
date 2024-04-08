package network

import "net"

type Connector struct {
	Host string
	conn net.Conn
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

func (connector *Connector) Close() {
	connector.conn.Close()
}
