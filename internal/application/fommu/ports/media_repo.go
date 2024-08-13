package ports

import (
	"app/internal/core/entities"
	"context"
)

type MediaRepo interface {
    CreateMedia(ctx context.Context, media *entities.MediaEntity) error
    FindMediaByID(ctx context.Context, id string) (*entities.MediaEntity, error)
    WriteFile(ctx context.Context, file []byte, fileName string) (string, error)
    ReadFile(ctx context.Context, fileUrl string) ([]byte, error)
}
