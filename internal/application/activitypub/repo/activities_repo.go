package repo

import (
	"app/internal/core/entity"
	"context"
)

type ActivitiesRepo interface {
    FindActivityByActivityId(ctx context.Context, activityId string) (*entity.ActivityEntity, error)
    FindActivityById(ctx context.Context, activityId string) (*entity.ActivityEntity, error)
    CreateActivity(ctx context.Context, activity *entity.ActivityEntity) error
    UpdateActivity(ctx context.Context, activity *entity.ActivityEntity) error
}
