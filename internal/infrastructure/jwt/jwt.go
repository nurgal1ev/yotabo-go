package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretKey string
}

func NewService(secretKey string) *Service {
	return &Service{secretKey: secretKey}
}

func (service *Service) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id missing")
	}

	return int(id), nil
}
