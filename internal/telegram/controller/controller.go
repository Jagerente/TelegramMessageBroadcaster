package controller

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/models"

	tele "gopkg.in/telebot.v3"

	"go.uber.org/zap"
)

type Controller struct {
	Bot      *tele.Bot
	Provider *repositories.Provider
	Logger   *zap.Logger
}

func (c *Controller) UpdateCache() error {
	if err := c.CreateUserService().UpdateCache(); err != nil {
		return err
	}

	if err := c.CreateChatService().UpdateCache(); err != nil {
		return err
	}

	if err := c.CreateLanguageService().UpdateCache(); err != nil {
		return err
	}

	if err := c.CreateGroupService().UpdateCache(); err != nil {
		return err
	}

	return nil
}

func (c *Controller) ClearUserState(user *models.User) {
	user.State = ""
	c.CreateUserService().Update(user)
}

func (c *Controller) SetUserState(user *models.User, state string) {
	user.State = state
	c.CreateUserService().Update(user)
}

func (c *Controller) SendText(chatId int64, text string) error {
	_, err := c.Bot.Send(&tele.User{ID: chatId}, text, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	return err
}

func (c *Controller) SendPhotoByID(chatId int64, photoId string, caption string) error {
	msg := &tele.Photo{File: tele.File{FileID: photoId}}
	msg.Caption = caption
	_, err := c.Bot.Send(&tele.User{ID: chatId}, msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
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
	ClearCache()
	UpdateCache() error
	FindByName(string) (*T, error)
	FindById(ID) (*T, error)
	FindBy(string, ...string) ([]T, error)
	FindAll() ([]T, error)
	Add(*T) (*T, error)
	Update(*T) (*T, error)
	Remove(ID) error
}
