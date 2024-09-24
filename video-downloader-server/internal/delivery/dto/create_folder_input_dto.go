package dto

type CreateFolderInputDto struct {
	FolderName  string `json:"folder_name" validate:"required,foldername"`
	ParentDirID string `json:"parent_dir_id" validate:"required,objectid"`
}
