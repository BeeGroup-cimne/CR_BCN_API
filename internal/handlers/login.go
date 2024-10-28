package handlers

import (
	"encoding/json"
	"go-api/internal/helpers/auth"
	"go-api/internal/helpers/dbOperations"
	"go-api/internal/helpers/jwtHelpers"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/timePrecision"
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

	neo4jUsr, err := dbOperations.FetchNeo4jUser(usr.Email)
	if err != nil {
		log.Print("Error fetching dbOperations:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = auth.UserAuthChecks(neo4jUsr, &usr.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var response *models.LoginResponse

	now := time.Now().UTC()
	lastConnection := utils.TernaryOperator(neo4jUsr.LastConnectionDate == nil, now, transform.ConvertIsoStringToTime(transform.ConvertNullableToString(neo4jUsr.LastConnectionDate)))
	isNewHour := utils.IsEqualDate(lastConnection, now, timePrecision.Hour)
	isNewDay := utils.IsEqualDate(lastConnection, now, timePrecision.Day)

	//If it doesn't already have a token it generates one
	if neo4jUsr.Token == nil {
		newTokens, err := jwtHelpers.CreateTokens(neo4jUsr.Email, neo4jUsr.Role)
		if err != nil {
			log.Print("Error creating token:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		response = newTokens
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
			newTokens, err := jwtHelpers.CreateTokens(neo4jUsr.Email, neo4jUsr.Role)
			if err != nil {
				log.Print("Error creating token:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			response = newTokens
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

	err = dbOperations.UpdateVisitCount(neo4jUsr, response, isNewHour, isNewDay)
	if err != nil {
		log.Print("Error updating database: ", err)
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
