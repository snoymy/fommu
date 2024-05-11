package repo

import (
	"app/internal/core/entity"
	"context"
)

type UsersRepo interface {
    FindUserByUsername(ctx context.Context, username string) (*entity.UserEntity, error)
    FindResource(ctx context.Context, resource string, domain string) (*entity.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
    CreateUser(ctx context.Context, user *entity.UserEntity) error
    UpdateUser(ctx context.Context, user *entity.UserEntity) error
}
