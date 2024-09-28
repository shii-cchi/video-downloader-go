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
	return err
}

func (r *VideosRepo) GetRealPath(ctx context.Context, videoID primitive.ObjectID) (string, error) {
	var video domain.Video

	err := r.db.FindOne(ctx, bson.M{"_id": videoID}, options.FindOne().SetProjection(bson.M{"realPath": 1, "_id": 0})).Decode(&video)
	if err != nil {
		return "", err
	}

	return video.RealPath, nil
}

func (r *VideosRepo) CheckExistByID(ctx context.Context, videoID primitive.ObjectID) error {
	return r.db.FindOne(ctx, bson.M{"_id": videoID}).Err()
}

func (r *VideosRepo) Rename(ctx context.Context, videoID primitive.ObjectID, newVideoName string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": videoID}, bson.M{"$set": bson.M{"video_name": newVideoName}})
	return err
}

func (r *VideosRepo) Move(ctx context.Context, videoID primitive.ObjectID, folderID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": videoID}, bson.M{"$set": bson.M{"folder_id": folderID}})
	return err
}

func (r *VideosRepo) Delete(ctx context.Context, videoID primitive.ObjectID) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"_id": videoID})
	return err
}

func (r *VideosRepo) GetPathsByFolders(ctx context.Context, foldersID []primitive.ObjectID) ([]string, []string, error) {
	cursor, err := r.db.Find(ctx, bson.M{"folder_id": bson.M{"$in": foldersID}}, options.Find().SetProjection(bson.M{"real_path": 1, "preview_path": 1, "_id": 0}))
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var video domain.Video
	var realPaths, previewPaths []string

	for cursor.Next(ctx) {
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

func (r *VideosRepo) DeleteVideos(ctx context.Context, foldersID []primitive.ObjectID) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"folder_id": bson.M{"$in": foldersID}})
	return err
}

func (r *VideosRepo) GetVideos(ctx context.Context, folderID primitive.ObjectID) ([]domain.Video, error) {
	cursor, err := r.db.Find(ctx, bson.M{"folder_id": folderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []domain.Video
	var video domain.Video

	for cursor.Next(ctx) {
		if err := cursor.Decode(&video); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
