package webresource

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/helpers/db/neo4jOperations"
	"log"
	"net/http"
	"strconv"
)

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format, must be an integer", http.StatusBadRequest)
		return
	}

	wr, err := neo4jOperations.GetWebResource(id)
	if err != nil {
		log.Print("Error fetching db:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(wr)
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
