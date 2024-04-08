package network

import (
	"crypto/rand"
	"encoding/base64"
	"mintalk/server/db"
	"time"
)

type Session struct {
	User   *db.User
	Expire time.Time
}

type SessionManager struct {
	Sessions    map[string]*Session
	Lifetime    time.Duration
	TokenLength int
}

func NewSessionManager(lifetime time.Duration, tokenLength int) *SessionManager {
	return &SessionManager{
		Sessions:    make(map[string]*Session),
		Lifetime:    lifetime,
		TokenLength: tokenLength,
	}
}

func randSessionToken(length int) (string, error) {
	numBytes := length / 4 * 3
	if length%4 != 0 {
		numBytes = (length/4 + 1) * 3
	}

	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	token = token[:length]

	return token, nil
}

func (manager *SessionManager) NewSession(user *db.User) (string, error) {
	sid, err := randSessionToken(manager.TokenLength)
	if err != nil {
		return "", err
	}
	for _, ok := manager.Sessions[sid]; ok; {
		sid, err = randSessionToken(manager.TokenLength)
		if err != nil {
			return "", err
		}
	}
	session := &Session{
		User:   user,
		Expire: time.Now().Add(manager.Lifetime),
	}
	manager.Sessions[sid] = session
	return sid, nil
}

func (manager *SessionManager) GetSession(sid string) *Session {
	session, ok := manager.Sessions[sid]
	if !ok {
		return nil
	}
	if time.Now().After(session.Expire) {
		delete(manager.Sessions, sid)
		return nil
	}
	return session
}
