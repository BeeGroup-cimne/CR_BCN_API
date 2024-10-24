package handlers

import (
	"encoding/json"
	"fmt"
	"go-api/internal/helpers/auth"
	"go-api/internal/helpers/dbOperations"
	"go-api/internal/helpers/jwt"
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
	var usr models.UserLogin

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Error when closing request json body:", err)
		}
	}(r.Body)

	neo4jUsr, err := dbOperations.FetchNeo4jUser(usr.Email, usr.Password)
	if err != nil {
		log.Fatal("Error fetching dbOperations:", err)
		//http.Error(w, "Internal server error", http.StatusInternalServerError)
		//return
	}

	if !auth.CheckPasswordHash(usr.Password, neo4jUsr.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if neo4jUsr.Blocked != nil && *neo4jUsr.Blocked {
		http.Error(w, "User is blocked", http.StatusUnauthorized)
		return
	}

	if neo4jUsr.ExpirationDate != nil && time.Now().UTC().After(transform.ConvertIsoStringToTime(*neo4jUsr.ExpirationDate)) {
		http.Error(w, "User no longer has access", http.StatusUnauthorized)
		return
	}

	if neo4jUsr.StartDate != nil && time.Now().UTC().Before(transform.ConvertIsoStringToTime(*neo4jUsr.StartDate)) {
		http.Error(w, fmt.Sprintf("User does not have access yet. It will start on %s", *neo4jUsr.StartDate), http.StatusUnauthorized)
		return
	}

	now := time.Now().UTC()
	lastConnection := utils.TernaryOperator(neo4jUsr.LastConnectionDate != nil, now, transform.ConvertIsoStringToTime(*neo4jUsr.LastConnectionDate))
	isNewHour := utils.IsEqualDate(lastConnection, now, timePrecision.Hour)
	isNewDay := utils.IsEqualDate(lastConnection, now, timePrecision.Day)
	if neo4jUsr.MaxHourlyConnections != nil && neo4jUsr.AccumulatedHourlyConnections != nil && isNewHour && *neo4jUsr.MaxHourlyConnections >= *neo4jUsr.AccumulatedHourlyConnections {
		http.Error(w, fmt.Sprintf("Access hourly limit of %s reached.", *neo4jUsr.MaxHourlyConnections), http.StatusUnauthorized)
		return
	}
	if neo4jUsr.MaxDailyConnections != nil && neo4jUsr.AccumulatedDailyConnections != nil && isNewDay && *neo4jUsr.MaxDailyConnections >= *neo4jUsr.AccumulatedDailyConnections {
		http.Error(w, fmt.Sprintf("Access daily limit of %s reached.", *neo4jUsr.MaxDailyConnections), http.StatusUnauthorized)
		return
	}

	var response *models.TokensResponse

	//If it doesn't already have a token it generates one
	if neo4jUsr.Token == nil {
		newTokens, err := jwt.CreateTokens(neo4jUsr.Email, neo4jUsr.Role)
		if err != nil {
			log.Fatal("Error creating token:", err)
		}
		response = newTokens
	} else {
		//If it already has a token it decodes it and checks the expiration date
		claims, err := jwt.DecodeToken(*neo4jUsr.Token)
		if err != nil {
			log.Fatal("Error decoding token:", err)
		}
		//If the token has expired it generates one
		if claims.ExpiresAt.Before(now) {
			newTokens, err := jwt.CreateTokens(neo4jUsr.Email, neo4jUsr.Role)
			if err != nil {
				log.Fatal("Error creating token:", err)
			}
			response = newTokens
		} else {
			//If the token has not expired it gets the token from the db
			response.Token = *neo4jUsr.Token
			response.RefreshToken = *neo4jUsr.RefreshToken
		}
	}
	jsonResponse, err := json.Marshal(&response)
	if err != nil {
		log.Fatal("Failed to marshal result into JSON: ", err)
	}

	err = dbOperations.UpdateVisitCount(*neo4jUsr, response, isNewHour, isNewDay)
	if err != nil {
		log.Fatal("Error updating database: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write json body in response.", http.StatusInternalServerError)
		return
	}
}
