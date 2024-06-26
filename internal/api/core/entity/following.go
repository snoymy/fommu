package entity

import (
	"app/internal/types"
	"time"
)

type FollowingEntity struct {
    ID                  string                              `db:"id"`
    Follower            string                              `db:"follower"`
    Following           string                              `db:"following"`
    Status              string                              `db:"status"`
    CreateAt            time.Time                           `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]           `db:"update_at"`
}

func NewFollowingEntity() *FollowingEntity {
    return &FollowingEntity{}
}
