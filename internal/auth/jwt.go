package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// create new access token for when admin log in
func MakeJWT(adminID uuid.UUID, serverSecretToken string) (string, error) {

	// create a new registered claim
	claim := jwt.RegisteredClaims{
		Issuer:    "jade-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		Subject:   adminID.String(),
	}

	// create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// sign the token wuth server secret key
	signedKey := []byte(serverSecretToken)
	signedToken, err := token.SignedString(signedKey)
	if err != nil {
		return "", fmt.Errorf("cannot sign token: %w", err)
	}
	return signedToken, nil
}

// check admin's token
func ValidateJWT(tokenString, serverSecretToken string) (uuid.UUID, error) {
	// create new empty claim struct to be filled
	claim := &jwt.RegisteredClaims{}

	// pass a pointer to that struct so the library can modify it
	_, err := jwt.ParseWithClaims(
		tokenString,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(serverSecretToken), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Token is expired or bad signature: %w", err)
	}

	// retrieve admin ID from claim's subject
	adminIDstr := claim.Subject
	adminID, err := uuid.Parse(adminIDstr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error converting string to uuid: %w", err)
	}
	return adminID, nil
}
