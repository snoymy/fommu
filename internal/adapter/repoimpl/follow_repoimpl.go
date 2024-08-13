package repoimpl

import (
	"app/internal/adapter/command"
	"app/internal/adapter/mapper"
	"app/internal/adapter/query"
	"app/internal/config"
	"app/internal/core/entities"
	"app/internal/log"
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/snoymy/activitypub"
)

type FollowRepoImpl struct {
    db *sqlx.DB                 `injectable:""`
    query *query.Query          `injectable:""`
    command *command.Command    `injectable:""`
}

func NewFollowRepoImpl() *FollowRepoImpl {
    return &FollowRepoImpl{}
}

func (r *FollowRepoImpl) CreateFollow(ctx context.Context, follow *entities.FollowEntity) error {
    followsCount, err := r.query.CountFollows(ctx, follow) 
    if err != nil {
        return err
    }
    if followsCount > 0 {
        return nil
    }

    pendingFollow := *follow
    pendingFollow.Status = "pending"
    if err := r.command.CreateFollow(ctx, &pendingFollow); err != nil {
        return err
    }

    if follow.Status == "followed" {
        err := r.command.AcceptFollow(ctx, follow)
        if err != nil {
            return err
        }
        go r.sendAcceptActivity(ctx, follow)
    }

    return nil
}

func (r *FollowRepoImpl) sendAcceptActivity(ctx context.Context, follow *entities.FollowEntity) error {
    activityEnitity, err := r.query.FindActivityById(ctx, follow.ActivityId.ValueOrZero())
    if err != nil {
        return err
    }

    entityId := uuid.New().String()
    acceptActivityId, err := url.JoinPath(config.Fommu.URL, "activities/accept", entityId)

    activity, err := mapper.MapToStruct[activitypub.Activity](activityEnitity.Activity)
    if err != nil {
        return err
    }
    acceptActivity := activitypub.AcceptNew(activitypub.IRI(acceptActivityId), activity)
    acceptActivity.Actor = activity.Object.GetLink()

    activityMap, err := mapper.StructToMap(acceptActivity)
    if err != nil {
        return err
    }

    log.Debug(ctx, "create activity")
    acceptActivityEntity := entities.NewActiActivityEntity()
    acceptActivityEntity.ID = entityId
    acceptActivityEntity.Type.Set(string(activitypub.AcceptType))
    acceptActivityEntity.Remote = false
    acceptActivityEntity.Activity = activityMap
    acceptActivityEntity.Status = "pending"
    acceptActivityEntity.CreateAt = time.Now().UTC()
    if err := r.command.CreateActivity(ctx, acceptActivityEntity); err != nil {
        return err
    }
    log.Debug(ctx, "send activity")

    targetUrl, err := url.JoinPath(string(activity.Actor.GetLink()), "inbox")
    if err != nil {
        return err
    }

    following, err := r.query.FindUserById(ctx, follow.Following)
    if err != nil {
        return err
    }

    keyId := following.ActorId + "#main-key"
    privateKey := following.PrivateKey.ValueOrZero()
    if err := r.command.SendActivity(ctx, targetUrl, privateKey, keyId, acceptActivity); err != nil {
        acceptActivityEntity.Status = "failed"
        if err := r.command.UpdateActivity(ctx, acceptActivityEntity); err != nil {
            return err
        }

        return err
    }

    acceptActivityEntity.Status = "succeed"
    acceptActivityEntity.UpdateAt.Set(time.Now().UTC())
    if err := r.command.UpdateActivity(ctx, acceptActivityEntity); err != nil {
        return err
    }

    return nil
}
