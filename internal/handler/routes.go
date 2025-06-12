package http

import (
	"github.com/evgeniySeleznev/person-enrichment-service/internal/service"
	"github.com/evgeniySeleznev/person-enrichment-service/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(service *service.PersonService, logger logger.Logger) *mux.Router {
	logger.Debug("ENTER: NewRouter")
	router := mux.NewRouter()
	handler := NewPersonHandler(service, logger)

	// Middleware для логирования запросов
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"Request received",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			next.ServeHTTP(w, r)
		})
	})

	// Маршруты API
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/persons", handler.CreatePerson).Methods("POST")
	api.HandleFunc("/persons/{id}", handler.GetPerson).Methods("GET")
	api.HandleFunc("/persons", handler.GetAllPersons).Methods("GET")
	api.HandleFunc("/persons/{id}", handler.UpdatePerson).Methods("PATCH")
	api.HandleFunc("/persons/{id}", handler.DeletePerson).Methods("DELETE")
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
