package network

import (
	"log/slog"
	"mintalk/server/db"
)

func (server *Server) HandleRequest(sid string, request NetworkData) {
	action, ok := request["action"]
	if !ok {
		return
	}
	switch action {
	case "new-message":
		server.ActionNewMessage(sid, request)
	case "messages":
		server.ActionMessages(sid, request)
	case "groups":
		server.ActionGroups(sid, request)
	case "channels":
		server.ActionChannels(sid, request)
	case "user":
		server.ActionUser(sid, request)
	case "users":
		server.ActionUsers(sid, request)
	}
}

func (server *Server) ActionNewMessage(sid string, data NetworkData) {
	session := server.sessionManager.GetSession(sid)
	if session == nil {
		return
	}
	contents, ok := data["contents"].(string)
	if !ok {
		slog.Debug("failed to fetch message contents", "err", "no contents")
		return
	}
	if len(contents) == 0 {
		slog.Debug("failed to create message text", "err", "empty text")
		return
	}
	cid, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to fetch message channel", "err", "no channel")
		return
	}
	if err := server.CreateMessage(session.User.ID, cid, contents); err != nil {
		slog.Debug("failed to create message", "err", err)
		return
	}
}

func (server *Server) ActionMessages(sid string, data NetworkData) {
	limit, ok := data["limit"].(int)
	if !ok {
		limit = 0
	}
	cid, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to fetch messages", "err", "no channel")
		return
	}
	var messages []db.Message
	query := server.database.Where(&db.Message{Channel: cid}).Order("time desc")
	if limit > 0 {
		query = query.Limit(limit + 1)
	}
	err := query.Find(&messages).Error
	if err != nil {
		slog.Debug("failed to fetch messages", "err", err)
		return
	}
	responseMessages := make([]string, len(messages))
	for i, message := range messages {
		if i == 0 && limit > 0 {
			continue
		}
		messageTime, err := message.Time.MarshalText()
		if err != nil {
			slog.Debug("failed to marshal message time", "err", err)
			return
		}
		messageData := NetworkData{
			"mid":      message.ID,
			"uid":      message.UID,
			"contents": message.Contents,
			"time":     messageTime,
		}
		rawMessageData, err := Encode(messageData)
		if err != nil {
			slog.Debug("failed to encode message data", "err", err)
			return
		}
		responseMessages[i] = string(rawMessageData)
	}
	lastMid := uint(0)
	if len(messages) > 0 {
		lastMid = messages[0].ID
	}
	response := NetworkData{
		"action":         "messages",
		"messages":       responseMessages,
		"cid":            cid,
		"last-mid":       lastMid,
		"check-last-mid": len(messages) > limit,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionGroups(sid string, data NetworkData) {
	var groups []db.ChannelGroup
	err := server.database.Find(&groups).Error
	if err != nil {
		slog.Debug("failed to fetch groups", "err", err)
		return
	}
	responseGroups := make([]string, len(groups))
	for i, group := range groups {
		groupData := NetworkData{
			"gid":        group.ID,
			"name":       group.Name,
			"parent":     group.Parent,
			"has-parent": group.HasParent,
		}
		rawGroupData, err := Encode(groupData)
		if err != nil {
			slog.Debug("failed to encode group data", "err", err)
			return
		}
		responseGroups[i] = string(rawGroupData)
	}
	response := NetworkData{
		"action": "groups",
		"groups": responseGroups,
	}
	server.senders[sid] <- response
}

func (server *Server) ActionChannels(sid string, data NetworkData) {
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
		"action":   "channels",
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

func (server *Server) ActionUsers(sid string, data NetworkData) {
	var users []db.User
	err := server.database.Find(&users).Error
	if err != nil {
		slog.Debug("failed to fetch users", "err", err)
		return
	}
	responseUsers := make([]string, len(users))
	for i, user := range users {
		userData := NetworkData{
			"uid":  user.ID,
			"name": user.Name,
		}
		rawUserData, err := Encode(userData)
		if err != nil {
			slog.Debug("failed to encode user data", "err", err)
			return
		}
		responseUsers[i] = string(rawUserData)
	}
	response := NetworkData{
		"action": "users",
		"users":  responseUsers,
	}
	server.senders[sid] <- response
}
