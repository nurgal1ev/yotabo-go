package subtask

type SubtaskResponse struct {
	Body SubtaskDTO
}

type SubtaskDTO struct {
	Name      string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	TaskID    uint   `path:"task_id"`
	Completed bool   `json:"completed"`
}

type CreateSubtaskOutput struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	}
}

type DeleteSubtaskInput struct {
	ID uint `path:"id"`
}

type DeleteSubtaskOutput struct {
	Status  int
	Message string
}

type UpdateSubtaskInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name      *string `json:"name,omitempty" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Completed *bool   `json:"completed,omitempty"`
	}
}

type UpdateSubtaskOutput struct {
	Status int
	Body   SubtaskResponse
}
