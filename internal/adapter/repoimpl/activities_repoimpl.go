package repoimpl

import (
	"app/internal/adapter/commands"
	"app/internal/adapter/queries"
	"app/internal/core/entities"
	"context"
)

type ActivitiesRepoImpl struct {
    queries *queries.Query `injectable:""`
    commands *commands.Command `injectable:""`
}

func NewActActivitiesRepoImpl() *ActivitiesRepoImpl {
    return &ActivitiesRepoImpl{}
}

func (r *ActivitiesRepoImpl) FindActivityByActivityId(ctx context.Context, activityId string) (*entities.ActivityEntity, error) {
    activity, err := r.queries.FindActivityByActivityId(ctx, activityId)
    if err != nil {
        return nil, err
    }

    return activity, nil
}

func (r *ActivitiesRepoImpl) FindActivityById(ctx context.Context, activityId string) (*entities.ActivityEntity, error) {
    activity, err := r.queries.FindActivityById(ctx, activityId)
    if err != nil {
        return nil, err
    }

    return activity, nil
}

func (r *ActivitiesRepoImpl) CreateActivity(ctx context.Context, activity *entities.ActivityEntity) error {
    err := r.commands.CreateActivity(ctx, activity)
    if err != nil {
        return err
    }

    r.commands.NotifyProcessActivity(ctx, activity)

    return nil
}

func (r *ActivitiesRepoImpl) UpdateActivity(ctx context.Context, activity *entities.ActivityEntity) error {
    err := r.commands.UpdateActivity(ctx, activity)
    if err != nil {
        return err
    }

    return nil
}
