package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"context"
	"strings"
)

type SearchUserUsecase struct {
    userRepo repo.UsersRepo
}

func NewSearchUserUsecase(userRepo repo.UsersRepo) *SearchUserUsecase {
    return &SearchUserUsecase{
        userRepo: userRepo,
    }
}

func (uc *SearchUserUsecase) Exec(ctx context.Context, username string) ([]*entity.UserEntity, error) {
    if username == "" {
        return nil, nil
    }

    s := strings.SplitN(username, "@", 2)

    username = strings.TrimSpace(s[0])
    domain := ""
    if len(s) > 1 {
        domain = strings.TrimSpace(s[1])
    }

    users, err := uc.userRepo.SearchUser(ctx, username, domain)
    if err != nil {
        return nil, err
    }
    if users == nil {
        return nil, nil
    }

    return users, nil
}
