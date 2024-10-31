package jwtHelpers

import (
	"github.com/golang-jwt/jwt/v5"
	"go-api/internal/helpers/env"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/jwtType"
	"go-api/internal/models/enums/userRoles"
	"strconv"
	"time"
)

func CreateTokens(email string, role models.UserRoles, id int) (*string, *string, error) {
	jwtKey := []byte(env.GetVar("SECRET_KEY"))
	JwtAdmDuration, err := time.ParseDuration(env.GetVar("JWT_ADM_DURATION"))
	if err != nil {
		return nil, nil, err
	}
	JwtUsrDuration, err := time.ParseDuration(env.GetVar("JWT_USR_DURATION"))
	rJwtUsrDuration, err := time.ParseDuration(env.GetVar("R_JWT_USR_DURATION"))

	if err != nil {
		return nil, nil, err
	}
	timespan := utils.TernaryOperator(role != userRoles.Registered, JwtAdmDuration, JwtUsrDuration)
	expirationTime := time.Now().UTC().Add(timespan)
	claims := &models.JwtClaims{
		Email: email,
		Role:  role,
		Type:  jwtType.Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        strconv.Itoa(id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, nil, err
	}

	rTimespan := utils.TernaryOperator(role != userRoles.Registered, JwtAdmDuration, rJwtUsrDuration)
	rExpirationTime := time.Now().UTC().Add(rTimespan)
	rClaims := &models.JwtClaims{
		Email: email,
		Role:  role,
		Type:  jwtType.Refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        strconv.Itoa(id),
		},
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	rTokenString, rErr := rToken.SignedString(jwtKey)
	if rErr != nil {
		return nil, nil, rErr
	}

	return &tokenString, &rTokenString, nil
}
