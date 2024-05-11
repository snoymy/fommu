package apapi

import (
	"app/internal/controller"
	"app/internal/core/usecase"
	"app/internal/handler"
	"app/internal/repoimpl"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewWellKnownRoute(db *sqlx.DB) func(chi.Router) {
    userRepo := repoimpl.NewUserRepoImpl(db)
    findresource := usecase.NewFindResourceUsecase(userRepo)
    wellknown := controller.NewWellKnownController(findresource) 

    return func(r chi.Router) {
       r.Get("/webfinger", handler.Handle(wellknown.WebFinger)) 
    }
}
