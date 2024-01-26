package handler

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"go-conversation/model"
)

func generateJWTToken(user model.User, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 day)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
