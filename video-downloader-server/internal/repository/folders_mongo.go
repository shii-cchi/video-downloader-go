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

func (r *FoldersRepo) CheckExistenceByID(ctx context.Context, folderID primitive.ObjectID) error {
	return r.db.FindOne(ctx, bson.M{"_id": folderID}).Err()
}

func (r *FoldersRepo) CheckExistenceByName(ctx context.Context, folderName string, parentDirID primitive.ObjectID) error {
	filter := bson.M{"folder_name": folderName}

	if parentDirID != primitive.NilObjectID {
		filter["parent_dir_id"] = parentDirID
	}

	return r.db.FindOne(ctx, filter).Err()
}

func (r *FoldersRepo) Create(ctx context.Context, folderName string, parentDirID primitive.ObjectID) (primitive.ObjectID, error) {
	doc := bson.M{"folder_name": folderName}

	if parentDirID != primitive.NilObjectID {
		doc["parent_dir_id"] = parentDirID
	}

	res, err := r.db.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *FoldersRepo) GetParentDirID(ctx context.Context, folderID primitive.ObjectID) (primitive.ObjectID, error) {
	var folder domain.Folder

	if err := r.db.FindOne(ctx, bson.M{"_id": folderID}, options.FindOne().SetProjection(bson.M{"parent_dir_id": 1, "_id": 0})).Decode(&folder); err != nil {
		return primitive.NilObjectID, err
	}

	return folder.ParentDirID, nil
}

func (r *FoldersRepo) UpdateName(ctx context.Context, folderID primitive.ObjectID, newFolderName string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": folderID}, bson.M{"$set": bson.M{"folder_name": newFolderName}})
	return err
}

func (r *FoldersRepo) GetName(ctx context.Context, folderID primitive.ObjectID) (string, error) {
	var folder domain.Folder

	if err := r.db.FindOne(ctx, bson.M{"_id": folderID}, options.FindOne().SetProjection(bson.M{"folder_name": 1, "_id": 0})).Decode(&folder); err != nil {
		return "", err
	}

	return folder.FolderName, nil
}

func (r *FoldersRepo) Move(ctx context.Context, folderID primitive.ObjectID, parentDirID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": folderID}, bson.M{"$set": bson.M{"parent_dir_id": parentDirID}})
	return err
}

func (r *FoldersRepo) GetAllNestedFolders(ctx context.Context, parentDirID primitive.ObjectID) ([]primitive.ObjectID, error) {
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"parent_dir_id", parentDirID}}},
		},
		{
			{"$graphLookup", bson.D{
				{"from", "folders"},
				{"startWith", "$_id"},
				{"connectFromField", "_id"},
				{"connectToField", "parent_dir_id"},
				{"as", "nestedFolders"},
			}},
		},
		{
			{"$project", bson.D{
				{"allFolderIds", bson.D{
					{"$setUnion", bson.A{
						bson.A{bson.D{{"$toObjectId", "$_id"}}},
						bson.D{
							{"$map", bson.D{
								{"input", "$nestedFolders"},
								{"as", "folder"},
								{"in", "$$folder._id"},
							}},
						},
					}},
				}},
			}},
		},
		{
			{"$unwind", bson.D{{"path", "$allFolderIds"}}},
		},
		{
			{"$project", bson.D{
				{"_id", "$allFolderIds"},
			}},
		},
	}

	cursor, err := r.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var folder domain.Folder
	var results []primitive.ObjectID

	for cursor.Next(ctx) {
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		results = append(results, folder.ID)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *FoldersRepo) DeleteAllNestedFolders(ctx context.Context, foldersID []primitive.ObjectID) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": foldersID}})
	return err
}

func (r *FoldersRepo) GetNestedFolders(ctx context.Context, folderID primitive.ObjectID) ([]domain.Folder, error) {
	cursor, err := r.db.Find(ctx, bson.M{"parent_dir_id": folderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var folder domain.Folder
	var folders []domain.Folder
	for cursor.Next(ctx) {
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}
