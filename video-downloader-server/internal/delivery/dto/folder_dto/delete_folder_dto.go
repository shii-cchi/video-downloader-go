package folder_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteFolderDto struct {
	ID primitive.ObjectID `json:"id" validate:"required,objectid"`
}
