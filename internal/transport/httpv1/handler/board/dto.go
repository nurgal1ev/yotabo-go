package board

type BoardResponse struct {
	Body BoardDTO
}

type BoardDTO struct {
	Name string `json:"name"`
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
	Body   BoardResponse
}

type DeleteBoardInput struct {
	ID uint `path:"id"`
}

type DeleteBoardOutput struct {
	Status  int
	Message string
}
