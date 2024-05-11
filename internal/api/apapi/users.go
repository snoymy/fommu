package apapi

import (
	"app/internal/controller/apcontroller"
	"app/internal/core/usecase"
	"app/internal/handler"
	"app/internal/repoimpl"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewUsersRoute(db *sqlx.DB) func(chi.Router) {
    userRepo := repoimpl.NewUserRepoImpl(db)

    getUser := usecase.NewGetUserUsecase(userRepo)

    userController := apcontroller.NewAPUsersController(getUser) 

    return func(r chi.Router) {
        r.Get("/{username}", handler.Handle(userController.GetUser)) 
    }
}

