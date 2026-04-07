package board

import "github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"

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
type GetAllBoardsInput struct {
	FolderID common.OptionalParam[uint] `query:"folder_id"`
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

// invite
type InvitationDTO struct {
	ID        uint   `json:"id"`
	BoardID   uint   `json:"boardId"`
	BoardName string `json:"boardName"`
	InviterID uint   `json:"inviterId"`
	Inviter   string `json:"inviter"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type CreateInvitationInput struct {
	BoardID uint `path:"board_id"`
	Body    struct {
		Username string `json:"username"`
	}
}

type GetInvitationInput struct{}

type AcceptInvitationInput struct {
	ID uint `path:"id"`
}
type RejectInvitationInput struct{}
