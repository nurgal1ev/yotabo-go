package board

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/boards",
		Tags:    []string{"Board"},
		Summary: "Create board",
	}, CreateBoardHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/boards/{id}",
		Tags:    []string{"Board"},
		Summary: "Delete board",
	}, DeleteBoardHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/boards/{id}",
		Tags:    []string{"Board"},
		Summary: "Get board",
	}, GetBoardHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/boards",
		Tags:    []string{"Board"},
		Summary: "Get all boards",
	}, GetAllBoardsHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPatch,
		Path:    "/boards/{id}",
		Tags:    []string{"Board"},
		Summary: "Update board",
	}, UpdateBoardHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/boards/{board_id}/invitations",
		Tags:    []string{"Invitations"},
		Summary: "Invite user to board",
	}, CreateInviteHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/invitations",
		Tags:    []string{"Invitations"},
		Summary: "Get my invitations",
	}, GetInviteHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/invitations/{id}/accept",
		Tags:    []string{"Invitations"},
		Summary: "Accept invitation",
	}, AcceptInviteHandler)
}
