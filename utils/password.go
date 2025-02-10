package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password.
// Returns the hashed password or an error.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err) // Keep this for internal debugging
		return "", err                                // Return error to the caller for further handling
	}
	return string(bytes), nil
}

// VerifyPassword compares a hashed password with a plain-text password.
// Returns nil if the password matches, otherwise an error.
func VerifyPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Printf("Password mismatch: %v", err) // Keep this log for failed verifications
	}
	return err
}
