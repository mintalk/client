package network

import (
	"io"
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

func NewProtocolExecutor(conn net.Conn, database *db.Connection, sessionManager *SessionManager) *ProtocolExecutor {
	return &ProtocolExecutor{Conn: conn, Database: database, SessionManager: sessionManager}
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
	if authed {
		var user *db.User
		username, ok := request["username"].(string)
		if !ok {
			return false, nil
		}
		err := executor.Database.Where(&db.User{Name: username}).First(&user).Error
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
	data, err = Encode(response)
	if err != nil {
		return false, err
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
			if err == io.EOF {
				executor.SessionManager.DeleteSession(executor.Session)
				close(received)
				return
			}
			slog.Debug("failed to receive data", "err", err)
			continue
		}
		if rawData == nil || len(rawData) == 0 {
			close(received)
			return
		}
		data, err := Decode(rawData)
		if err != nil {
			slog.Debug("failed to decode received data", "err", err)
			continue
		}
		received <- data
	}
}

func (executor *ProtocolExecutor) Send(data <-chan NetworkData) {
	for {
		sendData, ok := <-data
		if !ok {
			return
		}
		if sendData == nil {
			continue
		}
		rawData, err := Encode(sendData)
		if err != nil {
			slog.Debug("failed to encode data", "err", err)
			continue
		}
		if err := secure.SendAES(executor.Conn, rawData, executor.Session); err != nil {
			slog.Debug("failed to send data", "err", err)
			continue
		}
	}
}
