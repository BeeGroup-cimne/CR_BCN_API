package handlers

import (
	"encoding/json"
	"go-api/internal/helpers/auth"
	"go-api/internal/helpers/dbOperations"
	"go-api/internal/helpers/jwtHelpers"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/jwtType"
	"go-api/internal/models/enums/timePrecision"
	"log"
	"net/http"
	"time"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(*models.JwtClaims)
	if !ok {
		http.Error(w, "Could not retrieve claims", http.StatusUnauthorized)
		return
	}

	if claims.Type != jwtType.Refresh {
		http.Error(w, "Token must be a valid refresh token", http.StatusBadRequest)
		return
	}

	neo4jUsr, err := dbOperations.FetchNeo4jUser(claims.Email)
	if err != nil {
		log.Print("Error fetching dbOperations:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = auth.UserAuthChecks(neo4jUsr, nil)
	if err != nil {
		log.Print("Failed to authenticate user: ", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokens, err := jwtHelpers.CreateTokens(claims.Email, claims.Role)
	if err != nil {
		log.Print("Error creating token:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	lastConnection := utils.TernaryOperator(neo4jUsr.LastConnectionDate == nil, now, transform.ConvertIsoStringToTime(transform.ConvertNullableToString(neo4jUsr.LastConnectionDate)))
	isNewHour := utils.IsEqualDate(lastConnection, now, timePrecision.Hour)
	isNewDay := utils.IsEqualDate(lastConnection, now, timePrecision.Day)

	err = dbOperations.UpdateVisitCount(neo4jUsr, tokens, isNewHour, isNewDay)
	if err != nil {
		log.Print("Error updating database: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(tokens)
	if err != nil {
		log.Print("Failed to marshal result into JSON: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Print("Failed to write json body in response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
