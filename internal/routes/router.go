package routes

import (
	"github.com/gorilla/mux"
	"go-api/internal/handlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/collection", handlers.GetCollectionHandler).Methods("GET")
	r.HandleFunc("/test", handlers.GetNeo4jTestHandler).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	return r
}
