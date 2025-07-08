package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHashedPassword(password string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashed))
	return err == nil
}
