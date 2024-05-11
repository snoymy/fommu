package api

import (
	"app/internal/controller"
	"app/internal/core/usecase"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/repoimpl"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewUsersRoute(db *sqlx.DB) func(chi.Router) {
    userRepo := repoimpl.NewUserRepoImpl(db)
    sessionRepo := repoimpl.NewSessionReoImpl(db)

    signup := usecase.NewSignupUsecase(userRepo)
    getUser := usecase.NewGetUserUsecase(userRepo)
    editProfile := usecase.NewEditProfileUsecase(userRepo)
    auth := usecase.NewAuthUsecase(sessionRepo)

    userController := controller.NewUsersController(signup, getUser, editProfile)
    authMiddleware := middleware.NewAuthMiddleware(auth)

    return func(r chi.Router) {
        r.Post("/", handler.Handle(userController.SignUp)) 
        r.Get("/lookup", handler.Handle(userController.LookUp)) 

        r.Group(func(r chi.Router) {
            r.Use(authMiddleware)
            r.Patch("/{username}/settings/profiles", handler.Handle(userController.EditProfile)) 
        })
    }
}
