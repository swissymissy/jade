package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	secret := "test-secret"
	adminID := uuid.New()

	token, err := MakeJWT(adminID, secret)
	if err != nil {
		t.Fatalf("unexpected error making token: %v", err)
	}

	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if returnedID != adminID {
		t.Errorf("expected %v, got %v", adminID, returnedID)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	adminID := uuid.New()

	token, err := MakeJWT(adminID, "correct-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Error("expected error for wrong secret, but got nil")
	}
}

func TestValidateJWTBadToken(t *testing.T) {
	_, err := ValidateJWT("not.a.real.token", "secret")
	if err == nil {
		t.Error("expected error for garbage token, but got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		expectToken string
		expectError bool
	}{
		{
			name:        "valid token",
			headerValue: "Bearer abc123",
			expectToken: "abc123",
			expectError: false,
		},
		{
			name:        "missing header",
			headerValue: "",
			expectError: true,
		},
		{
			name:        "wrong prefix",
			headerValue: "Token abc123",
			expectError: true,
		},
		{
			name:        "empty token",
			headerValue: "Bearer ",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := http.Header{}
			if tc.headerValue != "" {
				headers.Set("Authorization", tc.headerValue)
			}

			token, err := GetBearerToken(headers)
			if tc.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != tc.expectToken {
				t.Errorf("expected %s, got %s", tc.expectToken, token)
			}
		})
	}
}
