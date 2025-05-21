package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/models"
	"io"
	"log"
	"net/http"
)

func UpdateUserAccessStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	email_id := vars["id"]

	user, err := neo4jOperations.GetUserAccessRequest(email_id)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Access request not found", http.StatusNotFound)
		return
	}

	var userUpdate models.UpdateStatusAccessRequest

	err = json.NewDecoder(r.Body).Decode(&userUpdate)
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

	//err = checks.UserUpdateParamsChecks(*user, userUpdate)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	err = neo4jOperations.UpdateUserAccessRequestStatus(userUpdate)
	if err != nil {
		log.Print("Error updating db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result := models.Message{
		Status:  "OK",
		Message: "Access request updated",
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
