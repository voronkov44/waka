package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	modelspkg "rest_waka/internal/models"
	"rest_waka/pkg/modelsutil"
)

type BackfillConfig struct {
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS" env-required:"true"`
}

func mustLoadBackfillConfig(configPath string) BackfillConfig {
	var cfg BackfillConfig

	if configPath == "" {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read env: %s", err)
		}
		return cfg
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		var pe *os.PathError
		if errors.As(err, &pe) {
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				log.Fatalf("cannot read env: %s", err)
			}
			return cfg
		}
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}

	return cfg
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.backfill.yaml", "config file")
	flag.Parse()

	cfg := mustLoadBackfillConfig(configPath)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := gorm.Open(postgres.Open(cfg.DBAddress), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect db", "err", err)
		os.Exit(1)
	}

	var recs []modelspkg.WakaModel
	if err := db.Find(&recs).Error; err != nil {
		log.Error("failed to load models", "err", err)
		os.Exit(1)
	}

	updated := 0

	for _, rec := range recs {
		flavors, err := modelsutil.UnmarshalFlavors(rec.Flavors)
		if err != nil {
			log.Error("failed to unmarshal flavors", "model_id", rec.ID, "err", err)
			continue
		}

		cleaned := modelsutil.CleanupFlavors(flavors)
		if equalStrings(flavors, cleaned) {
			continue
		}

		raw, err := modelsutil.MarshalFlavors(cleaned)
		if err != nil {
			log.Error("failed to marshal flavors", "model_id", rec.ID, "err", err)
			continue
		}

		if err := db.Model(&modelspkg.WakaModel{}).
			Where("id = ?", rec.ID).
			Update("flavors", raw).Error; err != nil {
			log.Error("failed to update model", "model_id", rec.ID, "err", err)
			continue
		}

		updated++
		log.Info("model flavors normalized", "model_id", rec.ID, "flavors", cleaned)
	}

	fmt.Printf("done, updated=%d\n", updated)
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
