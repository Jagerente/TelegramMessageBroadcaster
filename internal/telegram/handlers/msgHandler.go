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

	tele "gopkg.in/telebot.v3"

	"go.uber.org/zap"
)

type MessageHandler struct {
	controller *controller.Controller
}

func CreateMessageHandler(controller *controller.Controller) *MessageHandler {
	return &MessageHandler{controller: controller}
}

func (h *MessageHandler) HandleMessage(user *models.User, tctx tele.Context) {
	logger := h.controller.Logger.With(
		zap.String("function", "HandleMessage"),
		zap.Any("user", user.Id),
		zap.Any("message", tctx.Message().ID),
	)

	logger.Debug("Handling message")

	msgId, langId, text, photo := h.parseMessage(*tctx.Message())

	if text == "" {
		return
	}

	h.controller.ClearUserState(user)

	messageId, err := strconv.ParseUint(msgId, 10, 64)
	if err != nil {
		tctx.Send(cmdError("invalid message id.\nMust be positive number."))
		return
	}

	languageId, err := strconv.ParseUint(langId, 10, 64)
	if err != nil {
		tctx.Send(cmdError("invalid language id.\nUse /%s to verify.", commands.LanguageGroup.List))
		return
	}

	lang, err := h.controller.CreateLanguageService().FindById(languageId)
	if err != nil {
		tctx.Send(cmdError("language [%d] not found.\nUse /%s to verify.", languageId, commands.LanguageGroup.List))
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

	tctx.Send(fmt.Sprintf("Message stashed\nMessage ID: %d\nLanguage: [%d] %s", msg.Id, lang.Id, lang.Name))
}

func (h *MessageHandler) parseMessage(msg tele.Message) (string, string, string, string) {
	logger := h.controller.Logger.With(
		zap.String("function", "parseMessage"),
		zap.Any("message", msg.ID),
	)

	logger.Debug("Parsing message")

	if msg.Text == "" && msg.Photo == nil {
		return "", "", "", ""
	}

	var text string
	var photo string

	if msg.Photo != nil {
		text = msg.Caption
		photo = msg.Photo.FileID
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
