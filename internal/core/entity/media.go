package entity

import (
	"app/internal/types"
	"time"
)

type MediaEntity struct {
    ID                  string                      `db:"id"`
    Url                 string                      `db:"url"`
    Type                string                      `db:"type"`
    MimeType            string                      `db:"mime_type"`
    OriginalFileName    string                      `db:"original_file_name"`
    Description         types.Nullable[string]      `db:"description"`
    Owner               string                      `db:"owner"`
    Status              string                      `db:"status"`
    ReferenceCount      int                         `db:"refernce_count"`
    CreateAt            time.Time                   `db:"create_at"`
    UpdateAt            types.Nullable[time.Time]   `db:"update_at"`
}

func NewMediaEntity() *MediaEntity {
    return &MediaEntity{}
}
