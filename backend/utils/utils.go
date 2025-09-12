package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(account string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	SignedToken, err := token.SignedString([]byte("secret"))

	return "Bearer " + SignedToken, err
}
func PaserJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err //解析出现错误
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if account, ok := claims["account"].(string); ok {
			return account, nil //成功返回用户名
		}
		return "", errors.New("account not found in token")
	}
	return "", errors.New("invalid token") //无效token
}
