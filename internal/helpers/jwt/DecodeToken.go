package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-api/internal/helpers/env"
	"go-api/internal/models"
)

func DecodeToken(tokenString string) (*models.JwtClaims, error) {
	SecretKey := env.GetVar("SECRET_KEY")
	// Parse the token with a custom claims type
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the method used to sign the token is HMAC-SHA256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and the claims are of the correct type
	if claims, ok := token.Claims.(*models.JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
