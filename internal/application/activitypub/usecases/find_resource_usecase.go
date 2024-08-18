package usecases

import (
	"app/internal/application/activitypub/repos"
	"app/internal/config"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"context"
	"strings"
)

type FindResourceUsecase struct {
    userRepo repos.UsersRepo `injectable:""`
}

func NewFindResourceUsecase() *FindResourceUsecase {
    return &FindResourceUsecase{}
}

func (uc *FindResourceUsecase) Exec(ctx context.Context, resource string) (*entities.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    // validate resource
    log.Info(ctx, "Validate resource")
    if resource == "" {
        log.Info(ctx, "Resource value is empty")
        return nil, nil
    }

    ar := strings.SplitN(resource, ":", 2)
    if len(ar) < 2 {
        log.Info(ctx, "No resource name provided")
        return nil, appstatus.NotFound()
    }

    //typ := ar[0]
    handle := ar[1]
    if handle[0] == '@' && len(handle) > 1 {
        handle = handle[1:]
    }

    // find resource
    log.Info(ctx, "Find user by resource")
    user, err := uc.userRepo.FindResource(ctx, handle, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: ", err.Error())
        return nil, err
    }

    if user == nil {
        log.Info(ctx, "User not found")
        return nil, nil
    }

    return user, nil
}
