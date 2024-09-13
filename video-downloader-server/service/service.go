package service

import "video-downloader-server/internal/delivery/dto"

type Videos interface {
	DownloadVideo(input dto.DownloadInputDto) error
}

type Service struct {
	Videos Videos
}

func NewService() *Service {
	videoService := newVideoService()

	return &Service{
		Videos: videoService,
	}
}
