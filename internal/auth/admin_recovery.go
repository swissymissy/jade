package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// function to generate a one-time recovery code every time admin registered or reset
func GenerateRecoveryCode() (string, error) {
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate recovery code: %w", err)
	}
	hex := hex.EncodeToString(bytes)
	return fmt.Sprintf("JADE-%s-%s-%s", hex[:4], hex[4:8], hex[8:12]), nil
}
