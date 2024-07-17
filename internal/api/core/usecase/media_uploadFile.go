package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"app/internal/utils"
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type UploadFileUsecase struct {
    mediaRepo repo.MediaRepo
}

func NewUploadFileUsecase(mediaRepo repo.MediaRepo) *UploadFileUsecase {
    return &UploadFileUsecase{
        mediaRepo: mediaRepo,
    }
}

func (uc *UploadFileUsecase) Exec(ctx context.Context, fileBytes []byte, originalFileName string, fileSize int64, mimeType string, uploader string) (*entity.MediaEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    // check file if empty
    log.Info(ctx, "Check if file is empty")
    if len(fileBytes) == 0 {
        log.Info(ctx, "file size is 0")
        return nil, appstatus.BadValue("No file uploaded.")
    }
    // check file size
    log.Info(ctx, "Check if file is empty")
    if utils.Byte(fileSize) >= utils.MiB(30) {
        log.Info(ctx, "File is too large, file size: " + strconv.Itoa(int(utils.Byte(fileSize))) + " bytes")
        return nil, appstatus.BadValue("File size exceed limit.")
    }
    // get file extension from mimeType
    log.Info(ctx, "Get file extension from mimeType")
    extension := utils.GetExtensionFromMIME(mimeType)
    // generate uuid
    log.Info(ctx, "Generate uuid")
    id := uuid.New().String()
    // generate fileName
    log.Info(ctx, "Create file name")
    fileName := id + extension
    // write file
    log.Info(ctx, "Write file")
    fileUrl, err := uc.mediaRepo.WriteFile(ctx, fileBytes, fileName)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }
    // create media entity
    log.Info(ctx, "Create media entity")
    media := entity.NewMediaEntity()
    media.ID = id
    media.Url = fileUrl
    media.PreviewUrl.Set(fileUrl)
    media.Type = utils.GetMediaTypeFromMime(mimeType)
    media.MimeType = mimeType
    media.OriginalFileName = originalFileName
    media.Description.SetNull()
    media.Metadata.SetNull()
    media.Owner = uploader
    media.Status = "active"
    media.ReferenceCount = 0
    media.CreateAt = time.Now().UTC()
    // insert media entity to db
    log.Info(ctx, "Write media data")
    if err := uc.mediaRepo.CreateMedia(ctx, media); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }
    return media, nil
}
