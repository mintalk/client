package panels

import (
	"mintalk/client/cache"
	"mintalk/client/network"

	"github.com/rivo/tview"
)

type Channels struct {
	*tview.TreeView
	Connector    *network.Connector
	serverCache  *cache.ServerCache
	root *tview.TreeNode
	channel *Channel
	app *tview.Application
}

func NewChannels(app *tview.Application, connector *network.Connector, channel *Channel, serverCache *cache.ServerCache) *Channels {
	channels := &Channels{
		TreeView: tview.NewTreeView(),
		Connector: connector,
		channel: channel,
		root: tview.NewTreeNode("Remote"),
		serverCache: serverCache,
		app: app,
	}
	channels.SetBorder(true)
	channels.SetRoot(channels.root)
	
	channels.serverCache.AddListener(channels.updateTreeData)

	channels.SetSelectedFunc(func (node *tview.TreeNode) {
		switch node.GetReference().(type) {
		case uint:
			cid := node.GetReference().(uint)
			channels.channel.MoveChannel(cid)
		}
	})

	connector.LoadGroups()
	connector.LoadChannels()
	return channels
}

func (channels *Channels) updateTreeData() {
	channels.app.QueueUpdateDraw(func () {
		channels.root.SetText(channels.serverCache.Hostname)
		channels.root.ClearChildren()
	})
	nodes := make([]*tview.TreeNode, 0)
	groupsLeft := make([]uint, 0)
	for gid := range channels.serverCache.Groups {
		groupsLeft = append(groupsLeft, gid)
	}
	nodeAssignments := make(map[uint]*tview.TreeNode)
	for len(groupsLeft) > 0 {
		newGroupsLeft := make([]uint, 0)
		for _, gid := range groupsLeft {
			group := channels.serverCache.Groups[gid]
			if !group.HasParent {
				node := tview.NewTreeNode(group.String()).SetReference(group)
				nodes = append(nodes, node)
				nodeAssignments[gid] = node
				continue
			}
			parent, ok := nodeAssignments[group.Parent]
			if !ok {
				newGroupsLeft = append(newGroupsLeft, gid)
				continue
			}
			node := tview.NewTreeNode(group.String()).SetReference(group)
			parent.AddChild(node)
			nodeAssignments[gid] = node
		}
		groupsLeft = newGroupsLeft
	}
	for cid, channel := range channels.serverCache.Channels {
		if !channels.channel.ChannelOpened {
			channels.channel.MoveChannel(cid)
		}
		node := tview.NewTreeNode(channel.String()).SetReference(cid)
		parent, ok := nodeAssignments[channel.Group]
		if !ok {
			continue
		}
		parent.AddChild(node)
	}
	for _, node := range nodes {
		channels.app.QueueUpdateDraw(func () {
			channels.root.AddChild(node)
		})
	}
}
