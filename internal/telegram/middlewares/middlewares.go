package middlewares

import (
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func Whitelist(s controller.IService[models.User, int64]) tele.MiddlewareFunc {
	admins, _ := s.FindAll()

	var adminsIds []int64

	for _, admin := range admins {
		adminsIds = append(adminsIds, admin.Id)
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return middleware.Restrict(middleware.RestrictConfig{
			Chats: adminsIds,
			In:    next,
			Out:   func(c tele.Context) error { return nil },
		})(next)
	}
}
