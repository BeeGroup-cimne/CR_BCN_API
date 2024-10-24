package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"go-api/internal/helpers/env"
	"go-api/internal/models"
	"time"
)

func CreateTokens(email string, role models.UserRoles) (*models.TokensResponse, error) {
	jwtKey := []byte(env.GetVar("SECRET_KEY"))
	expirationTime := time.Now().UTC().Add(20 * time.Minute) // Expires in 20 mins
	claims := &models.JwtClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	rExpirationTime := time.Now().UTC().Add(7 * 24 * time.Hour)
	rClaims := &models.JwtClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	rTokenString, rErr := rToken.SignedString(jwtKey)
	if rErr != nil {
		return nil, rErr
	}

	return &models.TokensResponse{
		Token:        tokenString,
		RefreshToken: rTokenString,
	}, nil
}
