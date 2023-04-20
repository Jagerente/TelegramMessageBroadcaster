package commands

import (
	"DC_NewsSender/internal/telegram/commands/middlewares"
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"

	"context"
	"fmt"
)

type Command struct {
	Name        string
	Description string
	Arguments   models.Arguments
	Handler     func(ctx context.Context) (string, error)
	Middlewares []middlewares.Middleware
}

var (
	Commands []Command = []Command{
		{
			Name:        AdminGroup.Add,
			Description: fmt.Sprintf("Add %s", AdminGroup.Name),
			Handler:     addAdmin,
			Arguments:   constants.UserAddArgs,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        AdminGroup.Remove,
			Description: fmt.Sprintf("Remove %s", AdminGroup.Name),
			Handler:     removeAdmin,
			Arguments:   constants.UserRemoveArgs,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        AdminGroup.List,
			Description: fmt.Sprintf("List %s", AdminGroup.Name),
			Handler:     listAllAdmins,
			Arguments:   constants.UserListArgs,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        ChatGroup.Add,
			Description: fmt.Sprintf("Add %s", ChatGroup.Name),
			Handler:     addChat,
			Arguments:   constants.ChatAddArgs,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        ChatGroup.Remove,
			Description: fmt.Sprintf("Remove %s", ChatGroup.Name),
			Handler:     removeChat,
			Arguments:   constants.ChatRemoveArgs,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        ChatGroup.List,
			Description: fmt.Sprintf("List %s", ChatGroup.Name),
			Handler:     listAllChats,
			Arguments:   constants.ChatListArgs,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        LanguageGroup.Add,
			Description: fmt.Sprintf("Add %s", LanguageGroup.Name),
			Arguments:   constants.LanguageAddArgs,
			Handler:     addLanguage,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        LanguageGroup.Remove,
			Description: fmt.Sprintf("Remove %s", LanguageGroup.Name),
			Arguments:   constants.LanguageRemoveArgs,
			Handler:     removeLanguage,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        LanguageGroup.List,
			Description: fmt.Sprintf("List %s", LanguageGroup.Name),
			Arguments:   constants.LanguageListArgs,
			Handler:     listAllLanguages,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        GroupGroup.Add,
			Description: fmt.Sprintf("Add %s", GroupGroup.Name),
			Arguments:   constants.GroupAddArgs,
			Handler:     addGroup,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        GroupGroup.Remove,
			Description: fmt.Sprintf("Remove %s", GroupGroup.Name),
			Arguments:   constants.GroupRemoveArgs,
			Handler:     removeGroup,
			Middlewares: []middlewares.Middleware{middlewares.IsMaster, middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        GroupGroup.List,
			Description: fmt.Sprintf("List %s", GroupGroup.Name),
			Arguments:   constants.GroupListArgs,
			Handler:     listAllGroups,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        "listmessages",
			Description: fmt.Sprintf("List messages queue"),
			Handler:     listMessages,
			Arguments:   constants.MessageListArgs,
			Middlewares: []middlewares.Middleware{},
		},
		{
			Name:        "testmessage",
			Description: fmt.Sprintf("Preview messages"),
			Arguments:   constants.MessageTestArgs,
			Handler:     testMessages,
			Middlewares: []middlewares.Middleware{middlewares.HasInput, middlewares.ParseInput},
		},
		{
			Name:        "sendmessages",
			Description: fmt.Sprintf("Send messages to a specific group"),
			Arguments:   constants.MessageSendArgs,
			Handler:     sendMessages,
			Middlewares: []middlewares.Middleware{middlewares.HasInput, middlewares.ParseInput},
		},
	}
)

func (c *Command) Execute(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, constants.CtxArgsRequired, c.Arguments)

	for _, middleware := range c.Middlewares {
		if err := middleware(&ctx); err != nil {
			return "", err
		}
	}

	result, err := c.Handler(ctx)
	if err != nil {
		return "", err
	}

	return result, nil
}

var (
	AdminGroup    = createCommandGroup(constants.CmdAdmin)
	ChatGroup     = createCommandGroup(constants.CmdChat)
	LanguageGroup = createCommandGroup(constants.CmdLanguage)
	GroupGroup    = createCommandGroup(constants.CmdGroup)
)

type CommandGroup struct {
	Name   string
	Add    string
	Remove string
	List   string
}

func createCommandGroup(name string) *CommandGroup {
	result := &CommandGroup{
		Name:   name,
		Add:    constants.CmdAdd + name,
		Remove: constants.CmdRemove + name,
		List:   constants.CmdList + name,
	}

	return result
}
