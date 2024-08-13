package repo

import (
	"app/internal/core/entities"
	"context"
)

type ActivitiesRepo interface {
    FindActivityByActivityId(ctx context.Context, activityId string) (*entities.ActivityEntity, error)
    FindActivityById(ctx context.Context, activityId string) (*entities.ActivityEntity, error)
    CreateActivity(ctx context.Context, activity *entities.ActivityEntity) error
    UpdateActivity(ctx context.Context, activity *entities.ActivityEntity) error
}
