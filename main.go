package main

import (
	"fmt"
	"go-clean-api/config"
	"go-clean-api/entity"
	"go-clean-api/middleware"
	"os"

	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	userDelivery "go-clean-api/feature/user/delivery"
	userRepo "go-clean-api/feature/user/repository"
	userUseCase "go-clean-api/feature/user/usecase"

	shopDelivery "go-clean-api/feature/shop/delivery"
	shopRepo "go-clean-api/feature/shop/repository"
	shopUseCase "go-clean-api/feature/shop/usecase"

	courierDelivery "go-clean-api/feature/courier/delivery"
	courierRepo "go-clean-api/feature/courier/repository"
	courierUseCase "go-clean-api/feature/courier/usecase"

	orderDelivery "go-clean-api/feature/order/delivery"
	orderRepo "go-clean-api/feature/order/repository"
	orderUseCase "go-clean-api/feature/order/usecase"

	productDelivery "go-clean-api/feature/product/delivery"
	productRepo "go-clean-api/feature/product/repository"
	productUseCase "go-clean-api/feature/product/usecase"
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

	v1Public := e.Group("/v1")
	v1Auth := e.Group("/v1", middleware.NewJWTAuth(cfg.Tools.JWTSecret))

	//user
	userUC := userUseCase.NewUserUsecase(
		userRepo.NewPostgresUserRepository(db),
		cfg.Tools.JWTSecret,
	)

	userH := userDelivery.NewHandler(userUC, cfg.Tools)
	userDelivery.RegisterAuthUserRoutes(v1Auth, userH)
	userDelivery.RegisterPublicUserRoutes(v1Public, userH)

	//shop
	shopUC := shopUseCase.NewShopUsecase(
		shopRepo.NewPostgresShopRepository(db),
	)
	userFetcher := func(id uint) (*entity.User, error) {
		return userUC.GetUserByID(id)
	}
	shopDelivery.NewHandler(v1Auth, shopUC, cfg.Tools)

	// courier
	courierUC := courierUseCase.NewCourierUsecase(
		courierRepo.NewPostgresCourierRepository(db),
	)
	courierDelivery.NewHandler(v1Auth, courierUC)

	// product
	productUC := productUseCase.NewProductUsecase(
		productRepo.NewPostgresProductRepository(db),
		shopRepo.NewPostgresShopRepository(db),
	)
	productDelivery.NewHandler(v1Auth, productUC, cfg.Tools.JWTSecret, userFetcher)

	// order
	orderUC := orderUseCase.NewOrderUsecase(
		orderRepo.NewPostgresOrderRepository(db),
		shopRepo.NewPostgresShopRepository(db),
		courierRepo.NewPostgresCourierRepository(db),
		userRepo.NewPostgresUserRepository(db),
		productRepo.NewPostgresProductRepository(db),
		cfg.Tools,
	)
	orderDelivery.NewHandler(v1Auth, orderUC, cfg.Tools, userFetcher)

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
	return nil
}
