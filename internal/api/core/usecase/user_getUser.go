package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/config"
	"context"
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
    if username == "" {
        return nil, nil
    }

    s := strings.SplitN(username, "@", 2)

    username = strings.TrimSpace(s[0])
    domain := ""
    if len(s) > 1 {
        domain = strings.TrimSpace(s[1])
    }

    var user *entity.UserEntity = nil
    if domain == "" {
        domain = config.Fommu.Domain
    }

    user, err := uc.userRepo.FindUserByUsername(ctx, username, domain)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, nil
    }

    return user, nil
}
