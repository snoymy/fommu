package ports

import (
	"app/internal/core/entities"
	"context"
)

type FollowRepo interface {
    CreateFollow(ctx context.Context, following *entities.FollowEntity) error
}
