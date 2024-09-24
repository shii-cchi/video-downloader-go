package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
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
	RealPath    string             `bson:"real_path"`
	UserPath    string             `bson:"user_path"`
	PreviewPath string             `bson:"preview_path"`
}

type VideoRangeInfo struct {
	RangeStart int64
	RangeEnd   int64
	VideoInfo  VideoFileInfo
}

type VideoFileInfo struct {
	VideoName string
	FileSize  int64
	VideoFile *os.File
}
