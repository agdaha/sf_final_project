package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string `yaml:"env" env-default:"devel"`
	LogLevel           string `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"`
	HTTPServer         `yaml:"http_server"`
	NewsServiceURL     string `yaml:"news_service_URL" env-default:"localhost:8081/api/news"`
	CommentsServiceURL string `yaml:"comments_service_URL" env-default:"localhost:8082/api/comments"`
	CensorServiceURL   string `yaml:"censor_service_URL" env-default:"localhost:8083/api/censor"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"API_ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func New() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
