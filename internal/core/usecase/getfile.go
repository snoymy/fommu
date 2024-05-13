package usecase

import (
	"app/internal/appstatus"
	"app/internal/core/adapter"
	"app/internal/core/entity"
	"app/internal/core/repo"
	"context"
	"path/filepath"
	"strings"
)

type GetFileUsecase struct {
    mediaRepo repo.MediaRepo
    fileAdapter adapter.FileAdapter
}

func NewGetFileUsecase(mediaRepo repo.MediaRepo, fileAdapter adapter.FileAdapter) *GetFileUsecase {
    return &GetFileUsecase{
        mediaRepo: mediaRepo,
        fileAdapter: fileAdapter,
    }
}

func (uc *GetFileUsecase) Exec(ctx context.Context, fileName string) ([]byte, *entity.MediaEntity, error) {
    if fileName == "" {
        return nil, nil, appstatus.BadValue("No file name provided.")
    }

    mediaId := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
    media, err := uc.mediaRepo.FindMediaByID(ctx, mediaId)
    if err != nil {
        return nil, nil, err
    }

    if media == nil {
        return nil, nil, appstatus.NotFound("Media not found.")
    }
    
    fileBytes, err := uc.fileAdapter.ReadFile(ctx, fileName)
    if err != nil {
        return nil, nil, err
    }

    if len(fileBytes) == 0 {
        return nil, nil, appstatus.NotFound("Cannot get file.")
    }

    return fileBytes, media, nil
}
