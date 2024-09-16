package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"video-downloader-server/internal/config"
	"video-downloader-server/internal/delivery/handlers"
	"video-downloader-server/internal/repository"
	"video-downloader-server/service"
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
	repo := repository.NewRepository(db)

	s := service.NewService(repo)
	v := validator.New()

	r := chi.NewRouter()
	h := handlers.NewHandler(s, v, cfg.ExtensionURL)
	h.RegisterRoutes(r)

	log.Infof(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
