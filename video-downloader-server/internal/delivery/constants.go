package delivery

type contextKey string

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DownloadInputDto struct {
	VideoURL   string `json:"video_url" validate:"required,url"`
	FolderName string `json:"folder_name" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=youtube general"`
	Quality    string `json:"quality" validate:"omitempty,oneof=2160p 1440p 1080p 720p 480p 360p 240p 144p best"`
}

const (
	DownloadInputKey     contextKey = "downloadInput"
	VideoIDInputKey      contextKey = "videoIDInput"
	CreateFolderInputKey contextKey = "createFolderInput"

	ErrInvalidDownloadInput     = "invalid download input body"
	MesInvalidDownloadInput     = "fields url, folder_name and type are required and can't be empty, field url must be url format, field type can be 'general' or 'youtube', field quality can be empty or one of 2160p 1440p 1080p 720p 480p 360p 240p 144p best"
	MesInvalidJSON              = "invalid JSON body"
	ErrInvalidVideoIDInput      = "invalid video id input"
	MesInvalidVideoIDInput      = "video id is required"
	ErrInvalidCreateFolderInput = "invalid create folder input body"
	MesInvalidCreateFolderInput = "field folder_name is required, must be valid name and can't be empty, field parent_dir_id is required and must be valid objectid"

	ErrGettingID                = "error getting videoID from VideoURL"
	ErrDownloadingVideoToServer = "error downloading video to server"
	SuccessfulLoad              = "successful video load"

	ErrGettingVideoRange          = "error getting video range info"
	ErrGettingVideo               = "error getting video"
	ErrDownloadingVideoFromServer = "error downloading video from server"

	ErrMarshalingJSON = "failed to marshal JSON response"

	ErrCreatingFolder = "error creating new folder"
)
