package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Token    string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}
