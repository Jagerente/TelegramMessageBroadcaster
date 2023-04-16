package controller

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/models"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Controller struct {
	Api      *tgbotapi.BotAPI
	Provider *repositories.Provider
	Logger   *zap.Logger
}

func (c *Controller) SendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := c.Api.Send(msg)
	return err
}

func (c *Controller) FindUser(id int64) *models.User {
	c.Logger.Debug("FindUser",
		zap.String("msg", "Looking in cache"))

	if user := cache.Users.Find(id); user != nil {
		return user
	}

	c.Logger.Debug("FindUser",
		zap.String("msg", "Looking in db"))
	usr, err := c.Provider.CreateAdminsRepo().FindById(id)
	if err != nil {
		return nil
	}

	cache.Users.Add(usr.Id, models.User{Admin: *usr})
	c.Logger.Debug("FindUser",
		zap.String("msg", "Found in db"))

	return cache.Users.Find(id)
}

func (c *Controller) AddUser(user *models.User) error {
	c.Logger.Debug("AddUser",
		zap.String("msg", "Adding user"),
		zap.Int64("id", user.Id))

	if c.FindUser(user.Id) != nil {
		return errors.New("user already exists")
	}

	result, err := c.Provider.CreateAdminsRepo().Add(&user.Admin)
	if err != nil {
		return err
	}

	cache.Users.Add(result.Id, models.User{Admin: *result})

	return nil
}

func (c *Controller) ChangeUserState(user *models.User, state string) error {
	c.Logger.Debug("ChangeUserState",
		zap.Any("user", user),
		zap.String("state", state))

	user.State = state
	cache.Users.List.Store(user.Id, *user)
	return nil
}

func (c *Controller) RemoveUser(id int64) error {
	c.Logger.Debug("RemoveUser",
		zap.String("msg", "Removing user"),
		zap.Int64("id", id))

	usrToDelete := c.FindUser(id)
	if usrToDelete == nil {
		return errors.New("no such user")
	}

	if err := c.Provider.CreateAdminsRepo().Remove(usrToDelete.Id); err != nil {
		return err
	}

	if err := cache.Users.Remove(usrToDelete.Id); err != nil {
		return err
	}

	return nil
}

func (c *Controller) ListUsers() ([]models.User, error) {
	c.Logger.Debug("ListUsers",
		zap.String("msg", "Looking for admins"))

	admins, err := c.Provider.CreateAdminsRepo().FindAll()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug("ListUsers",
		zap.String("msg", "Found in db"),
		zap.Any("admins", admins))

	result := make([]models.User, 0, len(*admins))

	for _, admin := range *admins {
		result = append(result, models.User{Admin: admin})
	}

	c.Logger.Debug("ListUsers",
		zap.String("msg", "Built a slice"),
		zap.Any("slice", result))

	return result, nil
}

func (c *Controller) FindLanguage(name string) *models.Language {
	c.Logger.Debug("FindLanguage",
		zap.String("msg", "Looking for language"),
		zap.String("name", name))

	c.Logger.Debug("FindLanguage",
		zap.String("msg", "Looking in cache"))

	if lang := cache.Languages.Find(name); lang != nil {

		c.Logger.Debug("Looking for language",
			zap.String("msg", "Found in cache"))

		return lang
	}

	c.Logger.Debug("FindLanguage",
		zap.String("msg", "Looking in db"))

	lang, err := c.Provider.CreateLanguageRepo().FindBy("name", name)
	if err != nil {
		return nil
	}

	cache.Languages.Add(lang.Name, models.Language(*lang))

	c.Logger.Debug("FindLanguage",
		zap.String("msg", "Found in db"))

	return cache.Languages.Find(lang.Name)
}

func (c *Controller) AddLanguage(name string) error {
	c.Logger.Debug("AddLanguage",
		zap.String("msg", "Adding language"),
		zap.String("name", name))

	if c.FindLanguage(name) != nil {
		return errors.New("language already exists")
	}

	result, err := c.Provider.CreateLanguageRepo().Add(&db_models.Language{Name: name})
	if err != nil {
		return err
	}

	cache.Languages.Add(result.Name, models.Language(*result))

	return nil
}

func (c *Controller) RemoveLanguage(name string) error {
	c.Logger.Debug("RemoveLanguage",
		zap.String("msg", "Removing language"),
		zap.String("name", name))

	langToDelete := c.FindLanguage(name)
	if langToDelete == nil {
		return errors.New("no such language")
	}

	if err := c.Provider.CreateLanguageRepo().Remove(langToDelete.Id); err != nil {
		return err
	}

	if err := cache.Languages.Remove(langToDelete.Name); err != nil {
		return err
	}

	return nil
}

