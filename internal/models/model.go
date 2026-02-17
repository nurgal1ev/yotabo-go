package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Avatar   string

	Tasks []Task
}

type Task struct {
	gorm.Model
	Name        string
	Description string
	Status      string

	CreatedBy User
	UpdatedBy *User
}
