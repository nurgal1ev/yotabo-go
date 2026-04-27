package task

import "github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"

type TaskDTO struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	Description *string       `json:"description" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]*$"`
	Status      string        `json:"status" enum:"backlog,in_progress,review,done"`
	Priority    string        `json:"priority" enum:"easy,medium,hard"`
	DueDate     *string       `json:"dueDate,omitempty" format:"date"`
	Subtasks    []SubtasksDTO `json:"subtasks"`
	Comments    []CommentsDTO `json:"comments"`
	BoardID     uint
}
type SubtasksDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type CommentsDTO struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
type CreateTaskInput struct {
	Body struct {
		Name        string  `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Description *string `json:"description" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]*$"`
		Status      string  `json:"status" enum:"backlog,in_progress,review,done"`
		Priority    string  `json:"priority" enum:"easy,medium,hard"`
		DueDate     *string `json:"dueDate,omitempty" format:"date"`
		BoardID     uint    `json:"board_id"`
	}
}
type GetTaskInput struct {
	ID uint `path:"id"`
}

type GetAllTasksInput struct {
	BoardID common.OptionalParam[uint] `query:"board_id"`
}

type UpdateTaskInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name        *string `json:"name,omitempty" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Description *string `json:"description,omitempty" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Status      *string `json:"status,omitempty" enum:"backlog,in_progress,review,done"`
		Priority    *string `json:"priority,omitempty" enum:"easy,medium,hard"`
		DueDate     *string `json:"dueDate,omitempty"`
	}
}
type DeleteTaskInput struct {
	ID uint `path:"id"`
}
