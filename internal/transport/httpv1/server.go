package httpv1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/board"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/comment"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/folder"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/subtask"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/task"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/user"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
)

func StartServer() {
	r := chi.NewMux()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.Load().App.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

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

	task.Register(router)
	board.Register(router)
	folder.Register(router)
	subtask.Register(router)
	comment.Register(router)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
