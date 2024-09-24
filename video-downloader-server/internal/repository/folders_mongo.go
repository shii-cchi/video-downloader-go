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
	doc := bson.M{"folder_name": folder.FolderName}

	if folder.ParentDirID != primitive.NilObjectID {
		doc["parent_dir_id"] = folder.ParentDirID
	}

	folderID, err := r.db.InsertOne(ctx, doc)
	if err != nil {
		return domain.Folder{}, err
	}

	folder.ID = folderID.InsertedID.(primitive.ObjectID)

	return folder, nil
}

func (r *FoldersRepo) CheckExistByName(ctx context.Context, folderName string, parentDirID primitive.ObjectID) error {
	filter := bson.M{"folder_name": folderName}

	if parentDirID != primitive.NilObjectID {
		filter["parent_dir_id"] = parentDirID
	}

	res := r.db.FindOne(ctx, filter)
	return res.Err()
}

func (r *FoldersRepo) CheckExistByID(ctx context.Context, folderID primitive.ObjectID) error {
	filter := bson.M{"_id": folderID}

	res := r.db.FindOne(ctx, filter)
	return res.Err()
}

func (r *FoldersRepo) UpdateNameByID(ctx context.Context, folderID primitive.ObjectID, newFolderName string) error {
	filter := bson.M{"_id": folderID}
	update := bson.M{"$set": bson.M{"folder_name": newFolderName}}

	result, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *FoldersRepo) Move(ctx context.Context, folderID primitive.ObjectID, parentDirID primitive.ObjectID) error {
	filter := bson.M{"_id": folderID}
	update := bson.M{"$set": bson.M{"parent_dir_id": parentDirID}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *FoldersRepo) GetNameByID(ctx context.Context, folderID primitive.ObjectID) (string, error) {
	folder := domain.Folder{}
	filter := bson.M{"_id": folderID}

	res := r.db.FindOne(ctx, filter).Decode(&folder)
	if res != nil {
		return "", res
	}

	return folder.FolderName, nil
}
