package route

import (
	"app/internal/activitypub/controller"
	"app/internal/activitypub/core/usecase"
	"app/internal/activitypub/impl/repo"
	"app/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB) {
    // repo and adapter
    userRepo := repo.NewUserRepoImpl(db)

    // usecase
    getUser := usecase.NewGetUserUsecase(userRepo)
    findresource := usecase.NewFindResourceUsecase(userRepo)

    // controller and middleware
    wellknown := controller.NewWellKnownController(findresource) 
    userController := controller.NewAPUsersController(getUser) 

    r.Route("/", func(r chi.Router) {
        r.Get("/{username}", handler.Handle(userController.GetUser)) 
        r.Route("/.well-known", func(r chi.Router) {
            r.Get("/webfinger", handler.Handle(wellknown.WebFinger)) 
        })
    })
}
