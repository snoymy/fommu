package repoimpl

import (
	"app/internal/core/entity"
	"app/internal/config"
	"context"
	"net/url"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
)

type MediaRepoImpl struct {
    db *sqlx.DB `injectable:""`
}

func NewMediaRepoImpl() *MediaRepoImpl {
    return &MediaRepoImpl{}
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

func (a *MediaRepoImpl) WriteFile(ctx context.Context, file []byte, fileName string) (string, error) {
    filePath := path.Join("./media", fileName)

    outFile, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer outFile.Close()

    _, err = outFile.Write(file)
    if err != nil {
        return "", nil
    }

    fileUrl, err := url.JoinPath(config.Fommu.FileHost, fileName)
    if err != nil {
        return "", err
    }

    return fileUrl, nil
}

func (a *MediaRepoImpl) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
    filePath := path.Join("./media", fileName)
    fileBytes, err := os.ReadFile(filePath)
	if err != nil {
	    return nil, err
	}
    return fileBytes, nil
}
