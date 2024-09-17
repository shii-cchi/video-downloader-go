package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateParams struct {
	VideoName   string
	RealPath    string
	UserPath    string
	PreviewPath string
}

type VideoDownloadRepo struct {
	db *mongo.Collection
}

func NewVideoDownloadRepo(db *mongo.Database) *VideoDownloadRepo {
	return &VideoDownloadRepo{
		db: db.Collection(videosCollection),
	}
}

func (r *VideoDownloadRepo) Create(ctx context.Context, createParams CreateParams) error {
	doc := bson.D{{"video_name", createParams.VideoName}, {"real_path", createParams.RealPath}, {"user_path", createParams.UserPath}, {"preview_path", createParams.PreviewPath}}
	_, err := r.db.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}