func (c *Controller) ListLanguages() ([]models.Language, error) {
	c.Logger.Debug("ListLanguages",
		zap.String("msg", "Looking for languages"))

	languages, err := c.Provider.CreateLanguageRepo().FindAll()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug("ListLanguages",
		zap.String("msg", "Found in db"),
		zap.Any("languages", languages))

	result := make([]models.Language, 0, len(*languages))

	for _, language := range *languages {
		result = append(result, models.Language(language))
	}

	c.Logger.Debug("ListLanguages",
		zap.String("msg", "Built a slice"),
		zap.Any("languages", result))

	return result, nil
}

func (c *Controller) FindGroup(name string) *models.Group {
	c.Logger.Debug("FindGroup",
		zap.String("msg", "Looking for group in cache"),
		zap.String("name", name))

	if group := cache.Groups.Find(name); group != nil {
		c.Logger.Debug("FindGroup",
			zap.String("msg", "Found in cache"))
		return group
	}

	c.Logger.Debug("FindGroup",
		zap.String("msg", "Looking for group in db"),
		zap.String("name", name))
	group, err := c.Provider.CreateGroupRepo().FindBy("name", name)
	if err != nil {
		return nil
	}

	cache.Groups.Add(group.Name, models.Group(*group))

	c.Logger.Debug("FindGroup",
		zap.String("msg", "Found in db"))
	return cache.Groups.Find(group.Name)
}

func (c *Controller) AddGroup(name string) error {
	c.Logger.Debug("AddGroup",
		zap.String("msg", "Adding group"),
		zap.String("name", name))

	if c.FindGroup(name) != nil {
		return errors.New("group already exists")
	}

	result, err := c.Provider.CreateGroupRepo().Add(&db_models.Group{Name: name})
	if err != nil {
		return err
	}

	cache.Groups.Add(result.Name, models.Group(*result))

	return nil
}

func (c *Controller) RemoveGroup(name string) error {
	c.Logger.Debug("RemoveGroup",
		zap.String("msg", "Removing group"),
		zap.String("name", name))

	groupToDelete := c.FindGroup(name)
	if groupToDelete == nil {
		return errors.New("no such group")
	}

	if err := c.Provider.CreateGroupRepo().Remove(groupToDelete.Id); err != nil {
		return err
	}

	cache.Groups.Remove(groupToDelete.Name)

	return nil
}

func (c *Controller) ListGroups() ([]models.Group, error) {
	c.Logger.Debug("ListGroups",
		zap.String("msg", "Looking for groups"))

	groups, err := c.Provider.CreateGroupRepo().FindAll()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug("ListGroups",
		zap.String("msg", "Found in db"),
		zap.Any("groups", groups))

	result := make([]models.Group, 0, len(*groups))

	for _, group := range *groups {
		result = append(result, models.Group(group))
	}

	c.Logger.Debug("ListGroups",
		zap.String("msg", "Built a slice"),
		zap.Any("groups", result))

	return result, nil
}

func (c *Controller) FindChat(id int64) *models.Chat {

	c.Logger.Debug("FindChat",
		zap.String("msg", "Looking in cache"))

	if chat := cache.Chats.Find(id); chat != nil {
		c.Logger.Debug("FindChat",
			zap.String("msg", "Found in cache"))
		return chat
	}

	c.Logger.Debug("FindChat",
		zap.String("msg", "Looking in db"))
	chat, err := c.Provider.CreateChatRepo().FindById(id)
	if err != nil {
		return nil
	}

	cache.Chats.Add(chat.Id, models.Chat(*chat))

	c.Logger.Debug("FindChat",
		zap.String("msg", "Found in db"))
	return cache.Chats.Find(id)
}

func (c *Controller) AddChat(upd tgbotapi.Update) error {
	repo := c.Provider.CreateChatRepo()
	chat := upd.FromChat().ChatConfig()

	c.Logger.Debug("AddChat",
		zap.String("msg", "Registering chat"),
		zap.Any("chat", chat))

	if c.FindChat(chat.ChatID) == nil {
		_, err := repo.Add(&db_models.Chat{Id: chat.ChatID, Name: chat.SuperGroupUsername})
		if err != nil {
			return err
		}
	}

	c.Logger.Debug("AddChat",
		zap.String("msg", "Registered chat"),
		zap.Any("chat", chat))

	return nil
}

func (c *Controller) ListChats() ([]models.Chat, error) {
	c.Logger.Debug("ListChats",
		zap.String("msg", "Looking for chats"))

	chats, err := c.Provider.CreateChatRepo().FindAll()
	if err != nil {
		return nil, err
	}

	c.Logger.Debug("ListChats",
		zap.String("msg", "Found in db"),
		zap.Any("chats", chats))

	result := make([]models.Chat, 0, len(*chats))

	for _, chat := range *chats {
		result = append(result, models.Chat(chat))
	}

	c.Logger.Debug("ListChats",
		zap.String("msg", "Built a slice"),
		zap.Any("chats", result))

	return result, nil
}
