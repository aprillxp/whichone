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

func GenerateJWT(playerID uint) (string, error) {
	claims := jwt.MapClaims{
		"player_id": playerID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ParseKeyWithClaims(tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return []byte("YOUR_SECRET_KEY"), nil
	})
}
