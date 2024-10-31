package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {

	// Generate 32 bytes random data
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", fmt.Errorf("error in generating random 32 bytes data:%w", err)
	}

	// Convert the random data to a hex string
	refresh_token := hex.EncodeToString(randomData)

	return refresh_token, nil
}
