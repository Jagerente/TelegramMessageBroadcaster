package middlewares

import (
	"DC_NewsSender/internal/telegram/constants"
	"DC_NewsSender/internal/telegram/models"
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Middleware func(ctx *context.Context) error

func IsMaster(ctx *context.Context) error {
	var user = (*ctx).Value(constants.CtxInitiator).(*models.User)

	if !user.IsMaster {
		return errors.New("not enough privileges")
	}

	return nil
}

func ParseInput(ctx *context.Context) error {
	var args = strings.Split((*ctx).Value(constants.CtxArgs).(string), ";")
	var argsRequired = (*ctx).Value(constants.CtxArgsRequired).(models.Arguments)

	var argNames = argsRequired.Names
	var argTypes = argsRequired.Types

	if len(args) != len(argNames) || len(args) != len(argTypes) {
		return constants.ErrInvalidInput
	}

	for i, arg := range args {
		switch argTypes[i] {
		case reflect.Int, reflect.Int64:
			val, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return constants.ErrInvalidInput
			}
			*ctx = context.WithValue(*ctx, argNames[i], val)
		case reflect.Uint, reflect.Uint64:
			val, err := strconv.ParseUint(arg, 10, 64)
			if err != nil {
				return constants.ErrInvalidInput
			}
			*ctx = context.WithValue(*ctx, argNames[i], val)
		case reflect.String:
			*ctx = context.WithValue(*ctx, argNames[i], arg)
		default:
			return constants.ErrInvalidInput
		}
	}

	return nil
}

func HasInput(ctx *context.Context) error {
	var args = (*ctx).Value(constants.CtxArgs).(string)
	var argsSplited = strings.Split(args, ";")
	var argsRequired = (*ctx).Value(constants.CtxArgsRequired).(models.Arguments)

	if len(argsSplited) < len(argsRequired.Names) || args == "" {
		return errors.New("Input " + strings.Join(argsRequired.Names, ";"))
	}

	return nil
}
