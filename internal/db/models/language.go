package models

type Language struct {
	Id   uint64 `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}
