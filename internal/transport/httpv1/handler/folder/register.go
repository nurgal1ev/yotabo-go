package folder

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/folders",
		Tags:    []string{"Folder"},
		Summary: "Create folder",
	}, CreateFolderHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/folders/{id}",
		Tags:    []string{"Folder"},
		Summary: "Delete folder",
	}, DeleteFolderHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/folders/{id}",
		Tags:    []string{"Folder"},
		Summary: "Get folder",
	}, GetFolderHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/folders",
		Tags:    []string{"Folder"},
		Summary: "Get all folders",
	}, GetAllFoldersHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPatch,
		Path:    "/folders/{id}",
		Tags:    []string{"Folder"},
		Summary: "Update folder",
	}, UpdateFolderHandler)
}
