package cache

import (
	"fmt"
	"time"
)

type ServerGroup struct {
	Name      string
	Parent    uint
	HasParent bool
}

func (group ServerGroup) String() string {
	return fmt.Sprintf("%s", group.Name)
}

type ServerChannel struct {
	Name  string
	Group uint
}

func (channel ServerChannel) String() string {
	return channel.Name
}

type ServerCache struct {
	lastUpdated   time.Time
	Hostname      string
	Users         map[uint]string
	Groups        map[uint]ServerGroup
	Channels      map[uint]ServerChannel
	ChannelCaches map[uint]*ChannelCache
	Listeners     []func()
}

func NewServerCache() *ServerCache {
	return &ServerCache{
		lastUpdated:   time.Now(),
		Users:         make(map[uint]string),
		Groups:        make(map[uint]ServerGroup),
		Channels:      make(map[uint]ServerChannel),
		ChannelCaches: make(map[uint]*ChannelCache),
	}
}

func (cache *ServerCache) AddListener(listener func()) {
	cache.Listeners = append(cache.Listeners, listener)
}

func (cache *ServerCache) AddUser(uid uint, username string) {
	cache.Users[uid] = username
	for _, listener := range cache.Listeners {
		listener()
	}
}

func (cache *ServerCache) AddGroup(gid uint, group ServerGroup) {
	cache.Groups[gid] = group
	for _, listener := range cache.Listeners {
		listener()
	}
}

func (cache *ServerCache) AddChannel(cid uint, channel ServerChannel) {
	cache.Channels[cid] = channel
	for _, listener := range cache.Listeners {
		listener()
	}
}

func (cache *ServerCache) GetChannelCache(cid uint) *ChannelCache {
	channelCache, ok := cache.ChannelCaches[cid]
	if !ok {
		cache.ChannelCaches[cid] = NewChannelCache()
		channelCache = cache.ChannelCaches[cid]
	}
	return channelCache
}
