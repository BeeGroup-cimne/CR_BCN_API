package user

import (
	"encoding/json"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/mailing"
	"go-api/internal/helpers/security/password"
	"go-api/internal/models"
	"io"
	"log"
	"net/http"
)

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user models.UserRecoverPasswordRequest

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

	userExists, err := neo4jOperations.UserAccessRequestExists(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !userExists {
		http.Error(w, "Email does not exist.", http.StatusBadRequest)
		return
	}

	token, expiresAt, hashedToken, err := password.GenerateRecoveryToken(24)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = neo4jOperations.UpdateUser(models.UpdateUserRequest{PasswordResetTokenHash: &hashedToken, PasswordResetTokenExpiresAt: &expiresAt}, user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	message, err := mailing.CreatePasswordRecoveryMessage(user.Email, token, 24)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = mailing.SendEmail(user.Email, message, "Password Reset in ClimateReady")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	result := models.Message{
		Status:  "OK",
		Message: "Password recovery email sent",
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
