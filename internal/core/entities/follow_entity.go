package entities

import (
	"app/internal/core/types"
	"time"
)

type FollowEntity struct {
    ID                  string                              `db:"id"`
    Follower            string                              `db:"follower"`
    Following           string                              `db:"following"`
    ActivityId          types.Nullable[string]              `db:"activity_id"`
    Status              string                              `db:"status"`
    CreateAt            time.Time                           `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]           `db:"update_at"`
}

func NewFollowEntity() *FollowEntity {
    return &FollowEntity{}
}
