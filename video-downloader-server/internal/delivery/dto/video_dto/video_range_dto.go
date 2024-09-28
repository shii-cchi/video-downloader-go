package video_dto

import "os"

type VideoRangeInfoDto struct {
	RangeStart int64
	RangeEnd   int64
	VideoInfo  VideoFileInfoDto
}

type VideoFileInfoDto struct {
	VideoName string
	FileSize  int64
	VideoFile *os.File
}
