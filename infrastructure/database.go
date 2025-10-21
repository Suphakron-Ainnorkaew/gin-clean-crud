package infrastructure

import (
	"fmt"
	"go-clean-api/config"
	shopDomain "go-clean-api/internal/shop/domain"
	userDomain "go-clean-api/internal/user/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(cfg *config.DBConfig) *gorm.DB {

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
		&userDomain.User{},
		&shopDomain.Shop{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	log.Println("Database connection successful!")
	return db
}
