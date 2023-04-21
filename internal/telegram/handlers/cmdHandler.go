package handlers

import (
	"DC_NewsSender/internal/telegram/commands"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type CommandHandler struct {
	controller *controller.Controller
}

func CreateCommandHandler(controller *controller.Controller) *CommandHandler {
	return &CommandHandler{controller: controller}
}

func (h *CommandHandler) SetCommands() error {
	api := h.controller.Api
	var cmds = []tgbotapi.BotCommand{}
	for _, cmd := range commands.Commands {
		cmds = append(cmds, tgbotapi.BotCommand{Command: cmd.Name, Description: cmd.Description})
	}

	logger := h.controller.Logger.With(
		zap.String("function", "SetCommands"),
		zap.Any("commands", cmds),
	)

	logger.Debug("Setting commands list")

	cfg := tgbotapi.NewSetMyCommands(cmds...)

	users, err := h.controller.CreateUserService().FindAll()
	if err != nil {
		return err
	}

	for _, user := range users {
		cfg.Scope = &tgbotapi.BotCommandScope{Type: "chat", ChatID: user.Id}
		api.Request(cfg)
	}

	logger.Debug("Set commands list", zap.Any("users", users))

	return nil
}

func (h *CommandHandler) HandleCommand(user *models.User, update tgbotapi.Update) {
	logger := h.controller.Logger.With(
		zap.String("function", "HandleCommand"),
		zap.Any("user", user),
		zap.Any("message", update.Message.Text),
	)

	logger.Debug("Handling command")

	var cmd = h.parseCommand(user, update.Message)

	if cmd.Name != "" {
		var commandToExecute *commands.Command
		for _, command := range commands.Commands {
			if command.Name == cmd.Name {
				commandToExecute = &command
				break
			}
		}

		if commandToExecute != nil {
			logger.Debug("Executing command", zap.Any("cmd", commandToExecute.Name))

			ctx := context.WithValue(context.Background(), constants.CtxInitiator, user)
			ctx = context.WithValue(ctx, constants.CtxArgs, cmd.Arguments)
			ctx = context.WithValue(ctx, constants.CtxController, h.controller)
			ctx = context.WithValue(ctx, constants.CtxUser, user)

			result, err := commandToExecute.Execute(ctx)
			switch err {
			case constants.ErrEmptyInput:
				h.controller.ConfigureAndSendMessage(user.Id, fmt.Sprintf("Input %s", strings.Join(commandToExecute.Arguments.Names, ";")))
				return
			case nil:
				h.controller.ConfigureAndSendMessage(user.Id, result)
				h.controller.ClearUserState(user)
				return
			default:
				h.controller.ConfigureAndSendMessage(user.Id, cmdError(err.Error()))
				return
			}
		}

		h.controller.ConfigureAndSendMessage(user.Id, "Unknown command")
	}
}

func (h *CommandHandler) parseCommand(user *models.User, msg *tgbotapi.Message) *models.Command {
	logger := h.controller.Logger.With(
		zap.String("function", "parseCommand"),
		zap.Any("user", user),
		zap.Any("message", msg.Text),
	)

	logger.Debug("Parsing command")

	var command = &models.Command{
		Name:      user.State,
		Arguments: msg.Text,
	}

	if command.Name == "" || msg.IsCommand() {
		command.Name = msg.Command()
		command.Arguments = msg.CommandArguments()
	}

	h.controller.SetUserState(user, command.Name)

	command.Arguments = strings.TrimSpace(command.Arguments)

	logger.Debug("Parsed", zap.Any("command", command))

	return command
}

func cmdError(err string, fields ...any) string {
	return fmt.Sprintf("Error: %s", fmt.Sprintf(err, fields...))
}
