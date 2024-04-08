package network

import (
	"log/slog"
	"mintalk/server/db"
	"mintalk/server/secure"
	"net"
)

type ProtocolExecutor struct {
	Database       *db.Connection
	Conn           net.Conn
	SessionManager *SessionManager
	Session        string
}

func (executor *ProtocolExecutor) Run() error {
	_, err := executor.Auth()
	return err
}

func (executor *ProtocolExecutor) Auth() (bool, error) {
	data, err := secure.Receive3Pass(executor.Conn)
	if err != nil {
		return false, err
	}
	request, err := Decode(data)
	if err != nil {
		return false, err
	}
	authed, err := ValidateAuthRequest(executor.Database, request)
	if err != nil {
		return false, err
	}
	response := NetworkData{"authed": authed}
	data, err = Encode(response)
	if err != nil {
		return false, err
	}
	if authed {
		var user *db.User
		err := executor.Database.Where(&db.User{Name: request["username"].(string)}).First(&user).Error
		if err != nil {
			return false, err
		}
		session, err := executor.SessionManager.NewSession(user)
		if err != nil {
			return false, err
		}
		response["session"] = session
		executor.Session = session
	}
	prime, err := secure.RandomPrime(1024)
	if err != nil {
		return false, err
	}
	if err := secure.Send3Pass(executor.Conn, data, prime); err != nil {
		return false, err
	}
	return authed, nil
}

func (executor *ProtocolExecutor) Receive(received chan<- NetworkData) {
	for {
		rawData, err := secure.ReceiveAES(executor.Conn, executor.Session)
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

func (executor *ProtocolExecutor) Send(data <-chan NetworkData) {
	for {
		sendData := <-data
		rawData, err := Encode(sendData)
		if err != nil {
			slog.Error("failed to encode data", err)
			continue
		}
		if err := secure.SendAES(executor.Conn, rawData, executor.Session); err != nil {
			slog.Error("failed to send data", err)
			continue
		}
	}
}
