package user

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func Register(router *huma.Group) {
	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/users/me",
		Tags:    []string{"Users"},
		Summary: "Get user",
	}, GetUserHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPatch,
		Path:    "/users/me",
		Tags:    []string{"Users"},
		Summary: "Update user",
	}, UpdateUserHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/users/me/password",
		Tags:    []string{"Users"},
		Summary: "Change user password",
	}, ChangePasswordHandler)
}
