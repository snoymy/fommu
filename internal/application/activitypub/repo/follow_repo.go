package repo

import (
	"app/internal/core/entity"
	"context"
)

type FollowRepo interface {
    CreateFollowing(ctx context.Context, following *entity.FollowEntity) error
}
