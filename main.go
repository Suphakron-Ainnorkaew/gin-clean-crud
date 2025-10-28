package main

import (
	"fmt"
	"go-clean-api/config"
	"go-clean-api/entity"
	"os"

	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	userDelivery "go-clean-api/user/delivery"
	userRepo "go-clean-api/user/repository"
	userUseCase "go-clean-api/user/usecase"

	shopDelivery "go-clean-api/shop/delivery"
	shopRepo "go-clean-api/shop/repository"
	shopUseCase "go-clean-api/shop/usecase"
)

var runEnv string
var db *gorm.DB

var cfg *config.Config

func init() {
	runEnv = os.Getenv("RUN_ENV")
	if runEnv == "" {
		runEnv = "local"
	}

	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("[init]: Error loading .env file: ", err)
	}

	cfg = config.Load()

	if err := connectDB(); err != nil {
		log.Fatal("[init]: failed to connect to database: ", err)
	}
}

func main() {
	e := echo.New()

	v1 := e.Group("/v1")
	//user
	userUC := userUseCase.NewUserUsecase(
		userRepo.NewPostgresUserRepository(db),
		nil,
		nil,
		cfg.Server.JWTSecret,
	)
	userDelivery.NewHandler(v1, userUC, cfg.Server.JWTSecret)

	//shop
	shopUC := shopUseCase.NewShopUsecase(
		shopRepo.NewPostgresShopRepository(db),
		nil,
		nil,
	)
	// userFetcher will be used by role middleware to load user from DB (source-of-truth)
	userFetcher := func(id uint) (*entity.User, error) {
		return userUC.GetUserByID(id)
	}
	shopDelivery.NewHandler(v1, shopUC, cfg.Server.JWTSecret, userFetcher)

	addr := ":" + cfg.Server.Port
	log.Printf("🌐 starting HTTP server on %s (env=%s)", addr, runEnv)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}

func connectDB() error {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return errors.Wrap(err, "[Main.connectDB]: failed to connect to database")
	}

	createType := `
    DO $$
    BEGIN
      IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_type') THEN
        CREATE TYPE user_type AS ENUM ('general','shop','admin');
      END IF;
    END$$;
    `
	if execErr := db.Exec(createType).Error; execErr != nil {
		return errors.Wrap(execErr, "[Main.connectDB]: failed to create enum type user_type")
	}

	if migrateErr := db.AutoMigrate(&entity.User{}, &entity.Shop{}, &entity.Product{}); migrateErr != nil {
		return errors.Wrap(migrateErr, "[Main.connectDB]: auto migrate failed")
	}

	return nil
}

func NewGormDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	return db, nil
}
