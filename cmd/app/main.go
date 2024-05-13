package main

import (
	"app/internal/api"
	"app/internal/api/apapi"
	"app/internal/config"
	"app/internal/config/database"
	"app/internal/handler"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
    config.Init()
    db := database.NewConnection()
	defer db.Close()

    if err := database.TestConnection(db); err != nil {
        panic(err.Error())
    }

    handler.ErrorHandler = handler.HandleError

    r := chi.NewRouter()

    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        //ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, // Maximum value not ignored by any of major browsers
    }))

    r.Route("/", func(r chi.Router) {
        wellKnownRoute := apapi.NewWellKnownRoute(db)
        usersRoute := apapi.NewUsersRoute(db)

        r.Route("/.well-known", wellKnownRoute)
        r.Route("/users", usersRoute)
    })

    r.Route("/api", func(r chi.Router) {
        usersRoute := api.NewUsersRoute(db)
        sessionsRoute := api.NewSessionsRoute(db)
        mediaRoute := api.NewMediaRoute(db)

        r.Route("/users", usersRoute)
        r.Route("/sessions", sessionsRoute)
        r.Route("/media", mediaRoute)
    })

    http.ListenAndServe(":" + strconv.Itoa(config.Fommu.Port), r)
}
