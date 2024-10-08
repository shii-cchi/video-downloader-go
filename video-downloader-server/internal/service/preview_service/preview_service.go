package preview_service

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"video-downloader-server/internal/domain"
	"video-downloader-server/internal/service/common"
)

type PreviewService struct {
}

func NewPreviewService() *PreviewService {
	return &PreviewService{}
}

func (p *PreviewService) CreatePreview(videoName string, realPath string) (string, error) {
	previewDir, err := common.CreateRandomDir(domain.CommonPreviewDir)
	if err != nil {
		return "", err
	}

	videoPath := filepath.Join(domain.CommonVideoDir, realPath)
	videoName = common.ReplaceSpecialSymbols(videoName)
	previewPath := filepath.Join(domain.CommonPreviewDir, previewDir, videoName+domain.PreviewFormat)

	videoDuration, err := p.getVideoDuration(videoPath)
	if err != nil {
		return "", err
	}

	previewTime := p.generateRandomTime(videoDuration)

	if err := p.generatePreview(videoPath, previewPath, previewTime); err != nil {
		return "", err
	}

	return filepath.Join(previewDir, videoName+domain.PreviewFormat), nil
}

func (p *PreviewService) DeletePreviews(paths []string) error {
	for _, previewPath := range paths {
		if err := os.Remove(filepath.Join(domain.CommonPreviewDir, previewPath)); err != nil {
			return fmt.Errorf("%w (preview path: %s): %s", domain.ErrDeletingPreview, previewPath, err)
		}
	}

	return nil
}

func (p *PreviewService) getVideoDuration(videoPath string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("%w (video path: %s): %s", domain.ErrGettingVideoDuration, videoPath, err)
	}

	durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("%w (video path: %s): %s", domain.ErrParsingVideoDuration, videoPath, err)
	}

	return time.Duration(durationSeconds) * time.Second, nil
}

func (p *PreviewService) generateRandomTime(videoDuration time.Duration) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFraction := domain.MinTimeFraction + rng.Float64()*domain.MaxTimeFraction
	randomSeconds := int64(videoDuration.Seconds() * randomFraction)
	randomTime := time.Duration(randomSeconds) * time.Second
	return fmt.Sprintf("%d", int(randomTime.Seconds()))
}

func (p *PreviewService) generatePreview(videoPath, previewPath, previewTime string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", previewTime, "-vframes", "1", previewPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w (ffmpeg output: %s): %s", domain.ErrGeneratingPreview, string(output), err)
	}

	return nil
}
