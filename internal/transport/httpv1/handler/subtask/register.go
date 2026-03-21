package subtask

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/tasks/{task_id}/subtasks",
		Tags:    []string{"Subtask"},
		Summary: "Create subtask",
	}, CreateSubtaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPatch,
		Path:    "/subtasks/{id}",
		Tags:    []string{"Subtask"},
		Summary: "Update subtask",
	}, UpdateSubtaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/subtasks/{id}",
		Tags:    []string{"Subtask"},
		Summary: "Delete subtask",
	}, DeleteSubtaskHandler)
}
