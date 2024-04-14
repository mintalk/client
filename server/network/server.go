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
	config         *config.Config
}

func NewServer(database *db.Connection, conf *config.Config) *Server {
	return &Server{
		senders:        make(map[string]chan<- NetworkData),
		sessionManager: NewSessionManager(time.Duration(conf.SessionLifetime)*time.Minute, 32),
		database:       database,
		config:         conf,
	}
}

func (server *Server) Run() error {
	listener, err := net.Listen("tcp", server.config.Host)
	if err != nil {
		slog.Error("failed to create listener", "err", err)
		return err
	}
	defer listener.Close()
	slog.Info("listening on " + server.config.Host)
	ScheduleRepeatedTask(time.Duration(server.config.CleanupInterval)*time.Minute, server.Cleanup)
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Debug("failed to accept connection", "err", err)
			continue
		}
		go server.HandleClient(conn)
	}
}

func (server *Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	executor := NewProtocolExecutor(conn, server.database, server.sessionManager)
	authed, err := executor.Auth()
	if err != nil {
		slog.Debug("auth failed", "err", err)
		return
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
		request, ok := <-receiver
		if !ok {
			server.DeleteSender(executor.Session)
			break
		}
		if request == nil {
			continue
		}
		server.HandleRequest(executor.Session, request)
	}
}

func (server *Server) Broadcast(data NetworkData) {
	for _, sid := range server.sessionManager.GetSessions() {
		sender := server.GetSender(sid)
		if sender == nil {
			server.sessionManager.DeleteSession(sid)
			continue
		}
		sender <- data
	}
}

func (server *Server) GetSender(sid string) chan<- NetworkData {
	if server.sessionManager.GetSession(sid) == nil {
		delete(server.senders, sid)
		return nil
	}
	return server.senders[sid]
}

func (server *Server) DeleteSender(sid string) {
	if _, ok := server.senders[sid]; ok {
		delete(server.senders, sid)
	}
	if server.sessionManager.GetSession(sid) != nil {
		server.sessionManager.DeleteSession(sid)
	}
}

func (server *Server) Cleanup() {
	slog.Debug("cleaning up")
	for sid, sender := range server.senders {
		if server.sessionManager.GetSession(sid) == nil || sender == nil {
			server.DeleteSender(sid)
		}
	}
	var messages []db.Message
	err := server.database.Find(&messages).Error
	if err != nil {
		slog.Debug("failed to clean messages", "err", err)
	}
	for _, message := range messages {
		var channels []db.Channel
		err := server.database.Find(&channels).Error
		if err != nil {
			slog.Debug("failed to find channels", "err", err)
			continue
		}
		channelFound := false
		for _, channel := range channels {
			if channel.ID == message.Channel {
				channelFound = true
				break
			}
		}
		if !channelFound {
			err := server.database.Delete(&message).Error
			if err != nil {
				slog.Debug("failed to delete message", "err", err)
			}
		}
	}
}
