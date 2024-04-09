package db

type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}
