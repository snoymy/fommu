package usecase

import (
	"app/internal/activitypub/core/entity"
	"app/internal/activitypub/core/repo"
	"app/internal/appstatus"
	"app/internal/config"
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
    followingRepo repo.FollowingRepo `injectable:""`
}

func NewFollowUserUsecase() *FollowUserUsecase {
    return &FollowUserUsecase{}
}

func (uc *FollowUserUsecase) Exec(ctx context.Context, activity *activitypub.Activity) error {
    // check if activity is empty
    if activity == nil {
        return appstatus.BadValue("Invalid activity.")
    }
    // check if activity type is Follow
    if activity.Type != activitypub.FollowType {
        return appstatus.BadValue("Invalid activity type.")
    }
    // find follower user by actor id
    if !activity.Actor.IsLink() {
        appstatus.BadValue("Unsupport actor type.")
    }
    followerId := activity.Actor.GetLink().String()
    follower, err := uc.userRepo.FindUserByActorId(ctx, followerId)
    if err != nil {
        return err
    }
    if follower == nil {
        return appstatus.NotFound("Follower not found.")
    }

    // check if object is url
    if !activity.Object.IsLink() {
        appstatus.BadValue("Unsupport object type.")
    }
    targetId := activity.Object.GetLink().String()

    parsedUrl, err := url.Parse(targetId)
    if err != nil {
        return appstatus.BadValue("Invalid following ID.")
    }
    if strings.TrimPrefix(parsedUrl.Hostname(), "www.") != config.Fommu.Domain {
        return appstatus.BadValue("Invalid following ID.")
    }
    targetUsername := path.Base(parsedUrl.Path)
    target, err := uc.userRepo.FindUserByUsername(ctx, targetUsername, config.Fommu.Domain)
    if err != nil {
        return err
    }
    if target == nil {
        return appstatus.NotFound("Target user not found.")
    }

    following := entity.NewFollowingEntity()
    following.ID = uuid.New().String()
    following.Follower = follower.ID
    following.Following = target.ID
    if target.AutoApproveFollower == true {
        following.Status = "followed"
    } else {
        following.Status = "pending"
    }
    following.CreateAt = time.Now().UTC()
    // insert data to db
    if err := uc.followingRepo.CreateFollowing(ctx, following); err != nil {
        return err
    }

    return nil
}

