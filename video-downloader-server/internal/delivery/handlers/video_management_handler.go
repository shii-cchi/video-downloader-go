package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/service"
)

type videoManagementHandler struct {
	videoManagementService service.VideoManagement
}

func newVideoManagementHandler(videoManagementService service.VideoManagement) *videoManagementHandler {
	return &videoManagementHandler{
		videoManagementService: videoManagementService,
	}
}

func (h videoManagementHandler) streamVideo(w http.ResponseWriter, r *http.Request) {
	videoName := r.Context().Value(delivery.VideoNameInputKey).(string)
	rangeHeader := r.Header.Get("Range")

	videoRangeInfo, err := h.videoManagementService.GetVideoRange(videoName, rangeHeader)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingVideoRange)
		if strings.HasPrefix(err.Error(), delivery.ErrVideoNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrVideoNotFound})
			return
		}

		if strings.HasPrefix(err.Error(), delivery.ErrInvalidRangeHeader) {
			delivery.RespondWithJSON(w, http.StatusRequestedRangeNotSatisfiable, delivery.JsonError{Error: delivery.ErrInvalidRangeHeader})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingVideoRange})
		return
	}

	delivery.RespondWithVideo(w, videoRangeInfo)
}

func (h videoManagementHandler) downloadVideo(w http.ResponseWriter, r *http.Request) {

}
