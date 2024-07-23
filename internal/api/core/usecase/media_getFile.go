package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
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

func (uc *GetFileUsecase) Exec(ctx context.Context, fileName string) ([]byte, *entity.MediaEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if fileName == "" {
        log.Info(ctx, "File name is empty")
        return nil, nil, appstatus.BadValue("No file name provided.")
    }

    log.Info(ctx, "Get media id")
    mediaId := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
    log.Info(ctx, "Find media data")
    media, err := uc.mediaRepo.FindMediaByID(ctx, mediaId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, nil, appstatus.InternalServerError("Failed to get media data")
    }

    if media == nil {
        log.Info(ctx, "Media data not found")
        return nil, nil, appstatus.NotFound("Media not found.")
    }
    
    log.Info(ctx, "Read file")
    fileBytes, err := uc.mediaRepo.ReadFile(ctx, fileName)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, nil, appstatus.InternalServerError("Failed to get file")
    }

    if len(fileBytes) == 0 {
        log.Error(ctx, "File is size is 0")
        return nil, nil, appstatus.NotFound("Cannot get file.")
    }

    return fileBytes, media, nil
}
