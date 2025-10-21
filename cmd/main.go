// cmd/server/main.go
package main

import (
	"go-clean-api/config"
	"go-clean-api/infrastructure"
	"go-clean-api/internal/shop"
	"go-clean-api/internal/user"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadDBConfig()

	db := infrastructure.NewDatabase(cfg)

	userRepo := user.NewGormUserRepository(db)
	userUsecase := user.NewUserUsecase(userRepo)

	shopRepo := shop.NewGormShopRepository(db)
	shopUsecase := shop.NewShopUsecase(shopRepo)

	e := echo.New()

	user.NewUserHandler(e, userUsecase)

	shop.NewShopHandler(e, shopUsecase)

	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
