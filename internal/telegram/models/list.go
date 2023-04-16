package models

import (
	"errors"
	"sync"
)

type List[K, V any] struct {
	List sync.Map
}

func (list *List[K, V]) Find(key K) *V {
	if value, ok := list.List.Load(key); ok {
		v := value.(V)
		return &v
	}

	return nil
}

func (list *List[K, V]) Add(key K, value V) {
	list.List.Store(key, value)
}

func (list *List[K, V]) Remove(key K) error {
	if _, ok := list.List.Load(key); !ok {
		return errors.New("key not found")
	}

	list.List.Delete(key)
	return nil
}
