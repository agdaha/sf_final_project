package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"devel"`
	LogLevel   string `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"`
	HTTPServer `yaml:"http_server"`
	DB         DB
}

func (c Config) DbUrl() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
		c.DB.Host,
		c.DB.Name,
		c.DB.User,
		c.DB.Pass,
		c.DB.Port,
	)
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"API_ADDRESS" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type DB struct {
	Host string `yaml:"db_host" env:"DB_HOST"`
	Port int    `yaml:"db_port" env:"DB_PORT"`
	Name string `yaml:"db_name" env:"DB_NAME"`
	User string `yaml:"db_user" env:"DB_USER"`
	Pass string `yaml:"db_pass" env:"DB_PASSWORD"`
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
