package board

type BoardDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	FolderID *uint  `json:"folder_id"`
}
type CreateBoardInput struct {
	Body struct {
		Name     string `json:"name" minLength:"1" maxLength:"55" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
		FolderID *uint  `json:"folder_id,omitempty"`
	}
}
type GetBoardInput struct {
	ID uint `path:"id"`
}

type DeleteBoardInput struct {
	ID uint `path:"id"`
}
type UpdateBoardInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name     *string `json:"name"`
		FolderID *uint   `json:"folder_id"`
	}
}
