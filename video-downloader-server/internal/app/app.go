package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"video-downloader-server/internal/config"
	"video-downloader-server/internal/delivery/handlers"
	"video-downloader-server/service"
)

const (
	errLoadingConfig = "error loading config"

	successfulConfigLoad = "config has been loaded successfully"
	serverStart          = "server starting on port"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal(errLoadingConfig)
	}
	log.Info(successfulConfigLoad)

	s := service.NewService()
	v := validator.New()

	r := chi.NewRouter()
	h := handlers.NewHandler(s, v, cfg.ExtensionURL)
	h.RegisterRoutes(r)

	log.Infof(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
