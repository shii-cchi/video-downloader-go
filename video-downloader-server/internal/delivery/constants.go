package delivery

type contextKey string

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	DownloadInputKey  contextKey = "downloadInput"
	VideoNameInputKey contextKey = "videoNameInput"

	ErrInvalidDownloadInput  = "invalid download input body"
	MesInvalidDownloadInput  = "fields url and type are required and can't be empty, field url must be url format, field type can be 'general' or 'youtube'"
	MesInvalidJSON           = "invalid JSON body"
	ErrInvalidVideoNameInput = "invalid video name input"
	MesInvalidVideoNameInput = "video name is required"

	ErrNotFoundVideoID  = "videoID not found in url"
	ErrParsingURL       = "failed to parse VideoURL"
	ErrGettingID        = "error getting videoID from VideoURL"
	ErrDownloadingVideo = "error downloading video"
	SuccessfulLoad      = "successful video load"

	ErrMarshalingJSON = "failed to marshal JSON response"

	ErrGettingVideoRange  = "error getting video range info"
	ErrVideoNotFound      = "video not found"
	ErrGettingFileInfo    = "err getting file info"
	ErrInvalidRangeHeader = "invalid range header"
)
