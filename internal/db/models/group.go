package models

type Group struct {
	Id   uint64 `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}
