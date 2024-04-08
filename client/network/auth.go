package network

import "mintalk/client/secure"

func (connector *Connector) Auth(username, password string) error {
	prime, err := secure.RandomPrime(1024)
	if err != nil {
		return err
	}
	data := make(NetworkData)
	data["username"] = username
	data["password"] = password
	encodedData, err := Encode(data)
	if err != nil {
		return err
	}
	return secure.Send3Pass(connector.conn, encodedData, prime)
}
