package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	videosCollection = "videos"
)

type VideoDownload interface {
	Create(ctx context.Context, createParams CreateParams) error
}

type VideoManagement interface {
}

type Repository struct {
	VideoDownload   VideoDownload
	VideoManagement VideoManagement
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		VideoDownload:   NewVideoDownloadRepo(db),
		VideoManagement: NewVideoManagementRepo(db),
	}
}
