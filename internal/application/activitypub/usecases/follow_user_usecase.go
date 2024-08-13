package usecases

import (
	"app/internal/application/activitypub/ports"
	"app/internal/application/appstatus"
	"app/internal/config"
	"app/internal/core/entities"
	"app/internal/log"
	"app/internal/utils/structutil"
	"context"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/snoymy/activitypub"
)

type ProcessFollowActivityUsecase struct {
    userRepo        ports.UsersRepo          `injectable:""`
    followingRepo   ports.FollowRepo         `injectable:""`
    activitiesRepo  ports.ActivitiesRepo     `injectable:""`
}

func NewProcessFollowActivityUsecase() *ProcessFollowActivityUsecase {
    return &ProcessFollowActivityUsecase{}
}

func (uc *ProcessFollowActivityUsecase) Exec(ctx context.Context, activityId string) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    activityEntity, err := uc.getActivity(ctx, activityId)
    if err != nil {
        return err
    }

    activity, err := structutil.MapToStruct[activitypub.Activity](activityEntity.Activity)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }
    
    err = uc.process(ctx, activity, activityId)
    if err != nil {
        log.Debug(ctx, err.Error())
        activityEntity.Status = "failed"
    } else {
        log.Debug(ctx, "Succeed")
        activityEntity.Status = "succeed"
    }
    activityEntity.UpdateAt.Set(time.Now().UTC())

    if err := uc.activitiesRepo.UpdateActivity(ctx, activityEntity); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return err
    }

    return nil
}

func (uc *ProcessFollowActivityUsecase) process(ctx context.Context, activity *activitypub.Activity, activityId string) error {
    if err := uc.validateActivity(activity); err != nil {
        return err
    }

    log.Info(ctx, "Get follower")
    follower, err := uc.getFollower(ctx, activity)
    if err != nil {
        return err
    }

    log.Info(ctx, "Get Target")
    target, err := uc.getTarget(ctx, activity)
    if err != nil {
        return err
    }

    log.Info(ctx, "Create Follow Entity")
    follow := uc.createFollow(follower, target, activityId)

    // insert data to db
    log.Info(ctx, "Save Follow to db")
    if err := uc.followingRepo.CreateFollow(ctx, follow); err != nil {
        return err
    }

    return nil
}

func (uc *ProcessFollowActivityUsecase) getActivity(ctx context.Context, activityId string) (*entities.ActivityEntity, error) {
    log.Debug(ctx, activityId)
    activityEntity, err := uc.activitiesRepo.FindActivityById(ctx, activityId) 
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    if activityEntity == nil {
        log.Warn(ctx, "Activity not found")
        return nil, appstatus.NotFound("Activity not found")
    }

    return activityEntity, nil
}

func (uc *ProcessFollowActivityUsecase) validateActivity(activity *activitypub.Activity) error {
    if activity == nil {
        return appstatus.BadValue("Invalid activity.")
    }

    if activity.Type != activitypub.FollowType {
        return appstatus.BadValue("Invalid activity type.")
    }

    if !activity.Actor.IsLink() {
        appstatus.BadValue("Unsupport actor type.")
    }

    return nil
}

func (uc *ProcessFollowActivityUsecase) getFollower(ctx context.Context, activity *activitypub.Activity) (*entities.UserEntity, error) {
    followerId := activity.Actor.GetLink().String()
    follower, err := uc.userRepo.FindUserByActorId(ctx, followerId)
    if err != nil {
        return nil, err
    }
    if follower == nil {
        log.Warn(ctx, "Follower not found")
        return nil, appstatus.NotFound("Follower not found.")
    }

    return follower, nil
}

func (uc *ProcessFollowActivityUsecase) getTarget(ctx context.Context, activity *activitypub.Activity) (*entities.UserEntity, error) {
    if !activity.Object.IsLink() {
        log.Warn(ctx, "Unsupport object type.")
        appstatus.BadValue("Unsupport object type.")
    }
    targetId := activity.Object.GetLink().String()

    parsedUrl, err := url.Parse(targetId)
    if err != nil {
        log.Warn(ctx, "Invalid following ID.")
        return nil, appstatus.BadValue("Invalid following ID.")
    }

    if strings.TrimPrefix(parsedUrl.Hostname(), "www.") != config.Fommu.Domain {
        log.Warn(ctx, "Invalid following ID.")
        return nil, appstatus.BadValue("Invalid following ID.")
    }

    targetUsername := path.Base(parsedUrl.Path)
    target, err := uc.userRepo.FindUserByUsername(ctx, targetUsername, config.Fommu.Domain)
    if err != nil {
        return nil, err
    }

    if target == nil {
        log.Warn(ctx, "Target user not found")
        return nil, appstatus.NotFound("Target user not found.")
    }

    return target, nil
}

func (uc *ProcessFollowActivityUsecase) createFollow(follower *entities.UserEntity, target *entities.UserEntity, activityId string) *entities.FollowEntity {
    following := entities.NewFollowEntity()
    following.ID = uuid.New().String()
    following.Follower = follower.ID
    following.Following = target.ID
    if target.AutoApproveFollower == true {
        following.Status = "followed"
    } else {
        following.Status = "pending"
    }
    if activityId != "" {
        following.ActivityId.Set(activityId)
    }
    following.CreateAt = time.Now().UTC()

    return following
}
