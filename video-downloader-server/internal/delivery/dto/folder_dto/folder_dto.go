package folder_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FolderDto struct {
	ID          primitive.ObjectID  `json:"id"`
	FolderName  string              `json:"folder_name"`
	ParentDirID *primitive.ObjectID `json:"parent_dir_id,omitempty"`
}
