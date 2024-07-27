package entity

import (
	"app/internal/core/types"
	"time"
)

type MediaEntity struct {
    ID                  string                              `db:"id"`
    Url                 string                              `db:"url"`
    PreviewUrl          types.Nullable[string]              `db:"preview_url"`
    Type                string                              `db:"type"`
    MimeType            string                              `db:"mime_type"`
    OriginalFileName    string                              `db:"original_file_name"`
    Description         types.Nullable[string]              `db:"description"`
    Metadata            types.Nullable[types.JsonObject]    `db:"metadata"`
    Owner               string                              `db:"owner"`
    Status              string                              `db:"status"`
    ReferenceCount      int                                 `db:"refernce_count"`
    // Visibility          string                              `db:"visibility"`
    // Group               types.Nullable[string]              `db:"group"`
    CreateAt            time.Time                           `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]           `db:"update_at"`
}

func NewMediaEntity() *MediaEntity {
    return &MediaEntity{}
}
