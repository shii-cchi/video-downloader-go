package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto/folder_dto"
	"video-downloader-server/internal/delivery/dto/video_dto"
)

//func ApplyCors(extensionURL string) func(next http.Handler) http.Handler {
//	return cors.Handler(cors.Options{
//		AllowedOrigins:     []string{extensionURL},
//		AllowCredentials:   true,
//		AllowedMethods:     []string{"POST", "OPTIONS"},
//		AllowedHeaders:     []string{"Origin", "Content-Type", "Accept"},
//		OptionsPassthrough: true,
//	})
//}

type ValidatableDto interface {
	video_dto.DownloadVideoDto | video_dto.RenameVideoDto | video_dto.MoveVideoDto | video_dto.DeleteVideoDto | folder_dto.CreateFolderDto | folder_dto.RenameFolderDto | folder_dto.MoveFolderDto | folder_dto.DeleteFolderDto
}

func validateInput[V ValidatableDto](validate *validator.Validate, input V, ctxKey delivery.ContextKey, errInvalidInput, errMessage string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				log.WithError(err).Error(errInvalidInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: errInvalidInput, Message: delivery.MesInvalidJSON})
				return
			}

			if err := validate.Struct(input); err != nil {
				log.WithError(err).Error(errInvalidInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: errInvalidInput, Message: errMessage})
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey, input)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateDownloadVideoInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return validateInput(v, video_dto.DownloadVideoDto{}, delivery.DownloadVideoInputKey, delivery.ErrInvalidDownloadVideoInput, delivery.MesInvalidDownloadVideoInput)
}

func ValidateRenameVideoInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return validateInput(v, video_dto.RenameVideoDto{}, delivery.RenameVideoInputKey, delivery.ErrInvalidRenameVideoInput, delivery.MesInvalidRenameVideoInput)
}

func ValidateMoveVideoInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return validateInput(v, video_dto.MoveVideoDto{}, delivery.MoveVideoInputKey, delivery.ErrInvalidMoveVideoInput, delivery.MesInvalidMoveVideoInput)
}

func ValidateDeleteVideoInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return validateInput(v, video_dto.DeleteVideoDto{}, delivery.DeleteVideoInputKey, delivery.ErrInvalidDeleteVideoInput, delivery.MesInvalidDeleteVideoInput)
}

func ValidateCreateFolderInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return validateInput(v, folder_dto.CreateFolderDto{}, delivery.CreateFolderInputKey, delivery.ErrInvalidCreateFolderInput, delivery.MesInvalidCreateFolderInput)
}

func ValidateRenameFolderInput(v *validator.Validate) func(http.Handler) http.Handler {
	return validateInput(v, folder_dto.RenameFolderDto{}, delivery.RenameFolderInputKey, delivery.ErrInvalidRenameFolderInput, delivery.MesInvalidRenameFolderInput)
}

func ValidateMoveFolderInput(v *validator.Validate) func(http.Handler) http.Handler {
	return validateInput(v, folder_dto.MoveFolderDto{}, delivery.MoveFolderInputKey, delivery.ErrInvalidMoveFolderInput, delivery.MesInvalidMoveFolderInput)
}

func ValidateDeleteFolderInput(v *validator.Validate) func(http.Handler) http.Handler {
	return validateInput(v, folder_dto.DeleteFolderDto{}, delivery.DeleteFolderInputKey, delivery.ErrInvalidDeleteFolderInput, delivery.MesInvalidDeleteFolderInput)
}

func validateIDInput(paramName string, ctxKey delivery.ContextKey, errInvalidInput, errMessage string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paramValueStr := r.URL.Query().Get(paramName)
			if paramValueStr == "" {
				log.Error(errInvalidInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: errInvalidInput, Message: delivery.ErrEmptyIDParam})
				return
			}

			paramValueID, err := primitive.ObjectIDFromHex(paramValueStr)
			if err != nil {
				log.Error(errInvalidInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: errInvalidInput, Message: errMessage})
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey, paramValueID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateVideoIDInput(next http.Handler) http.Handler {
	return validateIDInput("video_id", delivery.VideoIDInputKey, delivery.ErrInvalidVideoIDInput, delivery.MesInvalidVideoIDInput)(next)
}

func ValidateFolderIDInput(next http.Handler) http.Handler {
	return validateIDInput("folder_id", delivery.FolderIDInputKey, delivery.ErrInvalidFolderIDInput, delivery.MesInvalidFolderIDInput)(next)
}
