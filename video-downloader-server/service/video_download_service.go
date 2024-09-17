package service

import (
	"context"
	crypto "crypto/rand"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/internal/repository"
)

const (
	commonVideoDir   = "videos"
	commonPreviewDir = "previews"

	videoFormat   = ".mp4"
	previewFormat = ".jpeg"

	errNotFoundVideoID      = "videoID not found in url"
	errParsingURL           = "failed to parse VideoURL"
	errFetchingMetadata     = "error fetching video metadata"
	errGettingStream        = "error getting stream"
	errCreatingFile         = "error creating file for saving video"
	errFillingFile          = "error filling file with video stream"
	errMerging              = "error merging video and audio"
	errSendingReq           = "error sending get request"
	errDownloadingVideo     = "error failed to download video with status code"
	errGettingVideoDuration = "error getting video duration"
	errParsingVideoDuration = "error parsing video duration"
	errGeneratingPreview    = "error generating preview"
	errCreatingDir          = "error creating directory"
	errGeneratingBytes      = "error generating random bytes"
)

type VideoDownloadService struct {
	repo repository.VideoDownload
}

func newVideoDownloadService(repo repository.VideoDownload) *VideoDownloadService {
	return &VideoDownloadService{
		repo: repo,
	}
}

func (v *VideoDownloadService) Download(input dto.DownloadInputDto) error {
	videoName, realPath, err := v.downloadVideo(input)
	if err != nil {
		return err
	}

	previewPath, err := v.createPreview(videoName, realPath)
	if err != nil {
		return err
	}

	return v.repo.Create(context.Background(), repository.CreateParams{
		VideoName:   videoName,
		RealPath:    realPath,
		UserPath:    input.FolderName,
		PreviewPath: previewPath,
	})
}

func (v *VideoDownloadService) downloadVideo(input dto.DownloadInputDto) (string, string, error) {
	switch input.Type {
	case "youtube":
		return v.downloadFromYouTube(input.VideoURL, input.Quality)
	default:
		return v.downloadFromOther(input.VideoURL)
	}
}

func (v *VideoDownloadService) downloadFromYouTube(videoURL string, quality string) (string, string, error) {
	videoID, err := v.getVideoID(videoURL)
	if err != nil {
		return "", "", err
	}

	video, err := v.fetchVideoMetadata(videoID)
	if err != nil {
		return "", "", err
	}

	videoName := v.replaceSpecialSymbols(video.Title)

	videoPath, audioPath, format, err := v.downloadAndPrepareFiles(video, quality, videoName)
	if err != nil {
		return "", "", err
	}
	defer v.deleteTmpFiles(videoPath, audioPath)

	realPath, err := v.createRandomDir(commonVideoDir)
	if err != nil {
		return "", "", err
	}

	mergedFilePath := filepath.Join(commonVideoDir, realPath, fmt.Sprintf("%s %s%s", videoName, format.QualityLabel, videoFormat))
	if err := v.mergeVideoAudio(videoPath, audioPath, mergedFilePath); err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s %s", videoName, format.QualityLabel), filepath.Join(realPath, fmt.Sprintf("%s %s%s", videoName, format.QualityLabel, videoFormat)), nil
}

func (v *VideoDownloadService) getVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf(errParsingURL+": %w", err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		return "", fmt.Errorf(errNotFoundVideoID+": %w", err)
	}

	return videoID, nil
}

func (v *VideoDownloadService) fetchVideoMetadata(videoID string) (*youtube.Video, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, fmt.Errorf(errFetchingMetadata+": %w", err)
	}

	return video, nil
}

func (v *VideoDownloadService) replaceSpecialSymbols(videoName string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	safeFileName := re.ReplaceAllString(videoName, "")

	return safeFileName
}

func (v *VideoDownloadService) downloadAndPrepareFiles(video *youtube.Video, quality string, videoName string) (string, string, *youtube.Format, error) {
	selectedVideoFormat := v.selectVideoFormat(video, quality)
	videoPath := filepath.Join(commonVideoDir, fmt.Sprintf("%s_video_%s%s", videoName, selectedVideoFormat.QualityLabel, videoFormat))
	if err := v.downloadStreamToFile(video, selectedVideoFormat, videoPath); err != nil {
		return "", "", nil, err
	}

	selectedAudioFormat := v.selectAudioFormat(video)
	audioPath := filepath.Join(commonVideoDir, fmt.Sprintf("%s_audio%s", videoName, videoFormat))
	if err := v.downloadStreamToFile(video, selectedAudioFormat, audioPath); err != nil {
		return "", "", nil, err
	}

	return videoPath, audioPath, selectedVideoFormat, nil
}

