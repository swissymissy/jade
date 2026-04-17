package storage

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// generate randome filename
func RandomFilename(extension string) (string, error) {
	randomSlice := make([]byte, 32)
	_, err := rand.Read(randomSlice)
	if err != nil {
		return "", fmt.Errorf("failed to generate random filename: %w", err)
	}
	randomName := base64.RawURLEncoding.EncodeToString(randomSlice) + extension
	return randomName, nil
}
