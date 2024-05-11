package api

import (
	"app/internal/controller"
	"app/internal/core/usecase"
	"app/internal/handler"
	"app/internal/repoimpl"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewSessionsRoute(db *sqlx.DB) func(chi.Router) {
    userRepo := repoimpl.NewUserRepoImpl(db)
    sessionRepo := repoimpl.NewSessionReoImpl(db)
    signin := usecase.NewSigninUsecase(userRepo, sessionRepo)
    sessionsController := controller.NewSessionsController(signin)
    return func(r chi.Router) {
        r.Post("/signin", handler.Handle(sessionsController.Signin))
    }
}
