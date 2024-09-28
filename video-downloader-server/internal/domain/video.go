package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CommonVideoDir         = "videos"
	VideoFormat            = ".mp4"
	YouTubeVideoType       = "youtube"
	DefaultRangePercentage = 0.05
)

type Video struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	VideoName   string             `bson:"video_name"`
	FolderID    primitive.ObjectID `bson:"folder_id"`
	RealPath    string             `bson:"real_path"`
	PreviewPath string             `bson:"preview_path"`
}
