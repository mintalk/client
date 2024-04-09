package network

import (
	"log/slog"
	"mintalk/server/db"
	"time"
)

func (server *Server) ActionMessage(sid string, data NetworkData) {
	session := server.sessionManager.GetSession(sid)
	if session == nil {
		return
	}
	message := &db.Message{
		UID:  session.User.ID,
		Text: data["text"].(string),
		Time: time.Now(),
	}
	server.database.Create(message)
	rawTime, err := message.Time.GobEncode()
	if err != nil {
		slog.Error("failed to encode time", "err", err)
	}
	broadcast := NetworkData{
		"action": "message",
		"mid":    message.ID,
		"text":   message.Text,
		"uid":    session.User.ID,
		"time":   rawTime,
	}
	server.Broadcast(broadcast)
}

func (server *Server) ActionFetch(sid string, data NetworkData) {
	limit, ok := data["limit"].(int)
	if !ok {
		limit = 0
	}
	var messages []db.Message
	query := server.database.Order("time desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&messages).Error
	if err != nil {
		slog.Error("failed to fetch messages", "err", err)
		return
	}
	responseMessages := make([]NetworkData, len(messages))
	for i, message := range messages {
		responseMessages[i] = NetworkData{
			"id":   message.ID,
			"uid":  message.UID,
			"text": message.Text,
			"time": message.Time,
		}
	}
	response := NetworkData{
		"action":   "fetch",
		"messages": responseMessages,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionUser(sid string, data NetworkData) {
	uid, ok := data["uid"].(uint)
	if !ok {
		return
	}
	var user db.User
	err := server.database.Where(&db.User{ID: uid}).First(&user).Error
	if err != nil {
		slog.Error("failed to find user", "err", err)
		return
	}
	response := NetworkData{
		"action": "user",
		"uid":    user.ID,
		"name":   user.Name,
	}
	server.senders[sid] <- response
}
