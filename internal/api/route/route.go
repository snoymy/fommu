package route

import (
	"app/internal/api/controller"
	"app/internal/api/core/usecase"
	"app/internal/api/impl/repo"
	"app/internal/api/middleware"
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
    container.Register(repo.NewSessionRepoImpl)
    container.Register(repo.NewMediaRepoImpl)

    // usecase
    container.Register(usecase.NewAuthUsecase)
    
    // users
    container.Register(usecase.NewSignupUsecase)
    container.Register(usecase.NewGetUserUsecase)
    container.Register(usecase.NewEditProfileUsecase)
    container.Register(usecase.NewEditAccountUsecase)
    container.Register(usecase.NewSearchUserUsecase)

    // sessions
    container.Register(usecase.NewSigninUsecase)
    container.Register(usecase.NewSignOutUsecase)
    container.Register(usecase.NewRefreshTokenUsecase)
    container.Register(usecase.NewRevokeSessionUsecase)
    container.Register(usecase.NewGetIdentityUsecase)

    // media
    container.Register(usecase.NewUploadFileUsecase)
    container.Register(usecase.NewGetFileUsecase)
    container.Register(usecase.NewGetTokenUsecase)

    // controller and middleware
    container.Register(middleware.NewRequestIDMiddleware)
    container.Register(middleware.NewAuthMiddleware)
    container.Register(controller.NewUsersController)
    container.Register(controller.NewSessionsController)
    container.Register(controller.NewMediaController)

    log.Info(ctx, "Init /api routes...")
    container.Resolve(resolveRoute)
    log.Info(ctx, "Init /api success")
}

func resolveRoute(
    r chi.Router,
    requestIdMiddleware middleware.RequestIdMiddleware, 
    authMiddleware middleware.AuthMiddleware,
    userController *controller.UsersController,
    sessionsController *controller.SessionsController,
    mediaController *controller.MediaController,
) {
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
}
