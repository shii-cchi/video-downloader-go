package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/service"
)

type videoHandler struct {
	videoService service.Videos
}

func newVideoHandler(videoService service.Videos) *videoHandler {
	return &videoHandler{
		videoService: videoService,
	}
}

func (h videoHandler) downloadVideo(w http.ResponseWriter, r *http.Request) {
	downloadInput := r.Context().Value(delivery.DownloadInputKey).(dto.DownloadInputDto)

	err := h.videoService.DownloadVideo(downloadInput)
	if err != nil {
		if strings.HasPrefix(err.Error(), delivery.ErrNotFoundVideoID) || strings.HasPrefix(err.Error(), delivery.ErrParsingURL) {
			log.WithError(err).Error(delivery.ErrGettingID)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingID})
			return
		}

		log.WithError(err).Error(delivery.ErrDownloadingVideo)
		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDownloadingVideo})
		return
	}

	log.Infof(delivery.SuccessfulLoad+": %s\n", downloadInput.VideoURL)
	delivery.RespondWithJSON(w, http.StatusOK, nil)
}
