package usecase

import (
	"app/internal/activitypub/core/entity"
	"app/internal/activitypub/core/repo"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/log"
	"context"
)

type GetUserUsecase struct {
    userRepo repo.UsersRepo
}

func NewGetUserUsecase(userRepo repo.UsersRepo) *GetUserUsecase {
    return &GetUserUsecase{
        userRepo: userRepo,
    }
}

func (uc *GetUserUsecase) Exec(ctx context.Context, username string) (*entity.UserEntity, error) {
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
