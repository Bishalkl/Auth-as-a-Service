package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/database"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/handlers"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/repositories"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/services"
	"gorm.io/gorm"
)

type Handlers struct {
	Auth *handlers.AuthHandler
}

type AppContainer struct {
	DB           *gorm.DB
	RedisService database.RedisService
	Handler      Handlers
}

func InitalizeApp() (*AppContainer, error) {
	// Load configuration
	log.Println("🔧 Loading configuration...")
	configs.LoadEnv()

	// Connect to the PostgreSQL db
	log.Println("💾 Connecting to the database...")
	dbService := database.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// Connect to Redis db
	log.Println("🔗 Connecting to Redis...")
	ctx := context.Background()
	redisService, err := database.NewRedisService(ctx)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Redis: %w", err)

	}

	// // AutoMigrate all models (TEMP use only)
	// log.Println("🔄 Running auto migration for models...")
	// err = migrateAutoModels(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("❌ Failed to auto-migrate authentication models: %w", err)
	// }

	// Initalize repositories
	log.Println("📦 Initializing repositories...")
	authRepo := repositories.NewAuthRepository(db)

	// Initalize service
	log.Println("🧠 Initializing services...")
	authService := services.NewAuthService(authRepo)

	// initalize handler
	log.Println("🧠 Initializing services...")
	authHandler := handlers.NewAuthHandler(authService)

	// Return app container
	return &AppContainer{
		DB:           db,
		RedisService: redisService,
		Handler: Handlers{
			Auth: authHandler,
		},
	}, nil

}

// // function for automigrate
// func migrateAutoModels(db *gorm.DB) error {
// 	// Auto migrate
// 	err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.UserRole{})
// 	if err != nil {
// 		return fmt.Errorf("❌ Error migrating auth models: %w", err)
// 	}
// 	log.Println("✅ Auth models migrated successfully")
// 	return nil
// }
