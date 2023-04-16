package configuration

import (
	"go.uber.org/zap"
)

func GetLogger(debug bool) *zap.Logger {
	logger, _ := zap.NewProduction()

	if debug {
		logger, _ = zap.NewDevelopment()
	}

	return logger
}
