package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email string    `json:"email"`
	Role  UserRoles `json:"role"`
	Type  JwtType   `json:"type"`
	jwt.RegisteredClaims
}
