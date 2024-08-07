package usecase

import (
	"app/internal/application/activitypub/repo"
	"app/internal/config"
	"app/internal/core/appstatus"
	"app/internal/core/entity"
	"context"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/snoymy/activitypub"
)

type FollowUserUsecase struct {
    userRepo      repo.UsersRepo     `injectable:""`
    followingRepo repo.FollowRepo `injectable:""`
}

func NewFollowUserUsecase() *FollowUserUsecase {
    return &FollowUserUsecase{}
}

func (uc *FollowUserUsecase) Exec(ctx context.Context, activity *activitypub.Activity) error {
    if err := uc.validateActivity(activity); err != nil {
        return err
    }

    follower, err := uc.getFollower(ctx, activity)
    if err != nil {
        return err
    }

    target, err := uc.getTarget(ctx, activity)
    if err != nil {
        return err
    }

    following := uc.createFollow(follower, target)

    // insert data to db
    if err := uc.followingRepo.CreateFollowing(ctx, following); err != nil {
        return err
    }

    return nil
}

func (uc *FollowUserUsecase) validateActivity(activity *activitypub.Activity) error {
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

func (uc *FollowUserUsecase) getFollower(ctx context.Context, activity *activitypub.Activity) (*entity.UserEntity, error) {
    followerId := activity.Actor.GetLink().String()
    follower, err := uc.userRepo.FindUserByActorId(ctx, followerId)
    if err != nil {
        return nil, err
    }
    if follower == nil {
        return nil, appstatus.NotFound("Follower not found.")
    }

    return follower, nil
}

func (uc *FollowUserUsecase) getTarget(ctx context.Context, activity *activitypub.Activity) (*entity.UserEntity, error) {
    if !activity.Object.IsLink() {
        appstatus.BadValue("Unsupport object type.")
    }
    targetId := activity.Object.GetLink().String()

    parsedUrl, err := url.Parse(targetId)
    if err != nil {
        return nil, appstatus.BadValue("Invalid following ID.")
    }

    if strings.TrimPrefix(parsedUrl.Hostname(), "www.") != config.Fommu.Domain {
        return nil, appstatus.BadValue("Invalid following ID.")
    }

    targetUsername := path.Base(parsedUrl.Path)
    target, err := uc.userRepo.FindUserByUsername(ctx, targetUsername, config.Fommu.Domain)
    if err != nil {
        return nil, err
    }

    if target == nil {
        return nil, appstatus.NotFound("Target user not found.")
    }

    return target, nil
}

func (uc *FollowUserUsecase) createFollow(follower *entity.UserEntity, target *entity.UserEntity) *entity.FollowEntity {
    following := entity.NewFollowEntity()
    following.ID = uuid.New().String()
    following.Follower = follower.ID
    following.Following = target.ID
    if target.AutoApproveFollower == true {
        following.Status = "followed"
    } else {
        following.Status = "pending"
    }
    following.CreateAt = time.Now().UTC()

    return following
}
