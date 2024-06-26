package main

import (
	aproute "app/internal/activitypub/route"
	apiroute "app/internal/api/route"
	"app/internal/config"
	"app/internal/config/database"
	"app/internal/handler"
	"app/internal/httpclient"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
    config.Init()
    db := database.NewConnection()
    defer db.Close()
    apClient := httpclient.NewActivitypubClient()

    if err := database.TestConnection(db); err != nil {
        panic(err.Error())
    }

    handler.ErrorHandler = handler.HandleError

    r := chi.NewRouter()

    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:4000", "https://fommu.loca.lt"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, // Maximum value not ignored by any of major browsers
    }))

    apiroute.InitRoute(r, db, apClient)
    aproute.InitRoute(r, db, apClient)

    http.ListenAndServe(":" + strconv.Itoa(config.Fommu.Port), r)
}
