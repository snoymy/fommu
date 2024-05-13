package api

import (
	"app/internal/adapterimpl"
	"app/internal/controller"
	"app/internal/core/usecase"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/repoimpl"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewMediaRoute(db *sqlx.DB) func(chi.Router) {
    mediaRepo := repoimpl.NewMediaRepoImpl(db)
    fileAdapter := adapterimpl.NewFileAdapterImpl()
    uploadFile := usecase.NewUploadFileUsecase(mediaRepo, fileAdapter)
    getFile := usecase.NewGetFileUsecase(mediaRepo, fileAdapter)
    mediaController := controller.NewMediaController(uploadFile, getFile)

    sessionRepo := repoimpl.NewSessionReoImpl(db)
    auth := usecase.NewAuthUsecase(sessionRepo)
    authMiddleware := middleware.NewAuthMiddleware(auth)

    return func(r chi.Router) {
        r.Group(func(r chi.Router) {
            r.Use(authMiddleware)
            r.Post("/upload", handler.Handle(mediaController.UploadFile))
        })
        r.Get("/{fileName}", handler.Handle(mediaController.GetFile))
    }
}
