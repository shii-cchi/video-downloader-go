package folders_service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"video-downloader-server/internal/delivery/dto/folder_dto"
	"video-downloader-server/internal/delivery/dto/video_dto"
	"video-downloader-server/internal/domain"
)

type FoldersRepo interface {
	CheckExistenceByID(ctx context.Context, folderID primitive.ObjectID) error
	CheckExistenceByName(ctx context.Context, folderName string, parentDirID primitive.ObjectID) error
	Create(ctx context.Context, folderName string, parentDirID primitive.ObjectID) (primitive.ObjectID, error)
	GetParentDirID(ctx context.Context, folderID primitive.ObjectID) (primitive.ObjectID, error)
	UpdateName(ctx context.Context, folderID primitive.ObjectID, newFolderName string) error
	GetName(ctx context.Context, folderID primitive.ObjectID) (string, error)
	Move(ctx context.Context, folderID primitive.ObjectID, parentDirID primitive.ObjectID) error
	GetAllNestedFolders(ctx context.Context, parentDirID primitive.ObjectID) ([]primitive.ObjectID, error)
	DeleteAllNestedFolders(ctx context.Context, foldersID []primitive.ObjectID) error
	GetNestedFolders(ctx context.Context, folderID primitive.ObjectID) ([]domain.Folder, error)
}

type Videos interface {
	DeleteVideos(foldersID []primitive.ObjectID) error
	GetVideos(folderID primitive.ObjectID) ([]video_dto.VideoDto, error)
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
	if createFolderInput.ParentDirID != primitive.NilObjectID {
		if err := f.checkFolderExistenceByID(createFolderInput.ParentDirID); err != nil {
			return folder_dto.FolderDto{}, err
		}
	}

	if err := f.checkFolderExistenceByName(createFolderInput.FolderName, createFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	folderID, err := f.repo.Create(context.Background(), createFolderInput.FolderName, createFolderInput.ParentDirID)
	if err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf("%w (folder name: %s, parent dir id: %s): %s", domain.ErrCreatingFolder, createFolderInput.FolderName, createFolderInput.ParentDirID, err)
	}

	return folder_dto.FolderDto{
		ID:         folderID,
		FolderName: createFolderInput.FolderName,
		ParentDirID: func() *primitive.ObjectID {
			if createFolderInput.ParentDirID != primitive.NilObjectID {
				return &createFolderInput.ParentDirID
			}
			return nil
		}(),
	}, nil
}

func (f *FoldersService) Rename(renameFolderInput folder_dto.RenameFolderDto) (folder_dto.FolderDto, error) {
	parentDirID, err := f.repo.GetParentDirID(context.Background(), renameFolderInput.ID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return folder_dto.FolderDto{}, fmt.Errorf("%w (folder id: %s)", domain.ErrFolderNotFound, renameFolderInput.ID)
		}
	}

	if err := f.checkFolderExistenceByName(renameFolderInput.FolderName, parentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	if err := f.repo.UpdateName(context.Background(), renameFolderInput.ID, renameFolderInput.FolderName); err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf("%w (folder id: %s, folder name: %s): %s", domain.ErrRenamingFolder, renameFolderInput.ID, renameFolderInput.FolderName, err)
	}

	return folder_dto.FolderDto{
		ID:          renameFolderInput.ID,
		FolderName:  renameFolderInput.FolderName,
		ParentDirID: &parentDirID,
	}, nil
}

