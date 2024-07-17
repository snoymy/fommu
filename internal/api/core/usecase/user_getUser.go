package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/log"
	"context"
	"fmt"
	"strings"
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
        return nil, nil
    }

    log.Info(ctx, "Split username and domain.")
    s := strings.SplitN(username, "@", 2)

    username = strings.TrimSpace(s[0])
    domain := config.Fommu.Domain
    if len(s) > 1 {
        log.Debug(ctx, "Has domain name.")
        domain = strings.TrimSpace(s[1])
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
