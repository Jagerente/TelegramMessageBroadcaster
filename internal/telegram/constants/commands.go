package constants

import (
	"DC_NewsSender/internal/telegram/models"
	"reflect"
)

const (
	CmdAdd    string = "add"
	CmdRemove string = "remove"
	CmdList   string = "list"

	CmdStart            string = "start"
	CmdConfigureMessage string = "configuremessage"
	CmdAdmin            string = "admin"
	CmdChat             string = "chat"
	CmdLanguage         string = "language"
	CmdGroup            string = "group"
)

var (
	UserAddArgs models.Arguments = models.Arguments{
		Names: []string{"user_id", "user_name"},
		Types: []reflect.Kind{reflect.Int64, reflect.String},
	}
	UserRemoveArgs models.Arguments = models.Arguments{
		Names: []string{"user_id"},
		Types: []reflect.Kind{reflect.Int64},
	}
	UserListArgs models.Arguments = models.Arguments{
		Names: []string{},
		Types: []reflect.Kind{},
	}

	ChatAddArgs models.Arguments = models.Arguments{
		Names: []string{"chat_id", "chat_name", "language_id", "group_id"},
		Types: []reflect.Kind{reflect.Int64, reflect.String, reflect.Uint64, reflect.Uint64},
	}
	ChatRemoveArgs models.Arguments = models.Arguments{
		Names: []string{"chat_id"},
		Types: []reflect.Kind{reflect.Int64},
	}
	ChatListArgs models.Arguments = models.Arguments{
		Names: []string{},
		Types: []reflect.Kind{},
	}

	LanguageAddArgs models.Arguments = models.Arguments{
		Names: []string{"language_name"},
		Types: []reflect.Kind{reflect.String},
	}
	LanguageRemoveArgs models.Arguments = models.Arguments{
		Names: []string{"language_id"},
		Types: []reflect.Kind{reflect.Uint64},
	}
	LanguageListArgs models.Arguments = models.Arguments{
		Names: []string{},
		Types: []reflect.Kind{},
	}

	GroupAddArgs models.Arguments = models.Arguments{
		Names: []string{"group_name"},
		Types: []reflect.Kind{reflect.String},
	}
	GroupRemoveArgs models.Arguments = models.Arguments{
		Names: []string{"group_id"},
		Types: []reflect.Kind{reflect.Uint64},
	}
	GroupListArgs models.Arguments = models.Arguments{
		Names: []string{},
		Types: []reflect.Kind{},
	}

	MessageTestArgs models.Arguments = models.Arguments{
		Names: []string{"message_id"},
		Types: []reflect.Kind{reflect.Uint64},
	}
	MessageListArgs models.Arguments = models.Arguments{
		Names: []string{},
		Types: []reflect.Kind{},
	}
	MessageSendArgs models.Arguments = models.Arguments{
		Names: []string{"message_id", "group_id"},
		Types: []reflect.Kind{reflect.Uint64, reflect.Uint64},
	}
)
