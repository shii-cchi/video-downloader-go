package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/service"
)

type Handler struct {
	downloadHandler        *downloadHandler
	videoManagementHandler *videoManagementHandler
	validator              *validator.Validate
	extensionURL           string
}

func NewHandler(service *service.Service, v *validator.Validate, extensionURL string) *Handler {
	dh := newDownloadHandler(service.VideoDownload)
	vm := newVideoManagementHandler(service.VideoManagement)

	return &Handler{
		downloadHandler:        dh,
		videoManagementHandler: vm,
		validator:              v,
		extensionURL:           extensionURL,
	}
}

func (h Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/extension", func(r chi.Router) {
		r.Use(middleware.CheckDownloadInput(h.validator))

		r.Post("/download-to-server", h.downloadHandler.downloadVideo)
	})

	r.Route("/videos", func(r chi.Router) {
		r.Use(middleware.CheckVideoNameInput)

		r.Get("/stream", h.videoManagementHandler.streamVideo)
		r.Get("/download-to-local", h.videoManagementHandler.downloadVideo)
	})
}
