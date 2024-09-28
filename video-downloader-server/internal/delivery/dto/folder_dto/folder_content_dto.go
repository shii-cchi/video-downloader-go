package folder_dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"video-downloader-server/internal/delivery/dto/video_dto"
)

type FolderContentDto struct {
	ID      primitive.ObjectID   `json:"id"`
	Folders []FolderDto          `json:"folders"`
	Videos  []video_dto.VideoDto `json:"videos"`
}
