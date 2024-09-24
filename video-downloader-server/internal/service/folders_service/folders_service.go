package folders_service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"video-downloader-server/internal/delivery/dto/folder_dto"
	"video-downloader-server/internal/domain"
)

type FoldersRepo interface {
	Create(ctx context.Context, folder domain.Folder) (domain.Folder, error)
	CheckExistByName(ctx context.Context, folderName string, parentDirID primitive.ObjectID) error
	CheckExistByID(ctx context.Context, folderID primitive.ObjectID) error
	UpdateNameByID(ctx context.Context, folderID primitive.ObjectID, newFolderName string) error
	Move(ctx context.Context, folderID primitive.ObjectID, parentDirID primitive.ObjectID) error
	GetNameByID(ctx context.Context, folderID primitive.ObjectID) (string, error)
}

type FoldersService struct {
	repo FoldersRepo
}

func NewFoldersService(repo FoldersRepo) *FoldersService {
	return &FoldersService{
		repo: repo,
	}
}

func (f *FoldersService) Create(createFolderInput folder_dto.CreateFolderDto) (folder_dto.FolderDto, error) {
	folder := domain.Folder{FolderName: createFolderInput.FolderName}

	if createFolderInput.ParentDirID != primitive.NilObjectID {
		folder.ParentDirID = createFolderInput.ParentDirID

		if err := f.checkParentDirExistence(folder.ParentDirID); err != nil {
			return folder_dto.FolderDto{}, err
		}
	}

	if err := f.checkFolderExistence(folder.FolderName, folder.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	newFolder, err := f.repo.Create(context.Background(), folder)
	if err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrCreatingFolder+": %w", err)
	}

	return f.createFolderOutput(newFolder), nil
}

func (f *FoldersService) Rename(renameFolderInput folder_dto.RenameFolderDto) (folder_dto.FolderDto, error) {
	if err := f.repo.UpdateNameByID(context.Background(), renameFolderInput.ID, renameFolderInput.FolderName); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrFolderNotFound+": %w", err)
		}

		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrUpdatingFolderName+": %w", err)
	}

	return folder_dto.FolderDto{
		ID:         renameFolderInput.ID,
		FolderName: renameFolderInput.FolderName,
	}, nil
}

func (f *FoldersService) Move(moveFolderInput folder_dto.MoveFolderDto) (folder_dto.FolderDto, error) {
	if err := f.checkParentDirExistence(moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	name, err := f.repo.GetNameByID(context.Background(), moveFolderInput.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrFolderNotFound+": %w", err)
		}

		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrGettingFolderName+": %w", err)
	}

	if err := f.checkFolderExistence(name, moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	if err := f.repo.Move(context.Background(), moveFolderInput.ID, moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrMovingFolder+": %w", err)
	}

	return folder_dto.FolderDto{
		ID:          moveFolderInput.ID,
		ParentDirID: &moveFolderInput.ParentDirID,
	}, nil
}

func (f *FoldersService) checkParentDirExistence(parentDirID primitive.ObjectID) error {
	err := f.repo.CheckExistByID(context.Background(), parentDirID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf(domain.ErrParentDirNotFound+": %w", err)
		}

		return fmt.Errorf(domain.ErrCheckingParentDir+": %w", err)
	}

	return nil
}

func (f *FoldersService) checkFolderExistence(folderName string, parentDirID primitive.ObjectID) error {
	err := f.repo.CheckExistByName(context.Background(), folderName, parentDirID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf(domain.ErrCheckingFolder+": %w", err)
	}

	if err == nil {
		return errors.New(domain.ErrFolderAlreadyExist)
	}

	return nil
}

func (f *FoldersService) createFolderOutput(folder domain.Folder) folder_dto.FolderDto {
	res := folder_dto.FolderDto{
		ID:         folder.ID,
		FolderName: folder.FolderName,
	}

	if folder.ParentDirID != primitive.NilObjectID {
		res.ParentDirID = &folder.ParentDirID
	}

	return res
}
