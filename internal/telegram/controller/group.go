package controller

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"

	"go.uber.org/zap"
)

type GroupService struct {
	logger *zap.Logger
	repo   repositories.IRepository[db_models.Group, uint64]
	cache  *models.List[string, models.Group]
}

func (s *GroupService) FindBy(selector string, values ...string) ([]models.Group, error) {
	panic("not implemented")
}

func (s *GroupService) FindByName(name string) (*models.Group, error) {
	logger := s.logger.With(
		zap.String("function", "FindByName"),
		zap.String("name", name),
	)

	logger.Debug("Started")

	if result := s.cache.Find(name); result != nil {
		logger.Debug("Found group in cache", zap.Any("group", result))
		return result, nil
	}

	dbResults, err := s.repo.FindBy("name = ?", name)
	if err != nil {
		logger.Error("Failed to find group in db", zap.Error(err))
		return nil, err
	}

	dbResult := (*dbResults)[0]

	logger.Debug("Found group in db", zap.Any("group", dbResult))

	s.cache.Add(name, models.Group(dbResult))

	return s.cache.Find(name), nil
}

func (s *GroupService) FindById(id uint64) (*models.Group, error) {
	logger := s.logger.With(
		zap.String("function", "FindById"),
		zap.Uint64("id", id),
	)

	logger.Debug("Started")

	dbResult, err := s.repo.FindById(id)
	if err != nil {
		logger.Error("Failed to find group in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found group in db", zap.Any("group", dbResult))

	result := models.Group(*dbResult)

	return &result, nil
}

func (s *GroupService) Add(group *models.Group) (*models.Group, error) {
	logger := s.logger.With(
		zap.String("function", "Add"),
		zap.Any("group", group),
	)

	logger.Debug("Started")

	value, _ := s.FindByName(group.Name)
	if value != nil {
		err := constants.ErrAlreadyExists
		logger.Error("Failed to add group", zap.Error(err))
		return nil, err
	}

	dbValue := db_models.Group{Name: group.Name}
	dbResult, err := s.repo.Add(&dbValue)
	if err != nil {
		logger.Error("Failed to add group", zap.Error(err))
		return nil, err
	}

	logger.Debug("Added group", zap.Any("result", dbResult))

	result := models.Group(*dbResult)

	s.cache.Add(dbResult.Name, result)

	return &result, nil
}

func (s *GroupService) Update(group *models.Group) (*models.Group, error) {
	logger := s.logger.With(
		zap.String("function", "Update"),
		zap.Any("group", group),
	)

	logger.Debug("Started")

	dbGroup := db_models.Group(*group)
	result, err := s.repo.Update(&dbGroup)
	if err != nil {
		logger.Error("Failed to update group", zap.Error(err))
		return nil, err
	}

	logger.Debug("Updated group", zap.Any("result", result))

	s.cache.List.Store(group.Id, *group)
	return group, nil
}

func (s *GroupService) Remove(id uint64) error {
	logger := s.logger.With(
		zap.String("function", "Remove"),
		zap.Uint64("id", id),
	)

	logger.Debug("Started")

	groupToDelete, _ := s.FindById(id)
	if groupToDelete == nil {
		err := constants.ErrNotFound
		logger.Error("Failed to remove group", zap.Error(err))
		return err
	}

	if err := s.repo.Remove(groupToDelete.Id); err != nil {
		logger.Error("Failed to remove group", zap.Error(err))
		return err
	}

	logger.Debug("Removed group")

	s.cache.Remove(groupToDelete.Name)

	return nil
}

func (s *GroupService) FindAll() ([]models.Group, error) {
	logger := s.logger.With(
		zap.String("function", "FindAll"),
	)

	logger.Debug("Started")

	dbResults, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Failed to find groups in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found groups in db", zap.Any("groups", dbResults))

	result := make([]models.Group, 0, len(*dbResults))

	for _, group := range *dbResults {
		result = append(result, models.Group(group))
	}

	return result, nil
}
