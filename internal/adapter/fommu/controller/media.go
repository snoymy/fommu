package controller

import (
	"app/internal/application/fommu/usecase"
	"app/internal/core/appstatus"
	"app/internal/log"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MediaController struct {
    uploadFileUsecase *usecase.UploadFileUsecase `injectable:""`
    getFileUsecase    *usecase.GetFileUsecase    `injectable:""`
}

func NewMediaController() *MediaController {
    return &MediaController{}
}

func (c *MediaController) UploadFile(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Parse Multi part form")
    if err := r.ParseMultipartForm(30 << 20); err != nil {
        log.Error(ctx, "Response with error: " + err.Error())
        return err
    }

    log.Info(ctx, "Get form file from request")
    file, handler, err := r.FormFile("file")
    if err != nil {
        log.Error(ctx, "Response with error: " + err.Error())
        return err
    }
    defer file.Close()

    log.Info(ctx, "Getting file metadata")
    fileHeader := make([]byte, 512) // Read first 512 bytes
    bytesRead, err := file.Read(fileHeader)
    if err != nil {
        return err
    }

    _, err = file.Seek(0, 0)
    if err != nil {
        return err
    }

	fileHeader = fileHeader[:bytesRead]

    fileBytes, err := io.ReadAll(file)
    fileName := handler.Filename
    size := handler.Size
    mimeType := http.DetectContentType(fileHeader)
    userId, ok := r.Context().Value("userId").(string)
    if !ok {
        log.Warn(ctx, "Cannot get user id")
        return appstatus.InternalServerError("Cannot get user id.")
    }

    media, err := c.uploadFileUsecase.Exec(ctx, fileBytes, fileName, size, mimeType, userId)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    res := map[string]interface{}{
        "url": media.Url,
        "type": media.Type,
        "mime_type": media.MimeType,
        "description": media.Description,
    }

    bytes, err := json.Marshal(res)
    if err != nil {
        log.Error(ctx, err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *MediaController) GetFile(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    fileName := chi.URLParam(r, "fileName")
    if fileName == "" {
        log.Info(ctx, "File name is empty")
        return appstatus.BadValue("No file name provided.")
    }

    fileBytes, media, err := c.getFileUsecase.Exec(ctx, fileName)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    if media == nil {
        err := appstatus.NotFound("Media not found.")
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    if fileBytes == nil {
        err := appstatus.NotFound("File not found.")
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

	w.Header().Set("Content-Type", media.MimeType)
	w.Write(fileBytes)

    return nil
}
