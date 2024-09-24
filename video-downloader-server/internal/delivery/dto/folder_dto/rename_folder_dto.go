package folder_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type RenameFolderDto struct {
	ID         primitive.ObjectID `json:"id" validate:"required,objectid"`
	FolderName string             `json:"folder_name" validate:"required,foldername"`
}
