package delivery

type СontextKey string

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	DownloadInputKey     СontextKey = "downloadInput"
	VideoIDInputKey      СontextKey = "videoIDInput"
	CreateFolderInputKey СontextKey = "createFolderInput"
	RenameFolderInputKey СontextKey = "renameFolderInput"
	MoveFolderInputKey   СontextKey = "moveFolderInput"
	DeleteFolderInputKey СontextKey = "deleteFolderInput"

	ErrInvalidDownloadInput     = "invalid download input body"
	MesInvalidDownloadInput     = "fields url, folder_id and type are required and can't be empty, field url must be url format, field folder_id must be object id, field type can be 'general' or 'youtube', field quality can be empty or one of 2160p 1440p 1080p 720p 480p 360p 240p 144p best"
	MesInvalidJSON              = "invalid JSON body"
	ErrInvalidVideoIDInput      = "invalid video id input"
	MesInvalidVideoIDInput      = "video id is required"
	ErrInvalidCreateFolderInput = "invalid create folder input body"
	MesInvalidCreateFolderInput = "field folder_name is required, must be valid name and can't be empty, field parent_dir_id must be valid object id"
	ErrInvalidRenameFolderInput = "invalid rename folder input body"
	MesInvalidRenameFolderInput = "fields id and folder_name are required, id must be valid object id, folder_name must be valid name and can't be empty"
	ErrInvalidMoveFolderInput   = "invalid move folder input body"
	MesInvalidMoveFolderInput   = "fields id and parent_dir_id are required, must be valid object id and can't be empty"
	ErrInvalidDeleteFolderInput = "invalid delete folder input body"
	MesInvalidDeleteFolderInput = "field id are required, must be valid object id and can't be empty"

	ErrGettingID                = "error getting videoID from VideoURL"
	ErrDownloadingVideoToServer = "error downloading video to server"
	SuccessfulLoad              = "successful video load"

	ErrGettingVideoRange          = "error getting video range info"
	ErrGettingVideo               = "error getting video"
	ErrDownloadingVideoFromServer = "error downloading video from server"

	ErrMarshalingJSON = "failed to marshal JSON response"

	ErrCreatingFolder = "error creating new folder"
	ErrRenameFolder   = "error rename folder"
	ErrMovingFolder   = "error moving folder"
	ErrDeletingFolder = "error deleting folder"
)
