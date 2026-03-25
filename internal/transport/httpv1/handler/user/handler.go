package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurgal1ev/yotabo-go/internal/config"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/handler/common"
	"github.com/nurgal1ev/yotabo-go/internal/transport/httpv1/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func GetUserHandler(ctx context.Context) (*common.HumaAPIResponse[UserDTO], error) {
	userID := middleware.GetUserID(ctx)
	user, err := gorm.G[models.User](postgres.Db).Where("id = ?", userID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("user not found")
		}
		slog.Error("failed get user", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	return common.NewHumaResponse(UserDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}), nil
}

func UpdateUserHandler(ctx context.Context, input *UpdateUserInput) (*common.HumaAPIResponse[UserDTO], error) {
	userID := middleware.GetUserID(ctx)
	user, err := gorm.G[models.User](postgres.Db).Where("id = ?", userID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("user not found")
		}
		slog.Error("failed get user", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	if input.Body.FirstName != nil {
		user.FirstName = *input.Body.FirstName
	}

	if input.Body.LastName != nil {
		user.LastName = *input.Body.LastName
	}

	if input.Body.Username != nil {
		user.Username = *input.Body.Username
	}

	if input.Body.Email != nil {
		user.Email = *input.Body.Email
	}

	_, err = gorm.G[models.User](postgres.Db).
		Where("id = ?", userID).
		Updates(ctx, user)

	if err != nil {
		slog.Error("failed update user", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("failed to update user")
	}

	return common.NewHumaResponse(UserDTO{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}), nil
}

func ChangePasswordHandler(ctx context.Context, input *ChangePasswordInput) (*common.HumaAPIResponse[any], error) {
	userID := middleware.GetUserID(ctx)

	user, err := gorm.G[models.User](postgres.Db).Where("id = ?", userID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("user not found")
		}
		slog.Error("failed get user", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Body.CurrentPassword))
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	user.Password = string(hashedPassword)
	err = postgres.Db.Save(&user).Error
	if err != nil {
		slog.Error("failed to save new password", slog.String("error", err.Error()))
		return nil, huma.Error500InternalServerError("internal server error")
	}

	return common.NewHumaResponse[any](nil, "password changed successfully"), nil
}
