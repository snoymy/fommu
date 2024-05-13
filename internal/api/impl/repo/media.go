package repo

import (
	"app/internal/api/core/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type MediaRepoImpl struct {
    db *sqlx.DB
}

func NewMediaRepoImpl(db *sqlx.DB) *MediaRepoImpl {
    return &MediaRepoImpl{
        db: db,
    }
}

func (r *MediaRepoImpl) CreateMedia(ctx context.Context, media *entity.MediaEntity) error {
    _, err := r.db.Exec(
        `
        insert into public.media
        (id, url, "type", mime_type, original_file_name, description, "owner", status, refernce_count, create_at, update_at)
        values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);
        `,
        media.ID, media.Url, media.Type, media.MimeType, media.OriginalFileName, media.Description, media.Owner, 
        media.Status, media.ReferenceCount, media.CreateAt, media.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (r *MediaRepoImpl) FindMediaByID(ctx context.Context, id string) (*entity.MediaEntity, error) {
    var media []*entity.MediaEntity = nil
    err := r.db.Select(&media, "select * from media where id=$1", id)
    if err != nil {
        return nil, err
    }

    if media == nil {
        return nil, nil
    }

    return media[0], nil
}
