package delivery

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ContextKey string

const (
	DownloadVideoInputKey ContextKey = "downloadVideoInput"
	RenameVideoInputKey   ContextKey = "renameVideoInput"
	MoveVideoInputKey     ContextKey = "modeVideoInput"
	DeleteVideoInputKey   ContextKey = "deleteVideoInput"
	CreateFolderInputKey  ContextKey = "createFolderInput"
	RenameFolderInputKey  ContextKey = "renameFolderInput"
	MoveFolderInputKey    ContextKey = "moveFolderInput"
	DeleteFolderInputKey  ContextKey = "deleteFolderInput"
	VideoIDInputKey       ContextKey = "videoIDInput"
	FolderIDInputKey      ContextKey = "folderIDInput"
)

const (
	ErrInvalidDownloadVideoInput = "invalid download video input body"
	MesInvalidDownloadVideoInput = "fields video_url, type and folder_id are required and can't be empty, field video_url must be url format, field type can be 'general' or 'youtube', field quality can be empty or one of 2160p 1440p 1080p 720p 480p 360p 240p 144p best, field folder_id must be object id"
	ErrInvalidRenameVideoInput   = "invalid rename video input body"
	MesInvalidRenameVideoInput   = "fields id and video_name are required and can't be empty, id must be valid object id, video_name must be valid name"
	ErrInvalidMoveVideoInput     = "invalid move video input body"
	MesInvalidMoveVideoInput     = "fields id and folder_id are required, can't be empty and must be valid object id"
	ErrInvalidDeleteVideoInput   = "invalid delete video input body"
	MesInvalidDeleteVideoInput   = "field id are required, can't be empty and must be valid object id"
	ErrInvalidCreateFolderInput  = "invalid create folder input body"
	MesInvalidCreateFolderInput  = "field folder_name is required, can't be empty and must be valid name, field parent_dir_id must be valid object id"
	ErrInvalidRenameFolderInput  = "invalid rename folder input body"
	MesInvalidRenameFolderInput  = "fields id and folder_name are required and can't be empty, id must be valid object id, folder_name must be valid name"
	ErrInvalidMoveFolderInput    = "invalid move folder input body"
	MesInvalidMoveFolderInput    = "fields id and parent_dir_id are required, can't be empty and must be valid object id"
	ErrInvalidDeleteFolderInput  = "invalid delete folder input body"
	MesInvalidDeleteFolderInput  = "field id are required, can't be empty and must be valid object id"
	ErrInvalidVideoIDInput       = "invalid video id input"
	MesInvalidVideoIDInput       = "video id param must be valid object id"
	ErrInvalidFolderIDInput      = "invalid folder id input"
	MesInvalidFolderIDInput      = "folder_id param must be valid object id"
	ErrEmptyIDParam              = "empty id param"
	MesInvalidJSON               = "invalid JSON body"
)

const (
	ErrGettingID                = "error getting videoID from VideoURL"
	ErrDownloadingVideoToServer = "error downloading video to server"
	SuccessfulLoad              = "successful video load"

	ErrGettingVideoRange          = "error getting video range info"
	ErrGettingVideo               = "error getting video"
	ErrDownloadingVideoFromServer = "error downloading video from server"
	ErrRenamingVideo              = "error renaming video"
	ErrMovingVideo                = "error moving video"
	ErrDeletingVideo              = "error deleting video"
)

const (
	ErrCreatingFolder = "error creating new folder"
	ErrRenamingFolder = "error renaming folder"
	ErrMovingFolder   = "error moving folder"
	ErrDeletingFolder = "error deleting folder"
	ErrGettingFolder  = "error getting folder content"
)
