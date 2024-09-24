package common

import (
	crypto "crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"video-downloader-server/internal/domain"
)

func CreateRandomDir(commonDir string) (string, error) {
	randBytes, err := generateRandomBytes(2)
	if err != nil {
		return "", err
	}

	firstRandDir := fmt.Sprintf("%02x", randBytes[0])
	secondRandDir := fmt.Sprintf("%02x", randBytes[1])

	fullDirPath := filepath.Join(commonDir, firstRandDir, secondRandDir)
	err = os.MkdirAll(fullDirPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(domain.ErrCreatingDir+": %w", err)
	}

	return filepath.Join(firstRandDir, secondRandDir), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	randBytes := make([]byte, length)
	_, err := crypto.Read(randBytes)
	if err != nil {
		return nil, fmt.Errorf(domain.ErrGeneratingBytes+": %w", err)
	}
	return randBytes, nil
}

func ReplaceSpecialSymbols(videoName string) string {
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	safeFileName := re.ReplaceAllString(videoName, "")

	return safeFileName
}

func CreateAndWriteFile(filePath string, data io.ReadCloser) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf(domain.ErrCreatingFile+": %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf(domain.ErrSavingDataToFile+": %w", err)
	}

	return nil
}
