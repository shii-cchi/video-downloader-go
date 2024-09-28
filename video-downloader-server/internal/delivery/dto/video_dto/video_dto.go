package video_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type VideoDto struct {
	ID          primitive.ObjectID `json:"id"`
	VideoName   string             `json:"video_name"`
	FolderID    primitive.ObjectID `json:"folder_id"`
	RealPath    string             `json:"real_path"`
	PreviewPath string             `json:"preview_path"`
}
