package httpv1

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/task"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/user"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"net/http"
)

func StartServer() {
	r := chi.NewMux()

	cors.New(cors.Options{
		AllowedOrigins:   config.Load().App.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	humaCfg := huma.DefaultConfig("Yotabo API", "1.0.0")
	humaCfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"jwt": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humachi.New(r, humaCfg)

	huma.Post(api, "/api/v1/auth/register", user.RegisterHandler)
	huma.Post(api, "/api/v1/auth/login", user.LoginHandler)

	router := huma.NewGroup(api, "/api/v1")
	router.UseMiddleware(middleware.HumaJWTMiddleware)
	router.UseSimpleModifier(func(op *huma.Operation) {
		op.Security = []map[string][]string{
			{"jwt": {}},
		}
	})

	huma.Register(router, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/tasks",
		Tags:    []string{"Task"},
		Summary: "Create task",
	}, task.CreateTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Get task",
	}, task.GetTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Update task",
	}, task.UpdateTaskHandler)

	huma.Register(router, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/tasks/{id}",
		Tags:    []string{"Task"},
		Summary: "Delete task",
	}, task.DeleteTaskHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
