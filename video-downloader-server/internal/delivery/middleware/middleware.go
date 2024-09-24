package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto"
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

func CheckDownloadInput(validate *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			downloadInput := dto.DownloadInputDto{}
			if err := json.NewDecoder(r.Body).Decode(&downloadInput); err != nil {
				log.WithError(err).Error(delivery.ErrInvalidDownloadInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidDownloadInput, Message: delivery.MesInvalidJSON})
				return
			}

			if err := validate.Struct(&downloadInput); err != nil {
				log.WithError(err).Error(delivery.ErrInvalidDownloadInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidDownloadInput, Message: delivery.MesInvalidDownloadInput})
				return
			}

			ctx := context.WithValue(r.Context(), delivery.DownloadInputKey, downloadInput)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CheckVideoIDInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		videoID := r.URL.Query().Get("id")
		if videoID == "" {
			log.Error(delivery.ErrInvalidVideoIDInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidVideoIDInput, Message: delivery.MesInvalidVideoIDInput})
			return
		}

		ctx := context.WithValue(r.Context(), delivery.VideoIDInputKey, videoID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CheckCreateFolderInput(validate *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			createFolderInput := dto.CreateFolderInputDto{}
			if err := json.NewDecoder(r.Body).Decode(&createFolderInput); err != nil {
				log.WithError(err).Error(delivery.ErrInvalidCreateFolderInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidCreateFolderInput, Message: delivery.MesInvalidJSON})
				return
			}

			if err := validate.Struct(&createFolderInput); err != nil {
				log.WithError(err).Error(delivery.ErrInvalidCreateFolderInput)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidCreateFolderInput, Message: delivery.MesInvalidCreateFolderInput})
				return
			}

			ctx := context.WithValue(r.Context(), delivery.CreateFolderInputKey, createFolderInput)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
