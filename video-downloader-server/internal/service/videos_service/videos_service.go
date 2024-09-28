package videos_service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	GetRealPath(ctx context.Context, videoID primitive.ObjectID) (string, error)
	CheckExistByID(ctx context.Context, videoID primitive.ObjectID) error
	Rename(ctx context.Context, videoID primitive.ObjectID, newVideoName string) error
	Move(ctx context.Context, videoID primitive.ObjectID, folderID primitive.ObjectID) error
	Delete(ctx context.Context, videoID primitive.ObjectID) error
	GetPathsByFolders(ctx context.Context, foldersID []primitive.ObjectID) ([]string, []string, error)
	DeleteVideos(ctx context.Context, foldersID []primitive.ObjectID) error
	GetVideos(ctx context.Context, folderID primitive.ObjectID) ([]domain.Video, error)
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

func (v *VideosService) DownloadToServer(downloadVideoInput video_dto.DownloadVideoDto) error {
	switch downloadVideoInput.Type {
	case domain.YouTubeVideoType:
		v.setVideoDownloadStrategy(strategies.YouTubeDownloadStrategy{})
	default:
		v.setVideoDownloadStrategy(strategies.GeneralDownloadStrategy{})
	}

	videoName, realPath, err := v.strategy.Download(downloadVideoInput.VideoURL, downloadVideoInput.Quality)
	if err != nil {
		return err
	}

	previewPath, err := v.previewService.CreatePreview(videoName, realPath)
	if err != nil {
		return err
	}

	if err := v.repo.Create(context.Background(), domain.Video{
		VideoName:   videoName,
		FolderID:    downloadVideoInput.FolderID,
		RealPath:    realPath,
		PreviewPath: previewPath,
	}); err != nil {
		return fmt.Errorf("%w (video name: %s): %s", domain.ErrSavingVideoToDb, videoName, err)
	}

	return nil
}

func (v *VideosService) GetVideoFileInfo(videoID primitive.ObjectID) (video_dto.VideoFileInfoDto, error) {
	videoRealPath, err := v.repo.GetRealPath(context.Background(), videoID)
	if err != nil {
		return video_dto.VideoFileInfoDto{}, fmt.Errorf("%w (video id: %s): %s", domain.ErrGettingRealVideoPath, videoID, err)
	}

	videoFile, err := os.Open(filepath.Join(domain.CommonVideoDir, videoRealPath))
	if err != nil {
		return video_dto.VideoFileInfoDto{}, fmt.Errorf("%w (video path: %s): %s", domain.ErrVideoNotFound, videoRealPath, err)
	}

	fileInfo, err := videoFile.Stat()
	if err != nil {
		return video_dto.VideoFileInfoDto{}, fmt.Errorf("%w (video path: %s): %s", domain.ErrGettingFileInfo, videoRealPath, err)
	}

	_, videoName := filepath.Split(videoRealPath)

	return video_dto.VideoFileInfoDto{
		VideoName: videoName,
		FileSize:  fileInfo.Size(),
		VideoFile: videoFile,
	}, nil
}

func (v *VideosService) GetVideoRangeInfo(videoID primitive.ObjectID, rangeHeader string) (video_dto.VideoRangeInfoDto, error) {
	videoFileInfo, err := v.GetVideoFileInfo(videoID)

	rangeStart, rangeEnd, err := v.parseRangeHeader(rangeHeader, videoFileInfo.FileSize)
	if err != nil {
		return video_dto.VideoRangeInfoDto{}, fmt.Errorf("%w: %s", domain.ErrInvalidRangeHeader, err)
	}

	return video_dto.VideoRangeInfoDto{
		RangeStart: rangeStart,
		RangeEnd:   rangeEnd,
		VideoInfo:  videoFileInfo,
	}, nil
}

func (v *VideosService) Rename(renameVideoInput video_dto.RenameVideoDto) (video_dto.VideoDto, error) {
	if err := v.checkVideoExistenceByID(renameVideoInput.ID); err != nil {
		return video_dto.VideoDto{}, err
	}

	if err := v.repo.Rename(context.Background(), renameVideoInput.ID, renameVideoInput.VideoName); err != nil {
		return video_dto.VideoDto{}, fmt.Errorf("%w (video id: %s): %s", domain.ErrRenamingVideo, renameVideoInput.ID, err)
	}

	return video_dto.VideoDto{
		ID:        renameVideoInput.ID,
		VideoName: renameVideoInput.VideoName,
	}, nil
}

