package commands

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func addLanguage(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var name string = ctx.Value(constants.LanguageAddArgs.Names[0]).(string)

	var langService = controller.CreateLanguageService()

	logger := controller.Logger.With(
		zap.String("function", "addGroup"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Started")

	if _, err := langService.Add(&models.Language{Name: name}); err != nil {
		logger.Error("Failed to add a language", zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("Language %s has been added!", name)
	logger.Debug("Finished", zap.String("result", result))

	return result, nil
}

func removeLanguage(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id uint64 = ctx.Value(constants.LanguageRemoveArgs.Names[0]).(uint64)

	var langService = controller.CreateLanguageService()

	logger := controller.Logger.With(
		zap.String("function", "listAllLanguages"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Started")

	if err := langService.Remove(id); err != nil {
		logger.Error("Failed to remove language", zap.Error(err))
		return "", err
	}

	response := fmt.Sprintf("Language %d has been removed!", id)
	logger.Debug("Finished", zap.String("response", response))

	return response, nil
}

func listAllLanguages(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var langService = controller.CreateLanguageService()

	logger := controller.Logger.With(
		zap.String("function", "listAllLanguages"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Started")

	var response strings.Builder
	response.WriteString("Language List:")

	languages, err := langService.FindAll()
	if err != nil {
		logger.Error("Failed to find all languages", zap.Error(err))
		return "", err
	}

	for _, language := range languages {
		response.WriteString(fmt.Sprintf("\n [%d] %s", language.Id, language.Name))
	}

	logger.Debug("Finished", zap.String("response", response.String()))

	return response.String(), nil
}
