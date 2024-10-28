package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Convert string type to []byte type
	bytePw := []byte(password)

	// Hashing the password
	hashedPw, err := bcrypt.GenerateFromPassword(bytePw, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate bcrypt password:%w", err)
	}

	// Convert []byte type to string type
	stringHashedPw := string(hashedPw)

	return stringHashedPw, nil
}

func CheckPasswordHash(password, hash string) error {

	// Compare hashed pw and pw in []byte type
	byteHashedPw := []byte(hash)
	byteProvidedPw := []byte(password)

	err := bcrypt.CompareHashAndPassword(byteHashedPw, byteProvidedPw)
	if err != nil {
		return fmt.Errorf("provided password is not same as stored password:%w", err)
	}
	return nil
}
