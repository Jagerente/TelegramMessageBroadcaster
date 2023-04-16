package telegram

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/commands"
	"DC_NewsSender/internal/telegram/controller"

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

	bot.handleUpdates(updates)
}

func (bot *BotCore) handleUpdates(updates tgbotapi.UpdatesChannel) {
	bot.controller.Logger.Info("Listening for updates")
	cmdHandler, err := commands.CreateHandler(bot.controller)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		bot.controller.Logger.Debug("Message Received",
			zap.String("user", update.Message.From.UserName),
			zap.String("message", update.Message.Text))

		if user := bot.controller.FindUser(update.Message.Chat.ID); user != nil {
			cmdHandler.HandleCommand(user, update)
		}
	}
}
