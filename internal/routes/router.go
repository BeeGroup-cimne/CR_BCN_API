package routes

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-api/internal/handlers/auth"
	"go-api/internal/handlers/data"
	"go-api/internal/handlers/user"
	"go-api/internal/handlers/webresource"
	"go-api/internal/middleware"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.CorsMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	dataRoutes := r.PathPrefix("/data").Subrouter()
	dataRoutes.Use(middleware.DataMiddleware)
	dataRoutes.HandleFunc("/climate_communications", data.GetClimateCommunicationsHandler).Methods(http.MethodPost)
	dataRoutes.HandleFunc("/collection/vulnerability_index", data.GetVulnerabilityIndex).Methods(http.MethodPost)
	dataRoutes.HandleFunc("/geometry_collection_vulnerability_index", data.GetGeometryVulnerabilityIndex).Methods(http.MethodPost)
	dataRoutes.HandleFunc("/climate_shelters", data.GetClimateShelters).Methods(http.MethodPost)
	dataRoutes.HandleFunc("/climate_shelters_details/{climate_shelter_id}", data.GetClimateSheltersDetails).Methods(http.MethodPost)
	dataRoutes.HandleFunc("/collection/{collection_name}", data.GetCollectionHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/geometry_collection/{collection_name}", data.GetGeometryCollectionHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/aggregation_levels/{lang}", data.GetAggregationLevelsHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators/{lang}", data.GetIndicatorsHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators_by_group/{aggregation_level}/{lang}", data.GetIndicatorsByGroupHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators_by_reference/{aggregation_level}/{reference}", data.GetIndicatorsByReferenceHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/collection/indicators_by_reference_data/{aggregation_level}/{lang}", data.GetIndicatorsByReferenceDataHandler).Methods(http.MethodGet)
	dataRoutes.HandleFunc("/indoor_sensor_data/{property}/building/{aggregation_time}/{start}:{end}", data.GetIndoorIndicators).Methods(http.MethodGet)

	adminDataRoutes := r.PathPrefix("/data").Subrouter()
	adminDataRoutes.Use(middleware.AdminMiddleware)
	adminDataRoutes.HandleFunc("/indoor_sensor_data/{property}/{aggregation_level}/{aggregation_time}/{start}:{end}", data.GetIndoorIndicators).Methods(http.MethodGet)

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	authRoutes.Handle("/refresh_token", middleware.RefreshTokenMiddleware(http.HandlerFunc(auth.RefreshToken))).Methods(http.MethodGet)

	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.AdminMiddleware)
	r.HandleFunc("/user/access_request", user.AccessRequest).Methods(http.MethodPost)
	r.HandleFunc("/user/change_password", user.ChangePassword).Methods(http.MethodPost)
	r.HandleFunc("/user/recover_password", user.RecoverPassword).Methods(http.MethodPost)

	userRoutes.HandleFunc("/create", user.Post).Methods(http.MethodPost)
	userRoutes.HandleFunc("/list", user.List).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}", user.Get).Methods(http.MethodGet)
	userRoutes.HandleFunc("/{id}", user.Put).Methods(http.MethodPut)
	userRoutes.HandleFunc("/{id}", user.Delete).Methods(http.MethodDelete)
	userRoutes.HandleFunc("/access_request/list", user.AccessRequestList).Methods(http.MethodGet)
	userRoutes.HandleFunc("/access_request/{id}/update", user.UpdateUserAccessStatus).Methods(http.MethodPut)

	webResourceUserRoutes := r.PathPrefix("/webresource").Subrouter()
	webResourceUserRoutes.Use(middleware.DataMiddleware)
	webResourceUserRoutes.HandleFunc("/list", webresource.List).Methods(http.MethodGet)
	webResourceUserRoutes.HandleFunc("/list/{type}", webresource.ListByType).Methods(http.MethodGet)

	webResourceAdminRoutes := r.PathPrefix("/webresource").Subrouter()
	webResourceAdminRoutes.Use(middleware.AdminMiddleware)
	webResourceAdminRoutes.HandleFunc("/create", webresource.Post).Methods(http.MethodPost)
	webResourceAdminRoutes.HandleFunc("/{id}", webresource.Get).Methods(http.MethodGet)
	webResourceAdminRoutes.HandleFunc("/{id}", webresource.Put).Methods(http.MethodPut)
	webResourceAdminRoutes.HandleFunc("/{id}", webresource.Delete).Methods(http.MethodDelete)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	return r
}
