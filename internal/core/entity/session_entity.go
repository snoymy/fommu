package entity

import (
	"app/internal/core/types"
	"time"
)

type SessionEntity struct {
    ID                  string                                          `db:"id"`
    AccessToken         string                                          `db:"access_token"`
    AccessExpireAt      time.Time                                       `db:"access_expire_at"`
    RefreshToken        string                                          `db:"refresh_token"`
    RefreshExpireAt     time.Time                                       `db:"refresh_expire_at"`
    Owner               string                                          `db:"owner"`
    Metadata            types.Nullable[types.JsonObject]                `db:"metadata"`
    LoginAt             time.Time                                       `db:"login_at"`
    LastRefresh         types.Nullable[time.Time]                       `db:"last_refresh"`
}

func NewSessionEntity() *SessionEntity {
    return &SessionEntity{}
}
