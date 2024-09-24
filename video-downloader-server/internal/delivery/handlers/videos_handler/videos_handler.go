package videos_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/internal/domain"
)

type VideosService interface {
	DownloadToServer(input dto.DownloadInputDto) error
	GetVideoFileInfo(videoID string) (domain.VideoFileInfo, error)
	GetVideoRangeInfo(videoID string, rangeHeader string) (domain.VideoRangeInfo, error)
}

type VideosHandler struct {
	videosService VideosService
	validator     *validator.Validate
}

func NewVideosHandler(videosService VideosService, validator *validator.Validate) *VideosHandler {
	return &VideosHandler{
		videosService: videosService,
		validator:     validator,
	}
}

func (h VideosHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/videos", func(r chi.Router) {
		r.With(middleware.CheckVideoIDInput).Get("/stream", h.streamVideo)
		r.With(middleware.CheckVideoIDInput).Get("/download-to-local", h.downloadVideoToLocal)

		r.With(middleware.CheckDownloadInput(h.validator)).Post("/download-to-server", h.downloadVideoToServer)
	})
}

func (h VideosHandler) downloadVideoToServer(w http.ResponseWriter, r *http.Request) {
	downloadInput := r.Context().Value(delivery.DownloadInputKey).(dto.DownloadInputDto)

	err := h.videosService.DownloadToServer(downloadInput)
	if err != nil {
		if strings.HasPrefix(err.Error(), domain.ErrNotFoundVideoID) || strings.HasPrefix(err.Error(), domain.ErrParsingURL) {
			log.WithError(err).Error(delivery.ErrGettingID)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrDownloadingVideoToServer, Message: delivery.ErrGettingID})
			return
		}

		log.WithError(err).Error(delivery.ErrDownloadingVideoToServer)
		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDownloadingVideoToServer})
		return
	}

	log.Infof(delivery.SuccessfulLoad+": %s\n", downloadInput.VideoURL)
	delivery.RespondWithJSON(w, http.StatusOK, nil)
}

func (h VideosHandler) downloadVideoToLocal(w http.ResponseWriter, r *http.Request) {
	videoID := r.Context().Value(delivery.VideoIDInputKey).(string)

	videoInfo, err := h.videosService.GetVideoFileInfo(videoID)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingVideo)
		if strings.HasPrefix(err.Error(), domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingVideo, Message: domain.ErrVideoNotFound})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingVideo})
		return
	}

	delivery.RespondWithVideo(w, videoInfo)
}

func (h VideosHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	videoID := r.Context().Value(delivery.VideoIDInputKey).(string)
	rangeHeader := r.Header.Get("Range")

	videoRangeInfo, err := h.videosService.GetVideoRangeInfo(videoID, rangeHeader)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingVideoRange)
		if strings.HasPrefix(err.Error(), domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingVideoRange, Message: domain.ErrVideoNotFound})
			return
		}

		if strings.HasPrefix(err.Error(), domain.ErrInvalidRangeHeader) {
			delivery.RespondWithJSON(w, http.StatusRequestedRangeNotSatisfiable, delivery.JsonError{Error: delivery.ErrGettingVideoRange, Message: domain.ErrInvalidRangeHeader})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingVideoRange})
		return
	}

	delivery.RespondWithVideoRange(w, videoRangeInfo)
}
