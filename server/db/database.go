package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mintalk/server/app"
)

type Connection struct {
	Database *gorm.DB
}

func NewConnetion(conf *app.Config) (*Connection, error) {
	connection := new(Connection)
	var err error
	connection.Database, err = gorm.Open(mysql.Open(conf.Database))
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (connection *Connection) Setup() {
	connection.Database.AutoMigrate(&User{})
}
