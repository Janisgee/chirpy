package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Create JSON Web Token
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	// Make sure input has tokenSecret
	if len(tokenSecret) == 0 {
		return "", errors.New("there is no tokenSecret provided.\n ")
	}

	// Get info
	currentTime := time.Now()
	ExpiresIn := currentTime.Add(expiresIn)
	signingKey := []byte(tokenSecret)

	// Convert info to jwt.NumericDate
	currentTimeNumericDate := jwt.NewNumericDate(currentTime)
	ExpiresInNumericDate := jwt.NewNumericDate(ExpiresIn)
	userIDString := userID.String()

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  currentTimeNumericDate,
		ExpiresAt: ExpiresInNumericDate,
		Subject:   userIDString,
	})

	// Sign the token using the signing Key to generte final signed JWT string
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error signing the token:%w", err)
	}

	return signedToken, nil
}

// Validate the signature of the JWT and extract the claims into a *jwt.Token struct
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var id uuid.UUID

	// Make sure input has tokenSecret
	if len(tokenSecret) == 0 {
		return id, errors.New("there is no tokenSecret provided.\n ")
	}

	// Create Claims object base on provided tokenString and token Secret.
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {

		return id, fmt.Errorf("error when validate the signature of Json Web Token:%w", err)
	}

	// Extract userID string from Claims
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return id, fmt.Errorf("error when access claims of user's id:%w", err)
	}
	// Convert string to UUID type
	id, err = uuid.Parse(userIDString)
	if err != nil {
		return id, fmt.Errorf("error when parsing string to UUID type:%w", err)
	}
	return id, nil
}
