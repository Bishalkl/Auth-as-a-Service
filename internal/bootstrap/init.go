package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/database"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB           *gorm.DB
	RedisService database.RedisService
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
	ctx := context.Background()
	redisService, err := database.NewRedisService(ctx)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Redis: %w", err)

	}

	// Return app container
	return &AppContainer{
		DB:           db,
		RedisService: redisService,
	}, nil

}
