package board

type BoardResponse struct {
	Body BoardDTO
}

type BoardDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	FolderID *uint  `json:"folder_id"`
}

type CreateBoardOutput struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	}
}

type GetBoardInput struct {
	ID uint `path:"id"`
}

type GetBoardOutput struct {
	Status int
	Body   BoardDTO
}

type GetAllBoardsInput struct{}

type GetAllBoardsOutput struct {
	Status int
	Body   struct {
		Boards []BoardDTO `json:"boards"`
	}
}
type DeleteBoardInput struct {
	ID uint `path:"id"`
}

type DeleteBoardOutput struct {
	Status  int
	Message string
}

type UpdateBoardInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name     *string `json:"name"`
		FolderID *uint   `json:"folder_id"`
	}
}

type UpdateBoardOutput struct {
	Status int
	Body   BoardDTO
}
