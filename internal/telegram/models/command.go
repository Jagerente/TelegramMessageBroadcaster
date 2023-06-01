package models

import "reflect"

type Command struct {
	Name      string
	Arguments []string
}

type Arguments struct {
	Names []string
	Types []reflect.Kind
}
