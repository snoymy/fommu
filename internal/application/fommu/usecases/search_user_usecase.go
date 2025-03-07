package usecases

import (
	"app/internal/core/entities"
	"app/internal/application/fommu/repos"
	"app/internal/application/appstatus"
	"app/internal/log"
	"context"
	"fmt"
)

type SearchUserUsecase struct {
    userRepo repos.UsersRepo `injectable:""`
}

func NewSearchUserUsecase() *SearchUserUsecase {
    return &SearchUserUsecase{}
}

func (uc *SearchUserUsecase) Exec(ctx context.Context, username string, domain string) ([]*entities.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if username == "" {
        return nil, nil
    }

    log.Info(ctx, fmt.Sprintf("Searching for username=\"%s\", domain=\"%s\"", username, domain))
    users, err := uc.userRepo.SearchUser(ctx, username, domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Somethig went wrong")
    }
    if users == nil {
        log.Info(ctx, "Users not found.")
        return nil, nil
    }

    return users, nil
}
