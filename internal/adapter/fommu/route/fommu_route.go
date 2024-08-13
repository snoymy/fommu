package route

import (
	"app/internal/adapter/command"
	"app/internal/adapter/fommu/controllers"
	"app/internal/adapter/fommu/middlewares"
	"app/internal/adapter/handler"
	"app/internal/adapter/httpclient"
	"app/internal/adapter/query"
	"app/internal/adapter/repoimpl"
	"app/internal/application/fommu/ports"
	"app/internal/application/fommu/usecases"
	"app/internal/log"
	"app/lib/di"
	"context"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoute(r chi.Router, db *sqlx.DB, apClient httpclient.ActivitypubClient) {
    ctx := context.Background()

    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    bus := EventBus.New()
    container := structdi.New()

    container.Register(func() chi.Router { return r })
    container.Register(func() *sqlx.DB { return db })
    container.Register(func() httpclient.ActivitypubClient { return apClient })
    container.Register(func() EventBus.Bus { return bus })
    container.Register(query.NewQuery)
    container.Register(command.NewCommand)

    // repo and adapter
    container.Register(func() ports.UsersRepo { return repoimpl.NewUserRepoImpl() })
    container.Register(func() ports.SessionsRepo { return repoimpl.NewSessionRepoImpl() })
    container.Register(func() ports.MediaRepo { return repoimpl.NewMediaRepoImpl() })

    // usecase
    container.Register(usecases.NewAuthUsecase)
    
    // users
    container.Register(usecases.NewSignupUsecase)
    container.Register(usecases.NewGetUserUsecase)
    container.Register(usecases.NewEditProfileUsecase)
    container.Register(usecases.NewEditAccountUsecase)
    container.Register(usecases.NewSearchUserUsecase)

    // sessions
    container.Register(usecases.NewSigninUsecase)
    container.Register(usecases.NewSignOutUsecase)
    container.Register(usecases.NewRefreshTokenUsecase)
    container.Register(usecases.NewRevokeSessionUsecase)
    container.Register(usecases.NewGetIdentityUsecase)

    // media
    container.Register(usecases.NewUploadFileUsecase)
    container.Register(usecases.NewGetFileUsecase)
    container.Register(usecases.NewGetTokenUsecase)

    // controller and middleware
    container.Register(middlewares.NewRequestIDMiddleware)
    container.Register(middlewares.NewAuthMiddleware)
    container.Register(controllers.NewUsersController)
    container.Register(controllers.NewSessionsController)
    container.Register(controllers.NewMediaController)

    log.Info(ctx, "Init /api routes...")
    container.Resolve(resolveRoute)
    log.Info(ctx, "Init /api success")
}

func resolveRoute(
    r chi.Router,
    requestIdMiddleware middlewares.RequestIdMiddleware, 
    authMiddleware middlewares.AuthMiddleware,
    userController *controllers.UsersController,
    sessionsController *controllers.SessionsController,
    mediaController *controllers.MediaController,
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
