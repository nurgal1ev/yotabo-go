package folder

type FolderResponse struct {
	Body FolderDTO
}

type FolderDTO struct {
	ID     uint             `json:"id"`
	Name   string           `json:"name"`
	Boards []FolderBoardDTO `json:"boards"`
}

type FolderBoardDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateFolderInput struct {
	Body struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
}
type CreateFolderOutput struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	}
}

type DeleteFolderInput struct {
	ID uint `path:"id"`
}

type DeleteFolderOutput struct {
	Status  int
	Message string
}

type GetFolderInput struct {
	ID uint `path:"id"`
}

type GetFolderOutput struct {
	Status int
	Body   FolderDTO
}

type GetAllFoldersInput struct{}

type GetAllFoldersOutput struct {
	Status int
	Body   struct {
		Folders []FolderDTO `json:"folders"`
	}
}

type UpdateFolderInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name *string `json:"name,omitempty" minLength:"1" maxLength:"15" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	}
}
type UpdateFolderOutput struct {
	Status int
	Body   FolderResponse
}
