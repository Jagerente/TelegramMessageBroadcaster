package commands

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/controller"
	"DC_NewsSender/internal/telegram/middlewares"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type command struct {
	Name        string
	Description string
	Inputs      []string
	Handler     func(ctx context.Context) (string, error)
	Middlewares []middlewares.Middleware
}

type commandGroup struct {
	Name   string
	Add    string
	Remove string
	List   string
}

func createCommandGroup(name string) *commandGroup {
	result := &commandGroup{
		Name:   name,
		Add:    constants.CmdAdd + name,
		Remove: constants.CmdRemove + name,
		List:   constants.CmdList + name,
	}

	return result
}

var (
	adminGroup    = createCommandGroup(constants.CmdAdmin)
	chatGroup     = createCommandGroup(constants.CmdChat)
	languageGroup = createCommandGroup(constants.CmdLanguage)
	groupGroup    = createCommandGroup(constants.CmdGroup)

	commands []command = []command{
		{
			Name:        adminGroup.Add,
			Description: fmt.Sprintf("Add %s", adminGroup.Name),
			Handler:     addAdmin,
			Inputs:      []string{"id", "name"},
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        adminGroup.Remove,
			Description: fmt.Sprintf("Remove %s", adminGroup.Name),
			Handler:     removeAdmin,
			Inputs:      []string{"id"},
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        adminGroup.List,
			Description: fmt.Sprintf("List %s", adminGroup.Name),
			Handler:     listAllAdmins,
			Inputs:      []string{},
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        chatGroup.Add,
			Description: fmt.Sprintf("Add %s", chatGroup.Name),
			Handler:     addChat,
			Inputs:      []string{"chat_id", "name", "language_id", "group_id"},
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        chatGroup.Remove,
			Description: fmt.Sprintf("Remove %s", chatGroup.Name),
			Handler:     removeChat,
			Inputs:      []string{"chat_id"},
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseChatAddInput},
		},
		{
			Name:        chatGroup.List,
			Description: fmt.Sprintf("List %s", chatGroup.Name),
			Handler:     listAllChats,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        languageGroup.Add,
			Description: fmt.Sprintf("Add %s", languageGroup.Name),
			Inputs:      []string{"name"},
			Handler:     addLanguage,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        languageGroup.Remove,
			Description: fmt.Sprintf("Remove %s", languageGroup.Name),
			Inputs:      []string{"id"},
			Handler:     removeLanguage,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        languageGroup.List,
			Description: fmt.Sprintf("List %s", languageGroup.Name),
			Inputs:      []string{},
			Handler:     listAllLanguages,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        groupGroup.Add,
			Description: fmt.Sprintf("Add %s", groupGroup.Name),
			Inputs:      []string{"name"},
			Handler:     addGroup,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        groupGroup.Remove,
			Description: fmt.Sprintf("Remove %s", groupGroup.Name),
			Inputs:      []string{"id"},
			Handler:     removeGroup,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput},
		},
		{
			Name:        groupGroup.List,
			Description: fmt.Sprintf("List %s", groupGroup.Name),
			Handler:     listAllGroups,
			Middlewares: []middlewares.Middleware{},
		},
	}
)

func (c *command) execute(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, constants.CtxArgsRequired, c.Inputs)

	for _, middleware := range c.Middlewares {
		if err := middleware(ctx); err != nil {
			return "", err
		}
	}

	result, err := c.Handler(ctx)
	if err != nil {
		return "", err
	}

	return result, nil
}

func addAdmin(ctx context.Context) (string, error) {
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	argsSlice := strings.Split(args, ";")
	if len(argsSlice) < 2 {
		return "", constants.ErrInvalidInput
	}

	id, err := strconv.ParseInt(argsSlice[0], 10, 64)
	if err != nil {
		return "", constants.ErrInvalidInput
	}
	name := argsSlice[1]

	if err := controller.AddUser(models.CreateUser(id, name)); err != nil {
		return "", err
	}

	return fmt.Sprintf("[%d] %s has been added!", id, name), nil
}