func (f *FoldersService) Move(moveFolderInput folder_dto.MoveFolderDto) (folder_dto.FolderDto, error) {
	if err := f.checkFolderExistenceByID(moveFolderInput.ID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	if err := f.checkFolderExistenceByID(moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	name, err := f.repo.GetName(context.Background(), moveFolderInput.ID)
	if err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf("%w (folder id: %s): %s", domain.ErrGettingFolderName, moveFolderInput.ID, err)
	}

	if err := f.checkFolderExistenceByName(name, moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, err
	}

	if err := f.repo.Move(context.Background(), moveFolderInput.ID, moveFolderInput.ParentDirID); err != nil {
		return folder_dto.FolderDto{}, fmt.Errorf("%w (folder id: %s, parent dir id: %s): %s", domain.ErrMovingFolder, moveFolderInput.ID, moveFolderInput.ParentDirID, err)
	}

	return folder_dto.FolderDto{
		ID:          moveFolderInput.ID,
		FolderName:  name,
		ParentDirID: &moveFolderInput.ParentDirID,
	}, nil
}

func (f *FoldersService) Delete(deleteFolderInput folder_dto.DeleteFolderDto) error {
	if err := f.checkFolderExistenceByID(deleteFolderInput.ID); err != nil {
		return err
	}

	allFolders, err := f.repo.GetAllNestedFolders(context.Background(), deleteFolderInput.ID)
	if err != nil {
		return fmt.Errorf("%w (folder id: %s): %s", domain.ErrGettingAllNestedFolders, deleteFolderInput.ID, err)
	}

	var foldersID []primitive.ObjectID
	foldersID = append(foldersID, deleteFolderInput.ID)
	foldersID = append(foldersID, allFolders...)

	if err := f.repo.DeleteAllNestedFolders(context.Background(), foldersID); err != nil {
		return fmt.Errorf("%w (folder id: %s): %s", domain.ErrDeletingAllNestedFolders, deleteFolderInput.ID, err)
	}

	if err := f.videosService.DeleteVideos(foldersID); err != nil {
		return err
	}

	return nil
}

func (f *FoldersService) Get(folderIDStr string) (folder_dto.FolderContentDto, error) {
	folderID, err := primitive.ObjectIDFromHex(folderIDStr)
	if err != nil {
		return folder_dto.FolderContentDto{}, fmt.Errorf("%w: %s", domain.ErrConvertingToObjectID, err)
	}

	if err := f.checkFolderExistenceByID(folderID); err != nil {
		return folder_dto.FolderContentDto{}, err
	}

	folders, err := f.repo.GetNestedFolders(context.Background(), folderID)
	if err != nil {
		return folder_dto.FolderContentDto{}, fmt.Errorf("%w (folder id: %s): %s", domain.ErrGettingNestedFolders, folderID, err)
	}

	videos, err := f.videosService.GetVideos(folderID)
	if err != nil {
		return folder_dto.FolderContentDto{}, err
	}

	return folder_dto.FolderContentDto{
		ID:      folderID,
		Folders: f.toFolderContentDto(folders),
		Videos:  videos,
	}, nil
}

func (f *FoldersService) checkFolderExistenceByID(folderID primitive.ObjectID) error {
	err := f.repo.CheckExistenceByID(context.Background(), folderID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("%w (folderID: %s)", domain.ErrFolderNotFound, folderID)
		}

		return fmt.Errorf("%w (folderID: %s): %s", domain.ErrCheckingFolder, folderID, err)
	}

	return nil
}

func (f *FoldersService) checkFolderExistenceByName(folderName string, parentDirID primitive.ObjectID) error {
	err := f.repo.CheckExistenceByName(context.Background(), folderName, parentDirID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%w (folder name: %s, parent dir id: %s): %s", domain.ErrCheckingFolder, folderName, parentDirID, err)
	}

	if err == nil {
		return fmt.Errorf("%w (folder name: %s, parent dir id: %s)", domain.ErrFolderAlreadyExist, folderName, parentDirID)
	}

	return nil
}

func (f *FoldersService) toFolderContentDto(folders []domain.Folder) []folder_dto.FolderDto {
	res := make([]folder_dto.FolderDto, 0, len(folders))

	for _, folder := range folders {
		res = append(res, folder_dto.FolderDto{
			ID:          folder.ID,
			FolderName:  folder.FolderName,
			ParentDirID: &folder.ParentDirID,
		})
	}

	return res
}
