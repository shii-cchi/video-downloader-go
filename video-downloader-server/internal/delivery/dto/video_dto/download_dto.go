package video_dto

type DownloadDto struct {
	VideoURL   string `json:"video_url" validate:"required,url"`
	FolderName string `json:"folder_name" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=youtube general"`
	Quality    string `json:"quality" validate:"omitempty,oneof=2160p 1440p 1080p 720p 480p 360p 240p 144p best"`
}
