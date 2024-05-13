package repo

import (
	"app/internal/activitypub/core/entity"
	"context"
)

type UsersRepo interface {
    FindUserByUsername(ctx context.Context, username string) (*entity.UserEntity, error)
    FindResource(ctx context.Context, resource string, domain string) (*entity.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
}