func removeAdmin(ctx context.Context) (string, error) {
	var user *models.User = ctx.Value(constants.CtxInitiator).(*models.User)
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	id, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		return "", constants.ErrInvalidInput
	}

	usrToDelete := controller.FindUser(id)
	if usrToDelete == nil {
		return "", errors.New("no such admin")
	}

	if usrToDelete.IsMaster {
		return "", errors.New("cannot remove Master")
	}

	if user.Id == usrToDelete.Id {
		return "", errors.New("cannot remove yourself")
	}

	if err = controller.RemoveUser(id); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d has been removed!", id), nil
}

func listAllAdmins(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var output strings.Builder
	output.WriteString("Admins List:")
	admins, err := controller.ListUsers()
	if err != nil {
		return "", err
	}

	for _, admin := range admins {
		output.WriteString(fmt.Sprintf("\n [%d] %s", admin.Id, admin.Name))
	}

	return output.String(), nil
}

// TODO: Chat Commands
func addChat(ctx context.Context) (string, error) {
	panic("not implemented")

	// var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	// var id string = ctx.Value(constants.CtxChatArgsId).(string)
	// var name string = ctx.Value(constants.CtxChatArgsName).(string)
	// var langId string = ctx.Value(constants.CtxChatArgsLangId).(string)
	// var groupId string = ctx.Value(constants.CtxChatArgsGroupId).(string)

	// if err := controller.AddChat(id, name, langId, groupId); err != nil {
	// 	return "", err
	// }

	// return fmt.Sprintf("%s has been added!", name), nil
}

func removeChat(ctx context.Context) (string, error) {
	panic("not implemented")
	// var args string = ctx.Value(constants.CtxArgs).(string)
	// var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	// argsSlice := strings.Split(args, ";")
	// if len(argsSlice) < 2 {
	// 	return "", constants.ErrInvalidInput
	// }

	// id, err := strconv.ParseInt(argsSlice[0], 10, 64)
	// if err != nil {
	// 	return "", constants.ErrInvalidInput
	// }
	// name := argsSlice[1]

	// name := args

	// if err := controller.RemoveGroup(name); err != nil {
	// 	return "", err
	// }

	// return fmt.Sprintf("%s has been removed!", name), nil
}

func listAllChats(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var output strings.Builder
	output.WriteString("Chat List:")
	chats, err := controller.ListChats()
	if err != nil {
		return "", err
	}

	for _, chat := range chats {
		output.WriteString(fmt.Sprintf("\n [%d] %s", chat.Id, chat.Name))
	}

	return output.String(), nil
}

func addLanguage(ctx context.Context) (string, error) {
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	name := args

	if err := controller.AddLanguage(name); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s has been added!", name), nil
}

func removeLanguage(ctx context.Context) (string, error) {
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	name := args

	if err := controller.RemoveLanguage(name); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s has been removed!", name), nil
}

func listAllLanguages(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var output strings.Builder
	output.WriteString("Language List:")
	languages, err := controller.ListLanguages()
	if err != nil {
		return "", err
	}

	for _, language := range languages {
		output.WriteString(fmt.Sprintf("\n [%d] %s", language.Id, language.Name))
	}

	return output.String(), nil
}

func addGroup(ctx context.Context) (string, error) {
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	name := args

	if err := controller.AddGroup(name); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s has been added!", name), nil
}

func removeGroup(ctx context.Context) (string, error) {
	var args string = ctx.Value(constants.CtxArgs).(string)
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)

	name := args

	if err := controller.RemoveGroup(name); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s has been removed!", name), nil
}

func listAllGroups(ctx context.Context) (string, error) {
	var controller = ctx.Value(constants.CtxController).(*controller.Controller)
	var output strings.Builder
	output.WriteString("Group List:")
	groups, err := controller.ListGroups()
	if err != nil {
		return "", err
	}

	for _, group := range groups {
		output.WriteString(fmt.Sprintf("\n [%d] %s", group.Id, group.Name))
	}

	return output.String(), nil
}
