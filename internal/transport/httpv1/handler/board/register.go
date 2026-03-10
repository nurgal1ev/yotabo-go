package board

import (
	"github.com/danielgtaylor/huma/v2"
	"net/http"
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
}
