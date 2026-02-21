package middleware

import (
	"context"
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/jwt"
	"log/slog"
	"net/http"
	"strings"
)

const UserIDKey = "userID"

func JWTMiddleware(next http.Handler) http.Handler {
	jwtService := jwt.NewService(config.Load().App.AuthToken)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		userID, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			slog.Warn("failed validate token", slog.Any("error", err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
