package mappers

import (
	db_models "DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/telegram/models"
	"reflect"
	"unsafe"
)

func MapFromDb[T, D any](value *D) *T {
	v := reflect.ValueOf(value)

	switch v.Interface().(type) {
	case db_models.Chat:
		mappedValue := v.Convert(reflect.TypeOf(models.Chat{})).Interface().(models.Chat)
		return (*T)(unsafe.Pointer(&mappedValue))
	case db_models.Group:
		mappedValue := v.Convert(reflect.TypeOf(models.Group{})).Interface().(models.Group)
		return (*T)(unsafe.Pointer(&mappedValue))
	case db_models.Language:
		mappedValue := v.Convert(reflect.TypeOf(models.Language{})).Interface().(models.Language)
		return (*T)(unsafe.Pointer(&mappedValue))
	case db_models.Admin:
		mappedValue := models.User{
			Admin: v.Interface().(db_models.Admin),
		}
		return (*T)(unsafe.Pointer(&mappedValue))
	default:
		return nil
	}
}

func MapToDb[T, D any](value *T) *D {
	v := reflect.ValueOf(value).Elem()

	switch v.Interface().(type) {
	case models.Chat:
		mappedValue := v.Convert(reflect.TypeOf(db_models.Chat{})).Interface().(db_models.Chat)
		return (*D)(unsafe.Pointer(&mappedValue))
	case models.Group:
		mappedValue := v.Convert(reflect.TypeOf(db_models.Group{})).Interface().(db_models.Group)
		return (*D)(unsafe.Pointer(&mappedValue))
	case models.Language:
		mappedValue := v.Convert(reflect.TypeOf(db_models.Language{})).Interface().(db_models.Language)
		return (*D)(unsafe.Pointer(&mappedValue))
	case models.User:
		mappedValue := v.FieldByName("Admin").Interface().(db_models.Admin)
		return (*D)(unsafe.Pointer(&mappedValue))
	default:
		return nil
	}
}
