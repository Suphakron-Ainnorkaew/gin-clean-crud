package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Tools    ToolsConfig
}

type ServerConfig struct {
	Port            string
	ContextTimeout  time.Duration
	ShutdownTimeout time.Duration
}

type ToolsConfig struct {
	JWTSecret string
	Logrus    *logrus.Logger
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ContextTimeout:  getDurationEnv("CONTEXT_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			User:     getEnv("DATABASE_USERNAME", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "supha"),
			DBName:   getEnv("DATABASE_NAME", "mydb"),
		},
		Tools: ToolsConfig{
			JWTSecret: getEnv("JWT_SECRET", "secret"),
			Logrus:    logrus.New(),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		} else {
			log.Printf("invalid duration for %s: %v — using default %s", key, err, defaultValue)
		}
	}
	return defaultValue
}

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		ForceColors:     true,
		DisableQuote:    true,
	})

	return log
}
