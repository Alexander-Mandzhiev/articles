package main

import (
	"articles/internal/apiserver"
	"articles/internal/apiserver/handlers"
	"articles/internal/config"
	"articles/internal/repository"
	"articles/internal/service"
	"articles/pkg/sl"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		if err := srv.Start(cfg.Address, handlers.InitRouters()); err != nil {
			logger.Error("error occured while running http server: %s", sl.Err(err))
		}
	}()
	logger.Info("starting application")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("application shutdown")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("error occured on server shutdown: %s", sl.Err(err))
	}

	logger.Info("application shutdown")
	if err := db.Close(); err != nil {
		logger.Error("error occured on db connection close: %s", sl.Err(err))
	}
}
