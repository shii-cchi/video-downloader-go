package folders_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto"
	"video-downloader-server/internal/delivery/middleware"
	"video-downloader-server/internal/domain"
)

type FoldersService interface {
	Create(createFolderInput dto.CreateFolderInputDto) (domain.Folder, error)
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
	})
}

func (f FoldersHandler) createFolder(w http.ResponseWriter, r *http.Request) {
	createFolderInput := r.Context().Value(delivery.CreateFolderInputKey).(dto.CreateFolderInputDto)

	folder, err := f.foldersService.Create(createFolderInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrCreatingFolder)
		if strings.HasPrefix(err.Error(), domain.ErrFolderAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingFolder, Message: domain.ErrFolderAlreadyExist})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrCreatingFolder})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, folder)
}
