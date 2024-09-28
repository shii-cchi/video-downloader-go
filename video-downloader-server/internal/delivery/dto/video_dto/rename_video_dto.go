package video_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type RenameVideoDto struct {
	ID        primitive.ObjectID `json:"id" validate:"required,objectid"`
	VideoName string             `json:"video_name" validate:"required,videoname"`
}
