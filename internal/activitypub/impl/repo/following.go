package repo

import (
	"app/internal/activitypub/core/entity"
	"app/internal/activitypub/core/repo"
	"context"

	"github.com/jmoiron/sqlx"
)

type FollowingRepoImpl struct {
    db *sqlx.DB `injectable:""`
}

func NewFollowingRepoImpl() repo.FollowingRepo {
    return &FollowingRepoImpl{}
}

func (r *FollowingRepoImpl) CreateFollowing(ctx context.Context, following *entity.FollowingEntity) error {
    var rowsCount []int = nil
    err := r.db.Select(&rowsCount, "select count(*) from following where follower = $1 and following = $2", following.Follower, following.Following)
    if err != nil {
        return err
    }

    if rowsCount != nil && rowsCount[0] > 0 {
        return nil
    }

    tx, err := r.db.Beginx()
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
        `
        insert into following (
            id, follower, following, status, create_at, update_at
        )
        values ($1,$2,$3,$4,$5,$6)
        `,
        following.ID, following.Follower, following.Following, following.Status, following.CreateAt, following.UpdateAt,
    )
    if err != nil {
        return err
    }

    if following.Status == "followed" {
        _, err = tx.Exec(
            `update users set follower_count = follower_count + 1 where id = $1`,
            following.Following,
        )
        if err != nil {
            return err
        }

        _, err = tx.Exec(
            `update users set following_count = following_count + 1 where id = $1`,
            following.Follower,
        )
        if err != nil {
            return err
        }

        // add activity object
        //go sent accept
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}
