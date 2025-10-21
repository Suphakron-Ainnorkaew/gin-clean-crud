package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" default:"127.0.0.1"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

func LoadDBConfig() *DBConfig {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	var cfg DBConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	return &cfg
}
