package db

type ChannelGroup struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Parent    uint   `gorm:"default:null"`
	HasParent bool
}
