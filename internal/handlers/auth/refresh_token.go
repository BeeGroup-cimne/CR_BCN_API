package auth

import (
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/security/checks"
	"go-api/internal/models"
	"log"
	"net/http"
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

	neo4jUsr, err := neo4jOperations.GetUser(claims.Email)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = checks.UserAuthChecks(neo4jUsr, nil)
	if err != nil {
		log.Print("Failed to authenticate user: ", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = neo4jOperations.UpdateVisitCount(neo4jUsr.Id, neo4jUsr.LastConnectionDate)
	if err != nil {
		log.Print("Error updating database: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := CreateTokensWrapper(neo4jUsr)
	if err != nil {
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
