package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FolderName  string             `bson:"folder_name"`
	ParentDirID primitive.ObjectID `bson:"parent_dir_id,omitempty"`
}
