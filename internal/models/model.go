package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  string
	Avatar    *string

	Tasks []Task `gorm:"foreignKey:CreatedByID"`
}

type Task struct {
	gorm.Model
	Name        string
	Description *string
	Status      string
	Priority    string

	CreatedByID uint
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`

	UpdatedByID uint
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID"`

	BoardID uint
	Board   Board

	SubtaskID *uint
	Subtasks  []Subtask

	CommentID *uint
	Comments  []Comment
}

type Board struct {
	gorm.Model
	Name        string
	Description string

	CreatedByID uint
	CreatedBy   User

	FolderID *uint
	Folder   Folder `gorm:"foreignKey:FolderID"`
}

type Folder struct {
	gorm.Model
	Name string

	CreatedByID uint
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`

	Boards []Board
}

type Subtask struct {
	gorm.Model
	Name      string
	Completed bool

	TaskID uint
	Task   Task `gorm:"foreignKey:TaskID"`
}

type Comment struct {
	gorm.Model

	TaskID uint
	Task   Task `gorm:"foreignKey:TaskID"`

	AuthorID uint
	Author   User `gorm:"foreignKey:AuthorID"`

	Message string
}

type BoardMember struct {
	gorm.Model
	BoardID uint
	UserID  uint

	Board Board `gorm:"foreignKey:BoardID"`
	User  User  `gorm:"foreignKey:UserID"`
}
