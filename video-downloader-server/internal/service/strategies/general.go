package strategies

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"video-downloader-server/internal/domain"
	"video-downloader-server/internal/service/common"
)

type GeneralDownloadStrategy struct{}

func (s GeneralDownloadStrategy) Download(videoURL string, quality string) (string, string, error) {
	res, err := http.Get(videoURL)
	if err != nil {
		return "", "", fmt.Errorf("%w (video url: %s): %s", domain.ErrSendingReq, videoURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("%w (video url: %s, status code: %d)", domain.ErrDownloadingVideo, videoURL, res.StatusCode)
	}

	realPath, err := common.CreateRandomDir(domain.CommonVideoDir)
	if err != nil {
		return "", "", err
	}

	videoName := common.ReplaceSpecialSymbols(filepath.Base(videoURL))
	filePath := filepath.Join(domain.CommonVideoDir, realPath, videoName)

	if err := common.CreateAndWriteFile(filePath, res.Body); err != nil {
		return "", "", err
	}

	return strings.TrimSuffix(videoName, filepath.Ext(videoName)), filepath.Join(realPath, videoName), nil
}
