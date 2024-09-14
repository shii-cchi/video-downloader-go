package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

type VideoManagementService struct {
}

func newVideoManagementService() *VideoManagementService {
	return &VideoManagementService{}
}

type VideoRangeInfo struct {
	RangeStart int64
	RangeEnd   int64
	FileSize   int64
	VideoFile  *os.File
}

func (s VideoManagementService) GetVideoRange(videoName string, rangeHeader string) (VideoRangeInfo, error) {
	videoFile, err := os.Open("videos/" + videoName)
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errVideoNotFound+": %s", err)
	}

	fileStat, err := videoFile.Stat()
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errGettingFileInfo+": %s", err)
	}

	rangeStart, rangeEnd, err := s.parseRangeHeader(rangeHeader, fileStat.Size())
	if err != nil {
		return VideoRangeInfo{}, fmt.Errorf(errInvalidRangeHeader+": %s", err)
	}

	return VideoRangeInfo{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		FileSize:   fileStat.Size(),
		VideoFile:  videoFile,
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
