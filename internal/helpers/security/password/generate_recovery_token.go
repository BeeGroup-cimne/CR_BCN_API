package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

func GenerateRecoveryToken(expiresInHours int) (string, time.Time, string, error) {
	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return "", time.Time{}, "", err
	}

	token := base64.URLEncoding.EncodeToString(raw)
	hash := sha256.Sum256([]byte(token))
	hashedToken := base64.URLEncoding.EncodeToString(hash[:])
	expiresAt := time.Now().Add(time.Duration(expiresInHours) * time.Hour)

	return token, expiresAt, hashedToken, nil
}
