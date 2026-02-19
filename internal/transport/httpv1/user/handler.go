package user

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type RegisterInput struct {
	Body struct {
		FirstName string `json:"firstname" minLength:"1" maxLength:"50" pattern:"^[\\p{L}]+$"`
		LastName  string `json:"lastname"  minLength:"1" maxLength:"50" pattern:"^[\\p{L}]+$"`
		Username  string `json:"username"  minLength:"3" maxLength:"12"`
		Email     string `json:"email"     format:"email"`
		Password  string `json:"password"  minLength:"7" maxLength:"12"`
	}
}

type RegisterOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

type LoginInput struct {
	Body struct {
		Username string `json:"username" minLength:"3" maxLength:"12"`
		Password string `json:"password" minLength:"7" maxLength:"12"`
	}
}

type LoginOutput struct {
	Body struct {
		AccessToken string `json:"accessToken"`
	}
}

var secretKey = os.Getenv("AUTH_TOKEN")

// TODO: /api/v1/auth/register
func RegisterHandler(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Body.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to hash password")
	}

	user := models.User{
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Username:  input.Body.Username,
		Email:     input.Body.Email,
		Password:  string(hashedPassword),
	}

	if err := postgres.Db.Create(&user).Error; err != nil {
		return nil, huma.Error500InternalServerError("failed to create user")
	}

	resp := &RegisterOutput{}
	resp.Body.Message = "successful registration"

	return resp, nil
}

// TODO: /api/v1/auth/login
func LoginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {

	var user models.User

	if err := postgres.Db.
		Where("username = ?", input.Body.Username).
		First(&user).Error; err != nil {

		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Body.Password),
	); err != nil {

		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	payload := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to generate token")
	}

	resp := &LoginOutput{}
	resp.Body.AccessToken = tokenString

	return resp, nil
}
