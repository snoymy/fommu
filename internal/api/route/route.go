package route

import (
	"app/internal/api/controller"
	"app/internal/api/core/usecase"
	"app/internal/api/impl/repo"
	"app/internal/api/middleware"
	"app/internal/handler"
	"app/internal/httpclient"
	"app/internal/log"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB, apClient *httpclient.ActivitypubClient) {
    ctx := context.Background()

    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)
    // repo and adapter
    userRepo := repo.NewUserRepoImpl(db, apClient)
    sessionRepo := repo.NewSessionReoImpl(db)
    mediaRepo := repo.NewMediaRepoImpl(db)

    // usecase
    auth := usecase.NewAuthUsecase(sessionRepo)
    
    // users
    signup := usecase.NewSignupUsecase(userRepo)
    getUser := usecase.NewGetUserUsecase(userRepo)
    editProfile := usecase.NewEditProfileUsecase(userRepo)
    editAccount := usecase.NewEditAccountUsecase(userRepo)
    searchUser := usecase.NewSearchUserUsecase(userRepo)

    // sessions
    signin := usecase.NewSigninUsecase(userRepo, sessionRepo)
    signout := usecase.NewSignOutUsecase(sessionRepo)
    refreshToken := usecase.NewRefreshTokenUsecase(sessionRepo)
    revokeSession := usecase.NewRevokeSessionUsecase(sessionRepo)
    verifySession := usecase.NewGetIdentityUsecase(sessionRepo, userRepo)

    // media
    uploadFile := usecase.NewUploadFileUsecase(mediaRepo)
    getFile := usecase.NewGetFileUsecase(mediaRepo)
    getToken := usecase.NewGetTokenUsecase(sessionRepo)

    // controller and middleware
    requestIdMiddleware := middleware.NewRequestIDMiddleware()
    authMiddleware := middleware.NewAuthMiddleware(auth)
    userController := controller.NewUsersController(signup, getUser, editProfile, editAccount, searchUser)
    sessionsController := controller.NewSessionsController(signin, signout, refreshToken, getToken, revokeSession, verifySession)
    mediaController := controller.NewMediaController(uploadFile, getFile)
    
    log.Info(ctx, "Init /api routes...")
    r.Route("/api", func(r chi.Router) {
        r.Use(requestIdMiddleware)
        r.Route("/users", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                r.Use(authMiddleware)
                r.Patch("/{username}/settings/profiles", handler.Handle(userController.EditProfile)) 
                r.Patch("/{username}/settings/account", handler.Handle(userController.EditAccount)) 
            })

            r.Get("/{username}", handler.Handle(userController.GetUser)) 
            r.Post("/", handler.Handle(userController.SignUp)) 
            r.Get("/lookup", handler.Handle(userController.LookUp)) 
            r.Get("/search", handler.Handle(userController.Search)) 
        })

        r.Route("/sessions", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                r.Use(authMiddleware)
                r.Delete("/revoke/{sessionId}", handler.Handle(sessionsController.RevokeSession))
                r.Delete("/signout", handler.Handle(sessionsController.SignOut))
                r.Get("/identity", handler.Handle(sessionsController.VerifySession))
            })

            r.Post("/signin", handler.Handle(sessionsController.Signin))
            r.Post("/token/refresh", handler.Handle(sessionsController.RefreshToken))
            r.Get("/token", handler.Handle(sessionsController.GetToken))
        })

        r.Route("/media", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                r.Use(authMiddleware)
                r.Post("/upload", handler.Handle(mediaController.UploadFile))
            })

            r.Get("/{fileName}", handler.Handle(mediaController.GetFile))
        })
    })
    log.Info(ctx, "Init /api success")
}
