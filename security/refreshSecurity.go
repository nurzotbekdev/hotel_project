package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ParseRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid claims")
	}

	if claims["type"] != "refresh" {
		return 0, errors.New("This is not a refresh token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return 0, errors.New("Refresh token has expired")
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("Invalid sub")
	}
	userID := uint(sub)

	return uint(userID), nil
}
