package entities

import (
	"app/internal/core/types"
	"time"
)

type ActivityEntity struct {
    ID                  string                              `db:"id"`
    Type                types.Nullable[string]              `db:"type"`              
    Remote              bool                                `db:"remote"`
    Activity            types.JsonObject                    `db:"activity"`
    Status              string                              `db:"status"`
    CreateAt            time.Time                           `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]           `db:"update_at"`
}

func NewActiActivityEntity() *ActivityEntity {
    return &ActivityEntity{}
}
