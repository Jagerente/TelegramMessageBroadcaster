package controller

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"

	"go.uber.org/zap"
)

type UserService struct {
	cache  *cache.Cache[int64, models.User]
	logger *zap.Logger
	repo   repositories.IRepository[db_models.Admin, int64]
}

func (s *UserService) ClearCache() {
	logger := s.logger.With(
		zap.String("function", "ClearCache"),
	)

	logger.Debug("Clearing cache")

	s.cache.Clear()

	logger.Debug("Cache cleared")
}

func (s *UserService) UpdateCache() error {
	logger := s.logger.With(
		zap.String("function", "UpdateCache"),
	)

	logger.Debug("Updating cache")

	s.ClearCache()
	results, err := s.findAllFromDb()
	if err != nil {
		return err
	}

	for _, result := range results {
		s.cache.Add(result.Id, result)
	}

	logger.Debug("Cache updated")

	return nil
}

func (s *UserService) FindBy(selector string, values ...string) ([]models.User, error) {
	panic("not implemented")
}

func (s *UserService) FindByName(name string) (*models.User, error) {
	panic("not implemented")
}

func (s *UserService) FindById(id int64) (*models.User, error) {
	logger := s.logger.With(
		zap.String("function", "FindById"),
		zap.Int64("id", id),
	)

	logger.Debug("Finding user")

	if result := s.cache.Find(id); result != nil {
		logger.Debug("Found user in cache", zap.Any("user", result))
		return result, nil
	}

	logger.Error("Failed to find user in cache")

	dbResult, err := s.repo.FindById(id)
	if err != nil {
		logger.Error("Failed to find user in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found user in db", zap.Any("chat", dbResult))

	s.cache.Add(dbResult.Id, models.User{Admin: *dbResult})

	return s.cache.Find(id), nil
}

func (s *UserService) Add(user *models.User) (*models.User, error) {
	logger := s.logger.With(
		zap.String("function", "Add"),
		zap.Any("user", user),
	)

	logger.Debug("Adding user")

	value, _ := s.FindById(user.Id)
	if value != nil {
		err := constants.ErrAlreadyExists
		logger.Error("Failed to add user", zap.Error(err))
		return nil, err
	}

	dbResult, err := s.repo.Add(&user.Admin)
	if err != nil {
		logger.Error("Failed to add user", zap.Error(err))
		return nil, err
	}

	logger.Debug("Added user", zap.Any("result", dbResult))

	result := models.User{Admin: *dbResult}

	s.cache.Add(dbResult.Id, result)

	return &result, nil
}

func (s *UserService) Update(user *models.User) (*models.User, error) {
	logger := s.logger.With(
		zap.String("function", "Update"),
		zap.Any("user", user),
	)

	logger.Debug("Updating user")

	result, err := s.repo.Update(&user.Admin)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	logger.Debug("Updated user", zap.Any("result", result))

	s.cache.List.Store(user.Id, *user)

	return user, nil
}

func (s *UserService) Remove(id int64) error {
	logger := s.logger.With(
		zap.String("function", "Remove"),
		zap.Int64("id", id),
	)

	logger.Debug("Removing user")

	usrToDelete, _ := s.FindById(id)
	if usrToDelete == nil {
		err := constants.ErrNotFound
		logger.Error("Failed to remove user", zap.Error(err))
		return err
	}

	if err := s.repo.Remove(usrToDelete.Id); err != nil {
		logger.Error("Failed to remove user", zap.Error(err))
		return err
	}

	logger.Debug("Removed user")

	s.cache.Remove(usrToDelete.Id)

	return nil
}

func (s *UserService) FindAll() ([]models.User, error) {
	logger := s.logger.With(
		zap.String("function", "FindAll"),
	)

	logger.Debug("Finding users")

	result := s.cache.FindAll()

	logger.Debug("Found users in cache")

	return result, nil
}

func (s *UserService) findAllFromDb() ([]models.User, error) {
	logger := s.logger.With(
		zap.String("function", "findAllFromDb"),
	)

	logger.Debug("Finding users")

	dbResults, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Failed to find users in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found users in db", zap.Any("chats", dbResults))

	result := make([]models.User, 0, len(*dbResults))

	for _, admin := range *dbResults {
		result = append(result, models.User{Admin: admin})
	}

	return result, nil
}
