package config

import (
	"fmt"
	"log"
	"go-clean-api/entity"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func NewDatabase() *gorm.DB {

	cfg := LoadDBConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Bangkok",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running database migrations Table User...")
	err = db.AutoMigrate(
		&entity.User{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database connection successful!")
	return db
}
