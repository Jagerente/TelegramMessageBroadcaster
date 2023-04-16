package repositories

import (
	"gorm.io/gorm"
)

type BaseRepository struct {
	gormConnection *gorm.DB
}

func (repo BaseRepository) GetConnection() *gorm.DB {
	return repo.gormConnection
}
