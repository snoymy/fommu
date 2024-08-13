package usecases

import (
	"app/internal/adapter/mapper"
	"app/internal/application/activitypub/ports"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/snoymy/activitypub"
)

type CreateActivityUsecase struct {
    activitiesRepo ports.ActivitiesRepo `injectable:""`
}

func NewCreateActivityUsecase() *CreateActivityUsecase {
    return &CreateActivityUsecase{}
}

func (uc *CreateActivityUsecase) Exec(ctx context.Context, activity *activitypub.Activity) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Validate activity")
    if err := uc.validateActivity(ctx, activity); err != nil {
        return err
    }

    log.Info(ctx, "Check if activity exist")
    exist, err := uc.isActivityExist(ctx, activity)
    if err != nil {
        return err
    }

    if exist {
        log.Info(ctx, "Activity exist")
        return nil
    }

    log.Info(ctx, "Create activity entity")
    activityEntity, err := uc.createActivity(ctx, activity)
    if err != nil {
        return err
    }

    log.Info(ctx, "Save activity")
    if err := uc.activitiesRepo.CreateActivity(ctx, activityEntity); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }

    return nil
}
func (uc *CreateActivityUsecase) isActivityExist(ctx context.Context, activity *activitypub.Activity) (bool, error) {
    activityEntity, err := uc.activitiesRepo.FindActivityByActivityId(ctx, activity.ID.String()) 
    if err != nil {
        return false, appstatus.InternalServerError("Something went wrong")
    }

    if activityEntity != nil {
        return true, nil
    }

    return false, nil
}

func (uc *CreateActivityUsecase) validateActivity(ctx context.Context, activity *activitypub.Activity) error {
    supportedType := []activitypub.ActivityVocabularyType{
        activitypub.FollowType,
    }

    if !slices.Contains(supportedType, activity.Type) {
        log.Warn(ctx, "Activity not supported")
        return appstatus.NotSupport("Service not support activity type")
    }

    return nil
}

func (uc *CreateActivityUsecase) createActivity(ctx context.Context, activity *activitypub.Activity) (*entities.ActivityEntity, error) {
    activityMap, err := mapper.StructToMap(activity)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }
    activityEntity := entities.NewActiActivityEntity()
    activityEntity.ID = uuid.New().String()
    activityEntity.Type.Set(string(activity.Type))
    activityEntity.Remote = true
    activityEntity.Activity = activityMap
    activityEntity.Status = "pending"
    activityEntity.CreateAt = time.Now().UTC()

    return activityEntity, nil
}
