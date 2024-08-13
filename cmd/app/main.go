package main

import (
	aproute "app/internal/adapter/activitypub/route"
	apiroute "app/internal/adapter/fommu/route"
	"app/internal/config"
	"app/internal/infrastructure/database"
	"app/internal/infrastructure/httpclient"
	"app/internal/infrastructure/router"
	"app/internal/log"
	"context"
	"net/http"
	"strconv"
)

func main() {
    config.Init()

    ctx := context.Background()

    log.Info(ctx, "Init database connection...")
    db := database.NewConnection()
    defer db.Close()

    apClient := httpclient.NewActivitypubClientImpl()

    if err := database.TestConnection(db); err != nil {
        log.Panic(ctx, err.Error())
    }
    log.Info(ctx, "Init database succeed")

    r := router.NewRouter()

    apiroute.InitRoute(r, db, apClient)
    aproute.InitRoute(r, db, apClient)

    log.Info(ctx, "Listening port: " + strconv.Itoa(config.Fommu.Port))
    http.ListenAndServe(":" + strconv.Itoa(config.Fommu.Port), r)
}
