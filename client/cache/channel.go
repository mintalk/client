package cache

import (
	"fmt"
	"sort"
	"time"
)

type Message struct {
	Sender   uint
	Contents string
	Time     time.Time
}

func (message Message) String() string {
	return fmt.Sprintf("%v@%v %v", message.Sender, message.Time.Format(time.Kitchen), message.Contents)
}

type ChannelCache struct {
	lastUpdated time.Time
	Messages    map[uint]Message
	Listeners   []func()
}

func NewChannelCache() *ChannelCache {
	return &ChannelCache{time.Now(), make(map[uint]Message, 0), make([]func(), 0)}
}

func (cache *ChannelCache) AddMessage(mid uint, message Message) {
	cache.Messages[mid] = message
	for _, listener := range cache.Listeners {
		listener()
	}
}

func (cache *ChannelCache) AddListener(listener func()) {
	cache.Listeners = append(cache.Listeners, listener)
}

func (cache *ChannelCache) GetMessages() []Message {
	mids := make([]uint, 0)
	for mid := range cache.Messages {
		mids = append(mids, mid)
	}
	sort.Slice(mids, func(i, j int) bool {
		return cache.Messages[mids[i]].Time.Before(cache.Messages[mids[j]].Time)
	})
	messages := make([]Message, len(mids))
	for idx, mid := range mids {
		messages[idx] = cache.Messages[mid]
	}
	return messages
}
