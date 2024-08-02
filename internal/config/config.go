package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"development" env:"ENV" env-default:"development"`
	DatabaseURL string `yaml:"database_url" env:"DATABASE_URL" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}
type HTTPServer struct {
	Address     string        `yaml:"address" env:"ADDRESS" env-default:"localhost:3000"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s" `
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := fetchConfigFlag()

	if configPath == "" {
		panic("config path to empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigFlag() string {

	var res string

	flag.StringVar(&res, "config", "config/development.yaml", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

const (
	envDevelopment = "development"
	envProduction  = "production"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
