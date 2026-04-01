package board

type BoardDTO struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	FolderID *uint            `json:"folder_id"`
	Members  []BoardMemberDTO `json:"members,omitempty"`
}

type BoardMemberDTO struct {
	UserID    uint   `json:"userId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	JoinedAt  string `json:"joinedAt"`
}

type InvitationDTO struct {
	ID        uint   `json:"id"`
	BoardID   uint   `json:"boardId"`
	BoardName string `json:"boardName"`
	InviterID uint   `json:"inviterId"`
	Inviter   string `json:"inviter"` // username или first_name
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
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
type GetAllBoardsInput struct{}
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

type AddMemberInput struct {
	BoardID uint `path:"board_id"`
	Body    struct {
		Username string `json:"username" minLength:"1" maxLength:"12"`
	}
}

type GetMembersInput struct {
	BoardID uint `path:"board_id"`
}

type RemoveMemberInput struct {
	BoardID uint `path:"board_id"`
	UserID  uint `path:"user_id"`
}

type GetMyInvitationsInput struct {
	// пустой, userID из токена
}

type GetMyInvitationsOutput struct {
	Status int
	Body   struct {
		Invitations []InvitationDTO `json:"invitations"`
	}
}

type AcceptInvitationInput struct {
	ID uint `path:"id"`
}

type RejectInvitationInput struct {
	ID uint `path:"id"`
}
