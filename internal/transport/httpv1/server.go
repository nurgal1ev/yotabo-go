package httpv1

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/tasks"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/user"
	"net/http"
)

func StartServer() {
	r := chi.NewMux()

	api := humachi.New(r, huma.DefaultConfig("Yotabo API", "1.0.0"))

	huma.Post(api, "/api/v1/auth/register", user.RegisterHandler)
	huma.Post(api, "/api/v1/auth/login", user.LoginHandler)
	huma.Post(api, "/api/v1/tasks/", tasks.CreateTaskHandler)
	//huma.Get(api, "/api/v1/tasks/", tasks.GetTaskHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
