package repo

import (
	"app/internal/core/entity"
	"context"
)

type UsersRepo interface {
    FindUserByUsername(ctx context.Context, username string, domain string) (*entity.UserEntity, error)
    FindUserByActorId(ctx context.Context, actorId string) (*entity.UserEntity, error)
    FindResource(ctx context.Context, resource string, domain string) (*entity.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string, domain string) (*entity.UserEntity, error)
}
