package service

import (
	"video-downloader-server/internal/delivery/dto"
)

type Download interface {
	DownloadVideo(input dto.DownloadInputDto) error
}

type VideoManagement interface {
	GetVideoRange(videoName string, rangeHeader string) (VideoRangeInfo, error)
	GetVideoToDownload(videoName string) (VideoInfo, error)
	//GetVideos() ([]string, error)
	//DeleteVideo(videoName string) error
	//DownloadVideoFile(videoName string) ([]byte, error)
}

type Service struct {
	Download        Download
	VideoManagement VideoManagement
}

func NewService() *Service {
	downloadService := newDownloadService()
	videoManagementService := newVideoManagementService()

	return &Service{
		Download:        downloadService,
		VideoManagement: videoManagementService,
	}
}
