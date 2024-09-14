package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"video-downloader-server/internal/delivery"
	"video-downloader-server/internal/delivery/dto"
)

func ApplyCors(extensionURL string) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{extensionURL},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept"},
	})
}

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

func CheckVideoNameInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		videoName := r.URL.Query().Get("name")
		if videoName == "" {
			log.Error(delivery.ErrInvalidVideoNameInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidVideoNameInput, Message: delivery.MesInvalidVideoNameInput})
			return
		}

		ctx := context.WithValue(r.Context(), delivery.VideoNameInputKey, videoName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
