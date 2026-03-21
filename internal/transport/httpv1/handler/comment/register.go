package comment

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/tasks/{task_id}/comments",
		Tags:    []string{"Comments"},
		Summary: "Create comment",
	}, CreateCommentHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPatch,
		Path:    "/comments/{id}",
		Tags:    []string{"Comments"},
		Summary: "Update comment",
	}, UpdateCommentHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/comments/{id}",
		Tags:    []string{"Comments"},
		Summary: "Delete comment",
	}, DeleteCommentHandler)
}
