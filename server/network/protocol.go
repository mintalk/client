package network

import (
	"mintalk/server/db"
	"mintalk/server/secure"
	"net"
)

type ProtocolExecutor struct {
	Database       *db.Connection
	Conn           net.Conn
	SessionManager *SessionManager
}

func (executor *ProtocolExecutor) Run() error {
	_, err := executor.Auth()
	return err
}

func (executor *ProtocolExecutor) Auth() (bool, error) {
	data, err := secure.Recieve3Pass(executor.Conn)
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
