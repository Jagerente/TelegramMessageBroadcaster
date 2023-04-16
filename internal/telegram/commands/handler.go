package commands

import (
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

func CreateHandler(controller *controller.Controller) (*CommandHandler, error) {
	if err := setCommands(controller); err != nil {
		return nil, err
	}
	return &CommandHandler{controller: controller}, nil
}

func setCommands(controller *controller.Controller) error {
	api := controller.Api
	var cmds = []tgbotapi.BotCommand{}
	for _, cmd := range commands {
		cmds = append(cmds, tgbotapi.BotCommand{Command: cmd.Name, Description: cmd.Description})
	}
	cfg := tgbotapi.NewSetMyCommands(cmds...)

	users, err := controller.ListUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		cfg.Scope = &tgbotapi.BotCommandScope{Type: "chat", ChatID: user.Id}
		api.Request(cfg)
	}

	return nil
}

func (h *CommandHandler) HandleCommand(user *models.User, update tgbotapi.Update) {
	var cmd = h.parseCommand(user, update.Message)

	if cmd.Name != "" {
		var commandToExecute *command
		for _, command := range commands {
			if command.Name == cmd.Name {
				commandToExecute = &command
				break
			}
		}

		if commandToExecute != nil {
			h.controller.Logger.Debug("Handling command",
				zap.String("msg", "Executing"),
				zap.Any("cmd", commandToExecute))

			ctx := context.WithValue(context.Background(), constants.CtxInitiator, user)
			ctx = context.WithValue(ctx, constants.CtxArgs, cmd.Arguments)
			ctx = context.WithValue(ctx, constants.CtxController, h.controller)

			h.controller.Logger.Sugar().Infow("HandleCommand",
				"user", ctx.Value(constants.CtxInitiator),
				"command", commandToExecute)

			result, err := commandToExecute.execute(ctx)
			if err != nil {
				if !strings.Contains(err.Error(), "Input") {
					h.controller.SendMessage(user.Id, cmdError(err.Error()))
					return
				}

				h.controller.SendMessage(user.Id, err.Error())
				return
			}

			h.controller.SendMessage(user.Id, result)
			h.controller.ChangeUserState(user, "")
			return
		}

		h.controller.SendMessage(user.Id, "Unknown command")
	}
}

func (h *CommandHandler) parseCommand(user *models.User, msg *tgbotapi.Message) *models.Command {
	h.controller.Logger.Debug("Parsing command",
		zap.String("msg", "Parsing"),
		zap.Any("user", user),
		zap.Any("message", msg.Text))

	var command = &models.Command{
		Name:      user.State,
		Arguments: msg.Text,
	}

	if command.Name == "" || msg.IsCommand() {
		command.Name = msg.Command()
		command.Arguments = msg.CommandArguments()
	}

	h.controller.ChangeUserState(user, command.Name)

	command.Arguments = strings.TrimSpace(command.Arguments)

	h.controller.Logger.Debug("Parsing command",
		zap.String("msg", "Parsed"),
		zap.Any("user", user),
		zap.Any("cmd", command))

	return command
}

func cmdError(err string, fields ...any) string {
	return fmt.Sprintf("Error: %s", fmt.Sprintf(err, fields...))
}
