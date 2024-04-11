package db

import (
	"time"
)

type Message struct {
	ID      uint `gorm:"primaryKey"`
	UID     uint
	Text    string
	Channel uint      `gorm:"not null"`
	Time    time.Time `gorm:"column:time"`
}
