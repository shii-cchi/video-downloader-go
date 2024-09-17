package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"video-downloader-server/internal/repository"
)

const (
	defaultRangePercentage = 0.05
	errVideoNotFound       = "video not found"
	errGettingFileInfo     = "err getting file info"
	errInvalidRangeHeader  = "invalid range header"
	errInvalidRangeFormat  = "invalid range format"
	errInvalidBytesFormat  = "invalid bytes format"
	errInvalidRangeStart   = "invalid start of range"
	errInvalidRangeEnd     = "invalid end of range"
	errInvalidRange        = "invalid range"
)

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

type VideoManagementService struct {
	repo repository.VideoManagement
}

func newVideoManagementService(repo repository.VideoManagement) *VideoManagementService {
	return &VideoManagementService{
		repo: repo,
	}
}

func (s VideoManagementService) GetVideoRange(videoName string, rangeHeader string) (VideoRangeInfo, error) {
	videoFile, err := os.Open("videos/" + videoName)
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errVideoNotFound+": %w", err)
	}

	fileInfo, err := videoFile.Stat()
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errGettingFileInfo+": %w", err)
	}

	rangeStart, rangeEnd, err := s.parseRangeHeader(rangeHeader, fileInfo.Size())
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errInvalidRangeHeader+": %w", err)
	}

	return VideoRangeInfo{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		VideoInfo: VideoFileInfo{
			VideoName: videoName,
			FileSize:  fileInfo.Size(),
			VideoFile: videoFile,
		},
	}, nil
}

func (s VideoManagementService) parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, 0, fmt.Errorf(errInvalidRangeFormat)
	}

	byteRanges := strings.Split(parts[1], "-")
	if len(byteRanges) != 2 {
		return 0, 0, fmt.Errorf(errInvalidBytesFormat)
	}

	start, err := strconv.ParseInt(byteRanges[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf(errInvalidRangeStart)
	}

	var end int64
	if byteRanges[1] == "" {
		end = min(fileSize-1, start+int64(float64(fileSize-1)*defaultRangePercentage))
	} else {
		end, err = strconv.ParseInt(byteRanges[1], 10, 64)
		if err != nil || end >= fileSize {
			return 0, 0, fmt.Errorf(errInvalidRangeEnd)
		}
	}

	if start > end || start < 0 || end >= fileSize {
		return 0, 0, fmt.Errorf(errInvalidRange)
	}

	return start, end, nil
}

func (s VideoManagementService) GetVideoToDownload(videoName string) (VideoFileInfo, error) {
	videoFile, err := os.Open("videos/" + videoName)
	if err != nil {
		return VideoFileInfo{}, fmt.Errorf(errVideoNotFound+": %w", err)
	}

	fileInfo, err := videoFile.Stat()
	if err != nil {
		return VideoFileInfo{}, fmt.Errorf(errGettingFileInfo+": %w", err)
	}

	return VideoFileInfo{
		VideoName: videoName,
		FileSize:  fileInfo.Size(),
		VideoFile: videoFile,
	}, nil
}
