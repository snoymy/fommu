package entity

import (
	"app/internal/core/types"
	"time"
)

type UserEntity struct {
    ID                  string                              `db:"id"`
    Email               types.Nullable[string]              `db:"email"`
    PasswordHash        types.Nullable[string]              `db:"password_hash"`
    Username            string                              `db:"username"`
    Displayname         string                              `db:"display_name"`
    Status              string                              `db:"status"`
    NamePrefix          types.Nullable[string]              `db:"name_prefix"`
    NameSuffix          types.Nullable[string]              `db:"name_suffix"`
    Avatar              types.Nullable[string]              `db:"avatar"`
    Banner              types.Nullable[string]              `db:"banner"`
    Bio                 types.Nullable[string]              `db:"bio"`
    Preference          types.Nullable[types.JsonObject]    `db:"preference"`
    Tag                 types.Nullable[types.JsonArray]     `db:"tag"`
    Attachment          types.Nullable[types.JsonArray]     `db:"attachment"`
    Remote              bool                                `db:"remote"`
    ActorId             string                              `db:"actor_id"`
    URL                 string                              `db:"url"`
    InboxURL            string                              `db:"inbox_url"`
    OutboxURL           string                              `db:"outbox_url"`
    FollowersURL        string                              `db:"followers_url"`
    FollowingURL        string                              `db:"following_url"`
    Domain              string                              `db:"domain"`
    Discoverable        bool                                `db:"discoverable"`
    AutoApproveFollower bool                                `db:"auto_approve_follower"`
    FollowerCount       int                                 `db:"follower_count"`
    FollowingCount      int                                 `db:"following_count"`
    PublicKey           string                              `db:"public_key"`
    PrivateKey          types.Nullable[string]              `db:"private_key"`
    RedirectURL         types.Nullable[string]              `db:"redirect_url"`
    JoinAt              types.Nullable[time.Time]           `db:"join_at"` 
    CreateAt            time.Time                           `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]           `db:"update_at"`
}

func NewUserEntity() *UserEntity {
    return &UserEntity{}
}
