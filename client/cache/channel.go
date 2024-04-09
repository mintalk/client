package cache

import (
	"fmt"
	"time"
)

type Message struct {
	Sender   string
	Contents string
	Time     time.Time
}

func (message Message) String() string {
	return fmt.Sprintf("%v@%v %v", message.Sender, message.Time.Format(time.Kitchen), message.Contents)
}

type ChannelCache struct {
	lastUpdated time.Time
	Messages    []Message
}

func NewChannelCache() *ChannelCache {
	return &ChannelCache{time.Now(), make([]Message, 0)}
}

func (cache *ChannelCache) AddMessage(message Message) {
	cache.Messages = append(cache.Messages, message)
}
