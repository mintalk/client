package db

type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"not null,unique"`
	Password string
	Operator bool `gorm:"not null"`
}
