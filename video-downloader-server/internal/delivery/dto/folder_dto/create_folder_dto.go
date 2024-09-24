package folder_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateFolderDto struct {
	FolderName  string             `json:"folder_name" validate:"required,foldername"`
	ParentDirID primitive.ObjectID `json:"parent_dir_id" validate:"omitempty,objectid"`
}
