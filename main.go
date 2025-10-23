package main

import (
	"context"
	"go-clean-api/config"
	"go-clean-api/domain"
	"go-clean-api/user/delivery"
	"go-clean-api/user/repository"
	"go-clean-api/user/usecase"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func main() {
	// Connect to PostgreSQL with AutoMigrate
	db := config.NewDatabase()

	// Connect to Real Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379", // ใช้ IPv4 แทน IPv6
	})

	// Test Redis connection and initialize cache repository
	ctx := context.Background()
	var userCacheRepo domain.UserCacheRepository

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("❌ Redis connection failed: %v", err)
		log.Println("Falling back to Mock Redis...")
		userCacheRepo = repository.NewRedisMockRepository()
	} else {
		log.Printf("✅ Redis connected successfully: %s", pong)
		userCacheRepo = repository.NewRedisUserRepository(rdb)
	}

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)
	userMessageRepo := repository.NewUserMessageRepository()

	// Initialize usecase
	userUsecase := usecase.NewUserUsecase(userRepo, userCacheRepo, userMessageRepo)

	// Start GraphQL server in background
	go startGraphQLServer(userUsecase)

	// Start HTTP server
	startHTTPServer(userUsecase)
}

func startHTTPServer(userUsecase domain.UserUsecase) {
	e := echo.New()

	// Setup routes
	v1 := e.Group("/v1")
	delivery.NewHandler(v1, userUsecase)

	// Start server
	log.Println("🌐 HTTP server starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startGraphQLServer(userUsecase domain.UserUsecase) {
	// Setup GraphQL handler
	h := delivery.NewGraphQLHandler(userUsecase)

	e := echo.New()
	e.POST("/graphql", h.GraphQLHandler)
	e.GET("/graphql", h.PlaygroundHandler)

	log.Println("🎮 GraphQL server starting on :8081")
	if err := e.Start(":8081"); err != nil {
		log.Fatalf("Failed to start GraphQL server: %v", err)
	}
}
