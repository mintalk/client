package network

import (
	"log/slog"
	"mintalk/client/cache"
	"time"
)

func (connector *Connector) HandleResponse(data NetworkData) {
	switch data["action"].(string) {
	case "messages":
		connector.ResponseMessages(data)
	case "user":
		connector.ResponseUser(data)
	case "users":
		connector.ResponseUsers(data)
	case "groups":
		connector.ResponseGroups(data)
	case "channels":
		connector.ResponseChannels(data)
	case "update":
		connector.ResponseUpdate(data)
	}
}

func (connector *Connector) ResponseMessages(data NetworkData) {
	messages, ok := data["messages"].([]string)
	if !ok {
		slog.Debug("failed to parse messages")
		return
	}
	cid, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to parse cid")
		return
	}
	lastMid, ok := data["last-mid"].(uint)
	if !ok {
		slog.Debug("failed to parse last-mid")
		return
	}
	checkLastMid, ok := data["check-last-mid"].(bool)
	if !ok {
		checkLastMid = false
	}
	channelCache := connector.serverCache.GetChannelCache(cid)

	if checkLastMid {
		lastMidFound := false
		for mid := range channelCache.Messages {
			if mid == lastMid {
				lastMidFound = true
				break
			}
		}
		if !lastMidFound {
			connector.LoadMessages(0, cid)
		}
	}

	messageMap := make(map[uint]cache.Message)
	for _, messageData := range messages {
		message, err := Decode([]byte(messageData))
		if err != nil {
			slog.Debug("failed to decode message", "err", err)
			continue
		}
		rawTime, ok := message["time"].([]byte)
		if !ok {
			slog.Debug("failed to parse time")
			continue
		}
		var messageTime time.Time
		err = messageTime.UnmarshalText(rawTime)
		if err != nil {
			slog.Debug("failed to parse time", "err", err)
			continue
		}
		uid, ok := message["uid"].(uint)
		if !ok {
			slog.Debug("failed to parse uid")
			continue
		}
		contents, ok := message["contents"].(string)
		if !ok {
			slog.Debug("failed to parse contents")
			continue
		}
		mid, ok := message["mid"].(uint)
		if !ok {
			slog.Debug("failed to parse mid")
			continue
		}
		messageItem := cache.Message{
			Sender:   uid,
			Contents: contents,
			Time:     messageTime,
			Username: "",
		}
		messageMap[mid] = messageItem
	}
	channelCache.AddMessages(messageMap)
}

func (connector *Connector) ResponseUser(data NetworkData) {
	uid, ok := data["uid"].(uint)
	if !ok {
		slog.Debug("failed to parse uid")
		return
	}
	name, ok := data["name"].(string)
	if !ok {
		slog.Debug("failed to parse name")
		return
	}
	connector.serverCache.AddUser(uid, name)
}

func (connector *Connector) ResponseUsers(data NetworkData) {
	users, ok := data["users"].([]string)
	if !ok {
		slog.Debug("failed to parse users")
		return
	}
	userMap := make(map[uint]string)
	for _, userData := range users {
		user, err := Decode([]byte(userData))
		if err != nil {
			slog.Debug("failed to decode user", "err", err)
			continue
		}
		uid, ok := user["uid"].(uint)
		if !ok {
			slog.Debug("failed to parse uid")
			continue
		}
		name, ok := user["name"].(string)
		if !ok {
			slog.Debug("failed to parse name")
			continue
		}
		userMap[uid] = name
	}
	connector.serverCache.AddUsers(userMap)
}

func (connector *Connector) ResponseGroups(data NetworkData) {
	groups, ok := data["groups"].([]string)
	if !ok {
		slog.Debug("failed to parse groups")
		return
	}
	groupMap := make(map[uint]cache.ServerGroup)
	for _, groupData := range groups {
		group, err := Decode([]byte(groupData))
		if err != nil {
			slog.Debug("failed to decode group", "err", err)
			continue
		}
		gid, ok := group["gid"].(uint)
		if !ok {
			slog.Debug("failed to parse gid")
			continue
		}
		name, ok := group["name"].(string)
		if !ok {
			slog.Debug("failed to parse name")
			continue
		}
		parent, ok := group["parent"].(uint)
		if !ok {
			slog.Debug("failed to parse parent")
			continue
		}
		hasParent, ok := group["has-parent"].(bool)
		if !ok {
			slog.Debug("failed to parse has-parent")
			continue
		}
		groupMap[gid] = cache.ServerGroup{
			Name: name, Parent: parent, HasParent: hasParent,
		}
	}
	connector.serverCache.AddGroups(groupMap)
}

func (connector *Connector) ResponseChannels(data NetworkData) {
	channels, ok := data["channels"].([]string)
	if !ok {
		slog.Debug("failed to parse channels")
		return
	}
	for _, channelData := range channels {
		channel, err := Decode([]byte(channelData))
		if err != nil {
			slog.Error("failed to decode channel", "err", err)
			continue
		}
		cid, ok := channel["cid"].(uint)
		if !ok {
			slog.Debug("failed to parse cid")
			continue
		}
		name, ok := channel["name"].(string)
		if !ok {
			slog.Debug("failed to parse name")
			continue
		}
		group, ok := channel["group"].(uint)
		if !ok {
			slog.Debug("failed to parse group")
			continue
		}
		connector.serverCache.AddChannel(cid, cache.ServerChannel{
			Name: name, Group: group,
		})
	}
}

func (connector *Connector) ResponseUpdate(data NetworkData) {
	updateType, ok := data["type"].(string)
	if !ok {
		slog.Debug("failed to parse update type")
		return
	}
	switch updateType {
	case "message":
		cid, ok := data["cid"].(uint)
		if !ok {
			slog.Debug("failed to parse cid")
			return
		}
		connector.LoadMessages(1, cid)
	case "user":
		connector.LoadUsers()
	case "group":
		connector.LoadGroups()
	case "channel":
		connector.LoadChannels()
	}
}
