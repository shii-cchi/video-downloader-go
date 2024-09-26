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
	Delete(ctx context.Context, folderID primitive.ObjectID) error
	GetFolders(ctx context.Context, parentDirID primitive.ObjectID) ([]domain.Folder, error)
}

type Videos interface {
	DeleteVideos(foldersID []primitive.ObjectID) error
}

type FoldersService struct {
	repo          FoldersRepo
	videosService Videos
}

func NewFoldersService(repo FoldersRepo, videosService Videos) *FoldersService {
	return &FoldersService{
		repo:          repo,
		videosService: videosService,
	}
}

func (f *FoldersService) Create(createFolderInput folder_dto.CreateFolderDto) (folder_dto.FolderDto, error) {
	folder := domain.Folder{FolderName: createFolderInput.FolderName}

	if createFolderInput.ParentDirID != primitive.NilObjectID {
		folder.ParentDirID = createFolderInput.ParentDirID

		if err := f.checkFolderExistenceByID(folder.ParentDirID); err != nil {
			return folder_dto.FolderDto{}, err
		}
	}

	if err := f.checkFolderExistenceByName(folder.FolderName, folder.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	newFolder, err := f.repo.Create(context.Background(), folder)
	if err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrCreatingFolder+": %w", err)
	}

	res := folder_dto.FolderDto{
		ID:         newFolder.ID,
		FolderName: newFolder.FolderName,
	}

	if newFolder.ParentDirID != primitive.NilObjectID {
		res.ParentDirID = &newFolder.ParentDirID
	}

	return res, nil
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
	if err := f.checkFolderExistenceByID(moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	if err := f.checkFolderExistenceByID(moveFolderInput.ID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	name, err := f.repo.GetNameByID(context.Background(), moveFolderInput.ID)
	if err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf(domain.ErrGettingFolderName+": %w", err)
	}

	if err := f.checkFolderExistenceByName(name, moveFolderInput.ParentDirID); err != nil {
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

func (f *FoldersService) Delete(deleteFolderInput folder_dto.DeleteFolderDto) error {
	if err := f.checkFolderExistenceByID(deleteFolderInput.ID); err != nil {
		return err
	}

	var foldersID []primitive.ObjectID

	if err := f.deleteFolder(deleteFolderInput.ID, &foldersID); err != nil {
		return err
	}

	if err := f.videosService.DeleteVideos(foldersID); err != nil {
		return err
	}

	return nil
}

func (f *FoldersService) checkFolderExistenceByID(folderID primitive.ObjectID) error {
	err := f.repo.CheckExistByID(context.Background(), folderID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf(domain.ErrFolderNotFound+": %w", err)
		}

		return fmt.Errorf(domain.ErrCheckingFolder+": %w", err)
	}

	return nil
}

func (f *FoldersService) checkFolderExistenceByName(folderName string, parentDirID primitive.ObjectID) error {
	err := f.repo.CheckExistByName(context.Background(), folderName, parentDirID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf(domain.ErrCheckingFolder+": %w", err)
	}

	if err == nil {
		return errors.New(domain.ErrFolderAlreadyExist)
	}

	return nil
}

func (f *FoldersService) deleteFolder(folderID primitive.ObjectID, foldersID *[]primitive.ObjectID) error {
	if err := f.repo.Delete(context.Background(), folderID); err != nil {
		return fmt.Errorf(domain.ErrDeletingFolder+": %w", err)
	}

	*foldersID = append(*foldersID, folderID)

	folders, err := f.repo.GetFolders(context.Background(), folderID)
	if err != nil {
		return fmt.Errorf(domain.ErrGettingFoldersList+": %w", err)
	}

	for _, folder := range folders {
		if err := f.deleteFolder(folder.ID, foldersID); err != nil {
			return err
		}
	}

	return nil
}
