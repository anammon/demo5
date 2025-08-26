package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	SignedToken, err := token.SignedString([]byte("secret"))

	return "Bearer" + SignedToken, err
}
