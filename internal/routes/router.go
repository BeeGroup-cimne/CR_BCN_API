package routes

import (
	"github.com/gorilla/mux"
	"go-api/internal/handlers"
	"go-api/internal/helpers/auth"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/collection", handlers.GetCollectionHandler).Methods("GET")
	r.HandleFunc("/test", handlers.GetNeo4jTestHandler).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.Handle("/refresh_token", auth.JWTAuthMiddleware(http.HandlerFunc(handlers.RefreshToken))).Methods("GET")
	r.Handle("/password", auth.JWTAuthMiddleware(http.HandlerFunc(handlers.Password))).Methods("POST")
	return r
}
