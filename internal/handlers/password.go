package handlers

import (
	"encoding/json"
	"go-api/internal/helpers/auth"
	"go-api/internal/models"
	"io"
	"log"
	"net/http"
)

type passwordRepsonse struct {
	Hash string `json:"hash"`
}

func Password(w http.ResponseWriter, r *http.Request) {
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
			log.Fatal("Error when closing request json body:", err)
		}
	}(r.Body)

	pass, err := auth.HashPassword(usr.Password)
	if err != nil {
		log.Fatal("Failed to hash password: ", err)
	}

	var response passwordRepsonse
	response.Hash = pass

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Failed to marshal result into JSON: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write json body in response.", http.StatusInternalServerError)
		return
	}
}
