package repo

import (
	"app/internal/core/entity"
	"context"
)

type MediaRepo interface {
    CreateMedia(ctx context.Context, media *entity.MediaEntity) error
    FindMediaByID(ctx context.Context, id string) (*entity.MediaEntity, error)
}
