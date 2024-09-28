package domain

import "errors"

// common service
const (
	ErrCreatingDir     = "error creating directory"
	ErrGeneratingBytes = "error generating random bytes"

	ErrCreatingFile     = "error creating file for saving video"
	ErrSavingDataToFile = "error saving data to file"
)

// preview service
const (
	ErrGettingVideoDuration = "error getting video duration"
	ErrParsingVideoDuration = "error parsing video duration"
	ErrGeneratingPreview    = "error generating preview"
)

// general strategy
const (
	ErrSendingReq       = "error sending get request by downloading from general player"
	ErrDownloadingVideo = "error failed to download video from general player with status code:"
)

// youtube strategy
const (
	ErrParsingURL       = "failed to parse VideoURL"
	ErrNotFoundVideoID  = "videoID not found in url"
	ErrFetchingMetadata = "error fetching video metadata"
	ErrGettingStream    = "error getting stream"
	ErrMerging          = "error merging video and audio"
)

// videos service
const (
	ErrVideoNotFound      = "video not found"
	ErrGettingFileInfo    = "err getting file info"
	ErrInvalidRangeHeader = "invalid range header"
	ErrInvalidRangeFormat = "invalid range format"
	ErrInvalidBytesFormat = "invalid bytes format"
	ErrInvalidRangeStart  = "invalid start of range"
	ErrInvalidRangeEnd    = "invalid end of range"
	ErrInvalidRange       = "invalid range"
	ErrGettingPaths       = "error getting real videos and previews paths"
	ErrDeletingVideo      = "error deleting video with path"
	ErrDeletingPreview    = "error deleting preview with path"
	ErrVideoAlreadyExist  = "video with this name already exist"
	ErrCheckingVideo      = "error checking video existence"
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
	ErrConvertingToObjectID     = errors.New("error converting str to object id")
	ErrGettingNestedFolders     = errors.New("error getting nested folders")
)
