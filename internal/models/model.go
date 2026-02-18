package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string //TODO: макс длина 50, а-Я, a-Z, только буквы
	LastName  string //TODO: макс длина 50, а-Я, a-Z, только буквы
	Username  string `gorm:"unique"` //TODO: символы a-z 0-9 _. длина мин 3 - 12 макс., не должен начинаться с числа или спец знаков
	Email     string `gorm:"unique"` //TODO: валидация
	Password  string //TODO: длина: мин 7, макс 12, разрешены все символы
	Avatar    *string

	Tasks []Task
}

type Task struct {
	gorm.Model
	Name        string //TODO: длина макс 55 символов, а-Я, a-Z, 1-10, спец символы
	Description string //TODO: длина макс 10000 символов, а-Я, a-Z, 1-10, спец символы
	Status      string //TODO: валидация статуса backlog | in_progress | review | done
	Priority    string //TODO: валидация статуса easy | medium | hard

	CreatedBy User
	UpdatedBy *User
}
