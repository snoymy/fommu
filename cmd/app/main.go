package main

import (
	aproute "app/internal/adapter/activitypub/route"
	apiroute "app/internal/adapter/fommu/route"
	"app/internal/config"
	"app/internal/adapter/database"
	"app/internal/adapter/handler"
	"app/internal/adapter/httpclient"
	"app/internal/log"
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
    config.Init()

    ctx := context.Background()

    log.Info(ctx, "Init database connection...")
    db := database.NewConnection()
    defer db.Close()

    apClient := httpclient.NewActivitypubClient()

    if err := database.TestConnection(db); err != nil {
        log.Panic(ctx, err.Error())
    }
    log.Info(ctx, "Init database succeed")

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

    log.Info(ctx, "Listening port: " + strconv.Itoa(config.Fommu.Port))
    http.ListenAndServe(":" + strconv.Itoa(config.Fommu.Port), r)
}
