package bootstrap

import (
	"fmt"
	"log"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func InitalizeApp() (*AppContainer, error) {
	// Load configuration
	// Load configuration
	log.Println("🔧 Loading configuration...")
	configs.LoadEnv()

	// Connect to the PostgreSQL
	log.Println("💾 Connecting to the database...")
	dbService := database.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// Connect to Redis
	log.Println("🔗 Connecting to Redis...")
	redisService := database.NewRedisService()
	redisClient := redisService.GetClient()

	// Return app container
	return &AppContainer{
		DB:    db,
		Redis: redisClient,
	}, nil

}
