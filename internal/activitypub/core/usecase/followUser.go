package usecase

import (
	"app/internal/activitypub/core/adapter"
	"app/internal/activitypub/core/repo"
	"app/internal/appstatus"
	"context"

	"github.com/go-ap/activitypub"
)

type FollowUserUsecase struct {
    userRepo repo.UsersRepo
    activitypubAdapter adapter.ActivitypubAdapter
}

func NewFollowUserUsecase(userRepo repo.UsersRepo, activitypubAdapter adapter.ActivitypubAdapter) *FollowUserUsecase {
    return &FollowUserUsecase{
        userRepo: userRepo,
        activitypubAdapter: activitypubAdapter,
    }
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
    followerId := activity.GetID().String()
    _, err := uc.userRepo.FindUserByActorId(ctx, followerId)
    if err != nil {
        return err
    }

    // check if object is url
    // extract username from url
    // find target user by username with domain
    // insert data to db
    return nil
}

