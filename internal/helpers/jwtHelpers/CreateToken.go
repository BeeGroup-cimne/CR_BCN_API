package jwtHelpers

import (
	"github.com/golang-jwt/jwt/v5"
	"go-api/internal/helpers/env"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/jwtType"
	"go-api/internal/models/enums/userRoles"
	"time"
)

func CreateTokens(email string, role models.UserRoles) (*models.LoginResponse, error) {
	jwtKey := []byte(env.GetVar("SECRET_KEY"))
	JwtAdmDuration, err := time.ParseDuration(env.GetVar("JWT_ADM_DURATION"))
	if err != nil {
		return nil, err
	}
	JwtUsrDuration, err := time.ParseDuration(env.GetVar("JWT_USR_DURATION"))
	rJwtUsrDuration, err := time.ParseDuration(env.GetVar("R_JWT_USR_DURATION"))

	if err != nil {
		return nil, err
	}
	timespan := utils.TernaryOperator(role != userRoles.Registered, JwtAdmDuration, JwtUsrDuration)
	rTimespan := utils.TernaryOperator(role != userRoles.Registered, JwtAdmDuration, rJwtUsrDuration)
	expirationTime := time.Now().UTC().Add(timespan)
	claims := &models.JwtClaims{
		Email: email,
		Role:  role,
		Type:  jwtType.Access,
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

	rExpirationTime := time.Now().UTC().Add(rTimespan)
	rClaims := &models.JwtClaims{
		Email: email,
		Role:  role,
		Type:  jwtType.Refresh,
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

	return &models.LoginResponse{
		Token:        tokenString,
		RefreshToken: rTokenString,
	}, nil
}
