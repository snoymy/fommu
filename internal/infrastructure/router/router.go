package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter() chi.Router {
    r := chi.NewRouter()
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:4000", "https://fommu.loca.lt"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, // Maximum value not ignored by any of major browsers
    }))

    ErrorHandler = HandleError

    return r
}
