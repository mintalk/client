package network

import (
	"log/slog"
	"mintalk/server/db"
	"time"
)

func (server *Server) ActionMessage(sid string, data map[string]interface{}) {
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
	messageTime, err := message.Time.MarshalText()
	if err != nil {
		slog.Error("failed to marshal message time", "err", err)
		return
	}
	broadcast := map[string]interface{}{
		"action": "message",
		"mid":    message.ID,
		"text":   message.Text,
		"uid":    session.User.ID,
		"time":   messageTime,
	}
	server.Broadcast(broadcast)
}

func (server *Server) ActionFetchMessages(sid string, data map[string]interface{}) {
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
	responseMessages := make([]string, len(messages))
	for i, message := range messages {
		messageTime, err := message.Time.MarshalText()
		if err != nil {
			slog.Error("failed to marshal message time", "err", err)
			return
		}
		messageData := map[string]interface{}{
			"mid":  message.ID,
			"uid":  message.UID,
			"text": message.Text,
			"time": messageTime,
		}
		rawMessageData, err := Encode(messageData)
		if err != nil {
			slog.Error("failed to encode message data", "err", err)
			return
		}
		responseMessages[i] = string(rawMessageData)
	}
	response := map[string]interface{}{
		"action":   "fetchmsg",
		"messages": responseMessages,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionFetchGroups(sid string, data map[string]interface{}) {
	var groups []db.ChannelGroup
	err := server.database.Find(&groups).Error
	if err != nil {
		slog.Error("failed to fetch groups", "err", err)
		return
	}
	responseGroups := make([]string, len(groups))
	for i, group := range groups {
		groupData := map[string]interface{}{
			"gid":       group.ID,
			"name":      group.Name,
			"parent":    group.Parent,
			"hasParent": group.HasParent,
		}
		rawGroupData, err := Encode(groupData)
		if err != nil {
			slog.Error("failed to encode group data", "err", err)
			return
		}
		responseGroups[i] = string(rawGroupData)
	}
	response := map[string]interface{}{
		"action": "fetchgroup",
		"groups": responseGroups,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionFetchChannels(sid string, data map[string]interface{}) {
	var channels []db.Channel
	err := server.database.Find(&channels).Error
	if err != nil {
		slog.Error("failed to fetch channels", "err", err)
		return
	}
	responseChannels := make([]string, len(channels))
	for i, channel := range channels {
		channelData := map[string]interface{}{
			"cid":   channel.ID,
			"name":  channel.Name,
			"group": channel.Group,
		}
		rawChannelData, err := Encode(channelData)
		if err != nil {
			slog.Error("failed to encode channel data", "err", err)
			return
		}
		responseChannels[i] = string(rawChannelData)
	}
	response := map[string]interface{}{
		"action":   "fetchchannel",
		"channels": responseChannels,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionUser(sid string, data map[string]interface{}) {
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
	response := map[string]interface{}{
		"action": "user",
		"uid":    user.ID,
		"name":   user.Name,
	}
	server.senders[sid] <- response
}
