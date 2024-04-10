package network

import (
	"mintalk/client/secure"
)

func (connector *Connector) Auth(username, password string) error {
	if err := connector.SendAuth(username, password); err != nil {
		return err
	}
	session, err := connector.ReceiveAuth()
	connector.session = session
	return err
}

func (connector *Connector) SendAuth(username, password string) error {
	prime, err := secure.RandomPrime(1024)
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	data["username"] = username
	data["password"] = password
	encodedData, err := Encode(data)
	if err != nil {
		return err
	}
	return secure.Send3Pass(connector.conn, encodedData, prime)
}

func (connector *Connector) ReceiveAuth() (string, error) {
	rawData, err := secure.Receive3Pass(connector.conn)
	if err != nil {
		return "", err
	}
	data, err := Decode(rawData)
	if err != nil {
		return "", err
	}
	if data["authed"].(bool) {
		return data["session"].(string), nil
	}
	return "", nil
}
