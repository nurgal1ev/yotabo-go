package user

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

func RegisterHandler(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Body.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		slog.Error(
			"failed generate hash from password",
			slog.String("error", err.Error()),
		)
		return nil, huma.Error500InternalServerError("failed to hash password")
	}

	user := models.User{
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Username:  input.Body.Username,
		Email:     input.Body.Email,
		Password:  string(hashedPassword),
	}

	err = postgres.Db.Create(&user).Error
	if err != nil {
		slog.Error("failed create user", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("failed to create user")
	}

	resp := &RegisterOutput{}
	resp.Body.Message = "successful registration"

	return resp, nil
}

func LoginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	var user models.User

	err := postgres.Db.Where("username = ?", input.Body.Username).First(&user).Error
	if err != nil {
		slog.Error(
			"failed get user by username",
			slog.String("error", err.Error()),
			slog.String("username", input.Body.Username),
		)
		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Body.Password),
	)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	payload := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(config.Load().App.AuthToken))
	if err != nil {
		slog.Error("failed generate token", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("failed to generate token")
	}

	resp := &LoginOutput{}
	resp.Body.AccessToken = tokenString

	return resp, nil
}
