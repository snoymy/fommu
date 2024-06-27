package route

import (
	"app/internal/activitypub/controller"
	"app/internal/activitypub/core/usecase"
	"app/internal/activitypub/impl/repo"
	"app/internal/activitypub/middleware"
	"app/internal/handler"
	"app/internal/httpclient"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB, apClient *httpclient.ActivitypubClient) {
    // repo and adapter
    userRepo := repo.NewUserRepoImpl(db, apClient)
    followingRepo := repo.NewFollowingRepoImpl(db)

    // usecase
    verifySignature := usecase.NewVerifySignatureUsecase(userRepo)
    getUser := usecase.NewGetUserUsecase(userRepo)
    findresource := usecase.NewFindResourceUsecase(userRepo)
    followUser := usecase.NewFollowUserUsecase(userRepo, followingRepo)

    // controller and middleware
    verifySignatureMiddleware := middleware.NewVerifyMiddleware(verifySignature)
    wellknown := controller.NewWellKnownController(findresource) 
    userController := controller.NewAPUsersController(getUser, followUser)

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
