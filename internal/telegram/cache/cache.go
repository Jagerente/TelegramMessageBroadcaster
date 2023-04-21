package cache

import (
	"DC_NewsSender/internal/telegram/models"
	"errors"
	"sync"
)

type Cache[K, V any] struct {
	List sync.Map
}

func (list *Cache[K, V]) Find(key K) *V {
	if value, ok := list.List.Load(key); ok {
		v := value.(V)
		return &v
	}

	return nil
}

func (list *Cache[K, V]) FindAll() []V {
	var values []V
	list.List.Range(func(key, value interface{}) bool {
		values = append(values, value.(V))
		return true
	})
	return values
}

func (list *Cache[K, V]) Add(key K, value V) {
	list.List.Store(key, value)
}

func (list *Cache[K, V]) Remove(key K) error {
	if _, ok := list.List.Load(key); !ok {
		return errors.New("key not found")
	}

	list.List.Delete(key)
	return nil
}

func (list *Cache[K, V]) Clear() {
	list.List = sync.Map{}
}

var Users Cache[int64, models.User]

var Chats Cache[int64, models.Chat]

var Languages Cache[string, models.Language]

var Groups Cache[string, models.Group]

var Messages Cache[uint64, models.Message]

type ICache[K comparable, T any] interface {
	Add(key K, value T)
	Find(key K) *T
	FindAll() []T
	Remove(key K) error
	Clear()
}
