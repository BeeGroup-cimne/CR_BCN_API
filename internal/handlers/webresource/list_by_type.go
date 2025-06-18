package webresource

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/models"
	"log"
	"net/http"
)

func ListByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	typeStr := vars["type"]

	requestList, err := neo4jOperations.GetWebResourceListByType(typeStr)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if requestList == nil {
		requestList = []*models.WebResource{}
	}

	jsonResponse, err := json.Marshal(requestList)
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
