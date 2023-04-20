package controller

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"

	"go.uber.org/zap"
)

type ChatService struct {
	logger *zap.Logger
	repo   repositories.IRepository[db_models.Chat, int64]
	cache  *models.List[int64, models.Chat]
}

func (s *ChatService) FindBy(selector string, values ...string) ([]models.Chat, error) {
	logger := s.logger.With(
		zap.String("function", "FindBy"),
		zap.String("selector", selector),
		zap.Strings("values", values),
	)

	logger.Debug("Started")

	dbResults, err := s.repo.FindBy(selector, values...)
	if err != nil {
		logger.Error("Failed to find chat in db", zap.Error(err))
		return nil, err
	}

	var result []models.Chat

	for _, dbRes := range *dbResults {
		result = append(result, models.Chat(dbRes))
	}

	logger.Debug("Found chat in db", zap.Any("chat", result))

	return result, nil
}

func (s *ChatService) FindByName(name string) (*models.Chat, error) {
	panic("not implemented")
}

func (s *ChatService) FindById(id int64) (*models.Chat, error) {
	logger := s.logger.With(
		zap.String("function", "FindById"),
		zap.Int64("id", id),
	)

	logger.Debug("Started")

	if result := s.cache.Find(id); result != nil {
		logger.Debug("Found chat in cache", zap.Any("chat", result))
		return result, nil
	}

	logger.Error("Failed to find chat in cache")

	dbResult, err := s.repo.FindById(id)
	if err != nil {
		logger.Error("Failed to find chat in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found chat in db", zap.Any("chat", dbResult))

	s.cache.Add(dbResult.Id, models.Chat(*dbResult))

	return s.cache.Find(id), nil
}

func (s *ChatService) Add(chat *models.Chat) (*models.Chat, error) {
	logger := s.logger.With(
		zap.String("function", "Add"),
		zap.Any("chat", chat),
	)

	logger.Debug("Started")

	value, _ := s.FindById(chat.Id)
	if value != nil {
		err := constants.ErrAlreadyExists
		logger.Error("Failed to add chat", zap.Error(err))
		return nil, err
	}

	dbChat := db_models.Chat(*chat)
	dbResult, err := s.repo.Add(&dbChat)
	if err != nil {
		logger.Error("Failed to add chat", zap.Error(err))
		return nil, err
	}

	logger.Debug("Added chat", zap.Any("result", dbResult))

	result := models.Chat(*dbResult)

	s.cache.Add(dbResult.Id, result)

	return &result, nil
}

func (s *ChatService) Update(chat *models.Chat) (*models.Chat, error) {
	logger := s.logger.With(
		zap.String("function", "Update"),
		zap.Any("chat", chat),
	)

	logger.Debug("Started")

	dbChat := db_models.Chat(*chat)
	result, err := s.repo.Update(&dbChat)
	if err != nil {
		logger.Error("Failed to update chat", zap.Error(err))
		return nil, err
	}

	logger.Debug("Updated chat", zap.Any("result", result))

	s.cache.List.Store(chat.Id, *chat)

	return chat, nil
}

func (s *ChatService) Remove(id int64) error {
	logger := s.logger.With(
		zap.String("function", "Remove"),
		zap.Int64("id", id),
	)

	logger.Debug("Started")

	chatToDelete, _ := s.FindById(id)
	if chatToDelete == nil {
		err := constants.ErrNotFound
		logger.Error("Failed to remove chat", zap.Error(err))
		return err
	}

	if err := s.repo.Remove(chatToDelete.Id); err != nil {
		logger.Error("Failed to remove chat", zap.Error(err))
		return err
	}

	logger.Debug("Removed chat")

	s.cache.Remove(chatToDelete.Id)

	return nil
}

func (s *ChatService) FindAll() ([]models.Chat, error) {
	logger := s.logger.With(
		zap.String("function", "FindAll"),
	)

	logger.Debug("Started")

	dbResults, err := s.repo.FindAll()
	if err != nil {
		logger.Error("Failed to find chats in db", zap.Error(err))
		return nil, err
	}

	logger.Debug("Found chats in db", zap.Any("chats", dbResults))

	result := make([]models.Chat, 0, len(*dbResults))

	for _, chat := range *dbResults {
		result = append(result, models.Chat(chat))
	}

	return result, nil
}
