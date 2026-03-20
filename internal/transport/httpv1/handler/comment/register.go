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

}
