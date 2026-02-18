package httpv1

import (
	"github.com/go-chi/chi/v5"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/user"
	"net/http"
)

func StartServer() {
	r := chi.NewRouter()

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", user.LoginHandler)
		r.Post("/register", user.RegisterHandler)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
