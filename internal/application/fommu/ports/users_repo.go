package ports

import (
	"app/internal/core/entities"
	"context"
)

type UsersRepo interface {
    FindUserByID(ctx context.Context, id string) (*entities.UserEntity, error)
    FindUserByUsername(ctx context.Context, username string, domain string) (*entities.UserEntity, error)
    FindUserByEmail(ctx context.Context, email string, domain string) (*entities.UserEntity, error)
    SearchUser(ctx context.Context, textSearch string, domain string) ([]*entities.UserEntity, error)
    CreateUser(ctx context.Context, user *entities.UserEntity) error
    UpdateUser(ctx context.Context, user *entities.UserEntity) error
}
