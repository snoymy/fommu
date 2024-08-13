package usecases

import (
	"app/internal/core/entities"
	"app/internal/application/activitypub/ports"
	"app/internal/application/appstatus"
	"app/internal/config"
	"app/internal/log"
	"context"
)

type GetUserUsecase struct {
    userRepo ports.UsersRepo `injectable:""`
}

func NewGetUserUsecase() *GetUserUsecase {
    return &GetUserUsecase{}
}

func (uc *GetUserUsecase) Exec(ctx context.Context, username string) (*entities.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if username == "" {
        log.Info(ctx, "User not found")
        return nil, nil
    }

    log.Info(ctx, "Find user")
    user, err := uc.userRepo.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Unable to find user")
    }

    if user == nil {
        log.Info(ctx, "User not found")
        return nil, nil
    }

    return user, nil
}
