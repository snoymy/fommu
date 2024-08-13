package usecases

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"context"
	"path/filepath"
	"strings"
)

type GetFileUsecase struct {
    mediaRepo repo.MediaRepo `injectable:""`
}

func NewGetFileUsecase() *GetFileUsecase {
    return &GetFileUsecase{}
}

func (uc *GetFileUsecase) Exec(ctx context.Context, filename string) ([]byte, *entities.MediaEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if filename == "" {
        log.Info(ctx, "File name is empty")
        return nil, nil, appstatus.BadValue("No file name provided.")
    }

    media, err := uc.getMedia(ctx, filename)
    if err != nil {
        return nil, nil, err
    }
    
    fileBytes, err := uc.getFile(ctx, filename)
    if err != nil {
        return nil, nil, err
    }

    return fileBytes, media, nil
}

func (uc *GetFileUsecase) getMedia(ctx context.Context, filename string) (*entities.MediaEntity, error) {
    log.Info(ctx, "Get media id")
    mediaId := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
    log.Info(ctx, "Find media data")
    media, err := uc.mediaRepo.FindMediaByID(ctx, mediaId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Failed to get media data")
    }

    if media == nil {
        log.Info(ctx, "Media data not found")
        return nil, appstatus.NotFound("Media not found.")
    }

    return media, nil
}

func (uc *GetFileUsecase) getFile(ctx context.Context, filename string) ([]byte, error) {
    log.Info(ctx, "Read file")
    fileBytes, err := uc.mediaRepo.ReadFile(ctx, filename)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Failed to get file")
    }

    if len(fileBytes) == 0 {
        log.Error(ctx, "File is size is 0")
        return nil, appstatus.NotFound("Cannot get file.")
    }

    return fileBytes, nil
}
