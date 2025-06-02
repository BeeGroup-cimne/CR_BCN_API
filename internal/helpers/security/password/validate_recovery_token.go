package password

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"go-api/internal/models"
	"time"
)

func ValidateRecoveryToken(user *models.User, inputToken string) error {
	if user.PasswordResetTokenExpiresAt == nil || time.Now().After(*user.PasswordResetTokenExpiresAt) {
		return errors.New("token has expired")
	}

	hash := sha256.Sum256([]byte(inputToken))
	hashedInput := base64.URLEncoding.EncodeToString(hash[:])

	if *user.PasswordResetToken != hashedInput {
		return errors.New("invalid token")
	}

	return nil
}
