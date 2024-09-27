package domain

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
)

// folder service
const (
	ErrCheckingFolder     = "error checking folder exist"
	ErrFolderAlreadyExist = "folder with this name already exist in this folder"
	ErrCreatingFolder     = "error creating folder"
	ErrFolderNotFound     = "folder with this ID not found"
	ErrUpdatingFolderName = "error updating folder name"
	ErrGettingFolderName  = "error getting folder name by id"
	ErrMovingFolder       = "error moving folder"
	ErrDeletingFolder     = "error deleting folder"
	ErrGettingFoldersList = "error getting folders list"
)
