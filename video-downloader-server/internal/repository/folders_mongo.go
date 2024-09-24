package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"video-downloader-server/internal/domain"
)

const (
	foldersCollection = "folders"
)

type FoldersRepo struct {
	db *mongo.Collection
}

func NewFoldersRepo(db *mongo.Database) *FoldersRepo {
	return &FoldersRepo{
		db: db.Collection(foldersCollection),
	}
}

func (r *FoldersRepo) Create(ctx context.Context, folder domain.Folder) (domain.Folder, error) {
	folderID, err := r.db.InsertOne(ctx, folder)
	if err != nil {
		return domain.Folder{}, err
	}

	folder.ID = folderID.InsertedID.(primitive.ObjectID)

	return folder, nil
}

func (r *FoldersRepo) CheckExist(ctx context.Context, folder domain.Folder) error {
	filter := bson.M{"folder_name": folder.FolderName, "parent_dir_id": folder.ParentDirID}
	res := r.db.FindOne(ctx, filter)
	return res.Err()
}
