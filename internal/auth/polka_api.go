package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	// Extract api key from Authorization header
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("there is no authorization in the header")
	}

	if !strings.Contains(authorization, "ApiKey") {
		return "", fmt.Errorf("authorization header does not contain ApiKey")
	}

	apiString := strings.TrimSpace(strings.Replace(authorization, "ApiKey", "", 1))

	return apiString, nil
}
