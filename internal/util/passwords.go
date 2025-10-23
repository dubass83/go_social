// Package util provides utility functions for password hashing and verification.
package util

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generate password hash or return error
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error: can not generate hash from password - %v", err)
	}
	return string(hash), nil
}

// CheckPassword check if provided password correct or not.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateToken(ID int64) string {
	return fmt.Sprintf("%d-%s", ID, uuid.New().String())
}
