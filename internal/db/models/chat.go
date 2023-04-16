package models

type Chat struct {
	Id         int64    `gorm:"primaryKey;autoIncrement:false"`
	Name       string   `gorm:"column:name"`
	Language   Language `gorm:"foreignKey:LanguageId"`
	LanguageId uint64   `gorm:"column:language_id"`
	Group      Group    `gorm:"foreignKey:GroupId"`
	GroupId    uint64   `gorm:"column:group_id"`
	IsActive   bool     `gorm:"column:is_active;default:false"`
}
