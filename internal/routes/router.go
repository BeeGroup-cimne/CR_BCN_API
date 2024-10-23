package routes

import (
	"github.com/gorilla/mux"
	"go-api/internal/handlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/collection", handlers.GetCollectionHandler).Methods("GET")
	return r
}
