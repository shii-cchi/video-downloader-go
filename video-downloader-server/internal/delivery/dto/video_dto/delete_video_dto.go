package video_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeleteVideoDto struct {
	ID primitive.ObjectID `json:"id" validate:"required,objectid"`
}
