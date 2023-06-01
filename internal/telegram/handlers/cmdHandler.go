package handlers

import (
	"DC_NewsSender/internal/telegram/commands"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type CommandHandler struct {
	controller *controller.Controller
}

func CreateCommandHandler(controller *controller.Controller) *CommandHandler {
	return &CommandHandler{controller: controller}
}

func (h *CommandHandler) SetCommands() error {
	bot := h.controller.Bot
	var cmds = []tele.Command{}
	for _, cmd := range commands.Commands {
		cmds = append(cmds, tele.Command{Text: cmd.Name, Description: cmd.Description})
	}

	logger := h.controller.Logger.With(
		zap.String("function", "SetCommands"),
	)

	logger.Debug("Setting commands list", zap.Any("commands", cmds))

	users, err := h.controller.CreateUserService().FindAll()
	if err != nil {
		return err
	}

	for _, user := range users {
		commands := tele.CommandParams{Commands: cmds, Scope: &tele.CommandScope{Type: tele.CommandScopeChat, ChatID: user.Id}}
		logger.Debug("Commands", zap.Any("obj", commands))
		if err := bot.SetCommands(commands); err != nil {
			return err
		}

	}

	logger.Debug("Set commands list", zap.Any("users", users))

	return nil
}

func (h *CommandHandler) HandleCommand(user *models.User, tctx tele.Context) {
	logger := h.controller.Logger.With(
		zap.String("function", "HandleCommand"),
		zap.Any("user", user.Id),
		zap.Any("message", tctx.Message().Text),
	)

	logger.Debug("Handling command")

	var cmd = h.parseCommand(user, tctx)

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
				tctx.Send(fmt.Sprintf("Input %s", strings.Join(commandToExecute.Arguments.Names, " ")))
				return
			case nil:
				tctx.Send(result)
				h.controller.ClearUserState(user)
				return
			default:
				tctx.Send(cmdError(err.Error()))
				return
			}
		}

		tctx.Send("Unknown command")
		h.controller.ClearUserState(user)
	}
}

func (h *CommandHandler) parseCommand(user *models.User, ctx tele.Context) *models.Command {
	logger := h.controller.Logger.With(
		zap.String("function", "parseCommand"),
		zap.Any("user", user),
		zap.Any("message", ctx.Message().Text),
	)

	logger.Debug("Parsing command")

	var cmd = &models.Command{
		Name:      user.State,
		Arguments: strings.Split(ctx.Message().Text, " "),
	}

	if cmd.Name == "" && strings.Contains(ctx.Message().Text, "/") {
		cmd.Name = strings.ReplaceAll(strings.Split(ctx.Message().Text, " ")[0], "/", "")
		cmd.Arguments = ctx.Args()
	}

	h.controller.SetUserState(user, cmd.Name)

	logger.Debug("Parsed", zap.Any("command", cmd))

	return cmd
}

func cmdError(err string, fields ...any) string {
	return fmt.Sprintf("Error: %s", fmt.Sprintf(err, fields...))
}
