package middleware

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"net/http"
)

func HumaJWTMiddleware(ctx huma.Context, next func(huma.Context)) {
	r, w := humachi.Unwrap(ctx)

	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)

		ctx = huma.WithValue(ctx, UserIDKey, userID)
		next(ctx)
	}))

	handler.ServeHTTP(w, r)
}

func GetUserID(ctx context.Context) int {
	if val, ok := ctx.Value(UserIDKey).(int); ok {
		return val
	}

	return -1
}
