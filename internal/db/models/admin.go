package models

type Admin struct {
	Id       int64  `gorm:"primaryKey;autoIncrement:false"`
	Name     string `gorm:"column:name"`
	IsMaster bool   `gorm:"column:is_master;default:false"`
}
