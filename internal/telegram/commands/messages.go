package commands

import (
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"

	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func listMessages(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	logger := controller.Logger.With(
		zap.String("function", "listMessages"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Listing all messages")

	var response strings.Builder
	response.WriteString("Messages List:")

	var keys []uint64

	cache.Messages.List.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(uint64))
		return true
	})

	for _, key := range keys {
		response.WriteString(fmt.Sprintf("\n [%d]", key))
	}

	logger.Debug("Listed all messages", zap.String("response", response.String()))

	return response.String(), nil
}

func testMessages(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id = ctx.Value(constants.MessageTestArgs.Names[0]).(uint64)

	logger := controller.Logger.With(
		zap.String("function", "testMessages"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Testing messages")

	msg := cache.Messages.Find(id)
	if msg == nil {
		logger.Warn("message not found")
		return "", constants.ErrNotFound
	}
	logger.Debug("Message found", zap.Any("message", msg))

	var response string
	counter := 0
	switch msg.GetType() {
	case models.PhotoMessage:
		for _, text := range msg.Text {
			msgToSend := tgbotapi.NewPhoto(user.Id, tgbotapi.FileID(msg.Photo))
			msgToSend.Caption = text
			if err := controller.Send(msgToSend); err != nil {
				logger.Error("Failed to send message", zap.Error(err))
				continue
			}
			logger.Debug("photo message sent")
			counter++
		}
	case models.TextMessage:
		for _, text := range msg.Text {
			msgToSend := tgbotapi.NewMessage(user.Id, text)
			if err := controller.Send(msgToSend); err != nil {
				logger.Error("Failed to send message", zap.Error(err))
				continue
			}
			logger.Debug("text message sent")
			counter++
		}
	default:
		logger.Warn("unknown message type")
		return "unknown message type", nil
	}

	response = fmt.Sprintf("%d messages sent.", counter)

	logger.Debug("Tested messages", zap.String("response", response))

	return response, nil
}

func sendMessages(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var msgId uint64 = ctx.Value(constants.MessageSendArgs.Names[0]).(uint64)
	var groupId uint64 = ctx.Value(constants.MessageSendArgs.Names[1]).(uint64)

	var groupService = controller.CreateGroupService()
	var chatService = controller.CreateChatService()

	logger := controller.Logger.With(
		zap.String("function", "sendMessages"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Sending message")

	group, err := groupService.FindById(groupId)
	if err != nil {
		return "", err
	}
	logger.Debug("Group found", zap.Any("group", group))

	msg := cache.Messages.Find(msgId)
	logger.Debug("Message found", zap.Any("msg", msg))
	chats, err := chatService.FindBy("group_id", fmt.Sprint(group.Id))
	if err != nil {
		return "", err
	}
	logger.Debug("Chats found", zap.Any("chats", chats))

	counter := 0
	failed := []string{}

	for _, chat := range chats {
		if chat.Group.Id == group.Id {
			switch msg.GetType() {
			case models.PhotoMessage:
				msgToSend := tgbotapi.NewPhoto(chat.Id, tgbotapi.FileID(msg.Photo))
				msgToSend.Caption = msg.Text[chat.LanguageId]
				if err := controller.Send(msgToSend); err != nil {
					logger.Error("Failed to send message", zap.Any("chat", chat), zap.Error(err))
					failed = append(failed, chat.Name)
					continue
				}
				counter++
				continue
			case models.TextMessage:
				msgToSend := tgbotapi.NewMessage(chat.Id, msg.Text[chat.LanguageId])
				if err := controller.Send(msgToSend); err != nil {
					logger.Error("Failed to send message", zap.Any("chat", chat), zap.Error(err))
					failed = append(failed, chat.Name)
					continue
				}
				counter++
				continue
			default:
				return "Error: unknown message type", nil
			}
		}
	}

	var response string
	if counter > 0 {
		response = fmt.Sprintf("%d messages sent to group [%d] %s.", counter, group.Id, group.Name)
	} else {
		response = fmt.Sprintf("No messages sent to group [%d] %s.", group.Id, group.Name)
	}
	if len(failed) > 0 {
		response += fmt.Sprintf("\nFailed to send to: %s", strings.Join(failed, ", "))
	}

	logger.Debug("Sent message")

	return response, nil
}
