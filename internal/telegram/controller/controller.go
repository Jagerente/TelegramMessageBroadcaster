package controller

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Controller struct {
	Api      *tgbotapi.BotAPI
	Provider *repositories.Provider
	Logger   *zap.Logger
}

func (c *Controller) ConfigureAndSendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := c.Api.Send(msg)
	return err
}

func (c *Controller) SendMessage(msg tgbotapi.MessageConfig) error {
	_, err := c.Api.Send(msg)
	return err
}

func (c *Controller) Send(msg tgbotapi.Chattable) error {
	_, err := c.Api.Send(msg)
	return err
}

func (c *Controller) CreateLanguageService() IService[models.Language, uint64] {
	s := &LanguageService{
		logger: c.Logger.With(zap.String("service", "LanguageService")),
		repo:   c.Provider.CreateLanguageRepo(),
		cache:  &cache.Languages,
	}

	return s
}
func (c *Controller) CreateGroupService() IService[models.Group, uint64] {
	s := &GroupService{
		logger: c.Logger.With(zap.String("service", "GroupService")),
		repo:   c.Provider.CreateGroupRepo(),
		cache:  &cache.Groups,
	}

	return s
}

func (c *Controller) CreateUserService() IService[models.User, int64] {
	s := &UserService{
		logger: c.Logger.With(zap.String("service", "UserService")),
		repo:   c.Provider.CreateAdminsRepo(),
		cache:  &cache.Users,
	}

	return s
}
func (c *Controller) CreateChatService() IService[models.Chat, int64] {
	s := &ChatService{
		logger: c.Logger.With(zap.String("service", "ChatService")),
		repo:   c.Provider.CreateChatRepo(),
		cache:  &cache.Chats,
	}

	return s
}

type IService[T any, ID comparable] interface {
	FindByName(string) (*T, error)
	FindById(ID) (*T, error)
	FindBy(string, ...string) ([]T, error)
	FindAll() ([]T, error)
	Add(*T) (*T, error)
	Update(*T) (*T, error)
	Remove(ID) error
}