func (v *VideosService) Move(moveVideoInput video_dto.MoveVideoDto) (video_dto.VideoDto, error) {
	if err := v.checkVideoExistenceByID(moveVideoInput.ID); err != nil {
		return video_dto.VideoDto{}, err
	}

	if err := v.repo.Move(context.Background(), moveVideoInput.ID, moveVideoInput.FolderID); err != nil {
		return video_dto.VideoDto{}, fmt.Errorf("%w (video id: %s): %s", domain.ErrMovingVideo, moveVideoInput.ID, err)
	}

	return video_dto.VideoDto{
		ID:       moveVideoInput.ID,
		FolderID: moveVideoInput.FolderID,
	}, nil
}

func (v *VideosService) Delete(deleteVideoInput video_dto.DeleteVideoDto) error {
	if err := v.checkVideoExistenceByID(deleteVideoInput.ID); err != nil {
		return err
	}

	realPath, err := v.repo.GetRealPath(context.Background(), deleteVideoInput.ID)
	if err != nil {
		return fmt.Errorf("%w (video id: %s): %s", domain.ErrGettingRealVideoPath, deleteVideoInput.ID, err)
	}

	if err := os.Remove(filepath.Join(domain.CommonVideoDir, realPath)); err != nil {
		return fmt.Errorf("%w (video path: %s): %s", domain.ErrDeletingVideo, realPath, err)
	}

	if err = v.repo.Delete(context.Background(), deleteVideoInput.ID); err != nil {
		return fmt.Errorf("%w (video id: %s): %s", domain.ErrDeletingVideoFromDB, deleteVideoInput.ID, err)
	}

	return nil
}

func (v *VideosService) DeleteVideos(foldersID []primitive.ObjectID) error {
	realPaths, previewPaths, err := v.repo.GetPathsByFolders(context.Background(), foldersID)
	if err != nil {
		return fmt.Errorf("%w (folders id: %s): %s", domain.ErrGettingPaths, foldersID, err)
	}

	if err = v.repo.DeleteVideos(context.Background(), foldersID); err != nil {
		return fmt.Errorf("%w (from folders: %s): %s", domain.ErrDeletingVideoFromDB, foldersID, err)
	}

	for _, realPath := range realPaths {
		if err := os.Remove(filepath.Join(domain.CommonVideoDir, realPath)); err != nil {
			return fmt.Errorf("%w (video path: %s): %s", domain.ErrDeletingVideo, realPath, err)
		}
	}

	if err := v.previewService.DeletePreviews(previewPaths); err != nil {
		return err
	}

	return nil
}

func (v *VideosService) GetVideos(folderID primitive.ObjectID) ([]video_dto.VideoDto, error) {
	videos, err := v.repo.GetVideos(context.Background(), folderID)
	if err != nil {
		return nil, fmt.Errorf("%w (folder id: %s): %s", domain.ErrGettingVideos, folderID, err)
	}

	return v.toVideoDto(videos), nil
}

func (v *VideosService) parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, 0, domain.ErrInvalidRangeFormat
	}

	byteRanges := strings.Split(parts[1], "-")
	if len(byteRanges) != 2 {
		return 0, 0, domain.ErrInvalidBytesFormat
	}

	start, err := strconv.ParseInt(byteRanges[0], 10, 64)
	if err != nil {
		return 0, 0, domain.ErrInvalidRangeStart
	}

	var end int64
	if byteRanges[1] == "" {
		end = min(fileSize-1, start+int64(float64(fileSize-1)*domain.DefaultRangePercentage))
	} else {
		end, err = strconv.ParseInt(byteRanges[1], 10, 64)
		if err != nil || end >= fileSize {
			return 0, 0, domain.ErrInvalidRangeEnd
		}
	}

	if start > end || start < 0 || end >= fileSize {
		return 0, 0, domain.ErrInvalidRange
	}

	return start, end, nil
}

func (v *VideosService) toVideoDto(videos []domain.Video) []video_dto.VideoDto {
	res := make([]video_dto.VideoDto, len(videos))

	for i, video := range videos {
		res[i] = video_dto.VideoDto{
			ID:          video.ID,
			VideoName:   video.VideoName,
			FolderID:    video.FolderID,
			RealPath:    video.RealPath,
			PreviewPath: video.PreviewPath,
		}
	}

	return res
}

func (v *VideosService) checkVideoExistenceByID(videoID primitive.ObjectID) error {
	err := v.repo.CheckExistByID(context.Background(), videoID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("%w (video id: %s): %s", domain.ErrVideoNotFound, videoID, err)
		}

		return fmt.Errorf("%w (video id: %s): %s", domain.ErrCheckingVideo, videoID, err)
	}

	return nil
}
