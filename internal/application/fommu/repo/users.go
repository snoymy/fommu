package repo

import (
	"app/internal/core/entity"
	"context"
)

type UsersRepo interface {
    FindUserByID(ctx context.Context, id string) (*entity.UserEntity, error)
    FindUserByUsername(ctx context.Context, username string, domain string) (*entity.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string, domain string) (*entity.UserEntity, error)
    SearchUser(ctx context.Context, textSearch string, domain string) ([]*entity.UserEntity, error)
    CreateUser(ctx context.Context, user *entity.UserEntity) error
    UpdateUser(ctx context.Context, user *entity.UserEntity) error
}
