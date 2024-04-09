package db

import (
	"time"
)

type Message struct {
	ID   uint `gorm:"primaryKey"`
	UID  uint
	Text string
	Time time.Time `gorm:"column:time"`
}
