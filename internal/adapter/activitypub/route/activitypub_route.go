package route

import (
	"app/internal/adapter/activitypub/controller"
	"app/internal/adapter/activitypub/listener"
	"app/internal/adapter/activitypub/middleware"
	"app/internal/adapter/command"
	"app/internal/adapter/handler"
	"app/internal/adapter/httpclient"
	"app/internal/adapter/query"
	"app/internal/adapter/repoimpl"
	"app/internal/application/activitypub/repo"
	"app/internal/application/activitypub/usecase"
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
    container.Register(func() repo.UsersRepo { return repoimpl.NewUserRepoImpl() })
    container.Register(func() repo.FollowRepo { return repoimpl.NewFollowRepoImpl() })
    container.Register(func() repo.ActivitiesRepo { return repoimpl.NewActActivitiesRepoImpl() })

    // usecase
    container.Register(usecase.NewVerifySignatureUsecase)
    container.Register(usecase.NewGetUserUsecase)
    container.Register(usecase.NewFindResourceUsecase)
    container.Register(usecase.NewProcessFollowActivityUsecase)
    container.Register(usecase.NewCreateActivityUsecase)

    // controller and middleware
    container.Register(middleware.NewVerifyMiddleware)
    container.Register(controller.NewWellKnownController)
    container.Register(controller.NewAPUsersController)

    container.Register(listener.NewProcessActivityListener)

    log.Info(ctx, "Init / routes...")
    container.Resolve(resolveRoute)
    log.Info(ctx, "Init / success")

    container.Resolve(registerEvent)
}

func registerEvent(bus EventBus.Bus, processActivityListener *listener.ProcessActivityListener) {
    bus.SubscribeAsync("topic:process_activity", processActivityListener.Handler, false)
}

func resolveRoute(r chi.Router, verifySignatureMiddleware middleware.VerifyMiddleware, userController *controller.APUsersController, wellknown *controller.WellKnown) {
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
