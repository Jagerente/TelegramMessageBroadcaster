package telegram

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type BotCore struct {
	controller *controller.Controller
}

type BotConfig struct {
	Token  string
	Logger *zap.Logger
	Db     *repositories.Provider
	Debug  bool
}

func CreateBotCore(cfg *BotConfig) (*BotCore, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	api.Debug = cfg.Debug

	db := cfg.Db

	logger := cfg.Logger

	controller := &controller.Controller{Api: api, Provider: db, Logger: logger}

	return &BotCore{controller: controller}, nil
}

func (bot *BotCore) Run() {
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.controller.Api.GetUpdatesChan(updateConfig)

	bot.controller.UpdateCache()

	bot.handleUpdates(updates)
}

func (bot *BotCore) handleUpdates(updates tgbotapi.UpdatesChannel) {
	bot.controller.Logger.Info("Listening for updates")

	cmdHandler := handlers.CreateCommandHandler(bot.controller)

	if err := cmdHandler.SetCommands(); err != nil {
		panic(err)
	}

	msgHandler := handlers.CreateMessageHandler(bot.controller)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if user, _ := bot.controller.CreateUserService().FindById(update.Message.Chat.ID); user != nil {
			msgHandler.HandleMessage(user, update)
			cmdHandler.HandleCommand(user, update)
		}
	}
}
