package routes

import (
	"github.com/gorilla/mux"
	"go-api/internal/handlers/auth"
	"go-api/internal/handlers/data"
	"go-api/internal/handlers/user"
	"go-api/internal/middleware"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	dataRoutes := r.PathPrefix("/data").Subrouter()
	dataRoutes.Use(middleware.DataMiddleware)
	dataRoutes.HandleFunc("/collection", data.GetCollectionHandler).Methods(http.MethodGet)

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	authRoutes.Handle("/refresh_token", middleware.RefreshTokenMiddleware(http.HandlerFunc(auth.RefreshToken))).Methods(http.MethodGet)

	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.AdminMiddleware)
	userRoutes.HandleFunc("/create", user.Post).Methods(http.MethodPost)
	userRoutes.HandleFunc("/{id}", user.Get).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}/update", user.Put).Methods(http.MethodPut)
	userRoutes.HandleFunc("/{id}/delete", user.Delete).Methods(http.MethodDelete)

	return r
}