func (v *VideoDownloadService) selectVideoFormat(video *youtube.Video, quality string) *youtube.Format {
	formats := video.Formats.Type("video")
	for _, format := range formats {
		if quality == "best" || format.QualityLabel == quality {
			return &format
		}
	}

	return &formats[0]
}

func (v *VideoDownloadService) selectAudioFormat(video *youtube.Video) *youtube.Format {
	formats := video.Formats.Type("audio")
	for _, format := range formats {
		if strings.Contains(format.MimeType, "audio/mp4") {
			return &format
		}
	}

	return &formats[0]
}

func (v *VideoDownloadService) downloadStreamToFile(video *youtube.Video, format *youtube.Format, fileName string) error {
	client := youtube.Client{}

	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf(errGettingStream+": %w", err)
	}
	defer stream.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf(errCreatingFile+": %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf(errFillingFile+": %w", err)
	}

	return nil
}

func (v *VideoDownloadService) deleteTmpFiles(videoPath, audioPath string) {
	os.Remove(videoPath)
	os.Remove(audioPath)
}

func (v *VideoDownloadService) mergeVideoAudio(videoFileName string, audioFileName string, mergedFileName string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFileName, "-i", audioFileName, "-c", "copy", mergedFileName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errMerging+": %w", err)
	}

	return nil
}

func (v *VideoDownloadService) downloadFromOther(videoURL string) (string, string, error) {
	resp, err := http.Get(videoURL)
	if err != nil {
		return "", "", fmt.Errorf(errSendingReq+": %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf(errDownloadingVideo+": %d", resp.StatusCode)
	}

	realPath, err := v.createRandomDir(commonVideoDir)
	if err != nil {
		return "", "", err
	}

	videoName := v.replaceSpecialSymbols(filepath.Base(videoURL))
	filePath := filepath.Join(commonVideoDir, realPath, videoName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf(errCreatingFile+": %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", "", fmt.Errorf(errFillingFile+": %w", err)
	}

	return videoName, filepath.Join(realPath, videoName), nil
}

func (v *VideoDownloadService) createPreview(videoName string, realPath string) (string, error) {
	previewDir, err := v.createRandomDir(commonPreviewDir)
	if err != nil {
		return "", err
	}

	videoPath := filepath.Join(commonVideoDir, realPath)
	previewPath := filepath.Join(commonPreviewDir, previewDir, videoName+previewFormat)

	videoDuration, err := v.getVideoDuration(videoPath)
	if err != nil {
		return "", err
	}

	previewTime := v.generateRandomTime(videoDuration)

	if err := v.generatePreview(videoPath, previewPath, previewTime); err != nil {
		return "", err
	}

	return filepath.Join(previewDir, videoName+previewFormat), nil
}

func (v *VideoDownloadService) getVideoDuration(videoPath string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf(errGettingVideoDuration+": %w", err)
	}

	durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf(errParsingVideoDuration+": %w", err)
	}

	return time.Duration(durationSeconds) * time.Second, nil
}

func (v *VideoDownloadService) generateRandomTime(videoDuration time.Duration) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFraction := 0.1 + rng.Float64()*0.8
	randomSeconds := int64(videoDuration.Seconds() * randomFraction)
	randomTime := time.Duration(randomSeconds) * time.Second
	return fmt.Sprintf("%d", int(randomTime.Seconds()))
}

func (v *VideoDownloadService) generatePreview(videoPath, previewPath, previewTime string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", previewTime, "-vframes", "1", previewPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(errGeneratingPreview+":%w\n%s", err, string(output))
	}

	return nil
}

func (v *VideoDownloadService) createRandomDir(commonDir string) (string, error) {
	randBytes, err := v.generateRandomBytes(2)
	if err != nil {
		return "", err
	}

	firstRandDir := fmt.Sprintf("%02x", randBytes[0])
	secondRandDir := fmt.Sprintf("%02x", randBytes[1])

	fullDirPath := filepath.Join(commonDir, firstRandDir, secondRandDir)
	err = os.MkdirAll(fullDirPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(errCreatingDir+": %w", err)
	}

	return filepath.Join(firstRandDir, secondRandDir), nil
}

func (v *VideoDownloadService) generateRandomBytes(length int) ([]byte, error) {
	randBytes := make([]byte, length)
	_, err := crypto.Read(randBytes)
	if err != nil {
		return nil, fmt.Errorf(errGeneratingBytes+": %w", err)
	}
	return randBytes, nil
}
