package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *VideosRepo) Delete(ctx context.Context, folderID primitive.ObjectID) error {
	filter := bson.M{"folder_id": folderID}

	_, err := r.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *VideosRepo) GetRealPaths(ctx context.Context, folderID primitive.ObjectID) ([]string, []string, error) {
	filter := bson.M{"folder_id": folderID}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var realPaths, previewPaths []string
	for cursor.Next(ctx) {
		var video domain.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, nil, err
		}
		realPaths = append(realPaths, video.RealPath)
		previewPaths = append(previewPaths, video.PreviewPath)
	}

	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}

	return realPaths, previewPaths, nil
}
