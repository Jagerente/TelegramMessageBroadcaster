package telegram

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/handlers"
	"DC_NewsSender/internal/telegram/middlewares"
	"time"

	tele "gopkg.in/telebot.v3"

	"go.uber.org/zap"
)

type Core struct {
	controller *controller.Controller
}

type BotConfig struct {
	Token  string
	Logger *zap.Logger
	Db     *repositories.Provider
	Debug  bool
}

func CreateBotCore(cfg *BotConfig) (*Core, error) {
	bot, err := tele.NewBot(tele.Settings{Token: cfg.Token, Poller: &tele.LongPoller{Timeout: 30 * time.Second}, Verbose: cfg.Debug})

	if err != nil {
		return nil, err
	}

	db := cfg.Db

	logger := cfg.Logger

	controller := &controller.Controller{Bot: bot, Provider: db, Logger: logger}

	return &Core{controller: controller}, nil
}

func (c *Core) Run() {
	c.controller.UpdateCache()

	c.handleUpdates()

	c.controller.Bot.Start()
}

func (bot *Core) handleUpdates() {
	bot.controller.Logger.Info("Listening for updates")

	cmdHandler := handlers.CreateCommandHandler(bot.controller)

	if err := cmdHandler.SetCommands(); err != nil {
		panic(err)
	}

	msgHandler := handlers.CreateMessageHandler(bot.controller)

	adminOnly := bot.controller.Bot.Group()

	adminOnly.Use(middlewares.Whitelist(bot.controller.CreateUserService()))

	adminOnly.Handle(tele.OnText, func(c tele.Context) error {
		user, err := bot.controller.CreateUserService().FindById(c.Sender().ID)
		if err != nil {
			bot.controller.Logger.Error(err.Error())

			return err
		}

		msgHandler.HandleMessage(user, c)
		cmdHandler.HandleCommand(user, c)
		return nil
	})

	adminOnly.Handle(tele.OnPhoto, func(c tele.Context) error {
		user, err := bot.controller.CreateUserService().FindById(c.Sender().ID)
		if err != nil {
			bot.controller.Logger.Error(err.Error())

			return err
		}

		msgHandler.HandleMessage(user, c)
		return nil
	})
}
