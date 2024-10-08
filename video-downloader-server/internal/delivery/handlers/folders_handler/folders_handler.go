package folders_handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto/folder_dto"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/internal/domain"
)

type FoldersService interface {
	Create(createFolderInput folder_dto.CreateFolderDto) (folder_dto.FolderDto, error)
	Rename(renameFolderInput folder_dto.RenameFolderDto) (folder_dto.FolderDto, error)
	Move(moveFolderInput folder_dto.MoveFolderDto) (folder_dto.FolderDto, error)
	Delete(deleteFolderInput folder_dto.DeleteFolderDto) error
	Get(folderID primitive.ObjectID) (folder_dto.FolderContentDto, error)
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
		r.With(middleware.ValidateCreateFolderInput(f.validator)).Post("/", f.createFolder)
		r.With(middleware.ValidateRenameFolderInput(f.validator)).Put("/rename", f.renameFolder)
		r.With(middleware.ValidateMoveFolderInput(f.validator)).Put("/move", f.moveFolder)
		r.With(middleware.ValidateDeleteFolderInput(f.validator)).Delete("/", f.deleteFolder)
		r.With(middleware.ValidateFolderIDInput).Get("/", f.getFolders)
	})
}

func (f FoldersHandler) createFolder(w http.ResponseWriter, r *http.Request) {
	createFolderInput := r.Context().Value(delivery.CreateFolderInputKey).(folder_dto.CreateFolderDto)

	folder, err := f.foldersService.Create(createFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrCreatingFolder)

		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingFolder, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingFolder, Message: domain.ErrFolderAlreadyExist.Error()})
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
		log.WithError(err).Error(delivery.ErrRenamingFolder)
		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrRenamingFolder, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrRenamingFolder, Message: domain.ErrFolderAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrRenamingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folder)
}

func (f FoldersHandler) moveFolder(w http.ResponseWriter, r *http.Request) {
	moveFolderInput := r.Context().Value(delivery.MoveFolderInputKey).(folder_dto.MoveFolderDto)

	folder, err := f.foldersService.Move(moveFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrMovingFolder)
		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingFolder, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrMovingFolder, Message: domain.ErrFolderAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrMovingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folder)
}

func (f FoldersHandler) deleteFolder(w http.ResponseWriter, r *http.Request) {
	deleteFolderInput := r.Context().Value(delivery.DeleteFolderInputKey).(folder_dto.DeleteFolderDto)

	err := f.foldersService.Delete(deleteFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrDeletingFolder)
		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrDeletingFolder, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDeletingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, nil)
}

func (f FoldersHandler) getFolders(w http.ResponseWriter, r *http.Request) {
	folderID := r.Context().Value(delivery.FolderIDInputKey).(primitive.ObjectID)

	folderContent, err := f.foldersService.Get(folderID)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingFolder)
		if errors.Is(err, domain.ErrFolderNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingFolder, Message: domain.ErrFolderNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folderContent)
}
