package task

import (
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/tasks",
		Tags:    []string{"Task"},
		Summary: "Create task",
	}, CreateTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Get task",
	}, GetTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/tasks",
		Tags:    []string{"Task"},
		Summary: "Get all tasks",
	}, GetAllTasksHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Update task",
	}, UpdateTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Delete task",
	}, DeleteTaskHandler)
}
