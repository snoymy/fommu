package repo

import (
	"app/internal/activitypub/core/entity"
	"context"
)

type FollowingRepo interface {
    CreateFollowing(ctx context.Context, following *entity.FollowingEntity) error
}
