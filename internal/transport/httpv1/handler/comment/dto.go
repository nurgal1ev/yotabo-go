package comment

type CommentDTO struct {
	ID       uint `json:"id"`
	Message  string
	TaskID   uint `path:"task_id"`
	AuthorID uint `json:"author_id"`
}

type CreateCommentInput struct {
	TaskID uint `path:"task_id"`
	Body   struct {
		Message string `json:"message"`
	}
}

type UpdateCommentInput struct {
	ID   uint `path:"id"`
	Body struct {
		Message *string `json:"message"`
	}
}

type DeleteCommentInput struct {
	ID uint `path:"id"`
}

type GetCommentInput struct {
	ID uint `path:"id"`
}

type GetCommentsByTaskInput struct {
	TaskID uint `path:"task_id"`
}
