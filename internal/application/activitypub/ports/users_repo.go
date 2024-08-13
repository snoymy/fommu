package ports

import (
	"app/internal/core/entities"
	"context"
)

type UsersRepo interface {
    FindUserByUsername(ctx context.Context, username string, domain string) (*entities.UserEntity, error)
    FindUserByActorId(ctx context.Context, actorId string) (*entities.UserEntity, error)
    FindResource(ctx context.Context, resource string, domain string) (*entities.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string, domain string) (*entities.UserEntity, error)
}
