package usecase

import (
	"app/internal/appstatus"
	"app/internal/core/adapter"
	"app/internal/core/entity"
	"app/internal/core/repo"
	"app/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

type UploadFileUsecase struct {
    mediaRepo repo.MediaRepo
    fileAdapter adapter.FileAdapter
}

func NewUploadFileUsecase(mediaRepo repo.MediaRepo, fileAdapter adapter.FileAdapter) *UploadFileUsecase {
    return &UploadFileUsecase{
        mediaRepo: mediaRepo,
        fileAdapter: fileAdapter,
    }
}

func (uc *UploadFileUsecase) Exec(ctx context.Context, fileBytes []byte, originalFileName string, fileSize int64, mimeType string, uploader string) (*entity.MediaEntity, error) {
    // check file if empty
    if len(fileBytes) == 0 {
        return nil, appstatus.BadValue("No file uploaded.")
    }
    // check file size
    if utils.Byte(fileSize) >= utils.MiB(30) {
        return nil, appstatus.BadValue("File size exceed limit.")
    }
    // get file extension from mimeType
    extension := utils.GetExtensionFromMIME(mimeType)
    // generate uuid
    id := uuid.New().String()
    // generate fileName
    fileName := id + extension
    // write file
    fileUrl, err := uc.fileAdapter.WriteFile(ctx, fileBytes, fileName)
    if err != nil {
        return nil, err
    }
    // create media entity
    media := entity.NewMediaEntity()
    media.ID = id
    media.Url = fileUrl
    media.Type = utils.GetMediaTypeFromMime(mimeType)
    media.MimeType = mimeType
    media.OriginalFileName = originalFileName
    media.Description.SetNull()
    media.Owner = uploader
    media.Status = "active"
    media.ReferenceCount = 0
    media.CreateAt = time.Now().UTC()
    // insert media entity to db
    if err := uc.mediaRepo.CreateMedia(ctx, media); err != nil {
        return nil, err
    }
    return media, nil
}
