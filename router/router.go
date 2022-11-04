package router

import (
	"github.com/gorilla/mux"
	"go-img-svc/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/images", middleware.GetAllMetadata).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/images/{id}", middleware.GetMetadata).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/images/{id}/data", middleware.GetImage).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/images", middleware.AddImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/images/{id}", middleware.UpdateImage).Methods("PUT", "OPTIONS")

	return router
}
