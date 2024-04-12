package network

import (
	"log/slog"
	"mintalk/client/cache"
	"time"
)

func (connector *Connector) HandleResponse(data NetworkData) {
	switch data["action"].(string) {
	case "message":
		connector.ResponseMessage(data)
	case "user":
		connector.ResponseUser(data)
	case "fetchmsg":
		connector.ResponseFetchMessages(data)
	case "fetchgroup":
		connector.ResponseFetchGroups(data)
	case "fetchchannel":
		connector.ResponseFetchChannels(data)
	}
}

func (connector *Connector) ResponseMessage(data NetworkData) {
	rawTime, ok := data["time"].([]byte)
	if !ok {
		slog.Debug("failed to parse time")
		return
	}
	var messageTime time.Time
	err := messageTime.UnmarshalText(rawTime)
	if err != nil {
		slog.Debug("failed to parse time", "err", err)
		return
	}
	uid, ok := data["uid"].(uint)
	if !ok {
		slog.Debug("failed to parse uid")
		return
	}
	contents, ok := data["text"].(string)
	if !ok {
		slog.Debug("failed to parse contents")
		return
	}
	mid, ok := data["mid"].(uint)
	if !ok {
		slog.Debug("failed to parse mid")
		return
	}
	cid, ok := data["cid"].(uint)
	if !ok {
		slog.Debug("failed to parse cid")
		return
	}
	message := cache.Message{
		Sender:   uid,
		Contents: contents,
		Time:     messageTime,
	}
	connector.serverCache.GetChannelCache(cid).AddMessage(mid, message)
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

func (connector *Connector) ResponseFetchMessages(data NetworkData) {
	messages, ok := data["messages"].([]string)
	if !ok {
		slog.Debug("failed to parse messages")
		return
	}
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
		contents, ok := message["text"].(string)
		if !ok {
			slog.Debug("failed to parse contents")
			continue
		}
		cid, ok := message["cid"].(uint)
		if !ok {
			slog.Debug("failed to parse cid")
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
		}
		connector.serverCache.GetChannelCache(cid).AddMessage(mid, messageItem)
	}
}

func (connector *Connector) ResponseFetchGroups(data NetworkData) {
	groups, ok := data["groups"].([]string)
	if !ok {
		slog.Debug("failed to parse groups")
		return
	}
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
		hasParent, ok := group["hasParent"].(bool)
		if !ok {
			slog.Debug("failed to parse hasParent")
			continue
		}
		connector.serverCache.AddGroup(gid, cache.ServerGroup{
			Name: name, Parent: parent, HasParent: hasParent,
		})
	}
}

func (connector *Connector) ResponseFetchChannels(data NetworkData) {
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
