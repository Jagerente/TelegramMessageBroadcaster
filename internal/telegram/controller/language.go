package controller

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"

	"go.uber.org/zap"
)

type LanguageService struct {
	cache  *cache.Cache[string, models.Language]
	logger *zap.Logger
	repo   repositories.IRepository[db_models.Language, uint64]
}

func (s *LanguageService) ClearCache() {
	logger := s.logger.With(
		zap.String("function", "ClearCache"),
	)

	logger.Debug("Clearing cache")

	s.cache.Clear()

	logger.Debug("Cache cleared")
}

func (s *LanguageService) UpdateCache() error {
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
		s.cache.Add(result.Name, result)
	}

	logger.Debug("Cache updated")

	return nil
}

func (s *LanguageService) FindBy(selector string, values ...string) ([]models.Language, error) {
	panic("not implemented")
}

func (s *LanguageService) FindByName(name string) (*models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "FindByName"),
		zap.String("name", name),
	)

	logger.Debug("Finding language")

	if result := s.cache.Find(name); result != nil {
		logger.Debug("Found language in cache", zap.Any("language", result))
		return result, nil
	}

	dbResults, err := s.repo.FindBy("name = ?", name)
	if err != nil {
		logger.Error("Failed to find language in db", zap.Error(err))
		return nil, err
	}

	dbResult := (*dbResults)[0]

	logger.Debug("Found language in db", zap.Any("language", dbResult))

	s.cache.Add(name, models.Language(dbResult))

	return s.cache.Find(name), nil
}

func (s *LanguageService) FindById(id uint64) (*models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "FindById"),
		zap.Uint64("id", id),
	)

	logger.Debug("Finding language")

	results := s.cache.FindAll()

	for _, result := range results {
		if result.Id == id {
			logger.Debug("Found language in cache", zap.Any("language", result))
			return &result, nil
		}
	}

	dbResult, err := s.repo.FindById(id)
	if err != nil {
		logger.Error("Failed to find language in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found language in db", zap.Any("language", dbResult))

	result := models.Language(*dbResult)

	return &result, nil
}

func (s *LanguageService) Add(language *models.Language) (*models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "Add"),
		zap.Any("language", language),
	)

	logger.Debug("Adding language")

	value, _ := s.FindByName(language.Name)
	if value != nil {
		err := constants.ErrAlreadyExists
		logger.Error("Failed to add language", zap.Error(err))
		return nil, err
	}

	dbValue := db_models.Language{Name: language.Name}
	dbResult, err := s.repo.Add(&dbValue)
	if err != nil {
		logger.Error("Failed to add language", zap.Error(err))
		return nil, err
	}

	logger.Debug("Added language", zap.Any("result", dbResult))

	result := models.Language(*dbResult)

	s.cache.Add(dbResult.Name, result)

	return &result, nil
}

func (s *LanguageService) Update(language *models.Language) (*models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "Update"),
		zap.Any("language", language),
	)

	logger.Debug("Updating language")

	dbLang := db_models.Language(*language)
	result, err := s.repo.Update(&dbLang)
	if err != nil {
		logger.Error("Failed to update language", zap.Error(err))
		return nil, err
	}

	logger.Debug("Updated language", zap.Any("result", result))

	s.cache.List.Store(language.Id, *language)
	return language, nil
}

func (s *LanguageService) Remove(id uint64) error {
	logger := s.logger.With(
		zap.String("function", "Remove"),
		zap.Uint64("id", id),
	)

	logger.Debug("Removing language")

	langToDelete, _ := s.FindById(id)
	if langToDelete == nil {
		err := constants.ErrNotFound
		logger.Error("Failed to remove language", zap.Error(err))
		return err
	}

	if err := s.repo.Remove(langToDelete.Id); err != nil {
		logger.Error("Failed to remove language", zap.Error(err))
		return err
	}

	logger.Debug("Removed language")

	s.cache.Remove(langToDelete.Name)

	return nil
}

func (s *LanguageService) FindAll() ([]models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "FindAll"),
	)

	logger.Debug("Finding languages")

	result := s.cache.FindAll()

	logger.Debug("Found languages in cache")

	return result, nil
}

func (s *LanguageService) findAllFromDb() ([]models.Language, error) {
	logger := s.logger.With(
		zap.String("function", "findAllFromDb"),
	)

	logger.Debug("Finding languages")

	dbResults, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Failed to find languages in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found languages in db", zap.Any("languages", dbResults))

	result := make([]models.Language, 0, len(*dbResults))

	for _, language := range *dbResults {
		result = append(result, models.Language(language))
	}

	return result, nil
}
