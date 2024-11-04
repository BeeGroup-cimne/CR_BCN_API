package auth

import (
	"encoding/json"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/jwtHelpers"
	"go-api/internal/helpers/security/checks"
	"go-api/internal/helpers/transform"
	"go-api/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var usr models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error when closing request json body:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}(r.Body)

	neo4jUsr, err := neo4jOperations.GetUser(usr.Email)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = checks.UserAuthChecks(neo4jUsr, &usr.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = neo4jOperations.UpdateVisitCount(neo4jUsr.Id, neo4jUsr.LastConnectionDate)
	if err != nil {
		log.Print("Error updating database: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response *models.LoginResponse

	now := time.Now().UTC()

	//If it doesn't already have a token it generates one
	if neo4jUsr.Token == nil {
		token, refreshToken, err := jwtHelpers.CreateTokens(neo4jUsr.Email, neo4jUsr.Role, neo4jUsr.Id)
		if err != nil {
			log.Print("Error creating token:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		response = &models.LoginResponse{Token: *token, RefreshToken: *refreshToken}
		err = neo4jOperations.UpdateUserTokens(neo4jUsr.Id, response.Token, response.RefreshToken)
		if err != nil {
			log.Print("Error updating database: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		//If it already has a token it decodes it and checks the expiration date
		claims, err := jwtHelpers.DecodeToken(*neo4jUsr.Token)
		if err != nil {
			log.Print("Error decoding token:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		//If the token has expired it generates one
		if claims.ExpiresAt.Before(now) {
			token, refreshToken, err := jwtHelpers.CreateTokens(neo4jUsr.Email, neo4jUsr.Role, neo4jUsr.Id)
			if err != nil {
				log.Print("Error creating token:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			response = &models.LoginResponse{Token: *token, RefreshToken: *refreshToken}
			err = neo4jOperations.UpdateUserTokens(neo4jUsr.Id, response.Token, response.RefreshToken)
			if err != nil {
				log.Print("Error updating database: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			response = &models.LoginResponse{Token: transform.ConvertNullableToString(neo4jUsr.Token),
				RefreshToken: transform.ConvertNullableToString(neo4jUsr.RefreshToken)}
		}
	}
	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		log.Print("Failed to marshal result into JSON: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write json body in response.", http.StatusInternalServerError)
		return
	}
}
