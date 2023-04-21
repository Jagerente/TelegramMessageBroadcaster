package commands

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func addAdmin(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id int64 = ctx.Value(constants.UserAddArgs.Names[0]).(int64)
	var name string = ctx.Value(constants.UserAddArgs.Names[1]).(string)

	var userService = controller.CreateUserService()

	logger := controller.Logger.With(
		zap.String("function", "addAdmin"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Adding user")

	if _, err := userService.Add(models.CreateUser(id, name)); err != nil {
		logger.Error("Failed to add user", zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("User [%d] %s has been added!", id, name)
	logger.Debug("User added", zap.String("result", result))

	return result, nil
}

func removeAdmin(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user *models.User = ctx.Value(constants.CtxInitiator).(*models.User)

	var id int64 = ctx.Value(constants.UserRemoveArgs.Names[0]).(int64)

	var userService = controller.CreateUserService()

	logger := controller.Logger.With(
		zap.String("function", "removeAdmin"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Removing user")

	usrToDelete, err := userService.FindById(id)
	if err != nil {
		err := constants.ErrNotFound
		logger.Error("Failed to remove a user", zap.Error(err))
		return "", err
	}

	if usrToDelete.IsMaster {
		err := errors.New("cannot remove Master")
		logger.Error("Failed to remove a user", zap.Error(err))
		return "", err
	}

	if user.Id == usrToDelete.Id {
		err := errors.New("cannot remove yourself")
		logger.Error("Failed to remove a user", zap.Error(err))
		return "", err
	}

	if err = userService.Remove(id); err != nil {
		logger.Error("Failed to remove a user", zap.Error(err))
		return "", err
	}

	logger.Debug("Removed user")

	response := fmt.Sprintf("%d has been removed!", id)

	return response, nil
}

func listAllAdmins(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var userService = controller.CreateUserService()

	logger := controller.Logger.With(
		zap.String("function", "listAllAdmins"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Listing all users")

	var response strings.Builder
	response.WriteString("Admins List:")
	admins, err := userService.FindAll()
	if err != nil {
		logger.Error("Failed to find all users", zap.Error(err))
		return "", err
	}

	for _, admin := range admins {
		response.WriteString(fmt.Sprintf("\n [%d] %s", admin.Id, admin.Name))
	}

	logger.Debug("Listed all users", zap.String("response", response.String()))

	return response.String(), nil
}
