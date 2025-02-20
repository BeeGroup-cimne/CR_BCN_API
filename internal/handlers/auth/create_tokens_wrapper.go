package auth

import (
	"encoding/json"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/jwtHelpers"
	"go-api/internal/models"
	"log"
)

func CreateTokensWrapper(neo4jUsr *models.User) ([]byte, error) {
	token, refreshToken, err := jwtHelpers.CreateTokens(neo4jUsr.Email, neo4jUsr.Role, neo4jUsr.Id)
	if err != nil {
		log.Print("Error creating token:", err)
		return nil, err
	}
	claims, errC := jwtHelpers.DecodeToken(*neo4jUsr.Token)
	if errC != nil {
		log.Print("Error decoding token:", errC)
		return nil, errC
	}
	rClaims, errRC := jwtHelpers.DecodeToken(*neo4jUsr.RefreshToken)
	if errRC != nil {
		log.Print("Error decoding refresh token:", errRC)
		return nil, errRC
	}
	response := &models.LoginResponse{
		Token:                 *token,
		RefreshToken:          *refreshToken,
		TokenExpiresAt:        claims.ExpiresAt.Time.Unix(),
		RefreshTokenExpiresAt: rClaims.ExpiresAt.Time.Unix(),
	}
	errDb := neo4jOperations.UpdateUserTokens(neo4jUsr.Id, response.Token, response.RefreshToken)
	if errDb != nil {
		log.Print("Error updating database: ", errDb)
		return nil, errDb
	}
	jsonResponse, errJson := json.Marshal(response)
	if errJson != nil {
		log.Print("Failed to marshal result into JSON: ", err)
		return nil, errJson
	}
	return jsonResponse, nil
}
