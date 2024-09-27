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

func (r *VideosRepo) GetPathsByFolders(ctx context.Context, foldersID []primitive.ObjectID) ([]string, []string, error) {
	filter := bson.M{"folder_id": bson.M{"$in": foldersID}}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.Video

	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	var realPaths, previewPaths []string

	for _, result := range results {
		realPaths = append(realPaths, result.RealPath)
		previewPaths = append(previewPaths, result.PreviewPath)
	}

	return realPaths, previewPaths, nil
}

func (r *VideosRepo) DeleteVideos(ctx context.Context, foldersID []primitive.ObjectID) error {
	filter := bson.M{"folder_id": bson.M{"$in": foldersID}}

	_, err := r.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

//func (r *VideosRepo) Delete(ctx context.Context, folderID primitive.ObjectID) error {
//	filter := bson.M{"folder_id": folderID}
//
//	_, err := r.db.DeleteMany(ctx, filter)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
