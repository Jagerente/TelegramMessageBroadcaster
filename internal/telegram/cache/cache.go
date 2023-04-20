package cache

import "DC_NewsSender/internal/telegram/models"

var Users models.List[int64, models.User]

var Chats models.List[int64, models.Chat]

var Languages models.List[string, models.Language]

var Groups models.List[string, models.Group]

var Messages models.List[uint64, models.Message]

type ICache[T any, K comparable] interface {
	Add(key K, value T)
	Find(key K) *T
	Remove(key K) error
}
