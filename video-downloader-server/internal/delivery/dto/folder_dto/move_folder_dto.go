package folder_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type MoveFolderDto struct {
	ID          primitive.ObjectID `json:"id" validate:"required,objectid"`
	ParentDirID primitive.ObjectID `json:"parent_dir_id" validate:"required,objectid"`
}
