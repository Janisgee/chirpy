package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {

	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("there is no authorization in the header")
	}
	if !strings.Contains(authorization, "Bearer") {
		return "", fmt.Errorf("authorization header does not contain 'Bearer'")
	}
	tokenString := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))

	return tokenString, nil
}
