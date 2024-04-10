package cache

import "time"

type ServerCache struct {
	lastUpdated time.Time
	Users       map[uint]string
	Listeners   []func()
}

func NewServerCache() *ServerCache {
	return &ServerCache{
		lastUpdated: time.Now(),
		Users:       make(map[uint]string),
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
