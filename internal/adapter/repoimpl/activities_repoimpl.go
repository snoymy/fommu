package repoimpl

import (
	"app/internal/adapter/command"
	"app/internal/adapter/query"
	"app/internal/core/entity"
	"context"
)

type ActivitiesRepoImpl struct {
    query *query.Query `injectable:""`
    command *command.Command `injectable:""`
}

func NewActActivitiesRepoImpl() *ActivitiesRepoImpl {
    return &ActivitiesRepoImpl{}
}

func (r *ActivitiesRepoImpl) FindActivityByActivityId(ctx context.Context, activityId string) (*entity.ActivityEntity, error) {
    activity, err := r.query.FindActivityByActivityId(ctx, activityId)
    if err != nil {
        return nil, err
    }

    return activity, nil
}

func (r *ActivitiesRepoImpl) FindActivityById(ctx context.Context, activityId string) (*entity.ActivityEntity, error) {
    activity, err := r.query.FindActivityById(ctx, activityId)
    if err != nil {
        return nil, err
    }

    return activity, nil
}

func (r *ActivitiesRepoImpl) CreateActivity(ctx context.Context, activity *entity.ActivityEntity) error {
    err := r.command.CreateActivity(ctx, activity)
    if err != nil {
        return err
    }

    r.command.NotifyProcessActivity(ctx, activity)

    return nil
}

func (r *ActivitiesRepoImpl) UpdateActivity(ctx context.Context, activity *entity.ActivityEntity) error {
    err := r.command.UpdateActivity(ctx, activity)
    if err != nil {
        return err
    }

    return nil
}
