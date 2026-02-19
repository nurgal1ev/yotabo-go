package httpv1

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/user"
	"net/http"
)

func StartServer() {
	r := chi.NewMux()

	api := huma.NewAPI(huma.DefaultConfig("Yotabo API", "1.0.0"), r)

	huma.Post(api, "/api/v1/auth/register", user.RegisterHandler)
	huma.Post(api, "/api/v1/auth/login", user.LoginHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
