package db

type Channel struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Group uint   `gorm:"default:null"`
}
