package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" `
	FolderName  string             `bson:"folder_name" json:"folder_name"`
	ParentDirID primitive.ObjectID `bson:"parent_dir_id" json:"parent_dir_id"`
}
