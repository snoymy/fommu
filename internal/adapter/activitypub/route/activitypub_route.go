package route

import (
	"app/internal/adapter/activitypub/controllers"
	"app/internal/adapter/activitypub/listeners"
	"app/internal/adapter/middlewares"
	"app/internal/adapter/commands"
	"app/internal/adapter/queries"
	"app/internal/adapter/repoimpl"
	"app/internal/application/activitypub/ports"
	"app/internal/application/activitypub/usecases"
	"app/internal/infrastructure/httpclient"
	"app/internal/infrastructure/router"
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
    container.Register(queries.NewQuery)
    container.Register(commands.NewCommand)

    // repo and adapter
    container.Register(func() ports.UsersRepo { return repoimpl.NewUserRepoImpl() })
    container.Register(func() ports.FollowRepo { return repoimpl.NewFollowRepoImpl() })
    container.Register(func() ports.ActivitiesRepo { return repoimpl.NewActActivitiesRepoImpl() })

    // usecase
    container.Register(usecases.NewVerifySignatureUsecase)
    container.Register(usecases.NewGetUserUsecase)
    container.Register(usecases.NewFindResourceUsecase)
    container.Register(usecases.NewProcessFollowActivityUsecase)
    container.Register(usecases.NewCreateActivityUsecase)

    // controller and middleware
    container.Register(middlewares.NewRequestIDMiddleware)
    container.Register(middlewares.NewVerifyMiddleware)
    container.Register(controllers.NewWellKnownController)
    container.Register(controllers.NewAPUsersController)

    container.Register(listeners.NewProcessActivityListener)

    log.Info(ctx, "Init / routes...")
    container.Resolve(resolveRoute)
    log.Info(ctx, "Init / success")

    container.Resolve(registerEvent)
}

func registerEvent(bus EventBus.Bus, processActivityListener *listeners.ProcessActivityListener) {
    bus.SubscribeAsync("topic:process_activity", processActivityListener.Handler, false)
}

func resolveRoute(
    r chi.Router, 
    requestIdMiddleware middlewares.RequestIdMiddleware, 
    verifySignatureMiddleware middlewares.VerifyMiddleware, 
    userController *controllers.APUsersController, 
    wellknown *controllers.WellKnown,
) {
    r.With().Route("/", func(r chi.Router) {
        r.Use(requestIdMiddleware)
        r.Get("/users/{username}", router.Handle(userController.GetUser)) 

        r.Group(func(r chi.Router) {
            r.Use(verifySignatureMiddleware)
            r.Post("/users/{username}/inbox", router.Handle(userController.Inbox)) 
        })

        r.Route("/.well-known", func(r chi.Router) {
            r.Get("/webfinger", router.Handle(wellknown.WebFinger)) 
        })
    })
}
