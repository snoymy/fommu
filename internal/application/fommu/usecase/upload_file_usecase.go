package usecase

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/core/types"
	"app/internal/log"
	"app/internal/utils/mimeutil"
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type UploadFileUsecase struct {
    mediaRepo repo.MediaRepo `injectable:""`
}

func NewUploadFileUsecase() *UploadFileUsecase {
    return &UploadFileUsecase{}
}

func (uc *UploadFileUsecase) Exec(ctx context.Context, fileBytes []byte, originalFileName string, fileSize int64, mimeType string, uploader string) (*entities.MediaEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if err := uc.checkFile(ctx, fileBytes, fileSize); err != nil {
        return nil, err
    }

    media, err := uc.createMedia(ctx, fileBytes, mimeType, originalFileName, uploader)
    if err != nil {
        return nil, err
    }

    log.Info(ctx, "Write media data")
    if err := uc.mediaRepo.CreateMedia(ctx, media); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }

    return media, nil
}

func (uc *UploadFileUsecase) checkFile(ctx context.Context, file []byte, fileSize int64) error {
    log.Info(ctx, "Check if file is empty")
    if len(file) == 0 {
        log.Info(ctx, "file size is 0")
        return appstatus.BadValue("No file uploaded.")
    }

    log.Info(ctx, "Check if file is empty")
    if types.Byte(fileSize) >= types.MiB(30) {
        log.Info(ctx, "File is too large, file size: " + strconv.Itoa(int(types.Byte(fileSize))) + " bytes")
        return appstatus.BadValue("File size exceed limit.")
    }

    return nil
}

func (uc *UploadFileUsecase) createMedia(ctx context.Context, file []byte, mimeType string, filename string, owner string) (*entities.MediaEntity, error) {
    log.Info(ctx, "Get file extension from mimeType")
    extension := mimeutil.GetExtensionFromMIME(mimeType)
    // generate uuid
    log.Info(ctx, "Generate uuid")
    id := uuid.New().String()
    // generate fileName
    log.Info(ctx, "Create file name")
    fileName := id + extension
    // write file
    log.Info(ctx, "Write file")
    fileUrl, err := uc.mediaRepo.WriteFile(ctx, file, fileName)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }
    // create media entity
    log.Info(ctx, "Create media entity")
    media := entities.NewMediaEntity()
    media.ID = id
    media.Url = fileUrl
    media.PreviewUrl.Set(fileUrl)
    media.Type = mimeutil.GetMediaTypeFromMime(mimeType)
    media.MimeType = mimeType
    media.OriginalFileName = filename 
    media.Description.SetNull()
    media.Metadata.SetNull()
    media.Owner = owner 
    media.Status = "active"
    media.ReferenceCount = 0
    media.CreateAt = time.Now().UTC()

    return media, nil
}
