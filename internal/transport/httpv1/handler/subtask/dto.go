package subtask

type SubtaskDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	TaskID    uint   `path:"task_id"`
	Completed bool   `json:"completed"`
}
type CreateSubtaskInput struct {
	TaskID uint `path:"task_id"`
	Body   struct {
		Name      string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		Completed bool   `json:"completed"`
	}
}
type DeleteSubtaskInput struct {
	ID uint `path:"id"`
}
type UpdateSubtaskInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name      *string `json:"name" minLength:"1" maxLength:"55"`
		Completed *bool   `json:"completed,omitempty"`
	}
}
