package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"video-downloader-server/internal/domain"
)

const (
	videosCollection = "videos"
)

type VideosRepo struct {
	db *mongo.Collection
}

func NewVideosRepo(db *mongo.Database) *VideosRepo {
	return &VideosRepo{
		db: db.Collection(videosCollection),
	}
}

func (r *VideosRepo) Create(ctx context.Context, video domain.Video) error {
	_, err := r.db.InsertOne(ctx, video)
	if err != nil {
		return err
	}

	return nil
}
