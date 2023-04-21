package handlers

import (
	"DC_NewsSender/internal/telegram/cache"
	"DC_NewsSender/internal/telegram/commands"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type MessageHandler struct {
	controller *controller.Controller
}

func CreateMessageHandler(controller *controller.Controller) *MessageHandler {
	return &MessageHandler{controller: controller}
}

func (h *MessageHandler) HandleMessage(user *models.User, update tgbotapi.Update) {
	logger := h.controller.Logger.With(
		zap.String("function", "HandleMessage"),
		zap.Any("user", user),
		zap.Any("message", update.Message),
	)

	logger.Debug("Handling message")

	msgId, langId, text, photo := h.parseMessage(update.Message)

	if text == "" {
		return
	}

	h.controller.ClearUserState(user)

	messageId, err := strconv.ParseUint(msgId, 10, 64)
	if err != nil {
		h.controller.ConfigureAndSendMessage(user.Id, cmdError("invalid message id.\nMust be positive number."))
		return
	}

	languageId, err := strconv.ParseUint(langId, 10, 64)
	if err != nil {
		h.controller.ConfigureAndSendMessage(user.Id, cmdError("invalid language id.\nUse /%s to verify.", commands.LanguageGroup.List))
		return
	}

	lang, err := h.controller.CreateLanguageService().FindById(languageId)
	if err != nil {
		h.controller.ConfigureAndSendMessage(user.Id, cmdError("language [%d] not found.\nUse /%s to verify.", languageId, commands.LanguageGroup.List))
		return
	}

	msg := cache.Messages.Find(messageId)

	if msg == nil {
		msg = models.CreateMessage()
	}

	msg.Id = messageId
	msg.Text[languageId] = text
	if photo != "" {
		msg.Photo = photo
	}

	cache.Messages.Add(msg.Id, *msg)

	h.controller.ConfigureAndSendMessage(user.Id, fmt.Sprintf("Message stashed\nMessage ID: %d\nLanguage: [%d] %s", msg.Id, lang.Id, lang.Name))
}

func (h *MessageHandler) parseMessage(msg *tgbotapi.Message) (string, string, string, string) {
	logger := h.controller.Logger.With(
		zap.String("function", "parseMessage"),
		zap.Any("message", msg),
	)

	logger.Debug("Parsing message")

	if msg == nil || (msg.Text == "" && msg.Photo == nil) {
		return "", "", "", ""
	}

	var text string
	var photo string

	if msg.Photo != nil {
		text = msg.Caption
		photo = msg.Photo[len(msg.Photo)-1].FileID
	} else {
		text = msg.Text
	}

	const pattern = `\${(.+?)}`
	const submatchIndex = 1

	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(text)
	if match == nil {
		return "", "", "", ""
	}

	params := strings.Split(match[submatchIndex], ";")
	if len(params) < 2 {
		return "", "", "", ""
	}

	msgId := params[0]
	langId := params[1]

	if match[submatchIndex] != "" {
		text = strings.Replace(text, match[0], "", 1)
	}

	logger.Debug("Parsed message", zap.Any("value", []string{msgId, langId, text, photo}))

	return msgId, langId, text, photo
}
