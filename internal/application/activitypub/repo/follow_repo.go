package repo

import (
	"app/internal/core/entity"
	"context"
)

type FollowRepo interface {
    CreateFollow(ctx context.Context, following *entity.FollowEntity) error
}
