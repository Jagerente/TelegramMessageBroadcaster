package commands

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func addChat(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id int64 = ctx.Value(constants.ChatAddArgs.Names[0]).(int64)
	var name string = ctx.Value(constants.ChatAddArgs.Names[1]).(string)
	var langId uint64 = ctx.Value(constants.ChatAddArgs.Names[2]).(uint64)
	var groupId uint64 = ctx.Value(constants.ChatAddArgs.Names[3]).(uint64)

	var langService = controller.CreateLanguageService()
	var groupService = controller.CreateGroupService()
	var chatService = controller.CreateChatService()

	logger := controller.Logger.With(
		zap.String("function", "addChat"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Adding chat")

	if lang, err := langService.FindById(langId); lang == nil {
		logger.Error("Failed to find a language", zap.Error(err))
		return "", errors.New("language not found")
	}

	if group, err := groupService.FindById(groupId); group == nil {
		logger.Error("Failed to find a group", zap.Error(err))
		return "", errors.New("group not found")
	}

	if _, err := chatService.Add(&models.Chat{Id: id, Name: name, LanguageId: langId, GroupId: groupId}); err != nil {
		logger.Error("Failed to add a chat", zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("Chat %s has been added!", name)
	logger.Debug("Added chat", zap.String("result", result))

	return result, nil
}

func removeChat(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id int64 = ctx.Value(constants.ChatRemoveArgs.Names[0]).(int64)

	var chatService = controller.CreateChatService()

	logger := controller.Logger.With(
		zap.String("function", "removeChat"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Removing chat")

	if err := chatService.Remove(id); err != nil {
		logger.Error("Failed to remove a group", zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("Chat %d has been removed!", id)

	logger.Debug("Removed chat")

	return result, nil
}

func listAllChats(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var chatService = controller.CreateChatService()

	logger := controller.Logger.With(
		zap.String("function", "listAllChats"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Listing all chats")

	var response strings.Builder
	response.WriteString("Chat List:")
	chats, err := chatService.FindAll()
	if err != nil {
		logger.Error("Failed to find all chats", zap.Error(err))
		return "", err
	}

	for _, chat := range chats {
		response.WriteString(fmt.Sprintf("\n [%d] %s", chat.Id, chat.Name))
	}

	logger.Debug("Listed all chats", zap.String("response", response.String()))

	return response.String(), nil
}
