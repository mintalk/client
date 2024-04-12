package panels

import (
	gc "github.com/rthornton128/goncurses"
	"mintalk/client/cache"
	"mintalk/client/network"
	"mintalk/client/ui/elements"
)

type ChannelsPanel struct {
	*elements.Panel
	Connector    *network.Connector
	serverCache  *cache.ServerCache
	tree         *elements.Tree
	channelPanel *ChannelPanel
}

func NewChannelsPanel(connector *network.Connector, channelPanel *ChannelPanel, serverCache *cache.ServerCache) (*ChannelsPanel, error) {
	panel, err := elements.NewPanel(1, 1)
	channelsPanel := &ChannelsPanel{Panel: panel, Connector: connector, channelPanel: channelPanel, serverCache: serverCache}
	channelsPanel.tree = elements.NewTree()
	channelsPanel.tree.Expand(1)
	channelsPanel.Add(channelsPanel.tree)
	channelsPanel.serverCache.AddListener(channelsPanel.updateTreeData)
	connector.LoadGroups()
	connector.LoadChannels()
	return channelsPanel, err
}

func (panel *ChannelsPanel) Resize() error {
	width, height := panel.RealWidth, panel.RealHeight
	panel.Panel.Resize(width, height)
	panel.tree.Move(1, 1)
	panel.tree.Resize(width-2, height-2)
	return nil
}

func (panel *ChannelsPanel) Update(key gc.Key) {
	panel.tree.Active = panel.Panel.Active
	panel.Panel.Update(key)
}

func (panel *ChannelsPanel) updateTreeData() {
	panel.tree.Nodes = make([]*elements.TreeNode, 0)
	groupsLeft := make([]uint, 0)
	for gid := range panel.serverCache.Groups {
		groupsLeft = append(groupsLeft, gid)
	}
	nodeAssignments := make(map[uint]*elements.TreeNode)
	for len(groupsLeft) > 0 {
		newGroupsLeft := make([]uint, 0)
		for _, gid := range groupsLeft {
			group := panel.serverCache.Groups[gid]
			if !group.HasParent {
				node := elements.NewTreeNode(elements.TreeItem{Value: group})
				panel.tree.Nodes = append(panel.tree.Nodes, node)
				nodeAssignments[gid] = node
				continue
			}
			parent, ok := nodeAssignments[group.Parent]
			if !ok {
				newGroupsLeft = append(newGroupsLeft, gid)
				continue
			}
			node := elements.NewTreeNode(elements.TreeItem{Value: group})
			parent.Children = append(parent.Children, node)
			nodeAssignments[gid] = node
		}
		groupsLeft = newGroupsLeft
	}
	for cid, channel := range panel.serverCache.Channels {
		if !panel.channelPanel.ChannelOpened {
			panel.channelPanel.MoveChannel(cid)
		}
		node := elements.NewTreeNode(elements.TreeItem{Value: channel, OnClick: func() {
			panel.channelPanel.MoveChannel(cid)
		}})
		parent, ok := nodeAssignments[channel.Group]
		if !ok {
			continue
		}
		parent.Children = append(parent.Children, node)
	}
}
