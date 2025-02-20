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
		errB := Body.Close()
		if errB != nil {
			log.Print("Error when closing request json body:", errB)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}(r.Body)

	neo4jUsr, errUsr := neo4jOperations.GetUser(usr.Email)
	if errUsr != nil {
		log.Print("Error fetching db:", errUsr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	errAuth := checks.UserAuthChecks(neo4jUsr, &usr.Password)
	if errAuth != nil {
		http.Error(w, errAuth.Error(), http.StatusUnauthorized)
		return
	}

	errDb := neo4jOperations.UpdateVisitCount(neo4jUsr.Id, neo4jUsr.LastConnectionDate)
	if errDb != nil {
		log.Print("Error updating database: ", errDb)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var jsonResponse []byte
	var errCreateToken error
	now := time.Now().UTC()

	//If it doesn't already have a token it generates one
	if neo4jUsr.Token == nil {
		jsonResponse, errCreateToken = CreateTokensWrapper(neo4jUsr)
		if errCreateToken != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		claims, errC := jwtHelpers.DecodeToken(*neo4jUsr.Token)
		if errC != nil {
			log.Print("Error decoding token:", errC)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		rClaims, errRC := jwtHelpers.DecodeToken(*neo4jUsr.RefreshToken)
		if errRC != nil {
			log.Print("Error decoding refresh token:", errRC)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		//If the token has expired it generates one
		if claims.ExpiresAt.Before(now) {
			jsonResponse, errCreateToken = CreateTokensWrapper(neo4jUsr)
			if errCreateToken != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		} else {
			response := models.LoginResponse{
				Token:                 transform.ConvertNullableToString(neo4jUsr.Token),
				RefreshToken:          transform.ConvertNullableToString(neo4jUsr.RefreshToken),
				TokenExpiresAt:        claims.ExpiresAt.Time.Unix(),
				RefreshTokenExpiresAt: rClaims.ExpiresAt.Time.Unix(),
			}
			jsonResponse, err = json.Marshal(response)
			if err != nil {
				log.Print("Failed to marshal result into JSON: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write json body in response.", http.StatusInternalServerError)
		return
	}
}
