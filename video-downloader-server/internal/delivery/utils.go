package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"video-downloader-server/internal/delivery/dto/video_dto"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithVideoRange(w http.ResponseWriter, info video_dto.VideoRangeInfoDto) {
	defer info.VideoInfo.VideoFile.Close()

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", info.RangeStart, info.RangeEnd, info.VideoInfo.FileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.RangeEnd-info.RangeStart+1))
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)

	info.VideoInfo.VideoFile.Seek(info.RangeStart, 0)
	if _, err := io.CopyN(w, info.VideoInfo.VideoFile, info.RangeEnd-info.RangeStart+1); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ErrGettingVideoRange)
	}
}

func RespondWithVideo(w http.ResponseWriter, info video_dto.VideoFileInfoDto) {
	defer info.VideoFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+info.VideoName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.FileSize))

	if _, err := io.Copy(w, info.VideoFile); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ErrDownloadingVideoFromServer)
	}
}
