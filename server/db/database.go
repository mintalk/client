package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mintalk/server/config"
)

type Connection struct {
	*gorm.DB
}

func NewConnection(conf *config.Config) (*Connection, error) {
	database, err := gorm.Open(mysql.Open(conf.Database))
	if err != nil {
		return nil, err
	}
	return &Connection{database}, nil
}

func (connection *Connection) Setup() error {
	return connection.AutoMigrate(&User{}, &Message{}, &Channel{}, &ChannelGroup{})
}
