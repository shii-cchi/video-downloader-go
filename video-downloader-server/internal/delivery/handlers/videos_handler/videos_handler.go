package videos_handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto/video_dto"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/internal/domain"
)

type VideosService interface {
	DownloadToServer(downloadVideoInput video_dto.DownloadVideoDto) error
	GetVideoFileInfo(videoID primitive.ObjectID) (video_dto.VideoFileInfoDto, error)
	GetVideoRangeInfo(videoID primitive.ObjectID, rangeHeader string) (video_dto.VideoRangeInfoDto, error)
	Rename(renameVideoInput video_dto.RenameVideoDto) (video_dto.VideoDto, error)
	Move(moveVideoInput video_dto.MoveVideoDto) (video_dto.VideoDto, error)
	Delete(deleteVideoInput video_dto.DeleteVideoDto) error
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
		r.With(middleware.ValidateDownloadVideoInput(h.validator)).Post("/download-to-server", h.downloadVideoToServer)
		r.With(middleware.ValidateVideoIDInput).Get("/download-to-local", h.downloadVideoToLocal)
		r.With(middleware.ValidateVideoIDInput).Get("/stream", h.streamVideo)
		r.With(middleware.ValidateRenameVideoInput(h.validator)).Put("/rename", h.renameVideo)
		r.With(middleware.ValidateMoveVideoInput(h.validator)).Put("/move", h.moveVideo)
		r.With(middleware.ValidateDeleteVideoInput(h.validator)).Delete("/", h.deleteVideo)
	})
}

func (h VideosHandler) downloadVideoToServer(w http.ResponseWriter, r *http.Request) {
	downloadVideoInput := r.Context().Value(delivery.DownloadVideoInputKey).(video_dto.DownloadVideoDto)

	err := h.videosService.DownloadToServer(downloadVideoInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrDownloadingVideoToServer)

		if errors.Is(err, domain.ErrNotFoundVideoID) || errors.Is(err, domain.ErrParsingURL) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrDownloadingVideoToServer, Message: delivery.ErrGettingID})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDownloadingVideoToServer})
		return
	}

	log.Infof(delivery.SuccessfulLoad+": %s\n", downloadVideoInput.VideoURL)
	delivery.RespondWithJSON(w, http.StatusOK, nil)
}

func (h VideosHandler) downloadVideoToLocal(w http.ResponseWriter, r *http.Request) {
	videoID := r.Context().Value(delivery.VideoIDInputKey).(primitive.ObjectID)

	videoInfo, err := h.videosService.GetVideoFileInfo(videoID)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingVideo)

		if errors.Is(err, domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingVideo, Message: domain.ErrVideoNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingVideo})
		return
	}

	delivery.RespondWithVideo(w, videoInfo)
}

func (h VideosHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	videoID := r.Context().Value(delivery.VideoIDInputKey).(primitive.ObjectID)
	rangeHeader := r.Header.Get("Range")

	videoRangeInfo, err := h.videosService.GetVideoRangeInfo(videoID, rangeHeader)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingVideoRange)

		if errors.Is(err, domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingVideoRange, Message: domain.ErrVideoNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrInvalidRangeHeader) {
			delivery.RespondWithJSON(w, http.StatusRequestedRangeNotSatisfiable, delivery.JsonError{Error: delivery.ErrGettingVideoRange, Message: domain.ErrInvalidRangeHeader.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingVideoRange})
		return
	}

	delivery.RespondWithVideoRange(w, videoRangeInfo)
}

func (h VideosHandler) renameVideo(w http.ResponseWriter, r *http.Request) {
	renameVideoInput := r.Context().Value(delivery.RenameVideoInputKey).(video_dto.RenameVideoDto)

	video, err := h.videosService.Rename(renameVideoInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrRenamingVideo)

		if errors.Is(err, domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrRenamingVideo, Message: domain.ErrVideoNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrRenamingVideo})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, video)
}

func (h VideosHandler) moveVideo(w http.ResponseWriter, r *http.Request) {
	moveVideoInput := r.Context().Value(delivery.MoveVideoInputKey).(video_dto.MoveVideoDto)

	video, err := h.videosService.Move(moveVideoInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrMovingVideo)

		if errors.Is(err, domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingVideo, Message: domain.ErrVideoNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingVideo, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrVideoAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingVideo, Message: domain.ErrVideoAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrMovingVideo})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, video)
}

func (h VideosHandler) deleteVideo(w http.ResponseWriter, r *http.Request) {
	deleteVideoInput := r.Context().Value(delivery.DeleteVideoInputKey).(video_dto.DeleteVideoDto)

	err := h.videosService.Delete(deleteVideoInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrDeletingVideo)

		if errors.Is(err, domain.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrDeletingVideo, Message: domain.ErrVideoNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDeletingVideo})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, nil)
}
