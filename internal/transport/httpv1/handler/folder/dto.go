package folder

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
		Name string `json:"name"`
	}
}
type DeleteFolderInput struct {
	ID uint `path:"id"`
}
type GetFolderInput struct {
	ID uint `path:"id"`
}
type GetAllFoldersInput struct{}
type UpdateFolderInput struct {
	ID   uint `path:"id"`
	Body struct {
		Name *string `json:"name,omitempty" minLength:"1" maxLength:"15" pattern:"^[a-zA-Zа-яА-Я0-9\\s]+$"`
	}
}
