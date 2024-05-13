package route

import (
	"app/internal/api/controller"
	"app/internal/api/core/usecase"
	"app/internal/api/impl/adapter"
	"app/internal/api/impl/repo"
	"app/internal/api/middleware"
	"app/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB) {
    // repo and adapter
    userRepo := repo.NewUserRepoImpl(db)
    sessionRepo := repo.NewSessionReoImpl(db)
    mediaRepo := repo.NewMediaRepoImpl(db)
    fileAdapter := adapter.NewFileAdapterImpl()

    // usecase
    signup := usecase.NewSignupUsecase(userRepo)
    getUser := usecase.NewGetUserUsecase(userRepo)
    editProfile := usecase.NewEditProfileUsecase(userRepo)
    auth := usecase.NewAuthUsecase(sessionRepo)
    signin := usecase.NewSigninUsecase(userRepo, sessionRepo)
    refreshToken := usecase.NewRefreshTokenUsecase(sessionRepo)
    uploadFile := usecase.NewUploadFileUsecase(mediaRepo, fileAdapter)
    getFile := usecase.NewGetFileUsecase(mediaRepo, fileAdapter)

    // controller and middleware
    authMiddleware := middleware.NewAuthMiddleware(auth)
    userController := controller.NewUsersController(signup, getUser, editProfile)
    sessionsController := controller.NewSessionsController(signin, refreshToken)
    mediaController := controller.NewMediaController(uploadFile, getFile)
    
    r.Route("/api", func(r chi.Router) {
        r.Route("/users", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                r.Use(authMiddleware)
                r.Patch("/{username}/settings/profiles", handler.Handle(userController.EditProfile)) 
            })

            r.Post("/", handler.Handle(userController.SignUp)) 
            r.Get("/lookup", handler.Handle(userController.LookUp)) 
        })

        r.Route("/sessions", func(r chi.Router) {
            r.Post("/signin", handler.Handle(sessionsController.Signin))
            r.Post("/token/refresh", handler.Handle(sessionsController.RefreshToken))
        })

        r.Route("/media", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                r.Use(authMiddleware)
                r.Post("/upload", handler.Handle(mediaController.UploadFile))
            })

            r.Get("/{fileName}", handler.Handle(mediaController.GetFile))
        })
    })
}
