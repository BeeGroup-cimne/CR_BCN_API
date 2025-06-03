package user

import (
	"encoding/json"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/security/checks"
	"go-api/internal/models"
	"io"
	"log"
	"net/http"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user models.UserChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&user)
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

	err = checks.UserUpdatePasswordParamsChecks(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = neo4jOperations.UpdateUser(models.UpdateUserRequest{Password: &user.Password}, user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result := models.Message{
		Status:  "OK",
		Message: "Password updated",
	}

	jsonResponse, err := json.Marshal(result)
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
