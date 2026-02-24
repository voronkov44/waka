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
	"rest_waka/internal/auth"
	"rest_waka/internal/models"
	"rest_waka/internal/users"
	"rest_waka/pkg/middleware"
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

	// repository
	modelsRepo := models.NewGormRepository(gormDB)
	authRepo := auth.NewGormRepository(gormDB)
	usersRepo := users.NewGormRepository(gormDB)

	// service
	modelsService := models.NewService(modelsRepo)
	authService, err := auth.NewService(authRepo, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
	if err != nil {
		log.Error("Failed to create auth service", "error", err)
	}
	usersService := users.NewService(usersRepo)

	//router
	router := http.NewServeMux()

	models.NewModelsHandler(router, models.HandlerDeps{Service: modelsService})
	auth.NewAuthHandler(router, auth.HandlerDeps{
		Service:   authService,
		JWTSecret: cfg.Auth.JWTSecret,
	})
	users.NewUsersHandler(router, users.HandlerDeps{Service: usersService})

	// Middlewares
	stack := middleware.Chain(
		middleware.Recover(log),
		middleware.Logging(log),
		middleware.CORS,
	)

	server := http.Server{
		Addr:              cfg.HTTP.Address,
		ReadHeaderTimeout: cfg.HTTP.Timeout,
		Handler:           stack(router),
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
