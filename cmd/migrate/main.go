package main

import (
	"flag"
	"log/slog"
	"os"
	"rest_waka/internal/favorites"
	"rest_waka/internal/users"

	"rest_waka/config"
	"rest_waka/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	db, err := gorm.Open(postgres.Open(cfg.DBAddress), &gorm.Config{})
	if err != nil {
		log.Error("cannot open db", "error", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&models.WakaModel{},
		&users.User{},
		&favorites.Favorite{},
	); err != nil {
		log.Error("automigrate failed", "error", err)
		os.Exit(1)
	}

	log.Info("automigrate done")
}
