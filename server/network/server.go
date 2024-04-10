package network

import (
	"log/slog"
	"mintalk/server/config"
	"mintalk/server/db"
	"net"
	"time"
)

type Server struct {
	senders        map[string]chan<- map[string]interface{}
	sessionManager *SessionManager
	database       *db.Connection
	host           string
}

func NewServer(database *db.Connection, conf *config.Config) *Server {
	return &Server{
		senders:        make(map[string]chan<- map[string]interface{}),
		sessionManager: NewSessionManager(time.Duration(conf.SessionLifetime)*time.Minute, 32),
		database:       database,
		host:           conf.Host,
	}
}

func (server *Server) Run() error {
	listener, err := net.Listen("tcp", server.host)
	if err != nil {
		slog.Error("failed to create listener", "err", err)
		return err
	}
	defer listener.Close()
	slog.Info("listening on " + server.host)
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Warn("failed to accept connection", err)
			continue
		}
		go server.handleClient(conn)
	}
}

func (server *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	executor := NewProtocolExecutor(conn, server.database, server.sessionManager)
	authed, err := executor.Auth()
	if err != nil {
		slog.Error("auth failed", "err", err)
	}
	if !authed {
		return
	}

	sender := make(chan map[string]interface{})
	server.senders[executor.Session] = sender
	go executor.Send(sender)
	receiver := make(chan map[string]interface{})
	go executor.Receive(receiver)
	for {
		request, ok := <-receiver
		if !ok {
			close(server.senders[executor.Session])
			delete(server.senders, executor.Session)
			break
		}
		if request == nil {
			continue
		}
		action, ok := request["action"]
		if !ok {
			continue
		}
		switch action {
		case "message":
			server.ActionMessage(executor.Session, request)
		case "fetch":
			server.ActionFetch(executor.Session, request)
		case "user":
			server.ActionUser(executor.Session, request)
		}
	}
}

func (server *Server) Broadcast(data map[string]interface{}) {
	for sid, sender := range server.senders {
		session := server.sessionManager.GetSession(sid)
		if session == nil {
			close(server.senders[sid])
			delete(server.senders, sid)
			continue
		}
		sender <- data
	}
}
