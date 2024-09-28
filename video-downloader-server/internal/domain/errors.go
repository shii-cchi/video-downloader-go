package domain

import "errors"

// common service
var (
	ErrCreatingDir      = errors.New("error creating directory")
	ErrGeneratingBytes  = errors.New("error generating random bytes")
	ErrCreatingFile     = errors.New("error creating file for saving video")
	ErrSavingDataToFile = errors.New("error saving data to file")
)

// preview service
var (
	ErrGettingVideoDuration = errors.New("error getting video duration")
	ErrParsingVideoDuration = errors.New("error parsing video duration")
	ErrGeneratingPreview    = errors.New("error generating preview")
	ErrDeletingPreview      = errors.New("error deleting preview")
)

// general strategy
var (
	ErrSendingReq       = errors.New("error sending get request for downloading from general player")
	ErrDownloadingVideo = errors.New("error failed to download video from general player")
)

// youtube strategy
var (
	ErrParsingURL       = errors.New("failed to parse VideoURL")
	ErrNotFoundVideoID  = errors.New("videoID not found in url")
	ErrFetchingMetadata = errors.New("error fetching video metadata")
	ErrGettingStream    = errors.New("error getting stream")
	ErrMerging          = errors.New("error merging video and audio")
	ErrDeletingTmpFiles = errors.New("error geleting tmp files")
)

// videos service
var (
	ErrSavingVideoToDb      = errors.New("error saving video info to db")
	ErrGettingRealVideoPath = errors.New("error getting real video path by id")
	ErrVideoNotFound        = errors.New("video not found")
	ErrGettingFileInfo      = errors.New("err getting file info")
	ErrInvalidRangeHeader   = errors.New("invalid range header")
	ErrInvalidRangeFormat   = errors.New("invalid range format")
	ErrInvalidBytesFormat   = errors.New("invalid bytes format")
	ErrInvalidRangeStart    = errors.New("invalid start of range")
	ErrInvalidRangeEnd      = errors.New("invalid end of range")
	ErrInvalidRange         = errors.New("invalid range")
	ErrCheckingVideo        = errors.New("error checking video existence")
	ErrRenamingVideo        = errors.New("error renaming video")
	ErrMovingVideo          = errors.New("error moving video")
	ErrDeletingVideo        = errors.New("error deleting video")
	ErrDeletingVideoFromDB  = errors.New("error deleting video from db")
	ErrGettingPaths         = errors.New("error getting real videos and previews paths")
	ErrGettingVideos        = errors.New("errors getting videos by folder id")

	ErrVideoAlreadyExist = errors.New("video with this name already exist")
)

// folder service
var (
	ErrCheckingFolder           = errors.New("error checking folder exist")
	ErrFolderNotFound           = errors.New("folder with this ID not found")
	ErrFolderAlreadyExist       = errors.New("folder with this name already exist in this folder")
	ErrCreatingFolder           = errors.New("error creating folder")
	ErrRenamingFolder           = errors.New("error renaming folder")
	ErrGettingFolderName        = errors.New("error getting folder name by id")
	ErrMovingFolder             = errors.New("error moving folder")
	ErrGettingAllNestedFolders  = errors.New("error getting all nested folders")
	ErrDeletingAllNestedFolders = errors.New("error deleting all nested folders")
	ErrGettingNestedFolders     = errors.New("error getting nested folders")
)
