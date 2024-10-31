package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {

	timeDuration := 1 * time.Second
	tokenSecret1 := "correctPassword123!"
	tokenSecret2 := "anotherPassword456!"
	userID1, _ := uuid.Parse("588679f1-e38f-404d-a1cf-233fb9cff445")
	userID2, _ := uuid.Parse("4a24b414-ccd8-4bfa-b680-04f0caf9654f")

	tokenString1, _ := MakeJWT(userID1, tokenSecret1, timeDuration)
	tokenString2, _ := MakeJWT(userID2, tokenSecret2, timeDuration)

	tests := []struct {
		name        string
		waitTime    time.Duration
		tokenSecret string
		signedToken string
		wantErr     bool
	}{
		{
			name:        "Correct signedToken and tokenSecret in validation",
			tokenSecret: tokenSecret1,
			signedToken: tokenString1,
			waitTime:    0 * time.Second,
			wantErr:     false,
		},
		{
			name:        "Empty signedToken in validation",
			tokenSecret: tokenSecret1,
			signedToken: "",
			waitTime:    0 * time.Second,
			wantErr:     true,
		},
		{
			name:        "Expired token in validation",
			tokenSecret: tokenSecret1,
			signedToken: tokenString1,
			waitTime:    3 * time.Second,
			wantErr:     true,
		},
		{
			name:        "Other signedToken in validation",
			tokenSecret: tokenSecret1,
			signedToken: tokenString2,
			waitTime:    0 * time.Second,
			wantErr:     true,
		},
		{
			name:        "Other tokenSecret in validation",
			tokenSecret: tokenSecret2,
			signedToken: "",
			waitTime:    0 * time.Second,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.waitTime)
			_, err := ValidateJWT(tt.signedToken, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error:%v, wantErr:%v", err, tt.wantErr)
			}
		})
	}

}
