package middlewares

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"errors"
	"strconv"
	"strings"
)

type Middleware func(ctx context.Context) error

func IsMaster(ctx context.Context) error {
	if !ctx.Value(constants.CtxInitiator).(*models.User).IsMaster {
		return errors.New("not enough privileges")
	}
	return nil
}

func ParseChatAddInput(ctx context.Context) error {
	args := ctx.Value(constants.CtxArgs).(string)
	fields := strings.Split(args, ";")
	name := fields[0]
	id, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return constants.ErrInvalidInput
	}

	langId, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return constants.ErrInvalidInput
	}

	groupId, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return constants.ErrInvalidInput
	}

	ctx = context.WithValue(ctx, constants.CtxChatArgsId, name)
	ctx = context.WithValue(ctx, constants.CtxChatArgsName, id)
	ctx = context.WithValue(ctx, constants.CtxChatArgsLangId, langId)
	ctx = context.WithValue(ctx, constants.CtxChatArgsGroupId, groupId)

	return nil
}

func HasInput(ctx context.Context) error {
	args := ctx.Value(constants.CtxArgs).(string)
	argsSplited := strings.Split(args, ";")
	argsRequired := ctx.Value(constants.CtxArgsRequired).([]string)

	if len(argsSplited) < len(argsRequired) || args == "" {
		return errors.New("Input " + strings.Join(argsRequired, ";"))
	}

	return nil
}
