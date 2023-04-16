package configuration

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

type Config[T any] struct {
	ENV T
}

func New[T any]() (*Config[T], error) {
	config := new(Config[T])

	if err := config.unmarshalConfig(); err != nil {
		val := reflect.ValueOf(&config.ENV).Elem()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			key := val.Type().Field(i).Tag.Get("mapstructure")
			envValue := os.Getenv(key)
			switch field.Kind() {
			case reflect.String:
				field.SetString(envValue)
			case reflect.Bool:
				field.SetBool(envValue == "true")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.ParseInt(envValue, 10, 64)
				if err != nil {
					return config, fmt.Errorf("failed to parse %s: %w", key, err)
				}
				field.SetInt(intValue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				uintValue, err := strconv.ParseUint(envValue, 10, 64)
				if err != nil {
					return config, fmt.Errorf("failed to parse %s: %w", key, err)
				}
				field.SetUint(uintValue)
			default:
				return config, fmt.Errorf("unsupported field type for %s", key)
			}
		}
	}

	return config, nil
}

func (config *Config[T]) unmarshalConfig() error {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config.ENV); err != nil {
		return err
	}

	return nil
}
