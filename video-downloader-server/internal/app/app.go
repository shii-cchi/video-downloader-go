package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"video-downloader-server/internal/config"
	"video-downloader-server/internal/delivery/handlers/folders_handler"
	"video-downloader-server/internal/delivery/handlers/videos_handler"
	"video-downloader-server/internal/repository"
	"video-downloader-server/internal/service/folders_service"
	"video-downloader-server/internal/service/preview_service"
	"video-downloader-server/internal/service/videos_service"
	"video-downloader-server/internal/validator"
)

const (
	errLoadingConfig    = "error loading config"
	errCreatingDbClient = "error creating mongo db client"
	errConnectingToDb   = "error connecting to mongo db"

	successfulConfigLoad     = "config has been loaded successfully"
	successfulConnectionToDb = "successfully connected to MongoDB"
	serverStart              = "server starting on port"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal(errLoadingConfig)
	}
	log.Info(successfulConfigLoad)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.vs9z4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", cfg.DbUser, cfg.DbPassword)).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.WithError(err).Fatal(errCreatingDbClient)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.WithError(err).Fatal(errConnectingToDb)
	}
	log.Info(successfulConnectionToDb)

	db := client.Database(cfg.DbName)
	videosRepo := repository.NewVideosRepo(db)
	foldersRepo := repository.NewFoldersRepo(db)

	previewService := preview_service.NewPreviewService()
	videosService := videos_service.NewVideosService(videosRepo, previewService)
	folderService := folders_service.NewFoldersService(foldersRepo, videosService)

	v := validator.Init()
	videosHandler := videos_handler.NewVideosHandler(videosService, v)
	foldersHandler := folders_handler.NewFoldersHandler(folderService, v)

	r := chi.NewRouter()
	videosHandler.RegisterRoutes(r)
	foldersHandler.RegisterRoutes(r)

	log.Infof(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
