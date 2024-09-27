package videos_service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"video-downloader-server/internal/delivery/dto/video_dto"
	"video-downloader-server/internal/domain"
	"video-downloader-server/internal/service/strategies"
)

type VideosRepo interface {
	Create(ctx context.Context, video domain.Video) error
	GetPathsByFolders(ctx context.Context, foldersID []primitive.ObjectID) ([]string, []string, error)
	DeleteVideos(ctx context.Context, foldersID []primitive.ObjectID) error

	//Delete(ctx context.Context, folderID primitive.ObjectID) error
}

type Preview interface {
	CreatePreview(videoName string, realPath string) (string, error)
	DeletePreviews(paths []string) error
}

type VideoDownloadStrategy interface {
	Download(videoURL string, quality string) (string, string, error)
}

type VideosService struct {
	repo           VideosRepo
	previewService Preview
	strategy       VideoDownloadStrategy
}

func NewVideosService(repo VideosRepo, previewService Preview) *VideosService {
	return &VideosService{
		repo:           repo,
		previewService: previewService,
	}
}

func (v *VideosService) setVideoDownloadStrategy(strategy VideoDownloadStrategy) {
	v.strategy = strategy
}

func (v *VideosService) DownloadToServer(input video_dto.DownloadVideoDto) error {
	switch input.Type {
	case domain.YouTubeVideoType:
		v.setVideoDownloadStrategy(strategies.YouTubeDownloadStrategy{})
	default:
		v.setVideoDownloadStrategy(strategies.GeneralDownloadStrategy{})
	}

	videoName, realPath, err := v.strategy.Download(input.VideoURL, input.Quality)
	if err != nil {
		return err
	}

	previewPath, err := v.previewService.CreatePreview(videoName, realPath)
	if err != nil {
		return err
	}

	return v.repo.Create(context.Background(), domain.Video{
		VideoName:   videoName,
		FolderID:    input.FolderID,
		RealPath:    realPath,
		PreviewPath: previewPath,
	})
}

func (v *VideosService) GetVideoFileInfo(videoID string) (domain.VideoFileInfo, error) {
	videoFile, err := os.Open(filepath.Join(domain.CommonVideoDir, videoID))
	if err != nil {
		return domain.VideoFileInfo{}, fmt.Errorf(domain.ErrVideoNotFound+": %w", err)
	}

	fileInfo, err := videoFile.Stat()
	if err != nil {
		return domain.VideoFileInfo{}, fmt.Errorf(domain.ErrGettingFileInfo+": %w", err)
	}

	return domain.VideoFileInfo{
		VideoName: videoID,
		FileSize:  fileInfo.Size(),
		VideoFile: videoFile,
	}, nil
}

func (v *VideosService) GetVideoRangeInfo(videoID string, rangeHeader string) (domain.VideoRangeInfo, error) {
	videoFileInfo, err := v.GetVideoFileInfo(videoID)

	rangeStart, rangeEnd, err := v.parseRangeHeader(rangeHeader, videoFileInfo.FileSize)
	if err != nil {
		return domain.VideoRangeInfo{}, fmt.Errorf(domain.ErrInvalidRangeHeader+": %w", err)
	}

	return domain.VideoRangeInfo{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		VideoInfo:  videoFileInfo,
	}, nil
}

func (v *VideosService) DeleteVideos(foldersID []primitive.ObjectID) error {
	realPaths, previewPaths, err := v.repo.GetPathsByFolders(context.Background(), foldersID)
	if err != nil {
		return fmt.Errorf(domain.ErrGettingPaths+": %w", err)
	}

	err = v.repo.DeleteVideos(context.Background(), foldersID)
	if err != nil {
		return fmt.Errorf(domain.ErrDeletingVideo+": %w", err)
	}

	for _, realPath := range realPaths {
		if err := os.Remove(filepath.Join(domain.CommonVideoDir, realPath)); err != nil {
			return fmt.Errorf(domain.ErrDeletingVideo+" %s: %w", realPath, err)
		}
	}

	if err := v.previewService.DeletePreviews(previewPaths); err != nil {
		return err
	}

	return nil
}

func (v *VideosService) parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, 0, fmt.Errorf(domain.ErrInvalidRangeFormat)
	}

	byteRanges := strings.Split(parts[1], "-")
	if len(byteRanges) != 2 {
		return 0, 0, fmt.Errorf(domain.ErrInvalidBytesFormat)
	}

	start, err := strconv.ParseInt(byteRanges[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf(domain.ErrInvalidRangeStart)
	}

	var end int64
	if byteRanges[1] == "" {
		end = min(fileSize-1, start+int64(float64(fileSize-1)*domain.DefaultRangePercentage))
	} else {
		end, err = strconv.ParseInt(byteRanges[1], 10, 64)
		if err != nil || end >= fileSize {
			return 0, 0, fmt.Errorf(domain.ErrInvalidRangeEnd)
		}
	}

	if start > end || start < 0 || end >= fileSize {
		return 0, 0, fmt.Errorf(domain.ErrInvalidRange)
	}

	return start, end, nil
}
