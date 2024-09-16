package repository

import "go.mongodb.org/mongo-driver/mongo"

type VideoManagementRepo struct {
	db *mongo.Collection
}

func NewVideoManagementRepo(db *mongo.Database) *VideoManagementRepo {
	return &VideoManagementRepo{
		db: db.Collection(videosCollection),
	}
}
