package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/service"
)

type Handler struct {
	videoHandler *videoHandler
	validator    *validator.Validate
	extensionURL string
}

func NewHandler(service *service.Service, v *validator.Validate, extensionURL string) *Handler {
	vh := newVideoHandler(service.Videos)

	return &Handler{
		videoHandler: vh,
		validator:    v,
		extensionURL: extensionURL,
	}
}

func (h Handler) RegisterRoutes(r *chi.Mux) {
	r.Use(middleware.ApplyCors(h.extensionURL))

	r.With(middleware.CheckDownloadInput(h.validator)).Post("/videos/download", h.videoHandler.downloadVideo)
}
