package service

import (
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/internal/repository"
)

type VideoDownload interface {
	Download(input dto.DownloadInputDto) error
}

type VideoManagement interface {
	GetVideoRange(videoName string, rangeHeader string) (VideoRangeInfo, error)
	GetVideoToDownload(videoName string) (VideoFileInfo, error)
	//GetVideos() ([]string, error)
	//DeleteVideo(videoName string) error
	//DownloadVideoFile(videoName string) ([]byte, error)
}

type Service struct {
	VideoDownload   VideoDownload
	VideoManagement VideoManagement
}

func NewService(repo *repository.Repository) *Service {
	videoDownloadService := newVideoDownloadService(repo.VideoDownload)
	videoManagementService := newVideoManagementService(repo.VideoManagement)

	return &Service{
		VideoDownload:   videoDownloadService,
		VideoManagement: videoManagementService,
	}
}
