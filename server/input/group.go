package input

import (
	"fmt"
	"mintalk/server/db"
	"strings"
)

type GroupListNode struct {
	Group    db.ChannelGroup
	Children []*GroupListNode
}

func (node GroupListNode) Print(indent int) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", indent), node.Group.Name)
	for _, child := range node.Children {
		child.Print(indent + 2)
	}
}

func (console *Console) group(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("group requires an argument")
	}
	switch args[0] {
	case "add":
		return console.groupadd(args[1:])
	case "del":
		return console.groupdel(args[1:])
	case "move":
		return console.groupmove(args[1:])
	case "root":
		return console.grouproot(args[1:])
	case "list":
		return console.grouplist(args[1:])
	}
	return fmt.Errorf("group subcommand not found: %s", args[0])
}

func (console *Console) groupadd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("group add requires 1 argument")
	}
	err := console.server.CreateGroup(args[0])
	if err != nil {
		err = fmt.Errorf("failed to create group: %v", err)
	}
	return err
}

func (console *Console) groupdel(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("group del requires 1 argument")
	}
	group := db.ChannelGroup{Name: args[0]}
	err := console.database.Where(group).First(&group).Error
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}
	err = console.server.RemoveGroup(group)
	if err != nil {
		err = fmt.Errorf("failed to delete group: %v", err)
	}
	return err
}

func (console *Console) groupmove(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("group move requires 2 arguments")
	}
	group := db.ChannelGroup{Name: args[0]}
	err := console.database.Where(group).First(&group).Error
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}
	parent := db.ChannelGroup{Name: args[1]}
	err = console.database.Where(parent).First(&parent).Error
	if err != nil {
		return fmt.Errorf("failed to find parent group: %v", err)
	}
	err = console.server.MoveGroup(group, parent)
	if err != nil {
		err = fmt.Errorf("failed to save group: %v", err)
	}
	return err
}

func (console *Console) grouproot(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("group move requires 1 argument")
	}
	group := db.ChannelGroup{Name: args[0]}
	err := console.database.Where(group).First(&group).Error
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}
	err = console.server.RootGroup(group)
	if err != nil {
		err = fmt.Errorf("failed to save group: %v", err)
	}
	return err
}

func (console *Console) grouplist(args []string) error {
	var groups []db.ChannelGroup
	err := console.database.Find(&groups).Error
	if err != nil {
		return fmt.Errorf("failed to list groups: %v", err)
	}
	groupMap := make(map[uint]*GroupListNode)
	rootGroup := make([]*GroupListNode, 0)
	groupsLeft := make([]db.ChannelGroup, len(groups))
	copy(groupsLeft, groups)
	for len(groupsLeft) > 0 {
		newGroupsLeft := make([]db.ChannelGroup, 0)
		for _, group := range groups {
			node := &GroupListNode{Group: group, Children: make([]*GroupListNode, 0)}
			if !group.HasParent {
				rootGroup = append(rootGroup, node)
			} else {
				parent, ok := groupMap[group.Parent]
				if ok {
					parent.Children = append(parent.Children, node)
				} else {
					newGroupsLeft = append(newGroupsLeft, group)
				}
			}
			groupMap[group.ID] = node
		}
		groupsLeft = newGroupsLeft
	}
	for _, group := range rootGroup {
		group.Print(0)
	}
	return nil
}
