package network

import (
	"log/slog"
	"mintalk/server/db"
	"time"
)

func (server *Server) CreateMessage(uid, cid uint, contents string) error {
	defer server.SendUpdate("message", NetworkData{"cid": cid})
	err := server.database.Where(&db.Channel{ID: cid}).First(&db.Channel{}).Error
	if err != nil {
		return err
	}
	message := &db.Message{
		UID:      uid,
		Contents: contents,
		Channel:  cid,
		Time:     time.Now(),
	}
	return server.database.Create(message).Error
}

func (server *Server) CreateChannel(name string, group db.ChannelGroup) error {
	defer server.SendUpdate("channel")
	channel := db.Channel{Name: name, Group: group.ID}
	return server.database.Create(&channel).Error
}

func (server *Server) MoveChannel(channel db.Channel, group db.ChannelGroup) error {
	defer server.SendUpdate("channel")
	channel.Group = group.ID
	return server.database.Save(&channel).Error
}

func (server *Server) RenameChannel(channel db.Channel, name string) error {
	defer server.SendUpdate("channel")
	channel.Name = name
	return server.database.Save(&channel).Error
}

func (server *Server) RemoveChannel(channel db.Channel) error {
	defer server.SendUpdate("channel")
	return server.database.Delete(&channel).Error
}

func (server *Server) CreateGroup(name string) error {
	defer server.SendUpdate("group")
	group := db.ChannelGroup{Name: name}
	return server.database.Create(&group).Error
}

func (server *Server) MoveGroup(group, parent db.ChannelGroup) error {
	defer server.SendUpdate("group")
	group.Parent = parent.ID
	return server.database.Save(&group).Error
}

func (server *Server) RootGroup(group db.ChannelGroup) error {
	defer server.SendUpdate("group")
	group.Parent = 0
	group.HasParent = false
	return server.database.Save(&group).Error
}

func (server *Server) RenameGroup(group db.ChannelGroup, name string) error {
	defer server.SendUpdate("group")
	group.Name = name
	return server.database.Save(&group).Error
}

func (server *Server) RemoveGroup(group db.ChannelGroup) error {
	defer server.SendUpdate("group")
	defer server.SendUpdate("channel")
	var children []db.ChannelGroup
	err := server.database.Where(&db.ChannelGroup{Parent: group.ID}).Find(&children).Error
	if err != nil {
		return err
	}
	for _, child := range children {
		err = server.RemoveGroup(child)
		if err != nil {
			slog.Debug("failed to remove group: %v", err)
			continue
		}
	}
	var childChannels []db.Channel
	err = server.database.Where(&db.Channel{Group: group.ID}).Find(&childChannels).Error
	if err != nil {
		return err
	}
	for _, childChannel := range childChannels {
		err = server.RemoveChannel(childChannel)
		if err != nil {
			slog.Debug("failed to remove channel: %v", err)
			continue
		}
	}
	return server.database.Delete(&group).Error
}

func (server *Server) CreateUser(name string) error {
	defer server.SendUpdate("user")
	user := db.User{Name: name, Password: ""}
	return server.database.Create(&user).Error
}

func (server *Server) RenameUser(user db.User, name string) error {
	defer server.SendUpdate("user")
	user.Name = name
	return server.database.Save(&user).Error
}

func (server *Server) UserResetPassword(user db.User) error {
	user.Password = ""
	return server.database.Save(&user).Error
}

func (server *Server) RemoveUser(user db.User) error {
	defer server.SendUpdate("user")
	server.database.Where(&db.Message{UID: user.ID}).Delete(&db.Message{})
	return server.database.Delete(&user).Error
}

func (server *Server) SendUpdate(updateType string, data ...NetworkData) {
	broadcast := NetworkData{
		"action": "update",
		"type":   updateType,
	}
	for _, d := range data {
		for k, v := range d {
			broadcast[k] = v
		}
	}
	go server.Broadcast(broadcast)
}
