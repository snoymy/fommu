package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"context"
	"fmt"
	"strings"
)

type SearchUserUsecase struct {
    userRepo repo.UsersRepo `injectable:""`
}

func NewSearchUserUsecase() *SearchUserUsecase {
    return &SearchUserUsecase{}
}

func (uc *SearchUserUsecase) Exec(ctx context.Context, username string) ([]*entity.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if username == "" {
        return nil, nil
    }

    log.Info(ctx, "Spliting username and domain.")
    s := strings.SplitN(username, "@", 2)

    username = strings.TrimSpace(s[0])
    domain := ""
    if len(s) > 1 {
        log.Debug(ctx, "Has domain name.")
        domain = strings.TrimSpace(s[1])
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
