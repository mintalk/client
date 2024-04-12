package network

import (
	"log/slog"
	"mintalk/server/db"
	"time"
)

func (server *Server) HandleRequest(sid string, request NetworkData) {
	action, ok := request["action"]
	if !ok {
		return
	}
	switch action {
	case "message":
		server.ActionMessage(sid, request)
	case "fetchmsg":
		server.ActionFetchMessages(sid, request)
	case "fetchgroup":
		server.ActionFetchGroups(sid, request)
	case "fetchchannel":
		server.ActionFetchChannels(sid, request)
	case "user":
		server.ActionUser(sid, request)
	}
}

func (server *Server) ActionMessage(sid string, data NetworkData) {
	session := server.sessionManager.GetSession(sid)
	if session == nil {
		return
	}
	text, ok := data["text"].(string)
	if !ok {
		slog.Debug("failed to fetch message text", "err", "no text")
		return
	}
	if len(text) == 0 {
		slog.Debug("failed to create message text", "err", "empty text")
		return
	}
	cid, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to fetch message channel", "err", "no channel")
		return
	}
	message := &db.Message{
		UID:     session.User.ID,
		Text:    text,
		Channel: cid,
		Time:    time.Now(),
	}
	server.database.Create(message)
	messageTime, err := message.Time.MarshalText()
	if err != nil {
		slog.Debug("failed to marshal message time", "err", err)
		return
	}
	broadcast := NetworkData{
		"action": "message",
		"mid":    message.ID,
		"text":   message.Text,
		"uid":    session.User.ID,
		"cid":    message.Channel,
		"time":   messageTime,
	}
	server.Broadcast(broadcast)
}

func (server *Server) ActionFetchMessages(sid string, data NetworkData) {
	limit, ok := data["limit"].(int)
	if !ok {
		limit = 0
	}
	channel, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to fetch messages", "err", "no channel")
		return
	}
	var messages []db.Message
	query := server.database.Where(&db.Message{Channel: channel}).Order("time desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&messages).Error
	if err != nil {
		slog.Debug("failed to fetch messages", "err", err)
		return
	}
	responseMessages := make([]string, len(messages))
	for i, message := range messages {
		messageTime, err := message.Time.MarshalText()
		if err != nil {
			slog.Debug("failed to marshal message time", "err", err)
			return
		}
		messageData := NetworkData{
			"mid":  message.ID,
			"uid":  message.UID,
			"cid":  message.Channel,
			"text": message.Text,
			"time": messageTime,
		}
		rawMessageData, err := Encode(messageData)
		if err != nil {
			slog.Debug("failed to encode message data", "err", err)
			return
		}
		responseMessages[i] = string(rawMessageData)
	}
	response := NetworkData{
		"action":   "fetchmsg",
		"messages": responseMessages,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionFetchGroups(sid string, data NetworkData) {
	var groups []db.ChannelGroup
	err := server.database.Find(&groups).Error
	if err != nil {
		slog.Debug("failed to fetch groups", "err", err)
		return
	}
	responseGroups := make([]string, len(groups))
	for i, group := range groups {
		groupData := NetworkData{
			"gid":       group.ID,
			"name":      group.Name,
			"parent":    group.Parent,
			"hasParent": group.HasParent,
		}
		rawGroupData, err := Encode(groupData)
		if err != nil {
			slog.Debug("failed to encode group data", "err", err)
			return
		}
		responseGroups[i] = string(rawGroupData)
	}
	response := NetworkData{
		"action": "fetchgroup",
		"groups": responseGroups,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionFetchChannels(sid string, data NetworkData) {
	var channels []db.Channel
	err := server.database.Find(&channels).Error
	if err != nil {
		slog.Debug("failed to fetch channels", "err", err)
		return
	}
	responseChannels := make([]string, len(channels))
	for i, channel := range channels {
		channelData := NetworkData{
			"cid":   channel.ID,
			"name":  channel.Name,
			"group": channel.Group,
		}
		rawChannelData, err := Encode(channelData)
		if err != nil {
			slog.Debug("failed to encode channel data", "err", err)
			return
		}
		responseChannels[i] = string(rawChannelData)
	}
	response := NetworkData{
		"action":   "fetchchannel",
		"channels": responseChannels,
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
		slog.Debug("failed to find user", "err", err)
		return
	}
	response := NetworkData{
		"action": "user",
		"uid":    user.ID,
		"name":   user.Name,
	}
	server.senders[sid] <- response
}
