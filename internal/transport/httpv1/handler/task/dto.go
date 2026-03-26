package task

type TaskDTO struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	Description *string       `json:"description" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]*$"`
	Status      string        `json:"status" enum:"backlog,in_progress,review,done"`
	Priority    string        `json:"priority" enum:"easy,medium,hard"`
	Subtasks    []SubtasksDTO `json:"subtasks"`
	Comments    []CommentsDTO `json:"comments"`
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
	}
}
type GetTaskInput struct {
	ID uint `path:"id"`
}

type GetAllTasksInput struct {
}

type UpdateTaskInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name        *string `json:"name,omitempty" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Description *string `json:"description,omitempty" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Status      *string `json:"status,omitempty" enum:"backlog,in_progress,review,done"`
		Priority    *string `json:"priority,omitempty" enum:"easy,medium,hard"`
	}
}
type DeleteTaskInput struct {
	ID uint `path:"id"`
}
