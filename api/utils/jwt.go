package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret []byte

func SetJWTSecret(secret string) {
	if secret == "" {
		panic("JWT_SECRET environment variable not set!")
	}
	JWTSecret = []byte(secret)
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ParseKey(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return JWTSecret, nil
	})
}
