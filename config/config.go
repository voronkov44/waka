package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AuthConfig struct {
	JWTSecret string        `yaml:"jwt_secret" env:"AUTH_JWT_SECRET" env-required:"true"`
	TokenTTL  time.Duration `yaml:"token_ttl" env:"AUTH_TOKEN_TTL" env-default:"24h"`
}

type HTTPConfig struct {
	Address string        `yaml:"address" env:"API_ADDRESS" env-default:":8080"`
	Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}

type S3Config struct {
	Endpoint        string `yaml:"endpoint" env:"S3_ENDPOINT" env-required:"true"`
	PresignEndpoint string `yaml:"presign_endpoint" env:"S3_PRESIGN_ENDPOINT" env-default:""`
	AccessKey       string `yaml:"access_key" env:"S3_ACCESS_KEY" env-required:"true"`
	SecretKey       string `yaml:"secret_key" env:"S3_SECRET_KEY" env-required:"true"`
	Bucket          string `yaml:"bucket" env:"S3_BUCKET" env-required:"true"`
	PublicBaseURL   string `yaml:"public_base_url" env:"S3_PUBLIC_BASE_URL" env-required:"true"`

	UsePresigned bool          `yaml:"use_presigned" env:"S3_USE_PRESIGNED" env-default:"false"`
	PresignTTL   time.Duration `yaml:"presign_ttl" env:"S3_PRESIGN_TTL" env-default:"15m"`
}

type Config struct {
	LogLevel  string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTP      HTTPConfig `yaml:"api_server"`
	DBAddress string     `yaml:"db_address" env:"DB_ADDRESS" env-required:"true"`

	S3   S3Config   `yaml:"s3"`
	Auth AuthConfig `yaml:"auth"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	// если путь пустой - просто env
	if configPath == "" {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read env: %s", err)
		}
		return cfg
	}

	// пробуем файл, если его нет - env
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
