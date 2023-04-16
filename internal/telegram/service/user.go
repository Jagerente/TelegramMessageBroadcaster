package service

import (
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram/cache"

	"go.uber.org/zap"
)

type Service[T, D any, K comparable] struct {
	cache      cache.ICache[T, K]
	repository repositories.IRepository[D, K]
	logger     *zap.Logger
}

func (s *Service[T, D, K]) Find(id K) *T {
	panic("Not implemented")

	// s.logger.Debug("Looking for user",
	// 	zap.String("msg", "Looking in cache"))

	// if value := s.cache.Find(id); value != nil {
	// 	s.logger.Debug("Looking for user",
	// 		zap.String("msg", "Found in cache"))

	// 	return value
	// }

	// s.logger.Debug("Looking for user",
	// 	zap.String("msg", "Looking in db"))

	// value, err := s.repository.FindById(id)
	// if err != nil {
	// 	return nil
	// }

	// s.cache.Add(id, *mappers.MapFromDb[T, D](value))

	// s.logger.Debug("Looking for user",
	// 	zap.String("msg", "Found in db"))

	// return s.cache.Find(id)
}

func (s *Service[T, D, K]) Add(key K) error {
	panic("Not implemented")
}

func (s *Service[T, D, K]) Remove() {
	panic("Not implemented")
}
