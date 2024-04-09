package network

import (
	"log/slog"
	"mintalk/server/config"
	"mintalk/server/db"
	"net"
	"time"
)

type Server struct {
	senders        map[string]chan<- NetworkData
	sessionManager *SessionManager
	database       *db.Connection
	host           string
}

func NewServer(database *db.Connection, conf *config.Config) *Server {
	return &Server{
		senders:        make(map[string]chan<- NetworkData),
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

	sender := make(chan NetworkData)
	server.senders[executor.Session] = sender
	go executor.Send(sender)
	receiver := make(chan NetworkData)
	go executor.Receive(receiver)
	for {
		request := <-receiver
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

func (server *Server) Broadcast(data NetworkData) {
	for sid, sender := range server.senders {
		session := server.sessionManager.GetSession(sid)
		if session == nil {
			delete(server.senders, sid)
			continue
		}
		sender <- data
	}
}
