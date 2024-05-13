package usecase

import (
	"app/internal/activitypub/core/entity"
	"app/internal/activitypub/core/repo"
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
    if username == "" {
        return nil, nil
    }

    user, err := uc.userRepo.FindUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, nil
    }

    return user, nil
}
