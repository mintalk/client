package input

import (
	"fmt"
	"mintalk/server/db"
)

func (console *Console) channel(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("channel requires an argument")
	}
	switch args[0] {
	case "add":
		return console.channeladd(args[1:])
	case "del":
		return console.channeldel(args[1:])
	case "move":
		return console.channelmove(args[1:])
	case "list":
		return console.channellist(args[1:])
	}
	return fmt.Errorf("channel subcommand not found: %s", args[0])
}

func (console *Console) channeladd(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("channel add requires 2 arguments")
	}
	var group db.ChannelGroup
	err := console.database.Where(db.ChannelGroup{Name: args[1]}).First(&group).Error
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}
	err = console.server.CreateChannel(args[0], group)
	if err != nil {
		err = fmt.Errorf("failed to create channel: %v", err)
	}
	return err
}

func (console *Console) channeldel(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("channel del requires 1 argument")
	}
	channel := db.Channel{Name: args[0]}
	err := console.database.Where(channel).First(&channel).Error
	if err != nil {
		return fmt.Errorf("failed to find channel: %v", err)
	}
	err = console.server.RemoveChannel(channel)
	if err != nil {
		err = fmt.Errorf("failed to delete channel: %v", err)
	}
	return err
}

func (console *Console) channelmove(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("channel move requires 2 arguments")
	}
	channel := db.Channel{Name: args[0]}
	err := console.database.Where(channel).First(&channel).Error
	if err != nil {
		return fmt.Errorf("failed to find channel: %v", err)
	}
	group := db.ChannelGroup{Name: args[1]}
	err = console.database.Where(group).First(&group).Error
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}
	err = console.server.MoveChannel(channel, group)
	if err != nil {
		err = fmt.Errorf("failed to save channel: %v", err)
	}
	return err
}

func (console *Console) channellist(args []string) error {
	var channels []db.Channel
	err := console.database.Find(&channels).Error
	if err != nil {
		return fmt.Errorf("failed to list channels: %v", err)
	}
	fmt.Printf("Name\tGroup\n")
	for _, channel := range channels {
		var group db.ChannelGroup
		err := console.database.Where(db.ChannelGroup{ID: channel.Group}).First(&group).Error
		if err != nil {
			return fmt.Errorf("failed to find group: %v", err)
		}
		fmt.Printf("%v\t%v\n", channel.Name, group.Name)
	}
	return nil
}
