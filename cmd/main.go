package main

import (
	"context"
	"errors"
	"flag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rest_waka/config"
	"rest_waka/internal/models"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	// DB (GORM)
	gormDB, err := gorm.Open(postgres.Open(cfg.DBAddress), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	repo := models.NewGormRepository(gormDB)
	svc := models.NewService(repo)

	mux := http.NewServeMux()
	models.NewModelsHandler(mux, models.HandlerDeps{Service: svc})

	server := http.Server{
		Addr:              cfg.HTTP.Address,
		ReadHeaderTimeout: cfg.HTTP.Timeout,
		Handler:           mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("http starting server", "address", server.Addr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down server")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped unexpectedly", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	_ = sqlDB.Close()
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}
