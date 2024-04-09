package cache

import "time"

type ServerCache struct {
	lastUpdated time.Time
	Users       map[uint]string
}

func NewServerCache() *ServerCache {
	return &ServerCache{
		lastUpdated: time.Now(),
		Users:       make(map[uint]string),
	}
}

func (cache *ServerCache) AddUser(uid uint, username string) {
	cache.Users[uid] = username
}
