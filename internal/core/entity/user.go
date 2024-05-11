package entity

import (
	"app/internal/types"
	"time"
)

type UserEntity struct {
    ID                  string                      `db:"id"`
    Email               string                      `db:"email"`
    Username            string                      `db:"username"`
    Displayname         string                      `db:"display_name"`
    PasswordHash        string                      `db:"password_hash"`
    Status              string                      `db:"status"`
    NamePrefix          types.Nullable[string]      `db:"name_prefix"`
    NameSuffix          types.Nullable[string]      `db:"name_suffix"`
    Avatar              types.Nullable[string]      `db:"avatar"`
    Banner              types.Nullable[string]      `db:"banner"`
    Bio                 types.Nullable[string]      `db:"bio"`
    Tag                 types.Nullable[[]string]    `db:"tag"`
    Remote              bool                        `db:"remote"`
    URL                 string                      `db:"url"`
    Discoverable        bool                        `db:"discoverable"`
    AutoApproveFollower bool                        `db:"auto_approve_follower"`
    FollowerCount       int                         `db:"follower_count"`
    FollowingCount      int                         `db:"following_count"`
    PublicKey           string                      `db:"public_key"`
    PrivateKey          string                      `db:"private_key"`
    RedirectURL         types.Nullable[string]      `db:"redirect_url"`
    CreateAt            time.Time                   `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]   `db:"update_at"`
}

func NewUserEntity() *UserEntity {
    return &UserEntity{}
}
