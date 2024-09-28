package video_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type MoveVideoDto struct {
	ID       primitive.ObjectID `json:"id" validate:"required,objectid"`
	FolderID primitive.ObjectID `json:"folder_id" validate:"required,objectid"`
}
