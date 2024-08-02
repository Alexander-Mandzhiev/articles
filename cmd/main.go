package main

import (
	"articles/internal/apiserver"
	"articles/internal/apiserver/handlers"
	"articles/internal/config"
	"articles/internal/repository"
	"articles/internal/service"
	"articles/pkg/sl"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	logger := config.SetupLogger(cfg.Env)
	logger.Info("Starting authors app", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to initialize db: %s", sl.Err(err))
	}
	logger.Info("initialize database")
	repo := repository.NewRepository(db)
	logger.Info("initialize repository")
	services := service.NewService(*repo)
	logger.Info("initialize services")
	handlers := handlers.NewHandler(services)
	logger.Info("initialize handlers")

	srv := new(apiserver.APIServer)
	if err := srv.Start(cfg.Address, handlers.InitRouters()); err != nil {
		logger.Error("error occured while running http server: %s", sl.Err(err))
	}
}
