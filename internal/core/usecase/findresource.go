package usecase

import (
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/core/entity"
	"app/internal/core/repo"
	"context"
	"strings"
)

type FindResourceUsecase struct {
    userRepo repo.UsersRepo
}

func NewFindResourceUsecase(userRepo repo.UsersRepo) *FindResourceUsecase {
    return &FindResourceUsecase{
        userRepo: userRepo,
    }
}

func (uc *FindResourceUsecase) Exec(ctx context.Context, resource string) (*entity.UserEntity, error) {
    // validate resource
    println(resource)
    if resource == "" {
        return nil, nil
    }

    ar := strings.SplitN(resource, ":", 2)
    if len(ar) < 2 {
        return nil, appstatus.NotFound()
    }

    //typ := ar[0]
    handle := ar[1]
    if handle[0] == '@' && len(handle) > 1 {
        handle = handle[1:]
    }

    // find resource
    user, err := uc.userRepo.FindResource(ctx, handle, config.Fommu.Domain)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, nil
    }

    return user, nil
}
