package delivery

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"video-downloader-server/service"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf(ErrMarshalingJSON+": %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithVideoRange(w http.ResponseWriter, info service.VideoRangeInfo) {
	defer info.VideoInfo.VideoFile.Close()

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", info.RangeStart, info.RangeEnd, info.VideoInfo.FileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.RangeEnd-info.RangeStart+1))
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)

	info.VideoInfo.VideoFile.Seek(info.RangeStart, 0)
	io.CopyN(w, info.VideoInfo.VideoFile, info.RangeEnd-info.RangeStart+1)
}

func RespondWithVideo(w http.ResponseWriter, info service.VideoFileInfo) {
	defer info.VideoFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+info.VideoName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.FileSize))

	if _, err := io.Copy(w, info.VideoFile); err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, ErrDownloadingVideoFromServer)
	}
}
