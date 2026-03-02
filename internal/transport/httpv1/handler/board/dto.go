package board

type BoardResponse struct {
	Body BoardDTO
}

type BoardDTO struct {
	Name string `json:"name"`
}

type CreateBoardOutput struct {
	Status int `status:"201"`
	Body   struct {
		Message string `json:"message"`
	}
}

type GetBoardInput struct {
	ID uint `path:"id"`
}

type GetBoardOutput struct {
	Status int `status:"200"`
	Body   BoardResponse
}

type DeleteBoardInput struct {
	ID uint `path:"id"`
}

type DeleteBoardOutput struct {
	Status  int `status:"201"`
	Message string
}
