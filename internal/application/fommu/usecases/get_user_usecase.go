package usecases

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"context"
	"fmt"
)

type GetUserUsecase struct {
    userRepo repo.UsersRepo `injectable:""`
}

func NewGetUserUsecase() *GetUserUsecase {
    return &GetUserUsecase{}
}

func (uc *GetUserUsecase) Exec(ctx context.Context, username string, domain string) (*entities.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if username == "" {
        return nil, nil
    }

    log.Info(ctx, fmt.Sprintf("Get username=\"%s\", domain=\"%s\"", username, domain))
    user, err := uc.userRepo.FindUserByUsername(ctx, username, domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Somethig went wrong")
    }
    if user == nil {
        log.Info(ctx, "User not found.")
        return nil, nil
    }

    return user, nil
}
