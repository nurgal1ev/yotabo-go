package task

type TaskResponse struct {
	Body TaskDTO
}

type TaskDTO struct {
	Name        string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	Description string `json:"description" maxLength:"10000" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	Status      string `json:"status" enum:"backlog,in_progress,review,done"`
	Priority    string `json:"priority" enum:"easy,medium,hard"`
}
type CreateTaskOutput struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	}
}

type GetTaskInput struct {
	ID uint `path:"id"` // {id} в URL
}
type GetTaskOutput struct {
	Status int `status:"201"`
	Body   TaskResponse
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

type UpdateTaskOutput struct {
	Status int `status:"201"`
	Body   TaskResponse
}

type DeleteTaskInput struct {
	ID uint `path:"id"`
}

type DeleteTaskOutput struct {
	Status  int `status:"201"`
	Message string
}
