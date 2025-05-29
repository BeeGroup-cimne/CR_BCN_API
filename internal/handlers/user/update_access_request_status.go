package user

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/security/checks"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/UserAccesRequestStatus"
	"go-api/internal/models/enums/userRoles"
	"io"
	"log"
	"net/http"
	"strconv"
)

func UpdateUserAccessStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	str := vars["id"]
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Error converting id to int64:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	if !utils.IsValidEnum(userUpdate.Status, userAccessRequestStatus.AllStatus) {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	uar, err := neo4jOperations.GetUserAccessRequest(id)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if uar == nil {
		http.Error(w, "Access request not found", http.StatusNotFound)
		return
	}

	err = neo4jOperations.UpdateUserAccessRequestStatus(id, userUpdate)
	if err != nil {
		log.Print("Error updating db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var result models.Message

	if userUpdate.Status == userAccessRequestStatus.Approved {
		err = checks.UserAccessRequestCreateParamsChecks(*uar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := models.CreateUserRequest{
			Email: uar.Email,
			//Password:   uar.Password,
			Role:       userRoles.Registered,
			Entity:     &uar.Entity,
			Department: &uar.Department,
			Name:       &uar.Name,
		}
		id, err := neo4jOperations.CreateUser(user)
		if err != nil {
			log.Print("Error creating user in db:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		result = models.Message{
			Status:  "OK",
			Message: *id,
		}
	} else {
		result = models.Message{
			Status:  "OK",
			Message: "Access request denied.",
		}
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
