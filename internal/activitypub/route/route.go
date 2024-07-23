package route

import (
	"app/internal/activitypub/controller"
	"app/internal/activitypub/core/usecase"
	"app/internal/activitypub/impl/repo"
	"app/internal/activitypub/middleware"
	"app/internal/handler"
	"app/internal/httpclient"
	"app/internal/log"
	"app/lib/di/structdi"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB, apClient *httpclient.ActivitypubClient) {
    ctx := context.Background()

    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    container := structdi.New()

    container.Register(func() chi.Router { return r })
    container.Register(func() *sqlx.DB { return db })
    container.Register(func() *httpclient.ActivitypubClient { return apClient })

    // repo and adapter
    container.Register(repo.NewUserRepoImpl)
    container.Register(repo.NewFollowingRepoImpl)

    // usecase
    container.Register(usecase.NewVerifySignatureUsecase)
    container.Register(usecase.NewGetUserUsecase)
    container.Register(usecase.NewFindResourceUsecase)
    container.Register(usecase.NewFollowUserUsecase)

    // controller and middleware
    container.Register(middleware.NewVerifyMiddleware)
    container.Register(controller.NewWellKnownController)
    container.Register(controller.NewAPUsersController)

    log.Info(ctx, "Init / routes...")
    container.Resolve(resolveRoute)
    log.Info(ctx, "Init / success")
}

func resolveRoute(r chi.Router, verifySignatureMiddleware middleware.VerifyMiddleware, userController *controller.APUsersController, wellknown *controller.WellKnown) {
    r.Route("/", func(r chi.Router) {
        r.Get("/users/{username}", handler.Handle(userController.GetUser)) 

        r.Group(func(r chi.Router) {
            r.Use(verifySignatureMiddleware)
            r.Post("/users/{username}/inbox", handler.Handle(userController.Inbox)) 
        })

        r.Route("/.well-known", func(r chi.Router) {
            r.Get("/webfinger", handler.Handle(wellknown.WebFinger)) 
        })
    })
}
