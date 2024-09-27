package video_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type DownloadVideoDto struct {
	VideoURL string             `json:"video_url" validate:"required,url"`
	Type     string             `json:"type" validate:"required,oneof=youtube general"`
	Quality  string             `json:"quality" validate:"omitempty,oneof=2160p 1440p 1080p 720p 480p 360p 240p 144p best"`
	FolderID primitive.ObjectID `json:"folder_id" validate:"required,objectid"`
}
