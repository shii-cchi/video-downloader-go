package service

import (
	"fmt"
	"github.com/kkdai/youtube/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"video-downloader-server/internal/delivery/dto"
)

const (
	errNotFoundVideoID  = "videoID not found in url"
	errParsingURL       = "failed to parse VideoURL"
	errFetchingMetadata = "error fetching video metadata"
	errGettingStream    = "error getting stream"
	errCreatingFile     = "error creating file for saving video"
	errFillingFile      = "error filling file with video stream"
	errMerging          = "error merging video and audio"
	errSendingReq       = "error sending get request"
	errDownloadingVideo = "error failed to download video with status code"
)

type VideoService struct {
}

func newVideoService() *VideoService {
	return &VideoService{}
}

func (v *VideoService) DownloadVideo(input dto.DownloadInputDto) error {
	if input.Type == "youtube" {
		return v.downloadYouTubeVideo(input.VideoURL, input.Quality)
	}

	return v.downloadGeneralVideo(input.VideoURL)
}

func (v *VideoService) downloadYouTubeVideo(videoURL string, quality string) error {
	videoID, err := v.getVideoID(videoURL)
	if err != nil {
		return err
	}

	client := youtube.Client{}
	video, err := client.GetVideo(videoID)
	if err != nil {
		return fmt.Errorf(errFetchingMetadata+": %s", err)
	}

	videoTitle := v.replaceSpecialSymbols(video.Title)

	selectedVideoFormat := v.selectVideoFormat(video, quality)
	videoFileName := fmt.Sprintf("videos/%s_video_%s.mp4", videoTitle, selectedVideoFormat.QualityLabel)
	if err := v.downloadStreamToFile(client, video, selectedVideoFormat, videoFileName); err != nil {
		return err
	}
	defer os.Remove(videoFileName)

	selectedAudioFormat := v.selectAudioFormat(video)
	audioFileName := fmt.Sprintf("videos/%s_audio.mp4", videoTitle)
	if err := v.downloadStreamToFile(client, video, selectedAudioFormat, audioFileName); err != nil {
		return err
	}
	defer os.Remove(audioFileName)

	mergedFileName := fmt.Sprintf("videos/%s %s.mp4", videoTitle, selectedVideoFormat.QualityLabel)
	if err := v.mergeVideoAudio(videoFileName, audioFileName, mergedFileName); err != nil {
		return err
	}

	return nil
}

func (v *VideoService) downloadGeneralVideo(videoURL string) error {
	resp, err := http.Get(videoURL)
	if err != nil {
		return fmt.Errorf(errSendingReq+": %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(errDownloadingVideo+": %d", resp.StatusCode)
	}

	pathParts := strings.Split(videoURL, "/")
	videoName := pathParts[len(pathParts)-1]
	videoName = v.replaceSpecialSymbols(videoName)

	file, err := os.Create("videos/" + videoName)
	if err != nil {
		return fmt.Errorf(errCreatingFile+": %s", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf(errFillingFile+": %s", err)
	}

	return nil
}

func (v *VideoService) getVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf(errParsingURL+": %s", err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		return "", fmt.Errorf(errNotFoundVideoID+": %s", err)
	}

	return videoID, nil
}

func (v *VideoService) replaceSpecialSymbols(videoName string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	safeFileName := re.ReplaceAllString(videoName, "")

	return safeFileName
}

func (v *VideoService) selectVideoFormat(video *youtube.Video, quality string) *youtube.Format {
	videoFormats := video.Formats.Type("video")
	for _, videoFormat := range videoFormats {
		if quality == "best" || videoFormat.QualityLabel == quality {
			return &videoFormat
		}
	}

	return &videoFormats[0]
}

func (v *VideoService) selectAudioFormat(video *youtube.Video) *youtube.Format {
	audioFormats := video.Formats.Type("audio")
	for _, audioFormat := range audioFormats {
		if strings.Contains(audioFormat.MimeType, "audio/mp4") {
			return &audioFormat
		}
	}

	return &audioFormats[0]
}

func (v *VideoService) downloadStreamToFile(client youtube.Client, video *youtube.Video, format *youtube.Format, fileName string) error {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf(errGettingStream+": %s", err)
	}
	defer stream.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf(errCreatingFile+": %s", err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf(errFillingFile+": %s", err)
	}

	return nil
}

func (v *VideoService) mergeVideoAudio(videoFileName string, audioFileName string, mergedFileName string) error {
	cmd := exec.Command("ffmpeg", "-i", videoFileName, "-i", audioFileName, "-c", "copy", mergedFileName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(errMerging+": %s", err)
	}

	return nil
}
