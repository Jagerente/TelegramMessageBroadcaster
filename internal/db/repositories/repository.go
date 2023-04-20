package repositories

import (
	"errors"
)

type Repository[T, K any] struct {
	BaseRepository
}

func (repo *Repository[T, K]) FindById(id K) (*T, error) {
	var selectedValue *T

	var connection = repo.gormConnection

	if result := connection.First(&selectedValue, id); result.Error != nil {
		return nil, result.Error
	}

	return selectedValue, nil
}

func (repo *Repository[T, K]) FindBy(selector string, values ...string) (*[]T, error) {
	var selectedValue *[]T

	var connection = repo.gormConnection

	if result := connection.Where(selector, values).Find(&selectedValue); result.Error != nil {
		return nil, result.Error
	}

	if len(*selectedValue) < 1 {
		return nil, errors.New("nothing found")
	}

	return selectedValue, nil
}

func (repo *Repository[T, K]) FindAll() (*[]T, error) {
	var values []T = make([]T, 0)

	var connection = repo.gormConnection

	if result := connection.Find(&values); result.Error != nil {
		return nil, result.Error
	}

	return &values, nil
}

func (repo *Repository[T, K]) Add(value *T) (*T, error) {
	if value == nil {
		return nil, errors.New("null value provided")
	}

	var connection = repo.gormConnection

	if err := connection.Create(value).Error; err != nil {
		return nil, err
	}

	return value, nil
}

func (repo *Repository[T, K]) Update(value *T) (*T, error) {
	if value == nil {
		return nil, errors.New("null value provided")
	}

	var connection = repo.gormConnection

	if err := connection.Save(&value).Error; err != nil {
		return nil, err
	}

	return value, nil
}

func (repo *Repository[T, K]) Remove(id K) error {
	var connection = repo.gormConnection
	var value = new(T)

	if err := connection.Delete(value, id).Error; err != nil {
		return err
	}

	return nil
}
