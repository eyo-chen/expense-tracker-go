package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateHashPassword is a function that generates a hash password.
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CompareHashPassword is a function that compares a password and a hash password.
func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
