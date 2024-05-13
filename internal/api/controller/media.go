package controller

import (
	"app/internal/appstatus"
	"app/internal/api/core/usecase"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MediaController struct {
    uploadFileUsecase *usecase.UploadFileUsecase
    getFileUsecase *usecase.GetFileUsecase
}

func NewMediaController(uploadFileUsecase *usecase.UploadFileUsecase, getFileUsecase *usecase.GetFileUsecase) *MediaController {
    return &MediaController{
        uploadFileUsecase: uploadFileUsecase,
        getFileUsecase: getFileUsecase,
    }
}

func (c *MediaController) UploadFile(w http.ResponseWriter, r *http.Request) error {
    if err := r.ParseMultipartForm(30 << 20); err != nil {
        return err
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        return err
    }
    defer file.Close()

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
        return appstatus.InternalServerError("Cannot get user id.")
    }

    media, err := c.uploadFileUsecase.Exec(r.Context(), fileBytes, fileName, size, mimeType, userId)
    if err != nil {
        return err
    }

    res := map[string]interface{}{
        "url": media.Url,
        "type": media.Type,
        "mime_type": media.MimeType,
        "description": media.Description,
    }

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *MediaController) GetFile(w http.ResponseWriter, r *http.Request) error {
    fileName := chi.URLParam(r, "fileName")
    if fileName == "" {
        return appstatus.BadValue("No file name provided.")
    }

    fileBytes, media, err := c.getFileUsecase.Exec(r.Context(), fileName)
    if err != nil {
        return err
    }

    if media == nil {
        return appstatus.NotFound("Media not found.")
    }

    if fileBytes == nil {
        return appstatus.NotFound("File not found.")
    }

	w.Header().Set("Content-Type", media.MimeType)
	w.Write(fileBytes)

    return nil
}
