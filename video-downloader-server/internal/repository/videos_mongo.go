package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *VideosRepo) GetVideos(ctx context.Context, folderID primitive.ObjectID) ([]domain.Video, error) {
	filter := bson.M{"folder_id": folderID}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []domain.Video
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}

func (r *VideosRepo) GetRealPath(ctx context.Context, videoID primitive.ObjectID) (string, error) {
	filter := bson.M{"_id": videoID}

	var video domain.Video

	err := r.db.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"realPath": 1, "_id": 0})).Decode(&video)
	if err != nil {
		return "", err
	}

	return video.RealPath, nil
}

func (r *VideosRepo) CheckExistByID(ctx context.Context, videoID primitive.ObjectID) error {
	filter := bson.M{"_id": videoID}

	res := r.db.FindOne(ctx, filter)
	return res.Err()
}

func (r *VideosRepo) Rename(ctx context.Context, videoID primitive.ObjectID, newVideoName string) error {
	filter := bson.M{"_id": videoID}
	update := bson.M{"$set": bson.M{"video_name": newVideoName}}

	result, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *VideosRepo) Delete(ctx context.Context, videoID primitive.ObjectID) error {
	filter := bson.M{"_id": videoID}

	_, err := r.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *VideosRepo) Move(ctx context.Context, videoID primitive.ObjectID, folderID primitive.ObjectID) error {
	filter := bson.M{"_id": videoID}
	update := bson.M{"$set": bson.M{"folder_id": folderID}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
