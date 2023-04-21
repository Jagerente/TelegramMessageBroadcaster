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

func addGroup(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var name string = ctx.Value(constants.GroupAddArgs.Names[0]).(string)

	var groupService = controller.CreateGroupService()

	logger := controller.Logger.With(
		zap.String("function", "addGroup"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Adding group")

	if _, err := groupService.Add(&models.Group{Name: name}); err != nil {
		logger.Error("Failed to add a group", zap.Error(err))
		return "", err
	}

	result := fmt.Sprintf("Group %s has been added!", name)

	logger.Debug("Added group", zap.String("result", result))

	return result, nil
}

func removeGroup(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var id uint64 = ctx.Value(constants.GroupRemoveArgs.Names[0]).(uint64)

	var groupService = controller.CreateGroupService()

	logger := controller.Logger.With(
		zap.String("function", "listAllGroups"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Removing group")

	if err := groupService.Remove(id); err != nil {
		logger.Error("Failed to remove a group", zap.Error(err))
		return "", err
	}

	response := fmt.Sprintf("Group %d has been removed!", id)

	logger.Debug("Removed group", zap.String("response", response))

	return response, nil
}

func listAllGroups(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var user = ctx.Value(constants.CtxUser).(*models.User)

	var groupService = controller.CreateGroupService()

	logger := controller.Logger.With(
		zap.String("function", "listAllGroups"),
		zap.Int64("userID", user.Id),
	)

	logger.Debug("Listing all groups")

	var response strings.Builder
	response.WriteString("Group List:")
	groups, err := groupService.FindAll()
	if err != nil {
		logger.Error("Failed to find all groups", zap.Error(err))
		return "", err
	}

	for _, group := range groups {
		response.WriteString(fmt.Sprintf("\n [%d] %s", group.Id, group.Name))
	}

	logger.Debug("Listed all groups", zap.String("response", response.String()))

	return response.String(), nil
}
