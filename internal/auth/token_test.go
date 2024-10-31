package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		tokenString string
		wantErr     bool
	}{
		{
			name:        "Valid Bearer Token",
			headers:     http.Header{"Authorization": []string{"Bearer my-valid-token"}},
			tokenString: "my-valid-token",
			wantErr:     false,
		}, {
			name:        "Empty Authorization Header",
			headers:     http.Header{},
			tokenString: "",
			wantErr:     true,
		}, {
			name:        "No Bearer Prefix",
			headers:     http.Header{"Authorization": []string{"my-valid-token"}},
			tokenString: "",
			wantErr:     true,
		}, {
			name:        "Whitespace Handling",
			headers:     http.Header{"Authorization": []string{"Bearer    my-valid-token"}},
			tokenString: "my-valid-token",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if token != tt.tokenString || ((err != nil) != tt.wantErr) {
				t.Errorf("GetBearerToken() error:%v token:%v, wantErr:%v , wantToken:%v", err, token, tt.wantErr, tt.tokenString)
			}
		})
	}
}
