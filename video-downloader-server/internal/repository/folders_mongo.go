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

	var folders []domain.Folder
	if err := cursor.All(ctx, &folders); err != nil {
		return nil, err
	}

	var results []primitive.ObjectID
	for _, folder := range folders {
		results = append(results, folder.ID)
	}

	return results, nil
}

func (r *FoldersRepo) DeleteAllNestedFolders(ctx context.Context, foldersID []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": foldersID}}

	_, err := r.db.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *FoldersRepo) GetNestedFolders(ctx context.Context, folderID primitive.ObjectID) ([]domain.Folder, error) {
	filter := bson.M{"parent_dir_id": folderID}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var folders []domain.Folder
	if err := cursor.All(ctx, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

//func (r *FoldersRepo) Delete(ctx context.Context, folderID primitive.ObjectID) error {
//	filter := bson.M{"_id": folderID}
//
//	result, err := r.db.DeleteOne(ctx, filter)
//	if err != nil {
//		return err
//	}
//
//	if result.DeletedCount == 0 {
//		return mongo.ErrNoDocuments
//	}
//
//	return nil
//}
//
//func (r *FoldersRepo) GetFolders(ctx context.Context, parentDirID primitive.ObjectID) ([]domain.Folder, error) {
//	filter := bson.M{"parent_dir_id": parentDirID}
//
//	cursor, err := r.db.Find(ctx, filter)
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(ctx)
//
//	var folders []domain.Folder
//	for cursor.Next(ctx) {
//		var folder domain.Folder
//		if err := cursor.Decode(&folder); err != nil {
//			return nil, err
//		}
//		folders = append(folders, folder)
//	}
//
//	if err := cursor.Err(); err != nil {
//		return nil, err
//	}
//
//	return folders, nil
//}
