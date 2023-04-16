package models

import (
	db_models "DC_NewsSender/internal/db/models"
)

type User struct {
	db_models.Admin
	State string
}

func CreateUser(id int64, name string) *User {
	return &User{Admin: db_models.Admin{Id: id, Name: name}}
}
