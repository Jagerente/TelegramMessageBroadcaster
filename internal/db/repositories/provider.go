package repositories

import (
	"DC_NewsSender/internal/db/models"

	"gorm.io/gorm"
)

type Provider struct {
	gormConnection *gorm.DB
}

func CreateProvider(connection *gorm.DB) *Provider {
	return &Provider{
		gormConnection: connection,
	}
}

func (provider *Provider) CreateGroupRepo() IRepository[models.Group, uint64] {
	repo := &Repository[models.Group, uint64]{
		BaseRepository{
			gormConnection: provider.gormConnection,
		},
	}

	return repo
}

func (provider *Provider) CreateChatRepo() IRepository[models.Chat, int64] {
	repo := &Repository[models.Chat, int64]{
		BaseRepository{
			gormConnection: provider.gormConnection,
		},
	}
	return repo
}

func (provider *Provider) CreateLanguageRepo() IRepository[models.Language, uint64] {
	repo := &Repository[models.Language, uint64]{
		BaseRepository{
			gormConnection: provider.gormConnection,
		},
	}
	return repo
}

func (provider *Provider) CreateAdminsRepo() IRepository[models.Admin, int64] {
	repo := &Repository[models.Admin, int64]{
		BaseRepository{
			gormConnection: provider.gormConnection,
		},
	}
	return repo
}

type IRepository[T any, K comparable] interface {
	FindById(id K) (*T, error)
	FindBy(key string, value string) (*T, error)
	FindAll() (*[]T, error)
	Add(value *T) (*T, error)
	Update(value *T) (*T, error)
	Remove(id K) error
}
