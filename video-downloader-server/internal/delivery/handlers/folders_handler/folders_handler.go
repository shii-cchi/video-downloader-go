package folders_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto/folder_dto"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/internal/domain"
)

type FoldersService interface {
	Create(createFolderInput folder_dto.CreateFolderDto) (folder_dto.FolderDto, error)
	Rename(renameFolderInput folder_dto.RenameFolderDto) (folder_dto.FolderDto, error)
	Move(moveFolderInput folder_dto.MoveFolderDto) (folder_dto.FolderDto, error)
}

type FoldersHandler struct {
	foldersService FoldersService
	validator      *validator.Validate
}

func NewFoldersHandler(foldersService FoldersService, validator *validator.Validate) *FoldersHandler {
	return &FoldersHandler{
		foldersService: foldersService,
		validator:      validator,
	}
}

func (f FoldersHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/folders", func(r chi.Router) {
		r.With(middleware.CheckCreateFolderInput(f.validator)).Post("/", f.createFolder)
		r.With(middleware.CheckRenameFolderInput(f.validator)).Put("/rename", f.renameFolder)
		r.With(middleware.CheckMoveFolderInput(f.validator)).Put("/move", f.moveFolder)
	})
}

func (f FoldersHandler) createFolder(w http.ResponseWriter, r *http.Request) {
	createFolderInput := r.Context().Value(delivery.CreateFolderInputKey).(folder_dto.CreateFolderDto)

	folder, err := f.foldersService.Create(createFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrCreatingFolder)
		if strings.HasPrefix(err.Error(), domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingFolder, Message: domain.ErrFolderAlreadyExist})
			return
		}

		if strings.HasPrefix(err.Error(), domain.ErrParentDirNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingFolder, Message: domain.ErrParentDirNotFound})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrCreatingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusCreated, folder)
}

func (f FoldersHandler) renameFolder(w http.ResponseWriter, r *http.Request) {
	renameFolderInput := r.Context().Value(delivery.RenameFolderInputKey).(folder_dto.RenameFolderDto)

	folder, err := f.foldersService.Rename(renameFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrRenameFolder)
		if strings.HasPrefix(err.Error(), domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrRenameFolder, Message: domain.ErrFolderNotFound})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrRenameFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folder)
}

func (f FoldersHandler) moveFolder(w http.ResponseWriter, r *http.Request) {
	moveFolderInput := r.Context().Value(delivery.MoveFolderInputKey).(folder_dto.MoveFolderDto)

	folder, err := f.foldersService.Move(moveFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrMovingFolder)
		if strings.HasPrefix(err.Error(), domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingFolder, Message: domain.ErrFolderNotFound})
			return
		}

		if strings.HasPrefix(err.Error(), domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingFolder, Message: domain.ErrFolderAlreadyExist})
			return
		}

		if strings.HasPrefix(err.Error(), domain.ErrParentDirNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingFolder, Message: domain.ErrParentDirNotFound})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrMovingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folder)
}
