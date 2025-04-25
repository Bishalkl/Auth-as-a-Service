package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassoword
func HashPassword(password string) (string, error) {
	// Generate a hashed password
	hashPassowrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashPassowrd), nil
}

// ComparePassword compares a plain password with a hashed password
func ComparePassword(hashPassword, password string) bool {
	// Compare the plain password with hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
