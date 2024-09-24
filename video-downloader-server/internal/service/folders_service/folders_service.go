package folders_service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/internal/domain"
)

type FoldersRepo interface {
	Create(ctx context.Context, folder domain.Folder) (domain.Folder, error)
	CheckExist(ctx context.Context, folder domain.Folder) error
}

type FoldersService struct {
	repo FoldersRepo
}

func NewFoldersService(repo FoldersRepo) *FoldersService {
	return &FoldersService{
		repo: repo,
	}
}

func (f *FoldersService) Create(createFolderInput dto.CreateFolderInputDto) (domain.Folder, error) {
	parentDirID, _ := primitive.ObjectIDFromHex(createFolderInput.ParentDirID)

	err := f.repo.CheckExist(context.Background(), domain.Folder{FolderName: createFolderInput.FolderName, ParentDirID: parentDirID})

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Folder{}, fmt.Errorf(domain.ErrCheckingFolder+": %w", err)
	}

	if err == nil {
		return domain.Folder{}, errors.New(domain.ErrFolderAlreadyExist)
	}

	newFolder, err := f.repo.Create(context.Background(), domain.Folder{FolderName: createFolderInput.FolderName, ParentDirID: parentDirID})
	if err != nil {
		return domain.Folder{}, fmt.Errorf(domain.ErrCreatingFolder+": %w", err)
	}

	return newFolder, nil
}
