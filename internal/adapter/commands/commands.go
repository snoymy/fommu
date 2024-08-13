package commands

import (
	"app/internal/infrastructure/httpclient"
	"app/internal/core/entities"
	"context"
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/jmoiron/sqlx"
	"github.com/snoymy/activitypub"
)

type Command struct {
    db *sqlx.DB `injectable:""`
    client httpclient.ActivitypubClient `injectable:""`
    bus EventBus.Bus `injectable:""`
}

func NewCommand() *Command {
    return &Command{}
}

func (c *Command) CreateUser(ctx context.Context, user *entities.UserEntity) error {
    _, err := c.db.Exec(
        `
        insert into users (
            id, email, password_hash, status, username, display_name, name_prefix, name_suffix, 
            bio, avatar, banner, attachment, tag, discoverable, auto_approve_follower, follower_count, following_count, 
            public_key, private_key, actor_id, url, inbox_url, outbox_url, followers_url, following_url, Domain, remote, redirect_url, 
            create_at, update_at
        )
        values
        ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
        `,
        user.ID, user.Email, user.PasswordHash, user.Status, user.Username, user.Displayname,
        user.NamePrefix, user.NameSuffix, user.Bio, user.Avatar, user.Banner, user.Attachment, user.Tag, user.Discoverable, 
        user.AutoApproveFollower, user.FollowerCount, user.FollowingCount, user.PublicKey, user.PrivateKey, user.ActorId,
        user.URL, user.InboxURL, user.OutboxURL, user.FollowersURL, user.FollowingURL, user.Domain, user.Remote, user.RedirectURL, 
        user.CreateAt, user.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) UpdateUser(ctx context.Context, user *entities.UserEntity) error {
    _, err := c.db.Exec(
        `
        update users set 
            display_name=$1, name_prefix=$2, name_suffix=$3, bio=$4, avatar=$5, banner=$6, 
            discoverable=$7, auto_approve_follower=$8, attachment=$9, tag=$10, follower_count=$11, following_count=$12, 
            url=$13, inbox_url=$14, outbox_url=$15, followers_url=$16, following_url=$17, 
            preference=$18, update_at=$19, email=$20, password_hash=$21
        where id = $22
        `,
        user.Displayname, user.NamePrefix, user.NameSuffix, user.Bio, user.Avatar, user.Banner, 
        user.Discoverable, user.AutoApproveFollower, user.Attachment, user.Tag, user.FollowerCount, user.FollowingCount,
        user.URL, user.InboxURL, user.OutboxURL, user.FollowersURL, user.FollowingURL,
        user.Preference, user.UpdateAt, user.Email, user.PasswordHash,
        user.ID,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) CreateSession(ctx context.Context, session *entities.SessionEntity) error {
    _, err := c.db.Exec(
        `
        insert into sessions (
            id, access_token, access_expire_at, refresh_token, refresh_expire_at, 
            owner, metadata, login_at, last_refresh
        )
        values
        ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        `,
        session.ID, session.AccessToken, session.AccessExpireAt, session.RefreshToken, session.RefreshExpireAt,
        session.Owner, session.Metadata, session.LoginAt, session.LastRefresh,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) UpdateSession(ctx context.Context, session *entities.SessionEntity) error {
    _, err := c.db.Exec(
        `
        update sessions 
        set access_token=$1, access_expire_at=$2, refresh_token=$3, refresh_expire_at=$4, last_refresh=$5
        where id=$6
        `,
        session.AccessToken, session.AccessExpireAt, session.RefreshToken, session.RefreshExpireAt, session.LastRefresh, session.ID,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) DeleteSession(ctx context.Context, id string) error {
    _, err := c.db.Exec("delete from sessions where id=$1", id)
    if err != nil {
        return err
    }

    return nil
}

func (c *Command) CreateActivity(ctx context.Context, activity *entities.ActivityEntity) error {
    _, err := c.db.Exec(
        `
        insert into activities (
            id, type, remote, activity, status, create_at, update_at
        )
        values
        ($1,$2,$3,$4,$5,$6,$7)
        `,
        activity.ID, activity.Type, activity.Remote, activity.Activity, activity.Status, activity.CreateAt, activity.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) UpdateActivity(ctx context.Context, activity *entities.ActivityEntity) error {
    _, err := c.db.Exec(
        `
        update activities set type=$2, remote=$3, activity=$4, status=$5, create_at=$6, update_at=$7 where id=$1`,
        activity.ID, activity.Type, activity.Remote, activity.Activity, activity.Status, activity.CreateAt, activity.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) CreateFollow(ctx context.Context, follow *entities.FollowEntity) error {
    _, err := c.db.Exec(
        `
        insert into follows (
            id, follower, following, status, create_at, update_at
        )
        values ($1,$2,$3,$4,$5,$6)
        `,
        follow.ID, follow.Follower, follow.Following, follow.Status, follow.CreateAt, follow.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (c *Command) AcceptFollow(ctx context.Context, follow *entities.FollowEntity) error {
    tx, err := c.db.Beginx()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
    }()

    _, err = tx.Exec(
        `update users set follower_count = follower_count + 1 where id = $1`,
        follow.Following,
    )
    if err != nil {
        return err
    }

    _, err = tx.Exec(
        `update users set following_count = following_count + 1 where id = $1`,
        follow.Follower,
    )
    if err != nil {
        return err
    }

    _, err = tx.Exec(
        `update follows set status = $2 where id = $1`,
        follow.ID, follow.Status,
    )
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

func (c *Command) SendActivity(ctx context.Context, url string, privateKey string, keyId string, activity *activitypub.Activity) error {
    if err := c.client.PublishActivity(ctx, url, privateKey, keyId, activity); err != nil {
        fmt.Println(err.Error())
        return err
    }

    return nil
}

func (c *Command) NotifyProcessActivity(ctx context.Context, activityEntity *entities.ActivityEntity) {
    c.bus.Publish("topic:process_activity", activityEntity.ID, activityEntity.Type.ValueOrZero())
}
