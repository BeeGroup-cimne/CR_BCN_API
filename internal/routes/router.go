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
	r.Use(middleware.CorsMiddleware)
	dataRoutes := r.PathPrefix("/data").Subrouter()
	dataRoutes.Use(middleware.DataMiddleware)
	//dataRoutes.HandleFunc("/swagger/*", httpSwagger.WrapHandler)
	dataRoutes.HandleFunc("/climate_communications", data.GetClimateCommunicationsHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/{collection_name}", data.GetCollectionHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/geometry_collection/{collection_name}", data.GetGeometryCollectionHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/aggregation_levels/{lang}", data.GetAggregationLevelsHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators/{lang}", data.GetIndicatorsHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators_by_group/{agregation_level}/{lang}", data.GetIndicatorsByGroupHandler).Methods(http.MethodGet)

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	authRoutes.Handle("/refresh_token", middleware.RefreshTokenMiddleware(http.HandlerFunc(auth.RefreshToken))).Methods(http.MethodGet)

	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.AdminMiddleware)
	userRoutes.HandleFunc("/create", user.Post).Methods(http.MethodPost)
	userRoutes.HandleFunc("/{id}", user.Get).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}/update", user.Put).Methods(http.MethodPut)
	userRoutes.HandleFunc("/{id}/delete", user.Delete).Methods(http.MethodDelete)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	return r
}
